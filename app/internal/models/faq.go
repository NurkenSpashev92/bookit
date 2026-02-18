package models

import "time"

type FAQ struct {
	ID         int       `json:"id"`
	QuestionKz string    `json:"question_kz,omitempty"`
	AnswerKz   string    `json:"answer_kz,omitempty"`
	QuestionRu string    `json:"question_ru,omitempty"`
	AnswerRu   string    `json:"answer_ru,omitempty"`
	QuestionEn string    `json:"question_en,omitempty"`
	AnswerEn   string    `json:"answer_en,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Inquiry struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Text        string    `json:"text"`
	IsApproved  bool      `json:"is_approved"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
