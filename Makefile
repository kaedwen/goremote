# Project Directories
BUILD_ROOT?=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	go build -mod=vendor -ldflags "-s -w" -trimpath -o service main.go

.PHONY: proto-setup
proto-setup:
	@go install \
	    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    	google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: proto-build
proto-build:
	$(call grpc_proto_build,pkg/api/v1/proto)

define grpc_proto_build
	$(eval REPO_PATH_PROTO := $(1))
	$(eval REPO_PATH_VERSION := $(shell dirname $(REPO_PATH_PROTO)))
	$(eval REPO_PATH_SERVICE := $(shell dirname $(REPO_PATH_VERSION)))
	$(eval REPO_PATH_GENBASE := $(shell dirname $(REPO_PATH_SERVICE)))

	@mkdir -p $(BUILD_ROOT)/$(REPO_PATH_VERSION)/gen
	@cd $(BUILD_ROOT);

	@protoc \
	  --proto_path=$(REPO_PATH_GENBASE) \
	  --go_out=$(BUILD_ROOT)/$(REPO_PATH_VERSION) \
	  --go-grpc_out=$(BUILD_ROOT)/$(REPO_PATH_VERSION) \
	  $(REPO_PATH_PROTO)/*.proto

	@protoc \
	  --proto_path=$(REPO_PATH_GENBASE) \
	  --grpc-gateway_out=$(BUILD_ROOT)/$(REPO_PATH_VERSION) \
	  --grpc-gateway_opt logtostderr=true \
	  --grpc-gateway_opt grpc_api_configuration=$(REPO_PATH_PROTO)/service.yaml \
	  $(REPO_PATH_PROTO)/*.proto

	@protoc \
	  --proto_path=$(REPO_PATH_GENBASE) \
	  --openapiv2_out $(BUILD_ROOT)/$(REPO_PATH_VERSION)/gen \
	  --openapiv2_opt allow_merge=true \
	  --openapiv2_opt merge_file_name=model \
	  --openapiv2_opt logtostderr=true \
	  --openapiv2_opt grpc_api_configuration=$(REPO_PATH_PROTO)/service.yaml \
	  --openapiv2_opt openapi_configuration=$(REPO_PATH_PROTO)/service.swagger.yaml \
	  $(REPO_PATH_PROTO)/*.proto
endef