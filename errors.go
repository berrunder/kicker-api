package main

import (
	"fmt"
)

// APIError to send in JSON format
type APIError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func getAPIError(status string, err error) *APIError {
	return &APIError{
		Status:  status,
		Message: fmt.Sprint(err),
	}
}

func getInternalError(err error) *APIError {
	return getAPIError("Internal server error", err)
}
