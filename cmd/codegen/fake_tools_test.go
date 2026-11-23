package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// fakeProtocEnv marks a re-exec of the test binary as the fake protoc. The
// tests put a symlink to the test binary on PATH under the name "protoc", so
// the generator's protoc call runs without the real compiler and without any
// external command the restricted PATH would not provide.
const fakeProtocEnv = "CODEGEN_FAKE_PROTOC"

// fakeProtocNoServicesArg marks an include path that requests no services.
const fakeProtocNoServicesArg = "no-services-for-fake-protoc"

// TestMain serves as the fake protoc when the test binary is re-executed
// through the symlink that writeFakeProtoc installs.
func TestMain(m *testing.M) {
	if os.Getenv(fakeProtocEnv) == "1" {
		os.Exit(runFakeProtoc(os.Args[1:]))
	}
	os.Exit(m.Run())
}

// runFakeProtoc writes the prepared descriptor set to the path named by
// --descriptor_set_out and reports any problem on stderr.
func runFakeProtoc(args []string) int {
	out := ""
	noServices := slices.Contains(args, fakeProtocNoServicesArg)
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-I" && args[i+1] == fakeProtocNoServicesArg {
			noServices = true
		}
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "--descriptor_set_out=") {
			out = strings.TrimPrefix(arg, "--descriptor_set_out=")
		}
	}
	if out == "" {
		fmt.Fprintln(os.Stderr, "fake protoc: no --descriptor_set_out")
		return 1
	}
	descriptorSet := testDescriptorSet()
	if noServices {
		descriptorSet.File[0].Service = nil
	}
	data, err := proto.Marshal(descriptorSet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fake protoc: %v\n", err)
		return 1
	}
	if err := os.WriteFile(out, data, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "fake protoc: %v\n", err)
		return 1
	}
	return 0
}

// testDescriptorSet describes the proto the generator is asked to process.
func testDescriptorSet() *descriptorpb.FileDescriptorSet {
	str := func(s string) *string { return &s }

	return &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{{
			Name:    str("query.proto"),
			Package: str("test.v1"),
			Syntax:  str("proto3"),
			Options: &descriptorpb.FileOptions{GoPackage: str("example.com/test/v1")},
			MessageType: []*descriptorpb.DescriptorProto{
				{Name: str("PingRequest")},
				{Name: str("PingResponse")},
			},
			Service: []*descriptorpb.ServiceDescriptorProto{{
				Name: str("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{{
					Name:       str("Ping"),
					InputType:  str(".test.v1.PingRequest"),
					OutputType: str(".test.v1.PingResponse"),
					Options:    &descriptorpb.MethodOptions{},
				}},
			}},
		}},
	}
}

// writeFakeProtoc installs a protoc that re-executes the test binary in fake
// mode, so the generator's protoc call succeeds without the real compiler.
func writeFakeProtoc(t *testing.T, dir string) {
	t.Helper()

	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(self, filepath.Join(dir, "protoc")); err != nil {
		t.Fatal(err)
	}
}
