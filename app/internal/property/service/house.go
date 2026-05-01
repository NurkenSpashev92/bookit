package service

import (
	"context"
	"strings"

	"github.com/gosimple/slug"
	"golang.org/x/sync/singleflight"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/port"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/cache"
)

// HouseRepository describes the persistence contract HouseService depends on.
type HouseRepository interface {
	GetAllPaginated(ctx context.Context, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error)
	GetByOwnerPaginated(ctx context.Context, ownerID, limit, offset int) ([]schema.HouseListItem, int, error)
	GetBySlug(ctx context.Context, slug string) (schema.HouseDetailResponse, error)
	RecordView(ctx context.Context, slug string, userID *int, ip string)
	Create(ctx context.Context, req schema.HouseCreateRequest) (model.House, error)
	Update(ctx context.Context, slug string, req schema.HouseUpdateRequest) (model.House, error)
	Delete(ctx context.Context, slug string) error
	SlugExists(ctx context.Context, slug string) (bool, error)
}

type HouseService struct {
	repository     HouseRepository
	likeRepository port.LikeChecker
	bookingRepo    port.BookingChecker
	cache          *cache.Cache
	sf             singleflight.Group
}

func NewHouseService(repo HouseRepository, likeRepo port.LikeChecker, bookingRepo port.BookingChecker, c *cache.Cache) *HouseService {
	return &HouseService{
		repository:     repo,
		likeRepository: likeRepo,
		bookingRepo:    bookingRepo,
		cache:          c,
	}
}

func (s *HouseService) GetAllPaginated(ctx context.Context, userID int, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error) {
	cacheKey := filter.CacheKey(limit, offset)

	// Cache hit — fast path. is_liked is per-user, so always re-attach after copy.
	var cached cachedHouseList
	if s.cache.Get(ctx, cacheKey, &cached) {
		houses := make([]schema.HouseListItem, len(cached.Houses))
		copy(houses, cached.Houses)
		if userID > 0 {
			if likedIDs, err := s.likeRepository.GetUserLikedHouseIDs(ctx, userID); err == nil {
				applyLiked(houses, likedIDs)
			}
		}
		return houses, cached.Total, nil
	}

	// Cache miss — singleflight collapses concurrent identical requests into one DB call.
	v, err, _ := s.sf.Do(cacheKey, func() (any, error) {
		// Check cache again — earlier flight may have populated it while we waited.
		var inner cachedHouseList
		if s.cache.Get(ctx, cacheKey, &inner) {
			return &inner, nil
		}
		houses, total, err := s.repository.GetAllPaginated(ctx, filter, limit, offset)
		if err != nil {
			return nil, err
		}
		result := &cachedHouseList{Houses: houses, Total: total}
		s.cache.Set(ctx, cacheKey, result)
		return result, nil
	})
	if err != nil {
		return nil, 0, err
	}
	loaded := v.(*cachedHouseList)

	houses := make([]schema.HouseListItem, len(loaded.Houses))
	copy(houses, loaded.Houses)

	if userID > 0 {
		if likedIDs, lErr := s.likeRepository.GetUserLikedHouseIDs(ctx, userID); lErr == nil {
			applyLiked(houses, likedIDs)
		}
	}

	return houses, loaded.Total, nil
}

func (s *HouseService) getActiveBooking(ctx context.Context, houseID, userID int) (*schema.HouseBooking, error) {
	if s.bookingRepo == nil {
		return nil, nil
	}
	return s.bookingRepo.GetUserActiveBooking(ctx, houseID, userID)
}

func applyLiked(houses []schema.HouseListItem, likedIDs []int) {
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
	Houses []schema.HouseListItem `json:"houses"`
	Total  int                    `json:"total"`
}

func (s *HouseService) GetMyHouses(ctx context.Context, ownerID, limit, offset int) ([]schema.HouseListItem, int, error) {
	return s.repository.GetByOwnerPaginated(ctx, ownerID, limit, offset)
}

func (s *HouseService) GetBySlug(ctx context.Context, slugVal string, userID int, ip string) (schema.HouseDetailResponse, error) {
	house, err := s.repository.GetBySlug(ctx, slugVal)
	if err != nil {
		return house, err
	}

	var uid *int
	if userID > 0 {
		uid = &userID
	}
	go s.repository.RecordView(context.Background(), slugVal, uid, ip)

	if userID > 0 {
		liked, _, lErr := s.likeRepository.StatusWithCount(ctx, userID, slugVal)
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

func (s *HouseService) Create(ctx context.Context, req schema.HouseCreateRequest, ownerID int) (model.House, error) {
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

func (s *HouseService) Update(ctx context.Context, slugVal string, req schema.HouseUpdateRequest) (model.House, error) {
	house, err := s.repository.Update(ctx, slugVal, req)
	if err == nil {
		s.cache.DeleteByPrefix("houses:")
	}
	return house, err
}

func (s *HouseService) Delete(ctx context.Context, slugVal string) error {
	err := s.repository.Delete(ctx, slugVal)
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
