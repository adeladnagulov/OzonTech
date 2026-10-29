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
	go generate ./...

.PHONY: run-product
run-product:
	go run ./cmd/product-service

.PHONY: run-cart
run-cart:
	go run ./cmd/cart

.PHONY: run-all
run-all:
	go run ./cmd/product-service &
	go run ./cmd/cart

.PHONY: test
test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html