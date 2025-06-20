package entities

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidPasswordLength        = errors.New("password length must be in 8-20 characters")
	ErrNoNormalCharacterInPassword  = errors.New("password need at least 1 normal character")
	ErrNoCapitalCharacterInPassword = errors.New("password need at least 1 capital character")
	ErrNoNumberInPassword           = errors.New("password need at least 1 number")
)

func (u *User) Validate() error {
	if length := len(u.Password); length < 8 || length > 20 {
		return ErrInvalidPasswordLength
	}

	patternLower := regexp.MustCompile(`[a-z]`)
	if !patternLower.MatchString(u.Password) {
		return ErrNoNormalCharacterInPassword
	}

	patternUpper := regexp.MustCompile(`[A-Z]`)
	if !patternUpper.MatchString(u.Password) {
		return ErrNoCapitalCharacterInPassword
	}

	patternNumber := regexp.MustCompile(`[0-9]`)
	if !patternNumber.MatchString(u.Password) {
		return ErrNoNumberInPassword
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
