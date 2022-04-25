proto:
	@protoc -I resources/proto/ resources/proto/deposit.proto  --go_out=resources/proto/  --go-grpc_out=require_unimplemented_servers=false:resources/proto/