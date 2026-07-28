package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrFAQNotFound     = shared.NotFound("FAQ not found")
	ErrInquiryNotFound = shared.NotFound("Inquiry not found")
)
