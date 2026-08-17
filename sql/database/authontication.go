package queries

import "golang.org/x/crypto/bcrypt"

/*
recieve a password and username as strings and returns a boolean and an error
true if passowrd match with password in the database and false if not
*/ 
func (app App) AuthonticatePassword(password, identifier string) (bool, error) {
	hashedPassword, err := app.GetPasswordByIdentifier(identifier)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err 
	}

	return true, nil
}