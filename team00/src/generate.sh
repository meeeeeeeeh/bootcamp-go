#!/bin/bash
PROTO_DIR="./proto"
OUT_DIR="./generated"

protoc --go_out=paths=source_relative:${OUT_DIR} \
       --go-grpc_out=paths=source_relative:${OUT_DIR} \
       --proto_path=${PROTO_DIR} \
       ${PROTO_DIR}/randomaliens.proto

protoc --doc_out=. --doc_opt=markdown,README.md --proto_path=${PROTO_DIR} ${PROTO_DIR}/randomaliens.proto