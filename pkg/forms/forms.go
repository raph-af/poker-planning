package forms

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type NewStory struct {
	Title    string
	Content  string
	Failures map[string]string
}

func (f *NewStory) IsValid() bool {
	f.Failures = make(map[string]string)

	ValidateTitle(f)
	ValidateContent(f)

	return len(f.Failures) == 0
}

func ValidateTitle(f *NewStory) {
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.Failures["Title"] = "Title cannot be longer than 100 characters"
	}
}

func ValidateContent(f *NewStory) {
	if strings.TrimSpace(f.Content) == "" {
		f.Failures["Content"] = "Content is required"
	}
}

type SignupUser struct {
	Name     string
	Email    string
	Password string
	Failures map[string]string
}

var regexEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (f *SignupUser) IsValid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		f.Failures["Name"] = "Name is required"
	}

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	} else if len(f.Email) > 254 || !regexEmail.MatchString(f.Email) {
		f.Failures["Email"] = "Email is not a valid address"
	}

	if utf8.RuneCountInString(f.Password) < 8 {
		f.Failures["Password"] = "Password cannot be shorter than 8 characters"
	}

	return len(f.Failures) == 0
}

type LoginUser struct {
	Email    string
	Password string
	Failures map[string]string
}

func (f *LoginUser) IsValid() bool {
	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	}

	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
	}

	return len(f.Failures) == 0
}
