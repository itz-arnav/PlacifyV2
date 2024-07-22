package validate

import (
	"errors"
	"html"
	"placify/backend/src/storage"
	"strings"

	"github.com/asaskevich/govalidator"
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (v *UserValidator) ValidateUser(user *storage.User, updating bool) error {
	user.Username = sanitizeInput(user.Username)
	user.Email = sanitizeInput(user.Email)

	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Email == "" || !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}
	if !updating && (user.Password == "") {
		return errors.New("password cannot be empty")
	}
	user.Password = sanitizeInput(user.Password)
	if !isValidAccessLevel(user.Access) {
		return errors.New("invalid access level")
	}
	return nil
}

func sanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return html.EscapeString(trimmed)
}

func isValidEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func isValidAccessLevel(access storage.AccessLevel) bool {
	return access >= storage.AccessViewer && access <= storage.AccessAdmin
}
