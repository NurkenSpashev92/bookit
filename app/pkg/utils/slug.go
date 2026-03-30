package utils

import (
	"regexp"
	"strings"

	"github.com/gosimple/slug"
)

func GenerateSlug(slug, nameEN, nameKZ, nameRU string) string {
	if strings.TrimSpace(slug) != "" {
		return normalize(slug)
	}

	if nameEN != "" {
		return normalize(nameEN)
	}
	if nameKZ != "" {
		return normalize(nameKZ)
	}

	return normalize(nameRU)
}

var reMultiDash = regexp.MustCompile(`-+`)

func normalize(s string) string {
	s = slug.Make(s)
	return strings.Trim(reMultiDash.ReplaceAllString(s, "-"), "-")
}
