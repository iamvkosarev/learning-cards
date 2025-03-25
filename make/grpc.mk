BIN_DIR := $(ROOT_DIR)/bin

PROTO_DIR := $(ROOT_DIR)/api/proto
PROTO_GEN_DIR := $(ROOT_DIR)/pkg/proto
OPENAPI_DIR := $(ROOT_DIR)/docs/openapiv2
GOOGLEAPIS_DIR := $(ROOT_DIR)/third_party/googleapis
VALIDATE_DIR := $(ROOT_DIR)/third_party/protoc-gen-validate

PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")


.PHONY: googleapis validate-proto grpc-deps grpc-generate grpc-gateway grpc-swagger grpc-validate
grpc-all: googleapis validate-proto grpc-deps grpc-generate grpc-gateway grpc-swagger grpc-validate

googleapis:
	@if [ ! -d "$(GOOGLEAPIS_DIR)/google/api" ]; then \
		echo "Cloning googleapis into $(GOOGLEAPIS_DIR)..."; \
		git clone --depth=1 https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR); \
	else \
		echo "googleapis already exists in $(GOOGLEAPIS_DIR)"; \
	fi

validate-proto:
	@if [ ! -d "$(VALIDATE_DIR)/validate" ]; then \
		echo "Cloning protoc-gen-validate protos into $(VALIDATE_DIR)..."; \
		git clone --depth=1 https://github.com/envoyproxy/protoc-gen-validate.git $(VALIDATE_DIR); \
	else \
		echo "protoc-gen-validate protos already exist in $(VALIDATE_DIR)"; \
	fi

grpc-deps:
	@echo "Installing gRPC dependencies into $(BIN_DIR)"
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@GOBIN=$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@latest
	@echo "gRPC dependencies installed into $(BIN_DIR)"

grpc-generate:
	@echo "Generating gRPC .pb.go files..."
	@mkdir -p $(PROTO_GEN_DIR)
	@PATH=$(BIN_DIR):$$PATH \
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS_DIR) \
		--proto_path=$(VALIDATE_DIR) \
		--go_out=$(PROTO_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "gRPC-Gateway .gw.go files generated"

grpc-gateway:
	@echo "Generating gRPC-Gateway .gw.go files..."
	@PATH=$(BIN_DIR):$$PATH \
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS_DIR) \
		--proto_path=$(VALIDATE_DIR) \
		--grpc-gateway_out=$(PROTO_GEN_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		$(PROTO_FILES)

grpc-swagger:
	@echo "Generating OpenAPI (Swagger)..."
	@mkdir -p $(OPENAPI_DIR)
	@PATH=$(BIN_DIR):$$PATH \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--proto_path=$(GOOGLEAPIS_DIR) \
			--proto_path=$(VALIDATE_DIR) \
			--openapiv2_out=$(OPENAPI_DIR) \
			$(PROTO_FILES)
	@echo "Swagger files generated in $(OPENAPI_DIR)"

grpc-validate:
	@echo "Generating validation code..."
	@PATH=$(BIN_DIR):$$PATH \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--proto_path=$(GOOGLEAPIS_DIR) \
			--proto_path=$(VALIDATE_DIR) \
			--validate_out=lang=go,paths=source_relative:$(PROTO_GEN_DIR) \
			$(PROTO_FILES)
	@echo "Validation code generated"
