package validation

import (
	"fmt"
	"service-user/model/dto/json"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
)

func GetValidationError(err error) []json.ValidationField {
	var validationFields []json.ValidationField
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range ve {
			log.Debug().Msg(fmt.Sprintf("validationError : %v", validationError))
			myField := convertFieldRequired(validationError.Namespace())
			validationFields = append(validationFields, json.ValidationField{
				FieldName: myField,
				Message:   formatMessage(validationError),
			})
		}
	}
	return validationFields
}

func convertFieldRequired(myValue string) string {
	log.Debug().Msg("convertFieldRequired: " + myValue)
	fieldSegmen := strings.Split(myValue, ".")
	myField := ""
	length := len(fieldSegmen)
	i := 1
	for _, val := range fieldSegmen {
		if i == 1 {
			i++
			continue
		}

		if i == length {
			myField += strcase.SnakeCase(val)
			break
		}

		myField += strcase.LowerCamelCase(val) + `/`
		i++
	}

	return myField
}

func formatMessage(err validator.FieldError) string {
	var message string

	switch err.Tag() {
	case "required":
		message = "required"
	case "number":
		message = "must be number"
	case "email":
		message = "invalid format email"
	case "DateOnly":
		message = "invalid format date"
	case "min":
		message = "minimum value is not exceed"
	case "max":
		message = "max value is exceed"
	case "uuid4":
		message = "format invalid"
	case "uppercase":
		message = "does not contains uppercase letter"
	case "lowercase":
		message = "does not contains lowercase letter"
	case "numeric":
		message = "does not contain numeric value"
	case "containsany":
		message = "does not contain required value"
	case "advancedpass":
		message = "password must contain at least 1 uppercase letter, 1 lowercase letter, 1 number and 1 special character"
	}

	return message
}
