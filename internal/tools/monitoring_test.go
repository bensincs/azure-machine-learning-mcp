package tools_test

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/tools"
)

func TestMonitoringTools_AddToServer(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	monitoringTools := tools.NewMonitoringTools()

	// Test that AddToServer doesn't panic
	monitoringTools.AddToServer(s)
}

func TestMonitoringTools_New(t *testing.T) {
	monitoringTools := tools.NewMonitoringTools()
	if monitoringTools == nil {
		t.Error("NewMonitoringTools() returned nil")
	}
}

func TestMonitoringToolsStructure(t *testing.T) {
	mt := tools.NewMonitoringTools()
	if mt == nil {
		t.Fatal("NewMonitoringTools() returned nil")
	}

	// Test that we can create multiple instances
	mt2 := tools.NewMonitoringTools()
	if mt2 == nil {
		t.Fatal("Second NewMonitoringTools() returned nil")
	}

	// They should be different instances
	if mt == mt2 {
		t.Error("NewMonitoringTools() returned the same instance")
	}
}
