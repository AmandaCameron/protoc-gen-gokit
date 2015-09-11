package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/golang/glog"

	"examples/foo-service"
)

func main() {
	go grpcServer()
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8675", grpc.WithInsecure())
	if err != nil {
		glog.Fatal(err)
	}
	cli := foo_service.NewFooServiceClient(conn)

	http.Handle("/hello", kithttp.Server{
		Context:  ctx,
		Endpoint: foo_service.MakeEndpoint_FooService_SayHello(cli),

		DecodeRequestFunc: foo_service.Decode_FooService_SayHello,
		EncodeResponseFunc: func(wr http.ResponseWriter, data interface{}) error {
			return json.NewEncoder(wr).Encode(data)
		},

		Logger: log.NewLogfmtLogger(os.Stderr),
	})

	http.Handle("/count/to/", kithttp.Server{
		Context:  ctx,
		Endpoint: foo_service.MakeEndpoint_FooService_CountTo(cli),

		DecodeRequestFunc: foo_service.Decode_FooService_CountTo,
		EncodeResponseFunc: func(wr http.ResponseWriter, data interface{}) error {
			return json.NewEncoder(wr).Encode(data)
		},

		Logger: log.NewLogfmtLogger(os.Stderr),
	})

	http.ListenAndServe("localhost:8000", nil)
}

func grpcServer() {
	lis, err := net.Listen("tcp", ":8675")
	if err != nil {
		glog.Fatal(err)
	}

	srv := grpc.NewServer()

	foo_service.RegisterFooServiceServer(srv, &impl{})

	srv.Serve(lis)
}

type impl struct{}

func (i *impl) SayHello(ctx context.Context, req *foo_service.HelloRequest) (*foo_service.HelloResponse, error) {
	return &foo_service.HelloResponse{
		Response: "Hello " + req.Who + "!",
	}, nil
}

func (i *impl) CountTo(ctx context.Context, req *foo_service.CountToRequest) (*foo_service.CountToResponse, error) {
	resp := &foo_service.CountToResponse{}

	for i := int32(0); i < req.Target; i++ {
		resp.Response += fmt.Sprintf(" %d", i+1)
	}

	resp.Response = resp.Response[1:]

	return resp, nil
}
