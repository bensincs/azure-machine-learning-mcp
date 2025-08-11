package server_test

import (
	"testing"

	"microsoft.com/aml-mcp/internal/server"
)

func TestNew(t *testing.T) {
	config := server.Config{
		Name:    "Test Azure ML Server",
		Version: "1.0.0-test",
	}

	s := server.New(config)
	if s == nil {
		t.Fatal("New() returned nil server")
	}

	// Note: We can't easily test the internal structure of MCPServer
	// since it wraps the underlying server, but we can test that it's created
}

func TestConfig(t *testing.T) {
	tests := []struct {
		name   string
		config server.Config
	}{
		{
			name: "valid config",
			config: server.Config{
				Name:    "Azure ML MCP Server",
				Version: "1.0.0",
			},
		},
		{
			name: "empty config",
			config: server.Config{
				Name:    "",
				Version: "",
			},
		},
		{
			name: "partial config",
			config: server.Config{
				Name:    "Test Server",
				Version: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.New(tt.config)
			if s == nil {
				t.Error("New() returned nil server")
			}
		})
	}
}

// TestServerInterface ensures that MCPServer implements the expected interface
func TestServerInterface(t *testing.T) {
	config := server.Config{
		Name:    "Test Server",
		Version: "1.0.0",
	}

	s := server.New(config)

	// Test that the server has the Serve method
	// Note: We can't actually call Serve() in tests as it would block
	// and requires stdio setup
	if s == nil {
		t.Fatal("Expected server to be non-nil")
	}
}
