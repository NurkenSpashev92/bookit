package test

import (
	"context"
	"testing"

	contentschema "github.com/nurkenspashev92/bookit/internal/content/schema"
	contentsvc "github.com/nurkenspashev92/bookit/internal/content/service"
)

func TestFAQService_CRUD(t *testing.T) {
	repo := newMockFAQRepo()
	svc := contentsvc.NewFAQService(repo)
	ctx := context.Background()

	// Create
	req := contentschema.FAQCreateRequest{
		QuestionKz: "Q?", AnswerKz: "A",
		QuestionRu: "В?", AnswerRu: "О",
		QuestionEn: "Q?", AnswerEn: "A",
	}
	faq, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if faq.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if faq.QuestionKz != "Q?" {
		t.Errorf("QuestionKz = %q", faq.QuestionKz)
	}

	// GetByID
	got, err := svc.GetByID(ctx, faq.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.QuestionEn != "Q?" {
		t.Errorf("QuestionEn = %q", got.QuestionEn)
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
	newQ := "Updated?"
	updated, err := svc.Update(ctx, faq.ID, contentschema.FAQUpdateRequest{QuestionKz: &newQ})
	if err != nil {
		t.Fatal(err)
	}
	if updated.QuestionKz != "Updated?" {
		t.Errorf("QuestionKz = %q", updated.QuestionKz)
	}

	// Delete
	if err := svc.Delete(ctx, faq.ID); err != nil {
		t.Fatal(err)
	}

	// Verify deleted
	_, err = svc.GetByID(ctx, faq.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestFAQService_GetByID_NotFound(t *testing.T) {
	repo := newMockFAQRepo()
	svc := contentsvc.NewFAQService(repo)

	_, err := svc.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent FAQ")
	}
}

func TestFAQService_Delete_NotFound(t *testing.T) {
	repo := newMockFAQRepo()
	svc := contentsvc.NewFAQService(repo)

	err := svc.Delete(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent FAQ")
	}
}
