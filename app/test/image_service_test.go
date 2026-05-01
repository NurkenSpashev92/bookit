package test

import (
	"testing"

	propertysvc "github.com/nurkenspashev92/bookit/internal/property/service"
)

func TestImageService_MaxImagesConstant(t *testing.T) {
	// Verify the constant is accessible and correct
	if propertysvc.ErrMaxImagesExceeded == nil {
		t.Fatal("ErrMaxImagesExceeded should not be nil")
	}
	if propertysvc.ErrMaxImagesExceeded.Error() != "maximum 15 images allowed" {
		t.Errorf("unexpected error message: %s", propertysvc.ErrMaxImagesExceeded.Error())
	}
}

func TestServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"slug exists", propertysvc.ErrSlugExists, "slug already exists"},
		{"max images", propertysvc.ErrMaxImagesExceeded, "maximum 15 images allowed"},
		{"image not found", propertysvc.ErrImageNotFound, "image not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("got %q, want %q", tt.err.Error(), tt.msg)
			}
		})
	}
}
