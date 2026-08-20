package validation

import (
	"errors"
	"net/mail"
	"regexp"
	"social/backend/models"
	"time"
)

// it recieve user data in a User struct and return nil if validated
func ValidateUserData(userData models.RegisterRequest) error {
	// validating email
	if err := validateEmail(userData.Email); err != nil {
		return err
	}

	//validate username
	if err := validateUsername(userData.Username); err != nil {
		return err
	}

	//validate first and last name
	if err := validateName(userData.FirstName); err != nil {
		return errors.Join(err, errors.New("first name"))
	}
	if err := validateName(userData.LastName); err != nil {
		return errors.Join(err, errors.New("last name"))
	}

	//validate DOB
	dob,err := time.Parse("2006-01-02", userData.DOB)
	if err != nil {
		return errors.New("invalid dob")
	}
	if err := validateDOB(dob); err != nil {
		return err
	}

	//validate password
	if err := validatePassword(userData.Password); err != nil {
		return err
	}
	return nil
}

// recieve email as a string and return a boolean if email format is valid and its not too short or long
func validateEmail(email string) error {
	if len(email) < 4 {
		return errors.New("invalid email. email length must be more than four characters")
	} else if len(email) > 70 {
		return errors.New("invalid email. email length must be less than 70 characters")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email format is invalid")
	}
	
	return nil
}

// recieve username as a string and return a boolean if username is not too short or long
// and if it does not contain forbidden charcters
func validateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("invalid username. username length must be at least three characters long")
	}
	if len(username) > 15 {
		return errors.New("invalid username. username length must be less than 15 charcters long")
	}

	matched, err := regexp.MatchString(`^[a-z0-9._]+$`, username)
	if err != nil || !matched {
		return errors.New("invalid username. forbidden characters")
	}

	return nil
}

// recieve name and validate the length and that it only contains alphabatic characters
func validateName(name string) error {
	if len(name) < 2 {
		errors.New("invalid name. length cannot be less than 3")
	}
	if len(name) < 2 {
		errors.New("invalid name. length cannot be more than 25")
	}

	matched, err := regexp.MatchString(`^[a-z]+$`, name)
	if err != nil || !matched {
		return errors.New("invalid name. forbidden characters")
	}

	return nil
}

// validate date of birth by checkint if its more than 100 years ago or less then 13 years ago or if it is in the future
func validateDOB(dob time.Time) error {

	if dob.Compare(time.Now()) >= 0 {
		return errors.New("invalid date of birth")
	}

	tooOldCheck := time.Now().AddDate(-120, 0, 0)
	if dob.Before(tooOldCheck) {
		return errors.New("invalid date of birth")
	}

	tooYoungCheck := time.Now().AddDate(-13, 0, 0)
	if dob.After(tooYoungCheck) {
		return errors.New("user must be at least 13 years old")
	}

	return nil
}

// validate password by checking it length and that in contains special characters,  numbers and alphabatic characters
func validatePassword(pass string) error {
	if len(pass) < 8 {
		return errors.New("invalid password: password must be at least 8 characters long")
	}

	if len(pass) > 75 {
		return errors.New("invalid password: password must be less than 75 characters long")
	}

	numberCheck, err := regexp.MatchString(`[0-9]`, pass)
	if err != nil {
		return errors.New("invalid password: password format is wrong")
	}

	if !numberCheck {
		return errors.New("invalid password: password must contain a number")
	}

	specialCharacterCheck, err := regexp.MatchString(`[!@#$%*&^?]`, pass)
	if err != nil {
		return errors.New("invalid password: password format is wrong")
	}

	if !specialCharacterCheck {
		return errors.New("invalid password: password must contain a special character: ! @ # $ % * & ^ ?")
	}

	alphaCheck, err := regexp.MatchString(`[a-zA-Z]`, pass)
	if err != nil {
		return errors.New("invalid password: password format is wrong")
	}

	if !alphaCheck {
		return errors.New("invalid password: password must contain alphabetic characters")
	}

	forbiddenCheck, err := regexp.MatchString(`[^a-zA-Z0-9!@#$%*&^?]`, pass)
	if err != nil {
		return errors.New("invalid password: password format is wrong")
	}

	if forbiddenCheck {
		return errors.New("invalid password: password contains forbidden characters")
	}

	return nil
}

