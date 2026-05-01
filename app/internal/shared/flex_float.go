package shared

import (
	"encoding/json"
	"strconv"
	"strings"
)

// FlexFloat64 accepts JSON number (43.12), string ("43.12"), empty string (""), or null.
// Empty/non-numeric strings are treated as 0 (zero value).
// Marshals always as a JSON number.
type FlexFloat64 float64

func (f *FlexFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexFloat64(n)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	if s == "" {
		*f = 0
		return nil
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		*f = 0
		return nil
	}
	*f = FlexFloat64(n)
	return nil
}

func (f FlexFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

func (f FlexFloat64) Float64() float64 {
	return float64(f)
}

func (f *FlexFloat64) Float64Ptr() *float64 {
	if f == nil {
		return nil
	}
	v := float64(*f)
	return &v
}
