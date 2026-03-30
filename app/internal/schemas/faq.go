package schemas

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
	v := newValidator()
	v.required("question_kz", r.QuestionKz)
	v.maxLen("question_kz", r.QuestionKz, 500)
	v.required("answer_kz", r.AnswerKz)
	v.required("question_ru", r.QuestionRu)
	v.maxLen("question_ru", r.QuestionRu, 500)
	v.required("answer_ru", r.AnswerRu)
	v.required("question_en", r.QuestionEn)
	v.maxLen("question_en", r.QuestionEn, 500)
	v.required("answer_en", r.AnswerEn)
	return v.result()
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
	v := newValidator()
	v.maxLenPtr("question_kz", r.QuestionKz, 500)
	v.maxLenPtr("question_ru", r.QuestionRu, 500)
	v.maxLenPtr("question_en", r.QuestionEn, 500)
	return v.result()
}

// Inquiry response DTO
type Inquiry struct {
	ID          int    `json:"id" example:"1"`
	Email       string `json:"email" example:"guest@example.com" format:"email"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+77001234567"`
	Text        string `json:"text" example:"I have a question about booking"`
	IsApproved  bool   `json:"is_approved" example:"false"`
}

// InquiryCreateRequest create inquiry request
// @Description Request body for creating an inquiry
type InquiryCreateRequest struct {
	Email       string `json:"email" format:"email" maxLength:"255" example:"guest@example.com" validate:"required"`
	PhoneNumber string `json:"phone_number,omitempty" maxLength:"20" example:"+77001234567"`
	Text        string `json:"text" example:"I have a question about booking" validate:"required"`
	IsApproved  bool   `json:"is_approved" example:"false"`
}

func (r InquiryCreateRequest) Validate() error {
	v := newValidator()
	v.required("email", r.Email)
	v.email("email", r.Email)
	v.maxLen("email", r.Email, 255)
	v.maxLen("phone_number", r.PhoneNumber, 20)
	v.required("text", r.Text)
	return v.result()
}

// InquiryUpdateRequest partial update inquiry
// @Description Request body for updating an inquiry (all fields optional)
type InquiryUpdateRequest struct {
	Email       *string `json:"email,omitempty" format:"email" maxLength:"255" example:"guest@example.com"`
	PhoneNumber *string `json:"phone_number,omitempty" maxLength:"20"`
	Text        *string `json:"text,omitempty"`
	IsApproved  *bool   `json:"is_approved,omitempty" example:"true"`
}

func (r InquiryUpdateRequest) Validate() error {
	v := newValidator()
	v.emailPtr("email", r.Email)
	v.maxLenPtr("email", r.Email, 255)
	v.maxLenPtr("phone_number", r.PhoneNumber, 20)
	return v.result()
}
