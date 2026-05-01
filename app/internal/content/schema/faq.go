package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

// FAQ response DTO
type FAQ struct {
	ID         int    `json:"id" example:"1"`
	QuestionKz string `json:"question_kz,omitempty" example:"Қалай брондау керек?"`
	AnswerKz   string `json:"answer_kz,omitempty" example:"Сайтта тіркеліңіз"`
	QuestionRu string `json:"question_ru,omitempty" example:"Как забронировать?"`
	AnswerRu   string `json:"answer_ru,omitempty" example:"Зарегистрируйтесь на сайте"`
	QuestionEn string `json:"question_en,omitempty" example:"How to book?"`
	AnswerEn   string `json:"answer_en,omitempty" example:"Register on the website"`
}

// FAQCreateRequest create FAQ request
// @Description Request body for creating a FAQ entry
type FAQCreateRequest struct {
	QuestionKz string `json:"question_kz" maxLength:"500" example:"Қалай брондау керек?" validate:"required"`
	AnswerKz   string `json:"answer_kz" example:"Сайтта тіркеліңіз" validate:"required"`
	QuestionRu string `json:"question_ru" maxLength:"500" example:"Как забронировать?" validate:"required"`
	AnswerRu   string `json:"answer_ru" example:"Зарегистрируйтесь на сайте" validate:"required"`
	QuestionEn string `json:"question_en" maxLength:"500" example:"How to book?" validate:"required"`
	AnswerEn   string `json:"answer_en" example:"Register on the website" validate:"required"`
}

func (r FAQCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("question_kz", r.QuestionKz)
	v.MaxLen("question_kz", r.QuestionKz, 500)
	v.Required("answer_kz", r.AnswerKz)
	v.Required("question_ru", r.QuestionRu)
	v.MaxLen("question_ru", r.QuestionRu, 500)
	v.Required("answer_ru", r.AnswerRu)
	v.Required("question_en", r.QuestionEn)
	v.MaxLen("question_en", r.QuestionEn, 500)
	v.Required("answer_en", r.AnswerEn)
	return v.Result()
}

// FAQUpdateRequest partial update FAQ
// @Description Request body for updating a FAQ entry (all fields optional)
type FAQUpdateRequest struct {
	QuestionKz *string `json:"question_kz,omitempty" maxLength:"500"`
	AnswerKz   *string `json:"answer_kz,omitempty"`
	QuestionRu *string `json:"question_ru,omitempty" maxLength:"500"`
	AnswerRu   *string `json:"answer_ru,omitempty"`
	QuestionEn *string `json:"question_en,omitempty" maxLength:"500"`
	AnswerEn   *string `json:"answer_en,omitempty"`
}

func (r FAQUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("question_kz", r.QuestionKz, 500)
	v.MaxLenPtr("question_ru", r.QuestionRu, 500)
	v.MaxLenPtr("question_en", r.QuestionEn, 500)
	return v.Result()
}
