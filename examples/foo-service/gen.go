package foo_service

//go:generate protoc --go_out=plugins=grpc:. --gokit_out=. simple-service.proto -I $HOME/Bazel/darkdna/infra/lib/protobuf:/usr/local/include:.
