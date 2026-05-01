package test

import (
	"context"
	"testing"

	locationschema "github.com/nurkenspashev92/bookit/internal/location/schema"
	locationsvc "github.com/nurkenspashev92/bookit/internal/location/service"
)

func TestCountryService_CRUD(t *testing.T) {
	repo := newMockCountryRepo()
	svc := locationsvc.NewCountryService(repo)
	ctx := context.Background()

	// Create
	country, err := svc.Create(ctx, locationschema.CountryCreateRequest{
		NameKZ: "Қазақстан", NameEN: "Kazakhstan", NameRU: "Казахстан", Code: "KZ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if country.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if country.NameEN != "Kazakhstan" {
		t.Errorf("NameEN = %q", country.NameEN)
	}

	// GetByID
	got, err := svc.GetByID(ctx, country.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Code != "KZ" {
		t.Errorf("Code = %q", got.Code)
	}

	// GetAll
	all, err := svc.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Errorf("count = %d, want 1", len(all))
	}

	// Update
	newName := "Updated"
	updated, err := svc.Update(ctx, country.ID, locationschema.CountryUpdateRequest{NameEN: &newName})
	if err != nil {
		t.Fatal(err)
	}
	if updated.NameEN != "Updated" {
		t.Errorf("NameEN = %q", updated.NameEN)
	}

	// Delete
	if err := svc.Delete(ctx, country.ID); err != nil {
		t.Fatal(err)
	}

	_, err = svc.GetByID(ctx, country.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestCountryService_GetByID_NotFound(t *testing.T) {
	svc := locationsvc.NewCountryService(newMockCountryRepo())
	_, err := svc.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCountryService_Delete_NotFound(t *testing.T) {
	svc := locationsvc.NewCountryService(newMockCountryRepo())
	err := svc.Delete(context.Background(), 999)
	if err == nil {
		t.Error("expected error")
	}
}
