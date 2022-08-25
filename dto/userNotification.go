package dto

import "fmt"

// UserNotification represents basic DTO notification to the user if requested
// Initially, it supposed to be simple - if err != nil => error happens, else all is ok
type UserNotification struct {
	Error *ProcessingError `json:"error"`
}

func NewUserNotification(procErr *ProcessingError) UserNotification {
	return UserNotification{Error: procErr}
}

// ProcessingError represents basic DTO to provide information about the error
// when the processing request contains a notification request
type ProcessingError struct {
	//Here is string because of https://github.com/golang/go/issues/5161 - error marshals to {}
	Cause string `json:"cause"`

	Message string `json:"message"`
}

func (pe *ProcessingError) Error() string {
	return fmt.Sprintf("%s: %v", pe.Message, pe.Cause)
}

func NewProcessingError(message string, err string) ProcessingError {
	return ProcessingError{
		Cause:   err,
		Message: message,
	}
}
