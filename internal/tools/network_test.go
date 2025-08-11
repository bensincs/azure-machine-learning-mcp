package tools_test

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/tools"
)

func TestNetworkTools_AddToServer(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	networkTools := tools.NewNetworkTools()

	// Test that AddToServer doesn't panic
	networkTools.AddToServer(s)
}

func TestNetworkTools_New(t *testing.T) {
	networkTools := tools.NewNetworkTools()
	if networkTools == nil {
		t.Error("NewNetworkTools() returned nil")
	}
}

func TestNetworkToolsStructure(t *testing.T) {
	nt := tools.NewNetworkTools()
	if nt == nil {
		t.Fatal("NewNetworkTools() returned nil")
	}

	// Test that we can create multiple instances
	nt2 := tools.NewNetworkTools()
	if nt2 == nil {
		t.Fatal("Second NewNetworkTools() returned nil")
	}

	// They should be different instances
	if nt == nt2 {
		t.Error("NewNetworkTools() returned the same instance")
	}
}
