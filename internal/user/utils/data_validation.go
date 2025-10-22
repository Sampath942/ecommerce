package utils

import (
	"net/mail"
	"regexp"
)

func IsValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	digit := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#$&*]`)

	// At least 2 uppercase, 3 lowercase, 2 digits, and 1 special char
	return upper.MatchString(pw) &&
		lower.MatchString(pw) &&
		digit.MatchString(pw) &&
		special.MatchString(pw)
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`((\+*)((0[ -]*)*|((91 )*))((\d{12})+|(\d{10})+))|\d{5}([- ]*)\d{6}`)
	return re.MatchString(phone)
}
