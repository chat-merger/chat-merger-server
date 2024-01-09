# Golang setup:
export GOARCH := "amd64"
export GOOS := "linux"

# Protoc plugins:
protoc_gen_go_version  := "v1.31.0"
protoc_gen_go_grpc_version  := "v1.3.0"

# Api protobuf scheme:
api_version_tag := "v0.0.1c"
api_scheme_destination := "./api"
api_file_name := "mergerapi.proto"

generated_pb_package_destination := "./internal/api"

# build settings:
build_output_file_or_directory := "./bin/server"

default: run

run:
    go run ./cmd/server/main.go \
    -clients-cfg ./clients.json \
    -grpc-port 32256 \
    -http-port 32255


build *FLAGS:
    go build -o {{build_output_file_or_directory}} {{FLAGS}}  ./cmd/server/main.go

# required:
#   0. go programming language;
#   1. proto compiler - protoc;
#   2. add to .bashrc PATH="$PATH:$(go env GOPATH)/bin"
init:
    go mod tidy
    just install-deps
    just get-api
    just gen-pb

gen-pb out=generated_pb_package_destination scheme=(api_scheme_destination+"/"+api_file_name):
    mkdir -p {{out}}
    protoc --go_out={{out}} --go-grpc_out={{out}} {{scheme}}

get-api tag=api_version_tag dest=api_scheme_destination file=api_file_name:
	./scripts/download-api-scheme.sh -t {{tag}} -d {{dest}} -f {{file}}

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@{{protoc_gen_go_version}}
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@{{protoc_gen_go_grpc_version}}