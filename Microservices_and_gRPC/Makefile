.PHONY: version
version:
	go version
	protoc --version

.PHONY: generate
generate:
	protoc -I ./ \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	api/product.proto

	protoc-go-inject-tag -input="./pkg/api/*.pb.go"	

.PHONY: run-product
run-product:
	go run ./cmd/product-service