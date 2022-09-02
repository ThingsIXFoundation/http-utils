package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/gddo/httputil/header"
	"github.com/sirupsen/logrus"
)

// BadHTTPRequest represents an invalid HTTP request.
type BadHTTPRequest struct {
	Status int    // HTTP error reason code
	Msg    string // Description what is wrong
}

// Error returns a brief description what was wrong with the request.
func (mr *BadHTTPRequest) Error() string {
	return mr.Msg
}

// DecodeHTTPJSONBody JSON decodes the request body into `to`.
// In case of problems an error is returned.
// It only accepts request payloads that are smaller than 1MB.
func DecodeHTTPJSONBody(w http.ResponseWriter, r *http.Request, to interface{}) *BadHTTPRequest {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &BadHTTPRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	// Accept max 1MiB JSON body payloads
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // pedantic but prevents errors that will otherwise probably not be detected

	if err := dec.Decode(&to); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &BadHTTPRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			logrus.WithError(err).Error("unable to decode json request")
			return &BadHTTPRequest{Status: http.StatusInternalServerError, Msg: http.StatusText(http.StatusInternalServerError)}
		}
	}

	// ensure that there is nothing left after the parsed payload
	if err := dec.Decode(&struct{}{}); err == io.EOF {
		return nil
	}

	msg := "Request body must only contain a single JSON object"
	return &BadHTTPRequest{Status: http.StatusBadRequest, Msg: msg}
}

func ReplyJSON(w http.ResponseWriter, r *http.Request, statusCode int, reply interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if isEmpty(reply) {
		if _, err := w.Write([]byte("[]")); err != nil {
			if rid, ok := r.Context().Value(middleware.RequestIDKey).(string); ok && rid != "" {
				logrus.WithField("request_id", rid).WithError(err).Error("could not send http JSON reply")
			} else {
				logrus.WithError(err).Error("could not send http JSON reply")
			}
		}
		return
	}
	if !isNil(reply) {
		if err := json.NewEncoder(w).Encode(reply); err != nil {
			if rid, ok := r.Context().Value(middleware.RequestIDKey).(string); ok && rid != "" {
				logrus.WithField("request_id", rid).WithError(err).Error("could not send http JSON reply")
			} else {
				logrus.WithError(err).Error("could not send http JSON reply")
			}
		}
		return
	}
}

// isNil returns an indication if the given value is nil.
// Or in case its a nullable type the value it is pointing to is nil.
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// isEmpty returns true if the given i is an empty array or empty slice.
func isEmpty(i interface{}) bool {
	if i == nil {
		return false
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Array, reflect.Slice:
		return reflect.ValueOf(i).IsNil() || reflect.ValueOf(i).Len() == 0
	}
	return false
}
