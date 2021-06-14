package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Code() int
	Body() interface{}
	SendSuccess(w http.ResponseWriter)
}

type response struct {
	httpCode int
	body     interface{}
}

func (r *response) Code() int {
	return r.httpCode
}

func (r *response) Body() interface{} {
	return r.body
}

func New(code int, body interface{}) Response {
	return &response{
		httpCode: code,
		body:     body,
	}
}

func (r *response) SendSuccess(w http.ResponseWriter) {
	body, err := json.Marshal(r.body)
	if err != nil {
		return
	}

	w.WriteHeader(r.httpCode)
	w.Write(body)
}
