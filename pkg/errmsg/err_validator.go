package errmsg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func errorValidationHandler[T any](err error, payload *T) (int, map[string][]string) {
	var (
		errorMessages = make(map[string][]string)
		code          = 400
	)

	for _, err := range err.(validator.ValidationErrors) {
		var (
			// Get the JSON tag name
			namespace  = err.Namespace()               // ex: UpdateInterestRequest.interest
			fieldParts = strings.Split(namespace, ".") // ex: [UpdateInterestRequest, interest]

			field      string
			fieldInMsg string
			message    string

			value     = err.Value()
			valueType = reflect.TypeOf(value)

			// Get the error message
		)
		lastField := fieldParts[len(fieldParts)-1]                                // get the last element
		fieldParts = fieldParts[1:]                                               // remove the first element
		field = strings.Join(fieldParts, ".")                                     // join the rest of the elements
		if strings.Contains(lastField, "_") && strings.Contains(lastField, "]") { // check if the last element contains "_" and "]", ex: interested_in[0]
			// fieldInMsg = field
			// remove characters between "[" and "]"
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
		} else {
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			if strings.Contains(fieldInMsg, "[") {
				fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
			}
		}

		if err.Param() != "" {
			message = fmt.Sprintf("field validation for '%s' failed on the '%s' tag with param '%s'", field, err.Tag(), err.Param())
		} else {
			message = fmt.Sprintf("field validation for '%s' failed on the '%s' tag", field, err.Tag())
		}

		// get validate tag that causes the error
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required.", fieldInMsg)
		case "email":
			message = fmt.Sprintf("%s is not a valid email address.", field)
		case "email_blacklist":
			message = fmt.Sprintf("email %v is not allowed.", value)
		case "strong_password":
			message = fmt.Sprintf("%s must be at least 12 characters and contain at least one uppercase letter, one lowercase letter, and one number.", fieldInMsg)
		case "exist":
			message = "resource is not exist."
		case "datetime":
			message = fmt.Sprintf("%s is not a valid datetime format (Ex: %s).", fieldInMsg, err.Param())
		case "ulid":
			message = fmt.Sprintf("%s is not a valid ULID.", fieldInMsg)
		case "uuid":
			message = fmt.Sprintf("%s is not a valid UUID.", fieldInMsg)
		case "min":
			// check if the field is a number or a string
			if valueType.Kind() == reflect.Int || valueType.Kind() == reflect.Int8 || valueType.Kind() == reflect.Int16 || valueType.Kind() == reflect.Int32 || valueType.Kind() == reflect.Int64 || valueType.Kind() == reflect.Float32 || valueType.Kind() == reflect.Float64 {
				message = fmt.Sprintf("%s must be at least %s.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.String {
				message = fmt.Sprintf("%s must be at least %s characters.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				message = fmt.Sprintf("%s must have at least %s items.", fieldInMsg, err.Param())
			}
		case "max":
			// check if the field is a number or a string
			if _, ok := value.(int); ok {
				message = fmt.Sprintf("%s must not be greater than %s.", fieldInMsg, err.Param())
			}
			if _, ok := value.(float64); ok {
				message = fmt.Sprintf("%s must not be greater than %s.", fieldInMsg, err.Param())
			}
			if _, ok := value.(string); ok {
				message = fmt.Sprintf("%s must not be greater than %s characters.", fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				message = fmt.Sprintf("%s must not have more than %s items.", fieldInMsg, err.Param())
			}

		case "eqfield":
			eqField := err.Param()
			eqFieldName := ""
			eqFieldTag, _ := reflect.TypeOf(payload).Elem().FieldByName(eqField)
			eqFieldJSONTag := eqFieldTag.Tag.Get("json")
			eqFieldQueryTag := eqFieldTag.Tag.Get("query")
			eqFieldFormTag := eqFieldTag.Tag.Get("form")
			eqFieldParamsTag := eqFieldTag.Tag.Get("params")

			if eqFieldJSONTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldJSONTag, "_", " ")
			}
			if eqFieldQueryTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldQueryTag, "_", " ")
			}
			if eqFieldFormTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldFormTag, "_", " ")
			}
			if eqFieldParamsTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldParamsTag, "_", " ")
			}

			message = fmt.Sprintf("%s must be equal to %s.", fieldInMsg, eqFieldName)
		case "oneof":
			message = fmt.Sprintf("%s must be one of %s.", fieldInMsg, err.Param())
		case "unique_in_slice":
			message = fmt.Sprintf("%s elements must be unique.", fieldInMsg)
		}

		errorMessages[field] = append(errorMessages[field], message)
	}

	return code, errorMessages
}
