# protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 auth.proto

# # reserve proxy (number 5 https://github.com/grpc-ecosystem/grpc-gateway)
# protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1 auth.proto
# protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 auth.proto

protoc --proto_path=api/proto/v1 \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
    --go_out ./pkg/api/v1\
    --go-grpc_out ./pkg/api/v1 service.proto

protoc --proto_path=api/proto/v1 \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
    --grpc-gateway_out ./pkg/api/v1 \
    --grpc-gateway_opt logtostderr=true\
    service.proto

protoc --proto_path=api/proto/v1 --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --swagger_out=logtostderr=true:api/swagger/v1 service.proto

# spectacle ./api/swagger/v1/file.swagger.json -t ./docs/api/swagger/