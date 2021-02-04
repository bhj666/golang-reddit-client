package test

import "net/http"

type ResponseWriterMock struct {
	Data string
	Code int
}

func (*ResponseWriterMock) Header() http.Header {
	return http.Header{}
}

func (m *ResponseWriterMock) Write(arg []byte) (int, error) {
	m.Data = string(arg)
	return 1, nil
}

func (m *ResponseWriterMock) WriteHeader(statusCode int) {
	m.Code = statusCode
}
