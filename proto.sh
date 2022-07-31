DOMAIN=.
PROTO_PATH=${DOMAIN}/api/v1
OUT_PATH=./${DOMAIN}/api/v1
PROTO_FILE=user.proto


protoc -I=$PROTO_PATH --go_out $OUT_PATH --go_opt paths=source_relative --go-grpc_out $OUT_PATH --go-grpc_opt=paths=source_relative $PROTO_FILE