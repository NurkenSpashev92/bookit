package service

import (
	"context"
	"time"

	"github.com/nurkenspashev92/bookit/internal/booking/model"
	"github.com/nurkenspashev92/bookit/internal/booking/schema"
)

type BookingRepository interface {
	GetHouseBySlug(ctx context.Context, slug string) (int, int, error)
	HasOverlap(ctx context.Context, houseID int, startDate, endDate string) (bool, error)
	Create(ctx context.Context, houseID, userID, guestCount, totalPrice int, startDate, endDate, message string) (int, error)
	GetByID(ctx context.Context, id int) (schema.BookingResponse, error)
	GetUserBookings(ctx context.Context, userID int) ([]schema.BookingResponse, error)
	GetOwnerBookings(ctx context.Context, ownerID int) ([]schema.BookingResponse, error)
	UpdateStatus(ctx context.Context, bookingID int, status string) error
	GetOwnerIDByBooking(ctx context.Context, bookingID int) (int, error)
	GetBookingUserID(ctx context.Context, bookingID int) (int, error)
}

type BookingService struct {
	repository BookingRepository
}

func NewBookingService(repo BookingRepository) *BookingService {
	return &BookingService{repository: repo}
}

func (s *BookingService) Create(ctx context.Context, userID int, req schema.BookingCreateRequest) (schema.BookingResponse, error) {
	houseID, pricePerDay, err := s.repository.GetHouseBySlug(ctx, req.HouseSlug)
	if err != nil {
		return schema.BookingResponse{}, err
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	if !end.After(start) {
		return schema.BookingResponse{}, model.ErrInvalidDateRange
	}
	if start.Before(time.Now().Truncate(24 * time.Hour)) {
		return schema.BookingResponse{}, model.ErrStartDateInPast
	}

	overlap, err := s.repository.HasOverlap(ctx, houseID, req.StartDate, req.EndDate)
	if err != nil {
		return schema.BookingResponse{}, err
	}
	if overlap {
		return schema.BookingResponse{}, model.ErrBookingOverlap
	}

	days := int(end.Sub(start).Hours() / 24)
	totalPrice := days * pricePerDay

	bookingID, err := s.repository.Create(ctx, houseID, userID, req.GuestCount, totalPrice, req.StartDate, req.EndDate, req.Message)
	if err != nil {
		return schema.BookingResponse{}, err
	}

	return s.repository.GetByID(ctx, bookingID)
}

func (s *BookingService) GetMyBookings(ctx context.Context, userID int) ([]schema.BookingResponse, error) {
	return s.repository.GetUserBookings(ctx, userID)
}

func (s *BookingService) GetOwnerBookings(ctx context.Context, ownerID int) ([]schema.BookingResponse, error) {
	return s.repository.GetOwnerBookings(ctx, ownerID)
}

func (s *BookingService) GetByID(ctx context.Context, id, userID int) (schema.BookingResponse, error) {
	booking, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return booking, err
	}

	ownerID, _ := s.repository.GetOwnerIDByBooking(ctx, id)
	if booking.GuestID != userID && ownerID != userID {
		return schema.BookingResponse{}, model.ErrBookingNotFound
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

	if status == "cancelled" && bookingUserID != userID {
		return model.ErrCancelNotAllowed
	}
	if (status == "confirmed" || status == "rejected") && ownerID != userID {
		return model.ErrDecisionNotAllowed
	}

	return s.repository.UpdateStatus(ctx, bookingID, status)
}
