package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/url"
	"strings"
)

func ToCamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i, p := range parts {
		parts[i] = cases.Title(language.Portuguese).String(p)
	}
	return strings.Join(parts, "")
}

func ValidateURL(input string) bool {
	u, err := url.ParseRequestURI(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	
	// Verifica esquemas permitidos (http, https)
	if !strings.Contains(strings.ToLower(u.Scheme), "http") &&
		!strings.Contains(strings.ToLower(u.Scheme), "https") {
		return false
	}
	
	return true
}
