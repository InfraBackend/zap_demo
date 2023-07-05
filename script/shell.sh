# 生成grpc代码
protoc --go_out=. --proto_path=. chat.proto   
protoc --go-grpc_out=. --proto_path=. chat.proto
protoc --go-grpc_out=require_unimplemented_servers=false:. --proto_path=. chat.proto

# 其他写法
protoc --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false[,other options...] \
