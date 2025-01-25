package config

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"

	"time"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("time_validator", timeValidator)
	validate.RegisterValidation("is_uri", uriValidator)
	return validate
}

func timeValidator(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Example: Check if time is after 1970-01-01T00:00:00Z
	return !t.IsZero() && t.After(time.Unix(0, 0))
}

func uriValidator(fl validator.FieldLevel) bool {
	uri, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Check for empty string
	if strings.TrimSpace(uri) == "" {
		return false
	}

	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// Check if scheme is present and valid
	if parsedURL.Scheme == "" {
		return false
	}

	// Validate host presence for network-based URIs
	if parsedURL.Scheme != "file" && parsedURL.Host == "" {
		return false
	}

	// Additional validation for specific schemes
	switch parsedURL.Scheme {
	case "http", "https":
		// For HTTP(S), ensure there's a valid host
		if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
			return false
		}
	case "file":
		// For file scheme, ensure there's a path
		if parsedURL.Path == "" {
			return false
		}
	}

	return true
}
