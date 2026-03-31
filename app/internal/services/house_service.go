package services

import (
	"context"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/cache"
)

type HouseService struct {
	repository     HouseRepository
	likeRepository HouseLikeRepository
	cache          *cache.Cache
}

func NewHouseService(repo HouseRepository, likeRepo HouseLikeRepository) *HouseService {
	return &HouseService{
		repository:     repo,
		likeRepository: likeRepo,
		cache:          cache.New(5 * time.Minute),
	}
}

func (s *HouseService) GetAllPaginated(ctx context.Context, userID int, filter schemas.HouseFilter, limit, offset int) ([]schemas.HouseListItem, int, error) {
	cacheKey := filter.CacheKey(limit, offset)

	// Try cache — houses data is the same for all users, only is_liked differs
	var houses []schemas.HouseListItem
	var total int
	var fromCache bool

	if cached, ok := s.cache.Get(cacheKey); ok {
		result := cached.(*cachedHouseList)
		// Copy slice so is_liked overlay doesn't mutate cache
		houses = make([]schemas.HouseListItem, len(result.houses))
		copy(houses, result.houses)
		total = result.total
		fromCache = true
	}

	if !fromCache {
		// Parallel: fetch houses + liked IDs
		g, gctx := errgroup.WithContext(ctx)

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

		s.cache.Set(cacheKey, &cachedHouseList{houses: houses, total: total})

		if userID > 0 {
			applyLiked(houses, likedIDs)
		}

		return houses, total, nil
	}

	// From cache — still need to fetch liked IDs for auth users
	if userID > 0 {
		likedIDs, err := s.likeRepository.GetUserLikedHouseIDs(ctx, userID)
		if err == nil {
			applyLiked(houses, likedIDs)
		}
	}

	return houses, total, nil
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
	houses []schemas.HouseListItem
	total  int
}

func (s *HouseService) GetMyHouses(ctx context.Context, ownerID, limit, offset int) ([]schemas.HouseListItem, int, error) {
	return s.repository.GetByOwnerPaginated(ctx, ownerID, limit, offset)
}

func (s *HouseService) GetBySlug(ctx context.Context, slug string) (models.House, error) {
	return s.repository.GetBySlug(ctx, slug)
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
