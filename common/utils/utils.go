package utils

import "strings"

func StrInList(s string, l []string) bool {
	if nil == l {
		return false
	}
	for _, item := range l {
		if s == item {
			return true
		}
	}
	return false
}

func FieldToVarName(fieldName string) string {
	r := strings.ToLower(fieldName[:1]) + fieldName[1:]
	// sensitive key words
	if r == "type" {
		r = "typ"
	}
	return r
}

// CamelToUnderline 驼峰式命名转化为下划线命名
func CamelToUnderline(name string) string {
	var (
		underlineName string
	)
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			if i > 0 && !(name[i-1] >= 'A' && name[i-1] <= 'Z') {
				underlineName += "_"
			}
			underlineName += string(c + 32)
		} else {
			underlineName += string(c)
		}
	}
	return underlineName
}
