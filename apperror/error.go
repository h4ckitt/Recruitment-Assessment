package apperror

import "net/http"

//AppError : Custom Error type which contains http status code and message
type AppError struct {
	Status  int
	Message string
}

var (
	BadRequest  = AppError{http.StatusBadRequest, "Bad request body received"}
	ServerError = AppError{http.StatusInternalServerError, "An error occurred while processing that request"}
	NotFound    = AppError{http.StatusNotFound, "The requested resource was not found"}
)

func NewError(status int, message string) AppError {
	return AppError{Status: status, Message: message}
}

func (e AppError) Error() string {
	return e.Message
}
