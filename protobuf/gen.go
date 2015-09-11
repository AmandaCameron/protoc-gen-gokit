package google_apimv

//go:generate protoc --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. -I /usr/local/include:. google/api/annotations.proto google/api/http.proto
