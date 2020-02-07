package app_errors

import "net/http"

type UserError struct {
	FnErr string `json:"fn_err"`
	LnErr string `json:"ln_err"`
	EmailErr string `json:"email_err"`
	StatusErr int `json:"status_err"`
}

func ValidationError(fn, ln, email string, status int) UserError {
	return UserError{
		FnErr:     fn,
		LnErr:     ln,
		EmailErr:  email,
		StatusErr: http.StatusUnprocessableEntity,
	}
}