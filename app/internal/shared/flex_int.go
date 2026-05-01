package shared

import (
	"encoding/json"
	"strconv"
	"strings"
)

// FlexInt accepts JSON number (5), string ("5"), empty string (""), or null.
// Empty/non-numeric strings are treated as 0 (zero value).
// Marshals always as a JSON number.
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexInt(n)
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

	n, err := strconv.Atoi(s)
	if err != nil {
		*f = 0
		return nil
	}
	*f = FlexInt(n)
	return nil
}

func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

func (f FlexInt) Int() int {
	return int(f)
}

func (f *FlexInt) IntPtr() *int {
	if f == nil {
		return nil
	}
	v := int(*f)
	return &v
}
