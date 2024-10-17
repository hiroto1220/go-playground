package server

import (
	"context"
	"log"

	pb "github.com/hiroto1220/go-playground/gRPC-gateway/gen"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// gRPC/gen/helloworld_grpc.pb.goに定義されている、GreeterServerインターフェースを実装する
func (s *GreeterServer) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())

	// 本来ならここから、別の層のメソッドを呼び出すなどする

	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}
