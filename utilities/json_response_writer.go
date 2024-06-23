package utilities

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ErrorDetails struct {
	Code    int    `json:"code,omitempty" example:"400001"`
	Title   string `json:"title,omitempty" example:"Validation error!"`
	Message string `json:"message,omitempty" example:"Invalid password"`
}

type ErrorResponse struct {
	Success bool         `json:"success"`
	Result  interface{}  `json:"result,omitempty"`
	Failure ErrorDetails `json:"failure"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type HttpJSONResponse interface {
	// success responses
	WriteJSONWithStatus(r *http.Request, w http.ResponseWriter, statusCode int, v interface{}) error
	WriteJSON(r *http.Request, w http.ResponseWriter, v interface{}) error
	WriteWithStatus(w http.ResponseWriter, statusCode int)

	// error responses
	WriteError(r *http.Request, w http.ResponseWriter, result interface{}, err error, messages ...string) error
	// use custom message instead of error.Error()
	WriteErrorWithMessage(r *http.Request, w http.ResponseWriter, result interface{}, err error, message string, messages ...string) error

	// error responses
	ErrorAsJSON(r *http.Request, w http.ResponseWriter, err error, messages ...string) error
	JSON(r *http.Request, w http.ResponseWriter, result interface{}, v interface{}, statusCode ...int) error
}

type httpJSONResponseWriter struct {
}

func NewHttpResponseService() HttpJSONResponse {
	return &httpJSONResponseWriter{}
}

func (c *httpJSONResponseWriter) JSON(r *http.Request, w http.ResponseWriter, result interface{}, v interface{}, statusCodeOpt ...int) error {
	statusCode := http.StatusOK
	if len(statusCodeOpt) != 0 {
		statusCode = statusCodeOpt[0]
	}
	customError, ok := v.(*ResponseError)
	if ok {
		return c.WriteErrorWithCode(r, w, result, customError, statusCode, customError.Key(), customError.Message())
	}
	err, ok := v.(error)
	if ok {
		return c.WriteErrorWithCode(r, w, result, err, statusCode)
	}

	w.WriteHeader(statusCode)
	return c.WriteJSON(r, w, v)
}

func (c *httpJSONResponseWriter) WriteJSONWithStatus(r *http.Request, w http.ResponseWriter, statusCode int, result interface{}) error {
	w.WriteHeader(statusCode)
	return c.marshalToJSON(r, w, result)
}

func (c *httpJSONResponseWriter) WriteJSON(r *http.Request, w http.ResponseWriter, v interface{}) error {
	var result interface{}
	resp := SuccessResponse{
		Success: true,
		Result:  v,
	}
	result = resp

	return c.marshalToJSON(r, w, result)
}

// Write interface value as json
func (c *httpJSONResponseWriter) marshalToJSON(r *http.Request, w http.ResponseWriter, v interface{}) error {

	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("RESPONSE_BODY_MARSHALLING_ERROR")
		return err
	}
	log.Printf("RESPONSE_DATA: %s", string(data))
	_, err = w.Write(data)
	return err
}

func (c *httpJSONResponseWriter) WriteWithStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

// argument error object need to be CustomError type.
// if not, this function with return with 500 status code as default.
// both messages title and description should be translated manually.
func (c *httpJSONResponseWriter) WriteError(r *http.Request, w http.ResponseWriter, result interface{}, err error,
	messages ...string) error {
	ctx := r.Context()
	log := logging.GetStructuredLogger(ctx)

	var resp interface{}

	statusCode := http.StatusInternalServerError
	code := getErrorCode(err)
	title, desc := getTitleAndDescription(messages)
	statusCode, resp = c.getErrorResponse(code, title, err.Error(), desc, result)
	w.WriteHeader(statusCode)
	return c.marshalToJSON(r, w, resp)
}

func (c *httpJSONResponseWriter) WriteErrorWithCode(r *http.Request, w http.ResponseWriter, result interface{},
	err error, code int, messages ...string) error {

	var resp interface{}

	statusCode := http.StatusInternalServerError
	title, desc := getTitleAndDescription(messages)
	statusCode, resp = c.getErrorResponse(code, title, desc, err.Error(), result)
	w.WriteHeader(statusCode)
	return c.marshalToJSON(r, w, resp)
}

// argument error object need to be CustomError type.
// if not, this function with return with 500 status code as default.
// both messages title and description should be translated manually.
func (c *httpJSONResponseWriter) ErrorAsJSON(r *http.Request, w http.ResponseWriter, err error, messages ...string) error {
	ctx := r.Context()
	log := logging.GetStructuredLogger(ctx)

	var resp interface{}
	statusCode := http.StatusInternalServerError
	errorWithStatusCode, ok := err.(*ResponseError)
	if ok {
		statusCode = errorWithStatusCode.Code()
	}

	message := err.Error()
	errorWithFormatter, ok := err.(cerrors.ErrorFormatter)
	if ok {
		message = errorWithFormatter.FormattedMessage()
	}

	title := ""
	errorWithName, ok := err.(cerrors.ErrorName)
	if ok {
		title = errorWithName.Name()
	}

	customTitle, desc := getTitleAndDescription(messages)
	if customTitle != "" {
		title = customTitle
	}

	resp = &ErrorResponse{
		Success: false,
		Failure: ErrorDetails{
			Title:   title,
			Message: message,
		},
	}
	w.WriteHeader(statusCode)
	log.WithError(err).Errorf("HTTP_ERROR_RESPONSE::[%d]", statusCode)
	return c.marshalToJSON(r, w, resp)
}

func (c *httpJSONResponseWriter) getErrorResponse(code int, title string, message string, description string, result interface{}) (int, interface{}) {
	statusCode := getHttpStatus(code)
	resp := &ErrorResponse{
		Success: false,
		Failure: ErrorDetails{
			Code:    code,
			Title:   title,
			Message: message,
		},
	}

	if result != nil {
		resp.Result = result
	}

	return statusCode, resp
}

func (c *httpJSONResponseWriter) WriteErrorWithMessage(r *http.Request, w http.ResponseWriter, result interface{},
	err error, message string, messages ...string) error {
	ctx := r.Context()
	log := logging.GetStructuredLogger(ctx)
	var resp interface{}

	statusCode := http.StatusInternalServerError
	code := getErrorCode(err)
	title, desc := getTitleAndDescription(messages)

	statusCode, resp = c.getErrorResponse(code, title, message, desc, result)

	w.WriteHeader(statusCode)
	return c.marshalToJSON(r, w, resp)
}

// error is not CustomError type return default error code (Internal Server Error)
func getErrorArgsAndMessage(err error) ([]interface{}, string) {
	errorWithFormat, ok := err.(cerrors.ErrorFormatter)
	if !ok {
		return []interface{}{}, err.Error()
	}
	return errorWithFormat.GetArgs(), errorWithFormat.GetMessage()
}

// error is not CustomError type return default error code (Internal Server Error)
func getErrorCode(err error) int {
	errorWithCode, ok := err.(cerrors.ErrorCode)
	if !ok {
		return errors.UnknownErrorCode
	}
	return errorWithCode.Code()
}

func getHttpStatus(code int) (status int) {
	firstThreeDigits := "500"
	codeArr := strconv.Itoa(code)
	if len(codeArr) >= 3 {
		firstThreeDigits = codeArr[:3]
	}
	switch firstThreeDigits {
	case "200":
		status = http.StatusOK
	case "202":
		status = http.StatusAccepted
	case "400":
		status = http.StatusBadRequest
	case "401":
		status = http.StatusUnauthorized
	case "403":
		status = http.StatusForbidden
	case "404":
		status = http.StatusNotFound
	case "405":
		status = http.StatusMethodNotAllowed
	case "406":
		status = http.StatusNotAcceptable
	case "408":
		status = http.StatusRequestTimeout
	case "422":
		status = http.StatusUnprocessableEntity
	case "423":
		status = http.StatusLocked
	case "424":
		status = http.StatusFailedDependency
	default:
		status = http.StatusInternalServerError
	}
	return
}

func getTitleAndDescription(messages []string) (string, string) {
	var ttl, desc string
	if len(messages) > 0 {
		ttl = messages[0]
	}
	if len(messages) > 1 {
		desc = messages[1]
	}
	return ttl, desc
}
