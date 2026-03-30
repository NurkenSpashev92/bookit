package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
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

func (m *mockHouseRepo) GetBySlug(_ context.Context, slug string) (models.House, error) {
	h, ok := m.houses[slug]
	if !ok {
		return models.House{}, fmt.Errorf("no rows in result set")
	}
	return h, nil
}

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

func (m *mockHouseRepo) Update(_ context.Context, id int, _ schemas.HouseUpdateRequest) (models.House, error) {
	for _, h := range m.houses {
		if h.ID == id {
			return h, nil
		}
	}
	return models.House{}, fmt.Errorf("house with id %d not found", id)
}

func (m *mockHouseRepo) Delete(_ context.Context, id int) error {
	for slug, h := range m.houses {
		if h.ID == id {
			delete(m.houses, slug)
			delete(m.slugs, slug)
			return nil
		}
	}
	return nil
}

func (m *mockHouseRepo) SlugExists(_ context.Context, slug string) (bool, error) {
	return m.slugs[slug], nil
}

func TestHouseService_CheckSlug(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo)
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
	svc := services.NewHouseService(repo)

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
	svc := services.NewHouseService(repo)
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
	svc := services.NewHouseService(repo)

	svc.Create(context.Background(), schemas.HouseCreateRequest{NameEN: "Beach", Slug: "beach"}, 1)

	house, err := svc.GetBySlug(context.Background(), "beach")
	if err != nil {
		t.Fatal(err)
	}
	if house.NameEN != "Beach" {
		t.Errorf("NameEN = %q", house.NameEN)
	}
}

func TestHouseService_GetBySlug_NotFound(t *testing.T) {
	svc := services.NewHouseService(newMockHouseRepo())
	_, err := svc.GetBySlug(context.Background(), "nope")
	if err == nil {
		t.Error("expected error")
	}
}

func TestHouseService_Delete(t *testing.T) {
	repo := newMockHouseRepo()
	svc := services.NewHouseService(repo)

	house, _ := svc.Create(context.Background(), schemas.HouseCreateRequest{NameEN: "Del", Slug: "del"}, 1)

	err := svc.Delete(context.Background(), house.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.GetBySlug(context.Background(), "del")
	if err == nil {
		t.Error("expected not found after delete")
	}
}
