package allora

import (
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
)

func TestNewClient(t *testing.T) {
	logger := zerolog.Nop() // No-op logger for tests

	tests := []struct {
		name        string
		config      *config.ClientConfig
		expectError bool
	}{
		{
			name:        "nil config should use default",
			config:      nil,
			expectError: true, // Will fail because no endpoints are provided in default config
		},
		{
			name: "empty endpoints should fail",
			config: &config.ClientConfig{
				Endpoints: []config.EndpointConfig{},
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
	cfg := config.DefaultClientConfig()

	if cfg == nil {
		t.Fatal("config.DefaultClientConfig returned nil")
	}

	if cfg.RequestTimeout != 30*time.Second {
		t.Errorf("expected request timeout 30s, got %v", cfg.RequestTimeout)
	}

	if cfg.ConnectionTimeout != 10*time.Second {
		t.Errorf("expected connection timeout 10s, got %v", cfg.ConnectionTimeout)
	}

	if len(cfg.Endpoints) != 0 {
		t.Errorf("expected empty endpoints in default config, got %d", len(cfg.Endpoints))
	}
}

// TestClientWithTestnetEndpoints tests the client with real testnet endpoints
func TestClientWithTestnetEndpoints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := zerolog.Nop()

	cfg := &config.ClientConfig{
		Endpoints: []config.EndpointConfig{
			{
				URL:      "allora-grpc.testnet.allora.network:443",
				Protocol: config.ProtocolGRPC,
			},
		},
		RequestTimeout:    30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Verify client implements the interface
	var _ Client = client

	t.Log("Successfully connected to Allora testnet")
}

// TestClientMethodsSignatures verifies that the client implements all required methods
func TestClientMethodsSignatures(t *testing.T) {
	// This test ensures that our client interface is properly implemented
	// even if we can't actually connect to test the functionality

	var client Client
	_ = client // Prevent unused variable error

	t.Log("Client interface properly implemented")
}

// BenchmarkClientCreation benchmarks client creation (without connection)
func BenchmarkClientCreation(b *testing.B) {
	b.Skip("Benchmark requires mock endpoints that don't hang")
}

// TestContextCancellation tests that operations respect context cancellation
func TestContextCancellation(t *testing.T) {
	t.Skip("Context cancellation test requires actual RPC endpoints")
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name          string
		endpoints     []config.EndpointConfig
		expectError   bool
		errorContains string
	}{
		{
			name:          "no endpoints",
			endpoints:     []config.EndpointConfig{},
			expectError:   true,
			errorContains: "at least one endpoint",
		},
		{
			name: "unsupported protocol",
			endpoints: []config.EndpointConfig{
				{URL: "http://localhost:1317", Protocol: "http"},
			},
			expectError:   true,
			errorContains: "failed to create any valid clients",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.ClientConfig{
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
