package box

import (
	"fmt"
)

type BoxError struct {
	StatusCode int
	Message    string
}

func (e *BoxError) Error() string {
	return fmt.Sprintf("%v : %v", e.StatusCode, e.Message)
}

var (
	SUCCESS    = &BoxError{200, "Success"}
	CREATED    = &BoxError{201, "Created"}
	ACCEPTED   = &BoxError{202, "Accepted"}
	NO_CONTENT = &BoxError{204, "No Content"}
)

var (
	REDIRECT     = &BoxError{302, "Redirect"}
	NOT_MODIFIED = &BoxError{304, "Not Modified"}
)

var (
	UNAUTHORIZED        = &BoxError{401, "Unauthorized"}        // Authorization failed
	FORBIDDEN           = &BoxError{403, "Forbidden"}           // Not enough permission for the operation
	NOT_FOUND           = &BoxError{404, "Not found"}           // Not Found
	NOT_ALLOWED         = &BoxError{405, "Not allowed"}         // Method not allowed
	CONFLICT            = &BoxError{409, "Conflict"}            // Same name item already exist
	PRECONDITION_FAILED = &BoxError{412, "Precondition failed"} // Precondition (If match) failed
	TOO_MANY_REQUESTS   = &BoxError{429, "Too many requests"}   // Too many requests
)

var (
	SERVER_ERROR = &BoxError{500, "Internal server error"} // Internal server error
	UNAVAILABLE  = &BoxError{503, "Unavailable"}           // Unavailable
)

func toError(status int) *BoxError {
	switch status {
	case 200:
		return SUCCESS
	case 201:
		return CREATED
	case 202:
		return ACCEPTED
	case 204:
		return NO_CONTENT
	case 302:
		return REDIRECT
	case 304:
		return NOT_MODIFIED
	case 401:
		return UNAUTHORIZED
	case 403:
		return FORBIDDEN
	case 404:
		return NOT_FOUND
	case 405:
		return NOT_ALLOWED
	case 409:
		return CONFLICT
	case 412:
		return PRECONDITION_FAILED
	case 429:
		return TOO_MANY_REQUESTS
	case 500:
		return SERVER_ERROR
	case 503:
		return UNAVAILABLE
	default:
		return &BoxError{status, "Unknown error"}
	}
}
