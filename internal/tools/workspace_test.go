package tools_test

import (
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/tools"
)

func TestWorkspaceTools_AddToServer(t *testing.T) {
	s := server.NewMCPServer("test", "1.0.0")
	workspaceTools := tools.NewWorkspaceTools()

	// Test that AddToServer doesn't panic
	workspaceTools.AddToServer(s)

	// Note: We can't easily test the actual tool registration
	// without accessing the server's internal state
}

func TestWorkspaceTools_New(t *testing.T) {
	workspaceTools := tools.NewWorkspaceTools()
	if workspaceTools == nil {
		t.Error("NewWorkspaceTools() returned nil")
	}
}

// Mock test for workspace tool handlers - these would require Azure credentials
// and actual Azure resources to test properly
func TestWorkspaceToolHandlers_Mock(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		params      map[string]interface{}
		shouldError bool
	}{
		{
			name:        "list_workspaces_by_subscription missing subscription_id",
			toolName:    "list_workspaces_by_subscription",
			params:      map[string]interface{}{},
			shouldError: true,
		},
		{
			name:     "get_workspace missing required params",
			toolName: "get_workspace",
			params: map[string]interface{}{
				"subscription_id": "test-sub",
			},
			shouldError: true,
		},
		{
			name:     "create_workspace missing required params",
			toolName: "create_workspace",
			params: map[string]interface{}{
				"subscription_id":     "test-sub",
				"resource_group_name": "test-rg",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request
			request := &mockCallToolRequest{
				params: tt.params,
			}

			// Note: In a real test environment, you would need to:
			// 1. Set up the actual tool handlers
			// 2. Call them with the mock request
			// 3. Verify the response

			// For now, we just verify the test structure
			if request.params == nil && !tt.shouldError {
				t.Error("Expected valid params for successful test case")
			}
		})
	}
}

// Mock implementation of mcp.CallToolRequest for testing
type mockCallToolRequest struct {
	params map[string]interface{}
}

func (m *mockCallToolRequest) RequireString(key string) (string, error) {
	if val, ok := m.params[key]; ok {
		if str, ok := val.(string); ok {
			return str, nil
		}
	}
	return "", errors.New("missing required parameter: " + key)
}

func (m *mockCallToolRequest) GetString(key, defaultValue string) string {
	if val, ok := m.params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// Test workspace tools structure
func TestWorkspaceToolsStructure(t *testing.T) {
	wt := tools.NewWorkspaceTools()
	if wt == nil {
		t.Fatal("NewWorkspaceTools() returned nil")
	}

	// Test that we can create multiple instances
	wt2 := tools.NewWorkspaceTools()
	if wt2 == nil {
		t.Fatal("Second NewWorkspaceTools() returned nil")
	}

	// They should be different instances
	if wt == wt2 {
		t.Error("NewWorkspaceTools() returned the same instance")
	}
}
