package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/content/model"
	"github.com/nurkenspashev92/bookit/internal/content/schema"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type FAQRepository struct {
	db *pgxpool.Pool
}

func NewFAQRepository(db *pgxpool.Pool) *FAQRepository {
	return &FAQRepository{db: db}
}

const faqSearchWhere = " WHERE (question_kz ILIKE $1 OR question_ru ILIKE $1 OR question_en ILIKE $1 OR answer_kz ILIKE $1 OR answer_ru ILIKE $1 OR answer_en ILIKE $1)"

func (r *FAQRepository) GetAll(ctx context.Context, search string) ([]schema.FAQ, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = faqSearchWhere
		args = append(args, "%"+search+"%")
	}

	rows, err := r.db.Query(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq`+where, args...)
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

func (r *FAQRepository) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.FAQ, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = faqSearchWhere
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM faq`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx,
		fmt.Sprintf(`SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq%s ORDER BY id LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var faqs []schema.FAQ
	for rows.Next() {
		var f schema.FAQ
		if err := rows.Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn); err != nil {
			return nil, 0, err
		}
		faqs = append(faqs, f)
	}
	return faqs, total, rows.Err()
}

func (r *FAQRepository) GetByID(ctx context.Context, id int) (schema.FAQ, error) {
	var f schema.FAQ
	err := r.db.QueryRow(ctx, `SELECT id, question_kz, answer_kz, question_ru, answer_ru, question_en, answer_en FROM faq WHERE id=$1`, id).
		Scan(&f.ID, &f.QuestionKz, &f.AnswerKz, &f.QuestionRu, &f.AnswerRu, &f.QuestionEn, &f.AnswerEn)
	if err != nil {
		return f, store.MapNoRows(err, model.ErrFAQNotFound)
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
		return model.ErrFAQNotFound
	}
	return nil
}
