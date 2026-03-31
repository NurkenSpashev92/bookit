package services

import (
	"context"
	"strings"

	"github.com/gosimple/slug"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/cache"
)

type HouseService struct {
	repository      HouseRepository
	likeRepository  HouseLikeRepository
	bookingRepo     *repositories.BookingRepository
	cache           *cache.Cache
}

func NewHouseService(repo HouseRepository, likeRepo HouseLikeRepository, bookingRepo *repositories.BookingRepository, c *cache.Cache) *HouseService {
	return &HouseService{
		repository:     repo,
		likeRepository: likeRepo,
		bookingRepo:    bookingRepo,
		cache:          c,
	}
}

func (s *HouseService) GetAllPaginated(ctx context.Context, userID int, filter schemas.HouseFilter, limit, offset int) ([]schemas.HouseListItem, int, error) {
	cacheKey := filter.CacheKey(limit, offset)

	// Try cache — houses data is the same for all users, only is_liked differs
	var cached cachedHouseList
	if s.cache.Get(cacheKey, &cached) {
		houses := make([]schemas.HouseListItem, len(cached.Houses))
		copy(houses, cached.Houses)

		if userID > 0 {
			likedIDs, err := s.likeRepository.GetUserLikedHouseIDs(ctx, userID)
			if err == nil {
				applyLiked(houses, likedIDs)
			}
		}

		return houses, cached.Total, nil
	}

	// Cache miss — parallel fetch houses + liked IDs
	g, gctx := errgroup.WithContext(ctx)

	var houses []schemas.HouseListItem
	var total int
	g.Go(func() error {
		var err error
		houses, total, err = s.repository.GetAllPaginated(gctx, filter, limit, offset)
		return err
	})

	var likedIDs []int
	if userID > 0 {
		g.Go(func() error {
			var err error
			likedIDs, err = s.likeRepository.GetUserLikedHouseIDs(gctx, userID)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	s.cache.Set(cacheKey, &cachedHouseList{Houses: houses, Total: total})

	if userID > 0 {
		applyLiked(houses, likedIDs)
	}

	return houses, total, nil
}

func (s *HouseService) getActiveBooking(ctx context.Context, houseID, userID int) (*schemas.HouseBooking, error) {
	if s.bookingRepo == nil {
		return nil, nil
	}
	return s.bookingRepo.GetUserActiveBooking(ctx, houseID, userID)
}

func applyLiked(houses []schemas.HouseListItem, likedIDs []int) {
	if len(likedIDs) == 0 {
		return
	}
	liked := make(map[int]bool, len(likedIDs))
	for _, id := range likedIDs {
		liked[id] = true
	}
	for i := range houses {
		houses[i].IsLiked = liked[houses[i].ID]
	}
}

type cachedHouseList struct {
	Houses []schemas.HouseListItem `json:"houses"`
	Total  int                     `json:"total"`
}

func (s *HouseService) GetMyHouses(ctx context.Context, ownerID, limit, offset int) ([]schemas.HouseListItem, int, error) {
	return s.repository.GetByOwnerPaginated(ctx, ownerID, limit, offset)
}

func (s *HouseService) GetBySlug(ctx context.Context, slug string, userID int, ip string) (schemas.HouseDetailResponse, error) {
	house, err := s.repository.GetBySlug(ctx, slug)
	if err != nil {
		return house, err
	}

	// Async view recording
	var uid *int
	if userID > 0 {
		uid = &userID
	}
	go s.repository.RecordView(context.Background(), slug, uid, ip)

	if userID > 0 {
		liked, _, lErr := s.likeRepository.StatusWithCount(ctx, userID, slug)
		if lErr == nil {
			house.IsLiked = liked
		}

		myBooking, _ := s.getActiveBooking(ctx, house.ID, userID)
		if myBooking != nil {
			house.IsBooked = true
			house.MyBooking = myBooking
		}
	}

	return house, nil
}

func (s *HouseService) Create(ctx context.Context, req schemas.HouseCreateRequest, ownerID int) (models.House, error) {
	req.OwnerID = ownerID
	house, err := s.repository.Create(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "slug already exists") {
			return house, ErrSlugExists
		}
		return house, err
	}
	s.cache.DeleteByPrefix("houses:")
	return house, nil
}

func (s *HouseService) Update(ctx context.Context, slug string, req schemas.HouseUpdateRequest) (models.House, error) {
	house, err := s.repository.Update(ctx, slug, req)
	if err == nil {
		s.cache.DeleteByPrefix("houses:")
	}
	return house, err
}

func (s *HouseService) Delete(ctx context.Context, slug string) error {
	err := s.repository.Delete(ctx, slug)
	if err == nil {
		s.cache.DeleteByPrefix("houses:")
	}
	return err
}

func (s *HouseService) CheckSlug(ctx context.Context, rawSlug string) (bool, string, error) {
	normalized := slug.Make(rawSlug)
	exists, err := s.repository.SlugExists(ctx, normalized)
	if err != nil {
		return false, "", err
	}
	return !exists, normalized, nil
}
