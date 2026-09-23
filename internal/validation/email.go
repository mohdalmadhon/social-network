package validation

func ValidateEmail(email string) error {
	return validateEmail(&email)
}
