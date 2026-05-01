package shared

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

type ValidationErrors []string

func (v ValidationErrors) Error() string {
	return strings.Join(v, "; ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

type Validator struct {
	Errs ValidationErrors
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.Errs = append(v.Errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *Validator) RequiredInt(field string, value int) {
	if value == 0 {
		v.Errs = append(v.Errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *Validator) RequiredFlexInt(field string, value FlexInt) {
	if value.Int() == 0 {
		v.Errs = append(v.Errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *Validator) MaxLen(field, value string, max int) {
	if utf8.RuneCountInString(value) > max {
		v.Errs = append(v.Errs, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
}

func (v *Validator) MinLen(field, value string, min int) {
	if strings.TrimSpace(value) != "" && utf8.RuneCountInString(value) < min {
		v.Errs = append(v.Errs, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
}

func (v *Validator) MinInt(field string, value int, min int) {
	if value < min {
		v.Errs = append(v.Errs, fmt.Sprintf("%s must be at least %d", field, min))
	}
}

func (v *Validator) Email(field, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	if _, err := mail.ParseAddress(value); err != nil {
		v.Errs = append(v.Errs, fmt.Sprintf("%s must be a valid email", field))
	}
}

func (v *Validator) Date(field, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	if _, err := time.Parse("2006-01-02", value); err != nil {
		v.Errs = append(v.Errs, fmt.Sprintf("%s must be a valid date (YYYY-MM-DD)", field))
	}
}

func (v *Validator) MaxLenPtr(field string, value *string, max int) {
	if value != nil {
		v.MaxLen(field, *value, max)
	}
}

func (v *Validator) EmailPtr(field string, value *string) {
	if value != nil {
		v.Email(field, *value)
	}
}

func (v *Validator) Append(msg string) {
	v.Errs = append(v.Errs, msg)
}

func (v *Validator) Result() error {
	if v.Errs.HasErrors() {
		return v.Errs
	}
	return nil
}
