package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
)

type InquiryRepository struct {
	db *pgxpool.Pool
}

func NewInquiryRepository(db *pgxpool.Pool) *InquiryRepository {
	return &InquiryRepository{db: db}
}

func (r *InquiryRepository) GetAll(ctx context.Context) ([]schema.Inquiry, error) {
	rows, err := r.db.Query(ctx, `SELECT id, email, phone_number, text, is_approved FROM inquiries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []schema.Inquiry
	for rows.Next() {
		var i schema.Inquiry
		if err := rows.Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (r *InquiryRepository) GetByID(ctx context.Context, id int) (schema.Inquiry, error) {
	var i schema.Inquiry
	err := r.db.QueryRow(ctx, `SELECT id, email, phone_number, text, is_approved FROM inquiries WHERE id=$1`, id).
		Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved)
	if err != nil {
		return i, fmt.Errorf("Inquiry not found: %w", err)
	}
	return i, nil
}

func (r *InquiryRepository) Create(ctx context.Context, req schema.InquiryCreateRequest) (schema.Inquiry, error) {
	var i schema.Inquiry
	err := r.db.QueryRow(ctx,
		`INSERT INTO inquiries (email, phone_number, text, is_approved) VALUES ($1,$2,$3,$4) RETURNING id,email,phone_number,text,is_approved`,
		req.Email, req.PhoneNumber, req.Text, req.IsApproved,
	).Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved)
	return i, err
}

func (r *InquiryRepository) Update(ctx context.Context, id int, req schema.InquiryUpdateRequest) (schema.Inquiry, error) {
	i, err := r.GetByID(ctx, id)
	if err != nil {
		return i, err
	}

	if req.Email != nil {
		i.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		i.PhoneNumber = *req.PhoneNumber
	}
	if req.Text != nil {
		i.Text = *req.Text
	}
	if req.IsApproved != nil {
		i.IsApproved = *req.IsApproved
	}

	_, err = r.db.Exec(ctx,
		`UPDATE inquiries SET email=$1, phone_number=$2, text=$3, is_approved=$4 WHERE id=$5`,
		i.Email, i.PhoneNumber, i.Text, i.IsApproved, id,
	)
	return i, err
}

func (r *InquiryRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM inquiries WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("Inquiry not found")
	}
	return nil
}
