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
	if e.Response.StatusCode == 401 {
		return "Unauthorized, please login again"
	}

	return fmt.Sprintf("API error, status: %s. \nResponse: %v", e.Status, e.Response)
}
