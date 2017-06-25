package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type (
	Arithmetic   struct{}
	MultiplyArgs struct {
		A, B int
	}
	MultiplyResult struct {
		Result int
	}
)

func (a *Arithmetic) Multiply(args *MultiplyArgs, rely *MultiplyResult) error {
	rely.Result = args.A * args.B
	return nil
}

func init() {
	s := rpc.NewServer()
	arithmetic := &Arithmetic{}
	s.Register(arithmetic)

	http.HandleFunc("/jrpc", func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			log.Fatal(err)
		}
		s.ServeCodec(jsonrpc.NewServerCodec(conn))
	})

	go http.ListenAndServe(":8080", nil)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(conn, "CONNECT "+"/jrpc"+" HTTP/1.0\n\n")
	if err != nil {
		log.Fatalln(err)
	}

	cli := jsonrpc.NewClient(conn)
	var ret MultiplyResult
	err = cli.Call("Arithmetic.Multiply", &MultiplyArgs{A: 5, B: 5}, &ret)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	log.Println(ret.Result)
}
