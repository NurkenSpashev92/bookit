package schemas

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

type validator struct {
	errs ValidationErrors
}

func newValidator() *validator {
	return &validator{}
}

func (v *validator) required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.errs = append(v.errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *validator) requiredInt(field string, value int) {
	if value == 0 {
		v.errs = append(v.errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *validator) requiredFlexInt(field string, value FlexInt) {
	if value.Int() == 0 {
		v.errs = append(v.errs, fmt.Sprintf("%s is required", field))
	}
}

func (v *validator) maxLen(field, value string, max int) {
	if utf8.RuneCountInString(value) > max {
		v.errs = append(v.errs, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
}

func (v *validator) minLen(field, value string, min int) {
	if strings.TrimSpace(value) != "" && utf8.RuneCountInString(value) < min {
		v.errs = append(v.errs, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
}

func (v *validator) minInt(field string, value int, min int) {
	if value < min {
		v.errs = append(v.errs, fmt.Sprintf("%s must be at least %d", field, min))
	}
}

func (v *validator) email(field, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	if _, err := mail.ParseAddress(value); err != nil {
		v.errs = append(v.errs, fmt.Sprintf("%s must be a valid email", field))
	}
}

func (v *validator) date(field, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	if _, err := time.Parse("2006-01-02", value); err != nil {
		v.errs = append(v.errs, fmt.Sprintf("%s must be a valid date (YYYY-MM-DD)", field))
	}
}

func (v *validator) maxLenPtr(field string, value *string, max int) {
	if value != nil {
		v.maxLen(field, *value, max)
	}
}

func (v *validator) emailPtr(field string, value *string) {
	if value != nil {
		v.email(field, *value)
	}
}

func (v *validator) result() error {
	if v.errs.HasErrors() {
		return v.errs
	}
	return nil
}
