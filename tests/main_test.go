package examples

import (
  "net/http"
  // "net/http/httptest"
  "testing"

  "github.com/gavv/httpexpect"
)

func TestFruits(t *testing.T) {
  // handler := FruitsHandler()

  // server := httptest.NewServer(handler)
  // defer server.Close()

  // e := httpexpect.New(t, server.URL)
  e := httpexpect.New(t, "http://192.168.21.90:8080")

  e.GET("/ping").
    Expect().
    Status(http.StatusOK).JSON().Object().ContainsKey("message").ValueEqual("message", "pong")
}