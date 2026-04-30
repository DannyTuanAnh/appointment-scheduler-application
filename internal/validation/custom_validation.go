package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("not_blank", func(fl validator.FieldLevel) bool {
		field := reflect.Indirect(fl.Field())

		switch field.Kind() {
		case reflect.String:
			return len(strings.TrimSpace(field.String())) > 0
		case reflect.Slice, reflect.Array:
			return field.Len() > 0
		default:
			return true
		}
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_fl", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return false
		}
		return fl.Field().Float() >= minVal
	})

	v.RegisterValidation("max_fl", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		maxVal, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return false
		}
		return fl.Field().Float() <= maxVal
	})
}
