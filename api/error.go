package api

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	Status   string
	Response *http.Response
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API error, status: %s. \nResponse: %v", e.Status, e.Response)
}
