package foo_service

//go:generate protoc --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/AmandaCameron/protoc-gen-gokit/protobuf/google/api:. --gokit_out=Mgoogle/api/annotations.proto=github.com/AmandaCameron/protoc-gen-gokit/protobuf/google/api:. -I ../../protobuf:. simple-service.proto
