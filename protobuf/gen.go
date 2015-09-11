package google_api

//go:generate protoc --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/annotations.proto google/api/http.proto
