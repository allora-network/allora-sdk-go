package allora

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestNewClient(t *testing.T) {
	logger := zerolog.Nop() // No-op logger for tests

	tests := []struct {
		name        string
		config      *ClientConfig
		expectError bool
	}{
		{
			name:        "nil config should use default",
			config:      nil,
			expectError: true, // Will fail because no endpoints are provided in default config
		},
		{
			name: "empty endpoints should fail",
			config: &ClientConfig{
				Endpoints: []EndpointConfig{},
			},
			expectError: true,
		},
		{
			name: "invalid endpoint should fail",
			config: &ClientConfig{
				Endpoints: []EndpointConfig{
					{
						URL:      "grpc://nonexistent:9090",
						Protocol: "grpc",
					},
				},
			},
			expectError: true, // Will fail to connect
		},
		{
			name: "unsupported protocol should fail",
			config: &ClientConfig{
				Endpoints: []EndpointConfig{
					{
						URL:      "http://localhost:1317",
						Protocol: "rest",
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config, logger)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					if client != nil {
						client.Close()
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client != nil {
					client.Close()
				}
			}
		})
	}
}

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	if config == nil {
		t.Fatal("DefaultClientConfig returned nil")
	}

	if config.RequestTimeout != 30*time.Second {
		t.Errorf("expected request timeout 30s, got %v", config.RequestTimeout)
	}

	if config.ConnectionTimeout != 10*time.Second {
		t.Errorf("expected connection timeout 10s, got %v", config.ConnectionTimeout)
	}

	if len(config.Endpoints) != 0 {
		t.Errorf("expected empty endpoints in default config, got %d", len(config.Endpoints))
	}
}

// TestClientWithMockEndpoint tests the client with a configuration that won't connect
// but demonstrates the structure without requiring a running Allora node
func TestClientConfiguration(t *testing.T) {
	logger := zerolog.Nop()

	config := &ClientConfig{
		Endpoints: []EndpointConfig{
			{
				URL:      "grpc://localhost:19090", // Non-standard port to avoid conflicts
				Protocol: "grpc",
			},
			{
				URL:      "grpc://localhost:19091", // Another non-standard port
				Protocol: "grpc",
			},
		},
		RequestTimeout:    5 * time.Second,
		ConnectionTimeout: 2 * time.Second,
	}

	// This should fail to connect but should not panic
	client, err := NewClient(config, logger)

	// We expect this to fail since there's no server running
	if err == nil {
		t.Log("Unexpectedly connected to mock endpoints")
		if client != nil {
			client.Close()
		}
	} else {
		t.Logf("Expected connection failure: %v", err)
	}
}

// TestClientMethodsSignatures verifies that the client implements all required methods
func TestClientMethodsSignatures(t *testing.T) {
	// This test ensures that our client interface is properly implemented
	// even if we can't actually connect to test the functionality

	var client AlloraClient
	_ = client // Prevent unused variable error

	// This will compile only if our GRPCClient properly implements AlloraClient
	var grpcClient *GRPCClient
	client = grpcClient
	_ = client

	t.Log("Client interface properly implemented")
}

// BenchmarkClientCreation benchmarks client creation (without connection)
func BenchmarkClientCreation(b *testing.B) {
	logger := zerolog.Nop()

	config := &ClientConfig{
		Endpoints: []EndpointConfig{
			{
				URL:      "grpc://localhost:19090",
				Protocol: "grpc",
			},
		},
		RequestTimeout:    5 * time.Second,
		ConnectionTimeout: 1 * time.Millisecond, // Very short to fail quickly
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client, _ := NewClient(config, logger) // Ignore errors for benchmark
		if client != nil {
			client.Close()
		}
	}
}

// TestContextCancellation tests that operations respect context cancellation
func TestContextCancellation(t *testing.T) {
	logger := zerolog.Nop()

	config := &ClientConfig{
		Endpoints: []EndpointConfig{
			{
				URL:      "grpc://localhost:19090",
				Protocol: "grpc",
			},
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 1 * time.Millisecond, // Will fail quickly
	}

	client, err := NewClient(config, logger)
	if err != nil {
		t.Skip("Cannot create client for context cancellation test:", err)
	}
	defer client.Close()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Any operation should respect the cancelled context
	_, err = client.GetNodeInfo(ctx)
	if err == nil {
		t.Error("Expected error due to cancelled context, but got none")
	}
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name          string
		endpoints     []EndpointConfig
		expectError   bool
		errorContains string
	}{
		{
			name:          "no endpoints",
			endpoints:     []EndpointConfig{},
			expectError:   true,
			errorContains: "at least one endpoint",
		},
		{
			name: "unsupported protocol",
			endpoints: []EndpointConfig{
				{URL: "http://localhost:1317", Protocol: "rest"},
			},
			expectError:   true,
			errorContains: "unsupported protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ClientConfig{
				Endpoints: tt.endpoints,
			}

			_, err := NewClient(config, logger)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', but got: %v", tt.errorContains, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
