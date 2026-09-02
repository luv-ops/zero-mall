package Regx

import "regexp"

func IsValidPhone(phone string) bool {
	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return reg.MatchString(phone)
}
