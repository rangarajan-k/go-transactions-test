package utilities

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Err        error  `json:"error"`
}

func (e *ResponseError) Error() string { return e.Message }
func (e *ResponseError) Unwrap() error { return e.Err }
func (e *ResponseError) Code() int     { return e.StatusCode }

func NewResponseError(statusCode int, message string, err error) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}
