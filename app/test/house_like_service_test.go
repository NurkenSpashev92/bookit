package test

import (
	"context"
	"testing"

	"github.com/nurkenspashev92/bookit/internal/services"
)

func TestHouseLikeService_Like(t *testing.T) {
	repo := newMockHouseLikeRepo()
	svc := services.NewHouseLikeService(repo)

	resp, err := svc.Like(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Liked {
		t.Error("expected liked=true")
	}
	if resp.LikeCount != 1 {
		t.Errorf("count = %d, want 1", resp.LikeCount)
	}

	// Like again — idempotent
	resp2, _ := svc.Like(context.Background(), 1, 10)
	if resp2.LikeCount != 1 {
		t.Errorf("double like count = %d, want 1", resp2.LikeCount)
	}
}

func TestHouseLikeService_Unlike(t *testing.T) {
	repo := newMockHouseLikeRepo()
	svc := services.NewHouseLikeService(repo)

	svc.Like(context.Background(), 1, 10)

	resp, err := svc.Unlike(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Liked {
		t.Error("expected liked=false")
	}
	if resp.LikeCount != 0 {
		t.Errorf("count = %d, want 0", resp.LikeCount)
	}
}

func TestHouseLikeService_Status(t *testing.T) {
	repo := newMockHouseLikeRepo()
	svc := services.NewHouseLikeService(repo)

	// Not liked
	resp, _ := svc.Status(context.Background(), 1, 10)
	if resp.Liked {
		t.Error("should not be liked")
	}

	// Like then check
	svc.Like(context.Background(), 1, 10)
	resp2, _ := svc.Status(context.Background(), 1, 10)
	if !resp2.Liked {
		t.Error("should be liked")
	}
	if resp2.LikeCount != 1 {
		t.Errorf("count = %d, want 1", resp2.LikeCount)
	}
}

func TestHouseLikeService_MultipleLikes(t *testing.T) {
	repo := newMockHouseLikeRepo()
	svc := services.NewHouseLikeService(repo)

	svc.Like(context.Background(), 1, 10)
	svc.Like(context.Background(), 2, 10)
	svc.Like(context.Background(), 3, 10)

	resp, _ := svc.Status(context.Background(), 1, 10)
	if resp.LikeCount != 3 {
		t.Errorf("count = %d, want 3", resp.LikeCount)
	}
}
