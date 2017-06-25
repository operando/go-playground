package main

import (
	"context"
	"encoding/json"
	"github.com/osamingo/jsonrpc"
	"log"
	"net/http"
)

type (
	EchoHandler struct{}
	EchoParams  struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}
)

var _ (jsonrpc.Handler) = (*EchoHandler)(nil)

func (h *EchoHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var p EchoParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func init() {
	jsonrpc.RegisterMethod("Main.Echo", &EchoHandler{}, EchoParams{}, EchoResult{})
}

func main() {
	http.HandleFunc("/jrpc", jsonrpc.HandlerFunc)
	http.HandleFunc("/jrpc/debug", jsonrpc.DebugHandlerFunc)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// curl -d '{"jsonrpc": "2.0","method":"Main.Echo","params":{"name":"test"},"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"}' -H "Content-Type: application/json" localhost:8080/jrpc


