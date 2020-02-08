package app_errors

import (
	"fmt"
	"net/http"
)

type UserError struct {
	FnErr string `json:"fn_err"`
	LnErr string `json:"ln_err"`
	EmailErr string `json:"email_err"`
	PassErr string `json:"pass_err"`
	NoUser string `json:"no_user"`
	OtherErr string `json:"other_error"`
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

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

//func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
//	return restErr{
//		ErrMessage: message,
//		ErrStatus:  status,
//		ErrError:   err,
//		ErrCauses:  causes,
//	}
//}

//func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
//	var apiErr restErr
//	if err := json.Unmarshal(bytes, &apiErr); err != nil {
//		return nil, errors.New("invalid json")
//	}
//	return apiErr, nil
//}

//func NewBadRequestError(message string) RestErr {
//	return restErr{
//		ErrMessage: message,
//		ErrStatus:  http.StatusBadRequest,
//		ErrError:   "bad_request",
//	}
//}

func NewUnprocessibleError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "unprocessible_request",
	}
}

//func NewNotFoundError(message string) RestErr {
//	return restErr{
//		ErrMessage: message,
//		ErrStatus:  http.StatusNotFound,
//		ErrError:   "not_found",
//	}
//}
//
//func NewUnauthorizedError(message string) RestErr {
//	return restErr{
//		ErrMessage: message,
//		ErrStatus:  http.StatusUnauthorized,
//		ErrError:   "unauthorized",
//	}
//}
//
//func NewInternalServerError(message string, err error) RestErr {
//	result := restErr{
//		ErrMessage: message,
//		ErrStatus:  http.StatusInternalServerError,
//		ErrError:   "internal_server_error",
//	}
//	if err != nil {
//		result.ErrCauses = append(result.ErrCauses, err.Error())
//	}
//	return result
//}