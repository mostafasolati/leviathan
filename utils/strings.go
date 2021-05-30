package utils

import (
	"strings"

	"github.com/dongri/phonenumber"
)

// NormalizePhoneNumber converts the phone number to standard format.
func NormalizePhoneNumber(number string) string {
	normal := phonenumber.Parse(number, "IR")
	if strings.Index(normal, "98") == 0 {
		return "0" + normal[2:]
	}
	if strings.Index(normal, "0098") == 0 {
		return "0" + normal[4:]
	}
	return normal
}

// InternationalPhoneNumber converts the phone number to international format.
func InternationalPhoneNumber(number string) string {
	normal := phonenumber.Parse(number, "IR")
	if strings.Index(normal, "0") == 0 {
		if strings.Index(normal, "00") == 0 {
			return normal[2:]
		}
		return "98" + normal[1:]
	}
	return normal
}
