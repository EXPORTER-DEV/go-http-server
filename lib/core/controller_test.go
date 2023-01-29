package server

import (
	"bytes"
	"net/http"
	"testing"
)

type MockResponseWriter struct{}

func (response MockResponseWriter) WriteHeader(status int) {

}

func (response MockResponseWriter) Status(status int) {

}

func (response MockResponseWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (response MockResponseWriter) Header() http.Header {
	return *new(http.Header)
}

func TestController(t *testing.T) {
	t.Run("Should send status and write response", func(t *testing.T) {
		mockResponse := MockResponseWriter{}

		controller := &Controller{
			response: mockResponse,
			request:  new(http.Request),
		}

		controller.Status(200)

		controller.Send("TEST")

		if controller.status != 200 {
			t.Fatalf("%v expected to be %v", controller.status, 200)
		}

		if !bytes.Equal(controller.content, []byte("TEST")) {
			t.Fatal("Response bytes are not correct")
		}
	})
}
