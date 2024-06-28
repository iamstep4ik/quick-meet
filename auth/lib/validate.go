package lib

import (
	"errors"
	"regexp"
	"unicode"
)

func ValidateInput(username, email, password string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if email == "" || !isValidEmail(email) {
		return errors.New("invalid email address")
	}
	if err := isValidPassword(password); err != nil {
		return err
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPassword(password string) error {
	var hasMinLen, hasUpper, hasLower, hasNumber bool
	const minLen = 8

	if len(password) >= minLen {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasMinLen {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}
