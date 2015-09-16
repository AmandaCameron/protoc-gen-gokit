// Generated by protoc-gen-gokit DO NOT EDIT.
package foo_service

import (
  "net/http"
  "errors"
  "strings"

  "golang.org/x/net/context"

  kithttp "github.com/go-kit/kit/transport/http"
  "github.com/go-kit/kit/endpoint"

  "github.com/AmandaCameron/protoc-gen-gokit/runtime"


)



// MakeMux_FooService creates a server mux for the FooService service, 
// using the passed kithttp.Server as a template for the parameters of the endpoints.
func MakeMux_FooService(cli FooServiceClient, mw endpoint.Middleware, responseEncoder kithttp.EncodeResponseFunc, error options ...kithttp.ServerOption) (http.Handler, error) {
  ret := runtime.NewMux()


  ret.AddEndpoint("GET", "/hello", kithttp.NewServer(
   context.Background(), 
   mw(MakeEndpoint_FooService_SayHello(cli)),
   Decode_FooService_SayHello,
   responseEncoder, options...)
  )
  ret.AddEndpoint("GET", "/count/to/{target}", kithttp.NewServer(
   context.Background(), 
   mw(MakeEndpoint_FooService_CountTo(cli)),
   Decode_FooService_CountTo,
   responseEncoder, options...)
  )

  return ret, nil
}


// Decode_FooService_SayHello decodes an http.Request into a HelloRequest.
func Decode_FooService_SayHello(req *http.Request) (interface{}, error) {
  var ret HelloRequest

  qry := req.URL.Query()
  _ = qry


  if val := qry.Get("who"); val != "" {
    if err := runtime.Decode(&ret.Who, val); err != nil {
      return nil, err
    }
  }

  parts := strings.Split(req.URL.Path, "/")
  if len(parts) < 2 {
    return nil, errors.New("Missing Parameters.")
  }



  return &ret, nil
}

// MakeEndpoint_FooService_SayHello creates an endpoint function for Go-kit 
// that runs the specified service / endpoint on the specified grpc endpoint.
func MakeEndpoint_FooService_SayHello(cli FooServiceClient) endpoint.Endpoint {
  endp := func (ctx context.Context, inp interface{}) (interface{}, error) {
    return cli.SayHello(ctx, inp.(*HelloRequest))
  }

  return endp
}

// Decode_FooService_CountTo decodes an http.Request into a CountToRequest.
func Decode_FooService_CountTo(req *http.Request) (interface{}, error) {
  var ret CountToRequest

  qry := req.URL.Query()
  _ = qry


  if val := qry.Get("target"); val != "" {
    if err := runtime.Decode(&ret.Target, val); err != nil {
      return nil, err
    }
  }

  parts := strings.Split(req.URL.Path, "/")
  if len(parts) < 4 {
    return nil, errors.New("Missing Parameters.")
  }


  if err := runtime.Decode(&ret.Target, parts[3]); err != nil {
    return nil, err
  }

  return &ret, nil
}

// MakeEndpoint_FooService_CountTo creates an endpoint function for Go-kit 
// that runs the specified service / endpoint on the specified grpc endpoint.
func MakeEndpoint_FooService_CountTo(cli FooServiceClient) endpoint.Endpoint {
  endp := func (ctx context.Context, inp interface{}) (interface{}, error) {
    return cli.CountTo(ctx, inp.(*CountToRequest))
  }

  return endp
}
