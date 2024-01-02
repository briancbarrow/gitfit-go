package ui

import "github.com/briancbarrow/gitfit-go/internal/validator"

type UserRegisterForm struct {
	FirstName           string `form:"first_name"`
	LastName            string `form:"last_name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
	Test                string `form:"test"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type TemplateData struct {
	Toast           string
	IsAuthenticated bool
	Form            any
	CSRFToken       string
}
