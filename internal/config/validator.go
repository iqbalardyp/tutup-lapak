package config

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"time"
)

var (
	sortByCache = make(map[string]bool)
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("time_validator", timeValidator)
	validate.RegisterValidation("is_uri", uriValidator)
	validate.RegisterValidation("sort_by", productSortByValidator)
	validate.RegisterValidation("contact_detail_validator", contactDetailValidation)
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

func productSortByValidator(fl validator.FieldLevel) bool {
	sortBy, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if sortBy == "" {
		return false
	} else if sortBy == "newest" || sortBy == "cheapest" {
		return true
	}

	if result, found := sortByCache[sortBy]; found {
		return result
	}
	isValid, err := regexp.MatchString("^sold-[0-9]+$", sortBy)
	if err != nil {
		return false
	}
	sortByCache[sortBy] = isValid

	return isValid
}

func contactDetailValidation(fl validator.FieldLevel) bool {
	contactType := fl.Parent().FieldByName("SenderContactType").String()
	contactDetail := fl.Field().String()

	if contactType == "phone" {
		phoneRegex := regexp.MustCompile(`\+\d{1,15}$`)
		return phoneRegex.MatchString(contactDetail)
	} else if contactType == "email" {
		return validator.New().Var(contactDetail, "email") == nil
	}
	return false
}
