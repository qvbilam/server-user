function userProto {
  DOMAIN=.
  PROTO_FILE=$1.proto
  VERSION=$2
  PROTO_PATH=${DOMAIN}/api/qvbilam/user/"$VERSION"
  OUT_PATH=./${DOMAIN}/api/qvbilam/user/"$VERSION"
  protoc -I="$PROTO_PATH" --go_out "$OUT_PATH" --go_opt paths=source_relative --go-grpc_out $OUT_PATH --go-grpc_opt=paths=source_relative "$PROTO_FILE"
}

userProto user v1
