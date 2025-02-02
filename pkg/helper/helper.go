package helper

import (
	"reflect"
	"time"
)

func DerefString(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}

func DerefInt(i *int, fallback int) int {
	if i == nil {
		return fallback
	}
	return *i
}

func DerefGeneric[T any](value interface{}, fallback T) T {
	val := reflect.ValueOf(value)

	if val.Kind() == reflect.Struct && val.NumField() >= 2 {
		validField := val.Field(1) // Assume "Valid" is always the second field
		if validField.Kind() == reflect.Bool && validField.Bool() {
			return val.Field(0).Interface().(T) // Type assertion directly
		}
	}

	return fallback
}

func FormatTimeToUTC(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

func MapToSlice[K comparable, V any](inputMap map[K]V) []V {
	slice := make([]V, 0, len(inputMap))
	for _, value := range inputMap {
		slice = append(slice, value)
	}
	return slice
}
