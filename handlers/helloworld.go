package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"os"
)

type HelloWorldHandler struct {
	log *zap.SugaredLogger
}

func NewHelloWorldHandler(log *zap.SugaredLogger) *HelloWorldHandler {
	return &HelloWorldHandler{log: log}
}

func (*HelloWorldHandler) Pattern() string {
	return "/hello"
}

func (h *HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	h.log.Debugf("Someone has came to say hello to us... fuck %v", body)
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		h.log.Errorf("Thank god we failed our response... i mean whoops Error writing response: %v %v", err, os.Stderr)
	}
}
