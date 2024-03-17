//go:build tools
// +build tools

package install

// Tool imports to define the dependency, independent of generated code
// see also: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//
// With that go install protoc ... will use the same import version as the code generated from protoc
//
// To install the protoc addons for golang, grpc gateway and swagger run (only one time):
//
//     make proto-setup
//
// To use protoc tools for code generation out of proto model, run:
//
//     make proto-build
//
// If the make tool is not available on your platform, you may find at least the correct command line statements
// in the Makefile.
//

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
