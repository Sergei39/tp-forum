package response

type Response interface {
	Code() int
	Body() interface{}
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
