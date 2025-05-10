.PHONY: gen

gen:
	@echo "Generating sso protobuf"
	protoc -I proto proto/sso.proto \
		--go_out=gen/sso --go_opt=paths=source_relative \
		--go-grpc_out=gen/sso --go-grpc_opt=paths=source_relative