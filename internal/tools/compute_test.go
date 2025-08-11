package tools_test

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/tools"
)

func TestComputeTools_AddToServer(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	computeTools := tools.NewComputeTools()

	// Test that AddToServer doesn't panic
	computeTools.AddToServer(s)
}

func TestComputeTools_New(t *testing.T) {
	computeTools := tools.NewComputeTools()
	if computeTools == nil {
		t.Error("NewComputeTools() returned nil")
	}
}

func TestComputeToolsStructure(t *testing.T) {
	ct := tools.NewComputeTools()
	if ct == nil {
		t.Fatal("NewComputeTools() returned nil")
	}

	// Test that we can create multiple instances
	ct2 := tools.NewComputeTools()
	if ct2 == nil {
		t.Fatal("Second NewComputeTools() returned nil")
	}

	// They should be different instances
	if ct == ct2 {
		t.Error("NewComputeTools() returned the same instance")
	}
}
