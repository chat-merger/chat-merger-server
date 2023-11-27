# Golang setup:
export GOARCH := "amd64"
export GOOS := "linux"

# Api protobuf scheme:
api_version_tag := "v0.0.1a"
api_scheme_destination := "./api"
api_file_name := "mergerapi.proto"

generated_pb_files_destination := "./pkg"

default: build

build:
    go build -o ./bin/main ./main.go

init:
    install-deps
    just install-deps
    just get-api
    just gen-pb

gen-pb out=generated_pb_files_destination scheme=(api_scheme_destination+"/"+api_file_name):
    mkdir -p {{out}}
    protoc --go_out={{out}} --go-grpc_out={{out}} {{scheme}}

get-api tag=api_version_tag dest=api_scheme_destination file=api_file_name:
	./scripts/download-api-scheme.sh -t {{tag}} -d {{dest}} -f {{file}}

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0