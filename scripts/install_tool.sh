go mod tidy
go install \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
  github.com/envoyproxy/protoc-gen-validate \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1

go install  google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0