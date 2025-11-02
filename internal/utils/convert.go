<<<<<<< HEAD
package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)[A-Z][a-z]+")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func NormalizeString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
func ConvertInt32ToPointer(v int32) *int32 {
	if v == 0 {
		return nil
	}
	//var i int32 = int32(v)
	return &v
}

func CapitalLizeFirtCharacter(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
=======
package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)[A-Z][a-z]+")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func NormalizeString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
func ConvertInt32ToPointer(v int32) *int32 {
	if v == 0 {
		return nil
	}
	//var i int32 = int32(v)
	return &v
}

func CapitalLizeFirtCharacter(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114
