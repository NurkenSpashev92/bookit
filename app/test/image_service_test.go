package test

import (
	"testing"

	"github.com/nurkenspashev92/bookit/internal/services"
)

func TestImageService_MaxImagesConstant(t *testing.T) {
	// Verify the constant is accessible and correct
	if services.ErrMaxImagesExceeded == nil {
		t.Fatal("ErrMaxImagesExceeded should not be nil")
	}
	if services.ErrMaxImagesExceeded.Error() != "maximum 15 images allowed" {
		t.Errorf("unexpected error message: %s", services.ErrMaxImagesExceeded.Error())
	}
}

func TestServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"slug exists", services.ErrSlugExists, "slug already exists"},
		{"max images", services.ErrMaxImagesExceeded, "maximum 15 images allowed"},
		{"image not found", services.ErrImageNotFound, "image not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("got %q, want %q", tt.err.Error(), tt.msg)
			}
		})
	}
}
