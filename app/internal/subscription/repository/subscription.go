package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/subscription/model"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

const columns = `id, user_id, type, status, start_date, end_date, created_at, updated_at`

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func scan(row interface {
	Scan(dest ...any) error
}, s *model.Subscription) error {
	return row.Scan(&s.ID, &s.UserID, &s.Type, &s.Status, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt)
}

func collect(rows pgx.Rows) ([]model.Subscription, error) {
	defer rows.Close()

	var result []model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := scan(rows, &s); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *SubscriptionRepository) Create(ctx context.Context, s model.Subscription) (model.Subscription, error) {
	var created model.Subscription
	err := scan(r.db.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, type, status, start_date, end_date, created_at, updated_at)
		 VALUES ($1, $2::subscription_type, $3::subscription_status, $4, $5, NOW(), NOW())
		 RETURNING `+columns,
		s.UserID, string(s.Type), string(s.Status), s.StartDate, s.EndDate,
	), &created)
	if err != nil {
		return created, store.MapForeignKey(err, model.ErrUserRefInvalid)
	}
	return created, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int) (model.Subscription, error) {
	var s model.Subscription
	err := scan(r.db.QueryRow(ctx, `SELECT `+columns+` FROM subscriptions WHERE id = $1`, id), &s)
	return s, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) GetLatestByUserID(ctx context.Context, userID int) (model.Subscription, error) {
	var s model.Subscription
	err := scan(r.db.QueryRow(ctx,
		`SELECT `+columns+` FROM subscriptions WHERE user_id = $1 ORDER BY id DESC LIMIT 1`, userID,
	), &s)
	return s, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) GetActiveByUserID(ctx context.Context, userID int) (model.Subscription, error) {
	var s model.Subscription
	err := scan(r.db.QueryRow(ctx,
		`SELECT `+columns+`
		 FROM subscriptions
		 WHERE user_id = $1 AND status = 'active' AND (end_date IS NULL OR end_date > NOW())
		 ORDER BY id DESC LIMIT 1`, userID,
	), &s)
	return s, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) Deactivate(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE subscriptions SET status = 'in_active', updated_at = NOW() WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepository) ListByUserID(ctx context.Context, userID int) ([]model.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+columns+` FROM subscriptions WHERE user_id = $1 ORDER BY id DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list user subscriptions: %w", err)
	}
	return collect(rows)
}

func (r *SubscriptionRepository) ListAll(ctx context.Context) ([]model.Subscription, error) {
	rows, err := r.db.Query(ctx, `SELECT `+columns+` FROM subscriptions ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	return collect(rows)
}

func (r *SubscriptionRepository) ListPaginated(ctx context.Context, limit, offset int) ([]model.Subscription, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM subscriptions`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+columns+` FROM subscriptions ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	result, err := collect(rows)
	return result, total, err
}

func (r *SubscriptionRepository) Update(ctx context.Context, id int, s model.Subscription) (model.Subscription, error) {
	var updated model.Subscription
	err := scan(r.db.QueryRow(ctx,
		`UPDATE subscriptions
		 SET type = $1::subscription_type, status = $2::subscription_status, end_date = $3, updated_at = NOW()
		 WHERE id = $4
		 RETURNING `+columns,
		string(s.Type), string(s.Status), s.EndDate, id,
	), &updated)
	return updated, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, `DELETE FROM subscriptions WHERE id = $1 RETURNING user_id`, id).Scan(&userID)
	return userID, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}
