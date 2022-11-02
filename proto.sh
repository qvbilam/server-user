function proto {
  # 脚本所在路径
  SERVER=$1
  SCRIPT_DIR=$(cd $(dirname "$0");pwd)
  DOMAIN=.
  PROTO_FILE=$2.proto
  VERSION=$3
  PROTO_PATH=${DOMAIN}/api/qvbilam/${SERVER}/"$VERSION"
  OUT_PATH=./${DOMAIN}/api/qvbilam/${SERVER}/"$VERSION"
  # 引入项目目录下不同级的 proto 需要指定参数 --proto_path=绝对路径
  protoc -I="$PROTO_PATH" --go_out "$OUT_PATH" --go_opt paths=source_relative --go-grpc_out "$OUT_PATH" --go-grpc_opt=paths=source_relative "$PROTO_FILE" --proto_path="$SCRIPT_DIR"
}

# 分页服务
proto page page v1
# 用户服务
proto user user v1
proto user account v1