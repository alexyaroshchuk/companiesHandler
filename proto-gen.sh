#!/bin/bash

rm -rf proto/pb/*.go

protoc --go_out=. \
    --go-grpc_out=require_unimplemented_servers=false:. \
    proto/company.proto

protoc --go_out=. \
    --go-grpc_out=require_unimplemented_servers=false:. \
    proto/auth_service.proto

protoc --go_out=. \
    --go-grpc_out=require_unimplemented_servers=false:. \
    proto/user.proto