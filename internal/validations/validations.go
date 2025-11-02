package validations

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to register validator")
	}
	RegisterCustomValidation(v)
	return nil

}
func HandleValidationErr(err error) gin.H {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErr {
			//root:=strings.Split(e.Namespace(), ".")[1]
			// log.Printf("%+v :", e.Tag())
			// log.Printf("%+v :", e.Field())
			// log.Printf("%+v :", e.Param())
			switch e.Tag() {
			case "gt":
				errors[e.Field()] = "vlue of '" + e.Field() + "' Must be greater than " + e.Param()
			case "lt":
				errors[e.Field()] = "vlue of '" + e.Field() + "' Must be less than " + e.Param()
			case "gte":
				errors[e.Field()] = "vlue of '" + e.Field() + "' Must be greater than or equal " + e.Param()
			case "lte":
				errors[e.Field()] = "vlue of " + e.Field() + " Must be less than or equal" + e.Param()
			case "uuid":
				errors[e.Field()] = "Invalid uuid :" + e.Param()
			case "slug_format":
				errors[e.Field()] = "Invalid slug format"
			case "min": // overide original message
				errors[e.Field()] = fmt.Sprintf("length '%s' is very short. It must be at least %s characters", e.Field(), e.Param())
			case "max": // overide original message
				errors[e.Field()] = fmt.Sprintf("length '%s'  is very long. It must be less than %s characters", e.Field(), e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("'%s' must be one of %s", e.Field(), allowedValues)
			case "required":
				errors[e.Field()] = fmt.Sprintf("'%s' field is required", e.Field())
			case "search_format":
				errors[e.Field()] = fmt.Sprintf("'%s' field is not valid. It can only contain a-z, A-Z, 0-9 and space", e.Field())
			case "min_price":
				errors[e.Field()] = fmt.Sprintf("'%s' field is not valid. It must be greater than %v", e.Field(), e.Param())
			case "max_price":
				errors[e.Field()] = fmt.Sprintf("'%s' field is not valid. It must be less than %v", e.Field(), e.Param())
			case "file_ext":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("'%s' format is not valid. It only accept  image formats like  %s", e.Field(), allowedValue)
			case "email_advanced": //validate email in blacklist domain
				errors[e.Field()] = fmt.Sprintf(" This '%s' is in blacklist domain %v", e.Field(), e.Param())
			case "password_strong":
				errors[e.Field()] = fmt.Sprintf("'%s' field is not strong enough. It's at least 6 characters long, contains at least one lowercase letter, one uppercase letter, one digit, and one special character", e.Field())
			}
		}
		return gin.H{
			"error": errors,
		}
	}
	return gin.H{
		"error": "Invalid request " + err.Error(),
	}
}
