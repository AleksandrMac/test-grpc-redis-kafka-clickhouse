USER=user

.PHONY: generate

generate: proto-gen

proto-gen: user-gen

user-gen:
	mkdir -p ./pkg/${USER}/grpc/userservice
	protoc -I/usr/local/include -I . \
	--go_out ./pkg/${USER}/grpc/userservice \
	--go_opt paths=import \
	--go-grpc_out=./pkg/${USER}/grpc/userservice \
	./api/protobuf/user.proto