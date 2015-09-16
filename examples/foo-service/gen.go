package foo_service

//go:generate protoc --go_out=plugins=grpc:. --gokit_out=. -I ../../protobuf:. simple-service.proto
