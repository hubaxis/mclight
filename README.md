# MemCached light
## Generate proto:

protoc --proto_path=protocol protocol/*.proto --go_out=./protocol --go-grpc_out=./protocol

## test
read https://github.com/vektra/mockery to generate mocks


