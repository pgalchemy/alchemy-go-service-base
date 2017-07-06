package errors

type (
	// HTTPError defines the structure of common HTTP errors
	HTTPError struct {
		StatusCode int                    `json:"statusCode"`
		Message    string                 `json:"message"`
		Code       string                 `json:"code"`
		Data       map[string]interface{} `json:"data,omitempty"`
	}
)

func newErr(s int, m, c string, d []map[string]interface{}) *HTTPError {
	err := &HTTPError{StatusCode: s, Message: m, Code: c}
	if len(d) > 0 {
		err.Data = d[0]
	}
	return err
}

// Error satisfies the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// InvalidContentError when a generic client error occurs
func InvalidContentError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(400, message, "Bad Request", data)
}

// UnauthorizedError should be used when a client is not authorized to perform the request
func UnauthorizedError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(401, message, "Unauthorized", data)
}

// ForbiddenError should be used when an authorized client is forbidden
func ForbiddenError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(403, message, "Forbidden", data)
}

// ResourceNotFoundError should be used when a requested resource cannot be located
func ResourceNotFoundError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(404, message, "Not Found", data)
}

// MissingParameterError should be used when a request is missing a required field
func MissingParameterError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(409, message, "Missing Parameter", data)
}

// ConflictError should be used when the requested resource already exists
func ConflictError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(409, message, "Conflict", data)
}

// InternalServerError should be used when a server error occurs
func InternalServerError(message string, data ...map[string]interface{}) *HTTPError {
	return newErr(500, message, "Internal Server Error", data)
}
