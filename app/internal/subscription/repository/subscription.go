package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/subscription/model"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

const writeColumns = `id, user_id, type, status, start_date, end_date, created_at, updated_at`

const fullNameExpr = `COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), u.email, '')`

const readQuery = `SELECT s.id, s.user_id, s.type, s.status, s.start_date, s.end_date, s.created_at, s.updated_at,
	` + fullNameExpr + `
	FROM subscriptions s LEFT JOIN users u ON u.id = s.user_id`

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func scanFull(row interface {
	Scan(dest ...any) error
}, s *model.Subscription) error {
	return row.Scan(&s.ID, &s.UserID, &s.Type, &s.Status, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt, &s.UserFullName)
}

func collect(rows pgx.Rows) ([]model.Subscription, error) {
	defer rows.Close()

	var result []model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := scanFull(rows, &s); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *SubscriptionRepository) Create(ctx context.Context, s model.Subscription) (model.Subscription, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if s.Status == model.StatusActive {
		if _, err := tx.Exec(ctx,
			`UPDATE subscriptions SET status = 'in_active', updated_at = NOW()
			 WHERE user_id = $1 AND status = 'active'`,
			s.UserID,
		); err != nil {
			return model.Subscription{}, fmt.Errorf("failed to deactivate existing subscriptions: %w", err)
		}
	}

	var created model.Subscription
	err = scanFull(tx.QueryRow(ctx,
		`WITH ins AS (
			INSERT INTO subscriptions (user_id, type, status, start_date, end_date, created_at, updated_at)
			VALUES ($1, $2::subscription_type, $3::subscription_status, $4, $5, NOW(), NOW())
			RETURNING `+writeColumns+`
		 )
		 SELECT ins.id, ins.user_id, ins.type, ins.status, ins.start_date, ins.end_date, ins.created_at, ins.updated_at,
			`+fullNameExpr+`
		 FROM ins LEFT JOIN users u ON u.id = ins.user_id`,
		s.UserID, string(s.Type), string(s.Status), s.StartDate, s.EndDate,
	), &created)
	if err != nil {
		return created, store.MapForeignKey(store.MapUnique(err, model.ErrActiveSubscriptionExists), model.ErrUserRefInvalid)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Subscription{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return created, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int) (model.Subscription, error) {
	var s model.Subscription
	err := scanFull(r.db.QueryRow(ctx, readQuery+` WHERE s.id = $1`, id), &s)
	return s, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) GetLatestByUserID(ctx context.Context, userID int) (model.Subscription, error) {
	var s model.Subscription
	err := scanFull(r.db.QueryRow(ctx,
		readQuery+` WHERE s.user_id = $1 ORDER BY s.id DESC LIMIT 1`, userID,
	), &s)
	return s, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}

func (r *SubscriptionRepository) GetActiveByUserID(ctx context.Context, userID int) (model.Subscription, error) {
	var s model.Subscription
	err := scanFull(r.db.QueryRow(ctx,
		readQuery+`
		 WHERE s.user_id = $1 AND s.status = 'active' AND (s.end_date IS NULL OR s.end_date > NOW())
		 ORDER BY s.id DESC LIMIT 1`, userID,
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
		readQuery+` WHERE s.user_id = $1 ORDER BY s.id DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list user subscriptions: %w", err)
	}
	return collect(rows)
}

func (r *SubscriptionRepository) ListAll(ctx context.Context) ([]model.Subscription, error) {
	rows, err := r.db.Query(ctx, readQuery+` ORDER BY s.id DESC`)
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
		readQuery+` ORDER BY s.id DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	result, err := collect(rows)
	return result, total, err
}

func (r *SubscriptionRepository) Update(ctx context.Context, id int, s model.Subscription) (model.Subscription, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if s.Status == model.StatusActive {
		if _, err := tx.Exec(ctx,
			`UPDATE subscriptions SET status = 'in_active', updated_at = NOW()
			 WHERE status = 'active' AND id <> $1
			   AND user_id = (SELECT user_id FROM subscriptions WHERE id = $1)`,
			id,
		); err != nil {
			return model.Subscription{}, fmt.Errorf("failed to deactivate existing subscriptions: %w", err)
		}
	}

	var updated model.Subscription
	err = scanFull(tx.QueryRow(ctx,
		`WITH upd AS (
			UPDATE subscriptions
			SET type = $1::subscription_type, status = $2::subscription_status, end_date = $3, updated_at = NOW()
			WHERE id = $4
			RETURNING `+writeColumns+`
		 )
		 SELECT upd.id, upd.user_id, upd.type, upd.status, upd.start_date, upd.end_date, upd.created_at, upd.updated_at,
			`+fullNameExpr+`
		 FROM upd LEFT JOIN users u ON u.id = upd.user_id`,
		string(s.Type), string(s.Status), s.EndDate, id,
	), &updated)
	if err != nil {
		return updated, store.MapUnique(store.MapNoRows(err, model.ErrSubscriptionNotFound), model.ErrActiveSubscriptionExists)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Subscription{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return updated, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, `DELETE FROM subscriptions WHERE id = $1 RETURNING user_id`, id).Scan(&userID)
	return userID, store.MapNoRows(err, model.ErrSubscriptionNotFound)
}
