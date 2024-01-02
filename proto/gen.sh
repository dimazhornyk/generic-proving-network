protoc --proto_path="$1" \
      --go_opt=paths=source_relative \
      --go_out="$2" \
      --go-grpc_out="$2" \
      --go-grpc_opt=paths=source_relative \
      "$1"/*.proto
