package common

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Response type
type Response struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Status  int               `json:"status,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Headers map[string]string `json:"-"`
	Cookies []*http.Cookie    `json:"-"`
	Raw     bool              `json:"-"`
}

// ResponseCtor function
type ResponseCtor func(context.Context) *Response

// ResponseSetter setter
type ResponseSetter func(*Response) *Response

// NewResponse constructor
func NewResponse(ctx context.Context, status int) *Response {
	var message = http.StatusText(status)
	return &Response{
		Status:  status,
		Code:    mapCode(message),
		Message: message,
		Headers: map[string]string{},
	}
}

func mapCode(code string) string {
	return strings.Join(strings.Split(strings.ToUpper(code), " "), "_")
}

// WrapResponse error constructor
func WrapResponse(ctx context.Context, err error) *Response {
	return WrapResponse2(ctx, err, ResponseInternalServerError)
}

// WrapResponse2 error constructor
func WrapResponse2(ctx context.Context, err error, ctor ResponseCtor, setters ...ResponseSetter) *Response {
	if err == nil {
		return nil
	}

	switch v := err.(type) {
	case *Response:
		return v
	default:
		log.Printf("%+v\r\n", err)
		var resp = ctor(ctx)
		for _, s := range setters {
			resp = s(resp)
		}
		return resp
	}
}

// SetData method
func (r *Response) SetData(value interface{}) *Response {
	r.Data = value
	return r
}

// SetRaw method
func (r *Response) SetRaw(value bool) *Response {
	r.Raw = value
	return r
}

// SetCode method
func (r *Response) SetCode(values ...string) *Response {
	if len(values) == 1 {
		values[0] = mapCode(values[0])
	} else {
		for k, s := range values {
			values[k] = mapCode(s)
		}
	}
	r.Code = strings.Join(values, "_")
	return r
}

// SetMessage method
func (r *Response) SetMessage(args ...interface{}) *Response {
	r.Message = fmt.Sprint(args...)
	return r
}

// SetMessagef method
func (r *Response) SetMessagef(format string, args ...interface{}) *Response {
	r.Message = fmt.Sprintf(format, args...)
	return r
}

// SetHeader method
func (r *Response) SetHeader(key string, value string) *Response {
	r.Headers[key] = value
	return r
}

// SetCookie method
func (r *Response) SetCookie(cookies ...*http.Cookie) *Response {
	r.Cookies = append(r.Cookies, cookies...)
	return r
}

// IsError getter
func (r *Response) IsError() bool {
	return r != nil && r.Status >= 400
}

// Error implementation
func (r *Response) Error() string {
	return r.String()
}

// String implementation
func (r *Response) String() string {
	if r == nil {
		return "<nil>"
	}
	if r.Code == r.Message {
		return fmt.Sprintf("%s (%d)", r.Code, r.Status)
	}
	return fmt.Sprintf("%s (%d) %s", r.Code, r.Status, r.Message)
}

// Clone without Data
func (r *Response) Clone() *Response {
	var c = &Response{}
	// Copy all properties
	*c = *r
	// Reset Data
	c.Data = nil
	return c
}

// ResponseOK constructor 200
func ResponseOK(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusOK)
}

// ResponseCreated constructor 201
func ResponseCreated(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusCreated)
}

// ResponseAccepted constructor 202
func ResponseAccepted(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusAccepted)
}

// ResponseNoContent constructor 204
func ResponseNoContent(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusNoContent)
}

// ResponseBadRequest constructor 400
func ResponseBadRequest(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusBadRequest)
}

// ResponseUnauthorized constructor 401 (Unauthenticated)
func ResponseUnauthorized(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusUnauthorized)
}

// ResponseForbidden constructor 403 (Unauthorized)
func ResponseForbidden(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusForbidden)
}

// ResponseNotFound constructor 404
func ResponseNotFound(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusNotFound)
}

// ResponseMethodNotAllowed constructor 405
func ResponseMethodNotAllowed(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusMethodNotAllowed)
}

// ResponseRequestTimeout constructor 408
func ResponseRequestTimeout(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusRequestTimeout)
}

// ResponseInternalServerError constructor 500
func ResponseInternalServerError(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusInternalServerError)
}

// ResponseNotImplemented constructor 501
func ResponseNotImplemented(ctx context.Context) *Response {
	return NewResponse(ctx, http.StatusNotImplemented)
}
