.DEFAULT_GOAL := gen

gen:
	protoc --go_out=. --go-grpc_out=. api-scheme-proto/mergerapi.proto

