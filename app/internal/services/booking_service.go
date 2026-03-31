package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type BookingService struct {
	repository *repositories.BookingRepository
}

func NewBookingService(repo *repositories.BookingRepository) *BookingService {
	return &BookingService{repository: repo}
}

func (s *BookingService) Create(ctx context.Context, userID int, req schemas.BookingCreateRequest) (schemas.BookingResponse, error) {
	houseID, pricePerDay, err := s.repository.GetHouseBySlug(ctx, req.HouseSlug)
	if err != nil {
		return schemas.BookingResponse{}, err
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	if !end.After(start) {
		return schemas.BookingResponse{}, fmt.Errorf("end_date must be after start_date")
	}
	if start.Before(time.Now().Truncate(24 * time.Hour)) {
		return schemas.BookingResponse{}, fmt.Errorf("start_date cannot be in the past")
	}

	overlap, err := s.repository.HasOverlap(ctx, houseID, req.StartDate, req.EndDate)
	if err != nil {
		return schemas.BookingResponse{}, err
	}
	if overlap {
		return schemas.BookingResponse{}, ErrBookingOverlap
	}

	days := int(end.Sub(start).Hours() / 24)
	totalPrice := days * pricePerDay

	bookingID, err := s.repository.Create(ctx, houseID, userID, req.GuestCount, totalPrice, req.StartDate, req.EndDate, req.Message)
	if err != nil {
		return schemas.BookingResponse{}, err
	}

	return s.repository.GetByID(ctx, bookingID)
}

func (s *BookingService) GetMyBookings(ctx context.Context, userID int) ([]schemas.BookingResponse, error) {
	return s.repository.GetUserBookings(ctx, userID)
}

func (s *BookingService) GetOwnerBookings(ctx context.Context, ownerID int) ([]schemas.BookingResponse, error) {
	return s.repository.GetOwnerBookings(ctx, ownerID)
}

func (s *BookingService) GetByID(ctx context.Context, id, userID int) (schemas.BookingResponse, error) {
	booking, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return booking, err
	}

	ownerID, _ := s.repository.GetOwnerIDByBooking(ctx, id)
	if booking.GuestID != userID && ownerID != userID {
		return schemas.BookingResponse{}, fmt.Errorf("booking not found")
	}

	return booking, nil
}

func (s *BookingService) UpdateStatus(ctx context.Context, bookingID, userID int, status string) error {
	ownerID, err := s.repository.GetOwnerIDByBooking(ctx, bookingID)
	if err != nil {
		return err
	}

	bookingUserID, err := s.repository.GetBookingUserID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Owner can confirm/reject, guest can cancel
	if status == "cancelled" && bookingUserID != userID {
		return fmt.Errorf("only the guest can cancel a booking")
	}
	if (status == "confirmed" || status == "rejected") && ownerID != userID {
		return fmt.Errorf("only the owner can confirm or reject a booking")
	}

	return s.repository.UpdateStatus(ctx, bookingID, status)
}
