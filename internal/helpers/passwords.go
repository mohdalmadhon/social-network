package helpers

import "golang.org/x/crypto/bcrypt"

/*
convert a password (or any string) to a hash using bcypt package and using DefaultCost (equals to 10)

Parameters:
	password string 

Returns:
	stinrg
		-> hashed password
	error
		-> nil if success
*/
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

/*
funtion used to compare the password input from the user to the hashed version stored in the database or provided when calling the function. 
Any errors that are caused would probably be out of hand and nothing to do about it

Parameters:
	pass string 		-> plain text
	hashedpassword string -> hashed text

Returns:
	bool
		-> true if matched
		-> false if error or unmatch
*/
func AuthonticateUser(pass string, hashedPassword string) bool {
	if err := (bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))); err != nil {
		return false
	}
	return true
}
