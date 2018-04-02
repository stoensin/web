package web

import (
	"testing"
)

func TestWeb(t *testing.T) {
	srv := NewServer("")
	srv.Get("/rest/hello", func(h *THandler) {
		h.RespondString("Hello, World")
		return
	})
	//srv.ShowRoute(true)
	srv.Listen(":8080")
}
