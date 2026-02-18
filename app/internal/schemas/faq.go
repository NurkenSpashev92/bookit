package schemas

type FAQ struct {
	ID         int    `json:"id"`
	QuestionKz string `json:"question_kz,omitempty"`
	AnswerKz   string `json:"answer_kz,omitempty"`
	QuestionRu string `json:"question_ru,omitempty"`
	AnswerRu   string `json:"answer_ru,omitempty"`
	QuestionEn string `json:"question_en,omitempty"`
	AnswerEn   string `json:"answer_en,omitempty"`
}

type FAQCreateRequest struct {
	QuestionKz string `json:"question_kz"`
	AnswerKz   string `json:"answer_kz"`
	QuestionRu string `json:"question_ru"`
	AnswerRu   string `json:"answer_ru"`
	QuestionEn string `json:"question_en"`
	AnswerEn   string `json:"answer_en"`
}

type FAQUpdateRequest struct {
	QuestionKz *string `json:"question_kz,omitempty"`
	AnswerKz   *string `json:"answer_kz,omitempty"`
	QuestionRu *string `json:"question_ru,omitempty"`
	AnswerRu   *string `json:"answer_ru,omitempty"`
	QuestionEn *string `json:"question_en,omitempty"`
	AnswerEn   *string `json:"answer_en,omitempty"`
}

type Inquiry struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Text        string `json:"text"`
	IsApproved  bool   `json:"is_approved"`
}

type InquiryCreateRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Text        string `json:"text"`
	IsApproved  bool   `json:"is_approved"`
}

type InquiryUpdateRequest struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Text        *string `json:"text,omitempty"`
	IsApproved  *bool   `json:"is_approved,omitempty"`
}
