# Golang setup:
export GOARCH := "amd64"
export GOOS := "linux"
export GOPATH := `go env GOPATH`

# Api protobuf scheme:
api_version_tag := "v0.0.1a"
api_scheme_destination := "./api"
api_file_name := "mergerapi.proto"

default: build

build:
    go build -o ./bin/main ./main.go

protogen file=(api_scheme_destination + api_file_name):
    protoc --go_out=. --go-grpc_out={{file}}

get-api tag=api_version_tag dest=api_scheme_destination file=api_file_name:
	./scripts/download-api-scheme.sh -t {{tag}} -d {{dest}} -f {{file}}

