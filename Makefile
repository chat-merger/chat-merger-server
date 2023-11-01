.DEFAULT_GOAL := gen

gen:
	protoc --go_out=. --go-grpc_out=. proto/mergerapi.proto

