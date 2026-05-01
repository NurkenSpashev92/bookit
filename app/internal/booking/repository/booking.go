package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/booking/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) GetHouseBySlug(ctx context.Context, slug string) (int, int, error) {
	var houseID, price int
	err := r.db.QueryRow(ctx, `SELECT id, price FROM houses WHERE slug=$1`, slug).Scan(&houseID, &price)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, fmt.Errorf("house not found")
	}
	return houseID, price, err
}

func (r *BookingRepository) HasOverlap(ctx context.Context, houseID int, startDate, endDate string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE house_id=$1 AND status IN ('pending','confirmed')
			  AND start_date < $3::date AND end_date > $2::date
		)
	`, houseID, startDate, endDate).Scan(&exists)
	return exists, err
}

func (r *BookingRepository) Create(ctx context.Context, houseID, userID, guestCount, totalPrice int, startDate, endDate, message string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `
		INSERT INTO bookings (house_id, user_id, start_date, end_date, guest_count, total_price, message)
		VALUES ($1, $2, $3::date, $4::date, $5, $6, $7)
		RETURNING id
	`, houseID, userID, startDate, endDate, guestCount, totalPrice, message).Scan(&id)
	return id, err
}

func (r *BookingRepository) GetByID(ctx context.Context, id int) (schema.BookingResponse, error) {
	return r.scanBooking(ctx, "b.id=$1", id)
}

func (r *BookingRepository) scanBooking(ctx context.Context, where string, arg interface{}) (schema.BookingResponse, error) {
	var b schema.BookingResponse
	var sd, ed, createdAt, updatedAt time.Time
	var ownerPhone, guestPhone *string

	query := fmt.Sprintf(`
		SELECT b.id, b.house_id, h.slug, h.name_en, h.name_kz, h.name_ru,
			   owner.id, CONCAT(owner.first_name, ' ', owner.last_name), owner.email, owner.phone_number,
			   guest.id, CONCAT(guest.first_name, ' ', guest.last_name), guest.email, guest.phone_number,
			   b.start_date, b.end_date, b.guest_count, b.status, b.total_price,
			   COALESCE(b.message,''), b.created_at, b.updated_at
		FROM bookings b
		INNER JOIN houses h ON h.id = b.house_id
		INNER JOIN users owner ON owner.id = h.owner_id
		INNER JOIN users guest ON guest.id = b.user_id
		WHERE %s
	`, where)

	err := r.db.QueryRow(ctx, query, arg).Scan(
		&b.ID, &b.HouseID, &b.HouseSlug, &b.HouseNameEN, &b.HouseNameKZ, &b.HouseNameRU,
		&b.OwnerID, &b.OwnerFullName, &b.OwnerEmail, &ownerPhone,
		&b.GuestID, &b.GuestFullName, &b.GuestEmail, &guestPhone,
		&sd, &ed, &b.GuestCount, &b.Status, &b.TotalPrice,
		&b.Message, &createdAt, &updatedAt,
	)
	if err != nil {
		return b, fmt.Errorf("booking not found")
	}

	if ownerPhone != nil {
		b.OwnerPhone = *ownerPhone
	}
	if guestPhone != nil {
		b.GuestPhone = *guestPhone
	}
	b.StartDate = sd.Format("2006-01-02")
	b.EndDate = ed.Format("2006-01-02")
	b.CreatedAt = createdAt.Format(time.RFC3339)
	b.UpdatedAt = updatedAt.Format(time.RFC3339)
	return b, nil
}

func (r *BookingRepository) GetUserBookings(ctx context.Context, userID int) ([]schema.BookingResponse, error) {
	return r.queryBookings(ctx, "b.user_id=$1", userID)
}

func (r *BookingRepository) GetOwnerBookings(ctx context.Context, ownerID int) ([]schema.BookingResponse, error) {
	return r.queryBookings(ctx, "h.owner_id=$1", ownerID)
}

func (r *BookingRepository) queryBookings(ctx context.Context, where string, arg interface{}) ([]schema.BookingResponse, error) {
	query := fmt.Sprintf(`
		SELECT b.id, b.house_id, h.slug, h.name_en, h.name_kz, h.name_ru,
			   owner.id, CONCAT(owner.first_name, ' ', owner.last_name), owner.email, owner.phone_number,
			   guest.id, CONCAT(guest.first_name, ' ', guest.last_name), guest.email, guest.phone_number,
			   b.start_date, b.end_date, b.guest_count, b.status, b.total_price,
			   COALESCE(b.message,''), b.created_at, b.updated_at
		FROM bookings b
		INNER JOIN houses h ON h.id = b.house_id
		INNER JOIN users owner ON owner.id = h.owner_id
		INNER JOIN users guest ON guest.id = b.user_id
		WHERE %s
		ORDER BY b.created_at DESC
	`, where)

	rows, err := r.db.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schema.BookingResponse
	for rows.Next() {
		var b schema.BookingResponse
		var sd, ed, createdAt, updatedAt time.Time
		var ownerPhone, guestPhone *string
		if err := rows.Scan(
			&b.ID, &b.HouseID, &b.HouseSlug, &b.HouseNameEN, &b.HouseNameKZ, &b.HouseNameRU,
			&b.OwnerID, &b.OwnerFullName, &b.OwnerEmail, &ownerPhone,
			&b.GuestID, &b.GuestFullName, &b.GuestEmail, &guestPhone,
			&sd, &ed, &b.GuestCount, &b.Status, &b.TotalPrice,
			&b.Message, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		if ownerPhone != nil {
			b.OwnerPhone = *ownerPhone
		}
		if guestPhone != nil {
			b.GuestPhone = *guestPhone
		}
		b.StartDate = sd.Format("2006-01-02")
		b.EndDate = ed.Format("2006-01-02")
		b.CreatedAt = createdAt.Format(time.RFC3339)
		b.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, b)
	}
	return items, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, bookingID int, status string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE bookings SET status=$1, updated_at=NOW() WHERE id=$2
	`, status, bookingID)
	return err
}

func (r *BookingRepository) GetOwnerIDByBooking(ctx context.Context, bookingID int) (int, error) {
	var ownerID int
	err := r.db.QueryRow(ctx, `
		SELECT h.owner_id FROM bookings b
		INNER JOIN houses h ON h.id = b.house_id
		WHERE b.id=$1
	`, bookingID).Scan(&ownerID)
	if err != nil {
		return 0, fmt.Errorf("booking not found")
	}
	return ownerID, nil
}

// GetUserActiveBooking returns the user's currently active booking on a house
// in the property-domain shape (HouseBooking). This implements the property.port.BookingChecker
// port — a downstream dependency from booking back to the property aggregate.
func (r *BookingRepository) GetUserActiveBooking(ctx context.Context, houseID, userID int) (*propertyschema.HouseBooking, error) {
	var b propertyschema.HouseBooking
	var sd, ed time.Time
	err := r.db.QueryRow(ctx, `
		SELECT id, start_date, end_date, status, guest_count, total_price
		FROM bookings
		WHERE house_id=$1 AND user_id=$2 AND status IN ('pending','confirmed')
		ORDER BY created_at DESC LIMIT 1
	`, houseID, userID).Scan(&b.ID, &sd, &ed, &b.Status, &b.GuestCount, &b.TotalPrice)
	if err != nil {
		return nil, nil
	}
	b.StartDate = sd.Format("2006-01-02")
	b.EndDate = ed.Format("2006-01-02")
	return &b, nil
}

func (r *BookingRepository) GetBookingUserID(ctx context.Context, bookingID int) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, `SELECT user_id FROM bookings WHERE id=$1`, bookingID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("booking not found")
	}
	return userID, nil
}
