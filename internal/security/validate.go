package security

import (
	"errors"
	"regexp"
)

type ContactInput interface {
	GetName() string
	GetEmail() string
	GetMessage() string
	GetWebsite() string
}

type ContactLike struct {
	Name    string
	Email   string
	Message string
	Website string
}

func (c ContactLike) GetName() string    { return c.Name }
func (c ContactLike) GetEmail() string   { return c.Email }
func (c ContactLike) GetMessage() string { return c.Message }
func (c ContactLike) GetWebsite() string { return c.Website }

func ValidateContact(req ContactInput) error {
	if req.GetWebsite() != "" {
		return errors.New("bot detected")
	}

	name := sanitize(req.GetName())
	email := sanitize(req.GetEmail())
	message := sanitize(req.GetMessage())

	if len(name) < 2 || len(name) > 80 {
		return errors.New("invalid name")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}

	if len(message) < 10 || len(message) > 2000 {
		return errors.New("invalid message")
	}

	if containsAttackPattern(name) ||
		containsAttackPattern(email) ||
		containsAttackPattern(message) {
		return errors.New("malicious pattern detected")
	}

	return nil
}
