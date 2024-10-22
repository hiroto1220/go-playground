package handler

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	greet "github.com/hiroto1220/go-playground/connect/gen"
)

type GreetHandler struct{}

func NewGreetHandler() *GreetHandler {
	return &GreetHandler{}
}

func (s *GreetHandler) SayHello(ctx context.Context, req *connect.Request[greet.HelloRequest]) (*connect.Response[greet.HelloResponse], error) {
	log.Println("Request headers: ", req.Header())
	// 本来ならここから、別の層のメソッドを呼び出すなどする

	res := connect.NewResponse(&greet.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
