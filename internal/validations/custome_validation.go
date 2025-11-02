<<<<<<< HEAD
package validations

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"user-management-api/internal/utils"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	///password strong
	v.RegisterValidation("password_strong", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		password := fl.Field().String()
		if len(password) < 6 {
			return false
		}
		hasLower := regexp.MustCompile(`[a-z]+`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]+`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]+`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>?/\\|]`).MatchString(password)

		return hasLower && hasUpper && hasDigit && hasSpecial
	})

	//blocked domain email
	var blockedDomains = map[string]bool{
		"blacklist.com": true,
		"edu.vn":        true,
		"abc.vn":        true,
		"templete.vn":   true,
	}
	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		email := fl.Field().String()
		part := strings.Split(email, "@")

		if len(part) != 2 {
			return false
		}
		domain := utils.NormalizeString(part[1])
		log.Printf("%+v", domain)
		return !blockedDomains[domain]
	})

	// define a domain tag validation

	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`) // regex for validating slug
	// define a slug tag validation
	v.RegisterValidation("slug_format", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		return slugRegex.MatchString(fl.Field().String())
	})
	// define a slug tag validation
	searchRegex := regexp.MustCompile(`^[a-zA-Z0-9\s]*$`) // regex for validating slug
	v.RegisterValidation("search_format", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		return searchRegex.MatchString(fl.Field().String())
	})
	// define a slug tag validation

	v.RegisterValidation("min_price", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		minStr := fl.Param()
		min, err := strconv.ParseFloat(minStr, 10)
		if err != nil {
			return false
		}
		return fl.Field().Float() > float64(min)
	})
	v.RegisterValidation("max_price", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		maxStr := fl.Param()
		max, err := strconv.ParseFloat(maxStr, 10)
		if err != nil {
			return false
		}
		return fl.Field().Float() <= max
	})
	// define a img_link tag validation
	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		fileName := fl.Field().String()
		allowStr := fl.Param()

		if allowStr != "" {
			allowExt := strings.Fields(allowStr)
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
			for _, al := range allowExt {
				if ext == strings.ToLower(al) {
					return true
				}
			}
		}
		return false
	})

}
=======
package validations

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"user-management-api/internal/utils"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	///password strong
	v.RegisterValidation("password_strong", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		password := fl.Field().String()
		if len(password) < 6 {
			return false
		}
		hasLower := regexp.MustCompile(`[a-z]+`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]+`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]+`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>?/\\|]`).MatchString(password)

		return hasLower && hasUpper && hasDigit && hasSpecial
	})

	//blocked domain email
	var blockedDomains = map[string]bool{
		"blacklist.com": true,
		"edu.vn":        true,
		"abc.vn":        true,
		"templete.vn":   true,
	}
	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		email := fl.Field().String()
		part := strings.Split(email, "@")

		if len(part) != 2 {
			return false
		}
		domain := utils.NormalizeString(part[1])
		log.Printf("%+v", domain)
		return !blockedDomains[domain]
	})

	// define a domain tag validation

	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`) // regex for validating slug
	// define a slug tag validation
	v.RegisterValidation("slug_format", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		return slugRegex.MatchString(fl.Field().String())
	})
	// define a slug tag validation
	searchRegex := regexp.MustCompile(`^[a-zA-Z0-9\s]*$`) // regex for validating slug
	v.RegisterValidation("search_format", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		return searchRegex.MatchString(fl.Field().String())
	})
	// define a slug tag validation

	v.RegisterValidation("min_price", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		minStr := fl.Param()
		min, err := strconv.ParseFloat(minStr, 10)
		if err != nil {
			return false
		}
		return fl.Field().Float() > float64(min)
	})
	v.RegisterValidation("max_price", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		maxStr := fl.Param()
		max, err := strconv.ParseFloat(maxStr, 10)
		if err != nil {
			return false
		}
		return fl.Field().Float() <= max
	})
	// define a img_link tag validation
	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		// handle slug format validation here
		fileName := fl.Field().String()
		allowStr := fl.Param()

		if allowStr != "" {
			allowExt := strings.Fields(allowStr)
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
			for _, al := range allowExt {
				if ext == strings.ToLower(al) {
					return true
				}
			}
		}
		return false
	})

}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114
