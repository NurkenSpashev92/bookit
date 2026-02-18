package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type FAQRepository struct {
	db *pgxpool.Pool
}

type InquiryRepository struct {
	db *pgxpool.Pool
}

func NewFAQRepository(db *pgxpool.Pool) *FAQRepository {
	return &FAQRepository{db: db}
}

func NewInquiryRepository(db *pgxpool.Pool) *InquiryRepository {
	return &InquiryRepository{db: db}
}

func (r *FAQRepository) GetAll(ctx context.Context) ([]schemas.FAQ, error) {
	rows, err := r.db.Query(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqs []schemas.FAQ
	for rows.Next() {
		var f schemas.FAQ
		if err := rows.Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn); err != nil {
			return nil, err
		}
		faqs = append(faqs, f)
	}
	return faqs, nil
}

func (r *FAQRepository) GetByID(ctx context.Context, id int) (schemas.FAQ, error) {
	var f schemas.FAQ
	err := r.db.QueryRow(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq WHERE id=$1`, id).
		Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn)
	if err != nil {
		return f, fmt.Errorf("FAQ not found: %w", err)
	}
	return f, nil
}

func (r *FAQRepository) Create(ctx context.Context, req schemas.FAQCreateRequest) (schemas.FAQ, error) {
	var f schemas.FAQ
	err := r.db.QueryRow(ctx,
		`INSERT INTO faq (question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en) 
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en`,
		req.QuestionKz, req.AnswerKz, req.QuestionRu, req.AnswerRu, req.QuestionEn, req.AnswerEn,
	).Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn)
	return f, err
}

func (r *FAQRepository) Update(ctx context.Context, id int, req schemas.FAQUpdateRequest) (schemas.FAQ, error) {
	f, err := r.GetByID(ctx, id)
	if err != nil {
		return f, err
	}

	if req.QuestionKz != nil {
		f.QuestionKz = *req.QuestionKz
	}
	if req.AnswerKz != nil {
		f.AnswerKz = *req.AnswerKz
	}
	if req.QuestionRu != nil {
		f.QuestionRu = *req.QuestionRu
	}
	if req.AnswerRu != nil {
		f.AnswerRu = *req.AnswerRu
	}
	if req.QuestionEn != nil {
		f.QuestionEn = *req.QuestionEn
	}
	if req.AnswerEn != nil {
		f.AnswerEn = *req.AnswerEn
	}

	_, err = r.db.Exec(ctx,
		`UPDATE faq SET question_kz=$1, answer_kz=$2, question_ru=$3, answer_ru=$4, question_en=$5, answer_en=$6 WHERE id=$7`,
		f.QuestionKz, f.AnswerKz, f.QuestionRu, f.AnswerRu, f.QuestionEn, f.AnswerEn, id,
	)
	return f, err
}

func (r *FAQRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM faq WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("FAQ not found")
	}
	return nil
}

func (r *InquiryRepository) GetAll(ctx context.Context) ([]schemas.Inquiry, error) {
	rows, err := r.db.Query(ctx, `SELECT id, email, phone_number, text, is_approved FROM inquiries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []schemas.Inquiry
	for rows.Next() {
		var i schemas.Inquiry
		if err := rows.Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (r *InquiryRepository) GetByID(ctx context.Context, id int) (schemas.Inquiry, error) {
	var i schemas.Inquiry
	err := r.db.QueryRow(ctx, `SELECT id, email, phone_number, text, is_approved FROM inquiries WHERE id=$1`, id).
		Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved)
	if err != nil {
		return i, fmt.Errorf("Inquiry not found: %w", err)
	}
	return i, nil
}

func (r *InquiryRepository) Create(ctx context.Context, req schemas.InquiryCreateRequest) (schemas.Inquiry, error) {
	var i schemas.Inquiry
	err := r.db.QueryRow(ctx,
		`INSERT INTO inquiries (email, phone_number, text, is_approved) VALUES ($1,$2,$3,$4) RETURNING id,email,phone_number,text,is_approved`,
		req.Email, req.PhoneNumber, req.Text, req.IsApproved,
	).Scan(&i.ID, &i.Email, &i.PhoneNumber, &i.Text, &i.IsApproved)
	return i, err
}

func (r *InquiryRepository) Update(ctx context.Context, id int, req schemas.InquiryUpdateRequest) (schemas.Inquiry, error) {
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
