package validation

import (
	"fmt"
	"strings"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("failed to register validator")
	}

	RegisterCustomValidations(v)

	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationErr {
			root := strings.Split(e.Namespace(), ".")[0]

			raw := strings.TrimPrefix(e.Namespace(), root+".")

			paths := strings.Split(raw, ".")

			for i, path := range paths {
				if strings.Contains(path, "[") {
					idx := strings.Index(path, "[")
					base := utils.CamelToSnake(path[:idx])

					index := path[idx:]
					paths[i] = base + index
				} else {
					paths[i] = utils.CamelToSnake(path)
				}
			}

			fieldPath := strings.Join(paths, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s must only contain lowercase letters, numbers, hyphens, or dots", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be at least %s characters", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must not exceed %s characters", fieldPath, e.Param())
			case "min_fl":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "max_fl":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "oneof":
				allowedValue := strings.ReplaceAll(e.Param(), " ", ", ")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the following values: %s", fieldPath, allowedValue)
			case "file_ext":
				allowedValue := strings.ReplaceAll(e.Param(), " ", ", ")
				errors[fieldPath] = fmt.Sprintf("%s only allows the following extensions: %s", fieldPath, allowedValue)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s must only contain letters, numbers, and spaces", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid email address", fieldPath)
			case "datetime":
				format := e.Param()
				if format == "15:04" {
					format = "HH:mm"
				}
				errors[fieldPath] = fmt.Sprintf("%s must match the format %s", fieldPath, format)
			case "not_blank":
				errors[fieldPath] = fmt.Sprintf("%s cannot be blank or contain only spaces", fieldPath)
			default:
				errors[fieldPath] = fmt.Sprintf("%s is invalid", fieldPath)
			}
		}

		return gin.H{"errors": errors}
	}

	return gin.H{
		"error":  "Invalid request",
		"detail": err.Error(),
	}
}
