package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/booking/model"
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
		return 0, 0, model.ErrHouseNotFound
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
		return b, model.ErrBookingNotFound
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

func (r *BookingRepository) GetUserBookings(ctx context.Context, userID int, search string) ([]schema.BookingResponse, error) {
	return r.queryBookings(ctx, "b.user_id=$1", userID, search)
}

func (r *BookingRepository) GetOwnerBookings(ctx context.Context, ownerID int, search string) ([]schema.BookingResponse, error) {
	return r.queryBookings(ctx, "h.owner_id=$1", ownerID, search)
}

func (r *BookingRepository) GetUserBookingsPaginated(ctx context.Context, userID int, search string, limit, offset int) ([]schema.BookingResponse, int, error) {
	return r.queryBookingsPaginated(ctx, "b.user_id=$1", userID, search, limit, offset)
}

func (r *BookingRepository) GetOwnerBookingsPaginated(ctx context.Context, ownerID int, search string, limit, offset int) ([]schema.BookingResponse, int, error) {
	return r.queryBookingsPaginated(ctx, "h.owner_id=$1", ownerID, search, limit, offset)
}

const bookingSelectColumns = `b.id, b.house_id, h.slug, h.name_en, h.name_kz, h.name_ru,
		   owner.id, CONCAT(owner.first_name, ' ', owner.last_name), owner.email, owner.phone_number,
		   guest.id, CONCAT(guest.first_name, ' ', guest.last_name), guest.email, guest.phone_number,
		   b.start_date, b.end_date, b.guest_count, b.status, b.total_price,
		   COALESCE(b.message,''), b.created_at, b.updated_at`

const bookingFromJoins = `FROM bookings b
		INNER JOIN houses h ON h.id = b.house_id
		INNER JOIN users owner ON owner.id = h.owner_id
		INNER JOIN users guest ON guest.id = b.user_id`

func bookingSearchClause(n int) string {
	p := fmt.Sprintf("$%d", n)
	return "(h.name_en ILIKE " + p + " OR h.name_kz ILIKE " + p + " OR h.name_ru ILIKE " + p +
		" OR (owner.first_name || ' ' || owner.last_name) ILIKE " + p +
		" OR (guest.first_name || ' ' || guest.last_name) ILIKE " + p +
		" OR owner.email ILIKE " + p + " OR guest.email ILIKE " + p +
		" OR b.status::text ILIKE " + p + ")"
}

func (r *BookingRepository) queryBookings(ctx context.Context, where string, arg interface{}, search string) ([]schema.BookingResponse, error) {
	args := []any{arg}
	searchSQL := ""
	if search != "" {
		args = append(args, "%"+search+"%")
		searchSQL = " AND " + bookingSearchClause(len(args))
	}

	query := fmt.Sprintf(`
		SELECT %s
		%s
		WHERE %s%s
		ORDER BY b.created_at DESC
	`, bookingSelectColumns, bookingFromJoins, where, searchSQL)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBookingRows(rows)
}

func (r *BookingRepository) queryBookingsPaginated(ctx context.Context, where string, arg interface{}, search string, limit, offset int) ([]schema.BookingResponse, int, error) {
	args := []any{arg}
	searchSQL := ""
	if search != "" {
		args = append(args, "%"+search+"%")
		searchSQL = " AND " + bookingSearchClause(len(args))
	}

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) %s WHERE %s%s`, bookingFromJoins, where, searchSQL)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT %s
		%s
		WHERE %s%s
		ORDER BY b.created_at DESC
		LIMIT $%d OFFSET $%d
	`, bookingSelectColumns, bookingFromJoins, where, searchSQL, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items, err := scanBookingRows(rows)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func scanBookingRows(rows pgx.Rows) ([]schema.BookingResponse, error) {
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
	return items, rows.Err()
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
		return 0, model.ErrBookingNotFound
	}
	return ownerID, nil
}

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
		return 0, model.ErrBookingNotFound
	}
	return userID, nil
}
