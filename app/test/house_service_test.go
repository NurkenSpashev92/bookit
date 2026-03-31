package test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
	"github.com/nurkenspashev92/bookit/pkg/cache"
)

type mockHouseRepo struct {
	houses map[string]models.House
	slugs  map[string]bool
	nextID int
}

func newMockHouseRepo() *mockHouseRepo {
	return &mockHouseRepo{
		houses: make(map[string]models.House),
		slugs:  make(map[string]bool),
		nextID: 1,
	}
}

func (m *mockHouseRepo) GetAll(_ context.Context) ([]schemas.HouseListItem, error) {
	var items []schemas.HouseListItem
	for _, h := range m.houses {
		items = append(items, schemas.HouseListItem{ID: h.ID, NameEN: h.NameEN, Slug: h.Slug})
	}
	return items, nil
}

func (m *mockHouseRepo) GetAllPaginated(_ context.Context, _ schemas.HouseFilter, _, _ int) ([]schemas.HouseListItem, int, error) {
	items, err := m.GetAll(nil)
	return items, len(items), err
}

func (m *mockHouseRepo) GetByOwnerPaginated(_ context.Context, _, _, _ int) ([]schemas.HouseListItem, int, error) {
	return []schemas.HouseListItem{}, 0, nil
}

func (m *mockHouseRepo) GetBySlug(_ context.Context, slug string) (schemas.HouseDetailResponse, error) {
	h, ok := m.houses[slug]
	if !ok {
		return schemas.HouseDetailResponse{}, fmt.Errorf("no rows in result set")
	}
	return schemas.HouseDetailResponse{ID: h.ID, NameEN: h.NameEN, Slug: h.Slug}, nil
}

func (m *mockHouseRepo) RecordView(_ context.Context, _ string, _ *int, _ string) {}

func (m *mockHouseRepo) Create(_ context.Context, req schemas.HouseCreateRequest) (models.House, error) {
	slug := "test-slug"
	if req.Slug != "" {
		slug = req.Slug
	}
	if m.slugs[slug] {
		return models.House{}, fmt.Errorf("slug already exists")
	}
	h := models.House{ID: m.nextID, NameEN: req.NameEN, Slug: slug, OwnerID: req.OwnerID}
	m.houses[slug] = h
	m.slugs[slug] = true
	m.nextID++
	return h, nil
}

func (m *mockHouseRepo) Update(_ context.Context, slug string, _ schemas.HouseUpdateRequest) (models.House, error) {
	h, ok := m.houses[slug]
	if !ok {
		return models.House{}, fmt.Errorf("house with slug '%s' not found", slug)
	}
	return h, nil
}

func (m *mockHouseRepo) Delete(_ context.Context, slug string) error {
	if _, ok := m.houses[slug]; ok {
		delete(m.houses, slug)
		delete(m.slugs, slug)
	}
	return nil
}

func (m *mockHouseRepo) SlugExists(_ context.Context, slug string) (bool, error) {
	return m.slugs[slug], nil
}

func TestHouseService_CheckSlug(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo, newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))
	ctx := context.Background()

	available, normalized, err := svc.CheckSlug(ctx, "Beach House")
	if err != nil {
		t.Fatal(err)
	}
	if !available {
		t.Error("slug should be available")
	}
	if normalized != "beach-house" {
		t.Errorf("normalized = %q, want beach-house", normalized)
	}

	// Create a house with that slug
	repo.slugs["beach-house"] = true

	available2, _, _ := svc.CheckSlug(ctx, "Beach House")
	if available2 {
		t.Error("slug should NOT be available after creation")
	}
}

func TestHouseService_Create_SetsOwnerID(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo, newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))

	house, err := svc.Create(context.Background(), schemas.HouseCreateRequest{
		NameEN: "Test", Slug: "test-house",
	}, 42)
	if err != nil {
		t.Fatal(err)
	}
	if house.OwnerID != 42 {
		t.Errorf("OwnerID = %d, want 42", house.OwnerID)
	}
}

func TestHouseService_Create_SlugExists(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo, newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))
	ctx := context.Background()

	svc.Create(ctx, schemas.HouseCreateRequest{NameEN: "A", Slug: "dup"}, 1)

	_, err := svc.Create(ctx, schemas.HouseCreateRequest{NameEN: "B", Slug: "dup"}, 2)
	if err == nil {
		t.Fatal("expected error for duplicate slug")
	}
	if !errors.Is(err, services.ErrSlugExists) {
		t.Errorf("expected ErrSlugExists, got %v", err)
	}
}

func TestHouseService_GetBySlug(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo, newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))

	svc.Create(context.Background(), schemas.HouseCreateRequest{NameEN: "Beach", Slug: "beach"}, 1)

	house, err := svc.GetBySlug(context.Background(), "beach", 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if house.NameEN != "Beach" {
		t.Errorf("NameEN = %q", house.NameEN)
	}
}

func TestHouseService_GetBySlug_NotFound(t *testing.T) {
	svc := services.NewHouseService(newMockHouseRepo(), newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))
	_, err := svc.GetBySlug(context.Background(), "nope", 0, "")
	if err == nil {
		t.Error("expected error")
	}
}

func TestHouseService_Delete(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo, newMockHouseLikeRepo(), cache.New(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute))

	svc.Create(context.Background(), schemas.HouseCreateRequest{NameEN: "Del", Slug: "del"}, 1)

	err := svc.Delete(context.Background(), "del")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.GetBySlug(context.Background(), "del", 0, "")
	if err == nil {
		t.Error("expected not found after delete")
	}
}
