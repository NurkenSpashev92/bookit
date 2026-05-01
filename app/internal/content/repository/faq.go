package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
)

type FAQRepository struct {
	db *pgxpool.Pool
}

func NewFAQRepository(db *pgxpool.Pool) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) GetAll(ctx context.Context) ([]schema.FAQ, error) {
	rows, err := r.db.Query(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqs []schema.FAQ
	for rows.Next() {
		var f schema.FAQ
		if err := rows.Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn); err != nil {
			return nil, err
		}
		faqs = append(faqs, f)
	}
	return faqs, nil
}

func (r *FAQRepository) GetByID(ctx context.Context, id int) (schema.FAQ, error) {
	var f schema.FAQ
	err := r.db.QueryRow(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq WHERE id=$1`, id).
		Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn)
	if err != nil {
		return f, fmt.Errorf("FAQ not found: %w", err)
	}
	return f, nil
}

func (r *FAQRepository) Create(ctx context.Context, req schema.FAQCreateRequest) (schema.FAQ, error) {
	var f schema.FAQ
	err := r.db.QueryRow(ctx,
		`INSERT INTO faq (question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en`,
		req.QuestionKz, req.AnswerKz, req.QuestionRu, req.AnswerRu, req.QuestionEn, req.AnswerEn,
	).Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn)
	return f, err
}

func (r *FAQRepository) Update(ctx context.Context, id int, req schema.FAQUpdateRequest) (schema.FAQ, error) {
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
