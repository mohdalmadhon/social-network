package validation

/*
this function exists because the original validate email (with small letter) was private idk why I didnt make it public

Parameters:
	email stirng

Returns:
	error
		-> nil if success
*/
func ValidateEmail(email string) error {
	return validateEmail(&email)
}
