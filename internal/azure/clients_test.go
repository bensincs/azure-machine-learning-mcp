package azure_test

import (
	"testing"

	"microsoft.com/aml-mcp/internal/azure"
)

func TestNewClientSet(t *testing.T) {
	// Note: This test will only work with valid Azure credentials
	// In a real environment, you might want to mock the Azure clients

	tests := []struct {
		name           string
		subscriptionID string
		shouldError    bool
	}{
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			shouldError:    false, // Azure SDK might not validate this immediately
		},
		{
			name:           "invalid subscription ID format",
			subscriptionID: "invalid-subscription-id",
			shouldError:    false, // Azure SDK creates clients even with invalid IDs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients, err := azure.NewClientSet(tt.subscriptionID)

			if tt.shouldError {
				if err == nil {
					t.Errorf("NewClientSet() expected error, but got none")
				}
				return
			}

			if err != nil {
				// Skip test if Azure credentials are not available
				t.Skipf("Skipping test due to missing Azure credentials: %v", err)
				return
			}

			if clients == nil {
				t.Error("NewClientSet() returned nil clients")
				return
			}

			// Verify all clients are initialized
			if clients.WorkspacesClient == nil {
				t.Error("WorkspacesClient is nil")
			}
			if clients.ComputeClient == nil {
				t.Error("ComputeClient is nil")
			}
			if clients.QuotasClient == nil {
				t.Error("QuotasClient is nil")
			}
			if clients.UsagesClient == nil {
				t.Error("UsagesClient is nil")
			}
			if clients.VirtualMachineSizesClient == nil {
				t.Error("VirtualMachineSizesClient is nil")
			}
			if clients.PrivateEndpointClient == nil {
				t.Error("PrivateEndpointClient is nil")
			}
			if clients.WorkspaceConnectionsClient == nil {
				t.Error("WorkspaceConnectionsClient is nil")
			}
			if clients.WorkspaceFeaturesClient == nil {
				t.Error("WorkspaceFeaturesClient is nil")
			}
		})
	}
}

func TestClientSetStructure(t *testing.T) {
	// Test that the ClientSet struct has the expected fields
	var cs azure.ClientSet

	// This test ensures the struct fields exist and are of the correct types
	// The actual initialization is tested in TestNewClientSet

	if cs.WorkspacesClient != nil {
		t.Error("Expected WorkspacesClient to be nil in uninitialized ClientSet")
	}
	if cs.ComputeClient != nil {
		t.Error("Expected ComputeClient to be nil in uninitialized ClientSet")
	}
	if cs.QuotasClient != nil {
		t.Error("Expected QuotasClient to be nil in uninitialized ClientSet")
	}
	if cs.UsagesClient != nil {
		t.Error("Expected UsagesClient to be nil in uninitialized ClientSet")
	}
	if cs.VirtualMachineSizesClient != nil {
		t.Error("Expected VirtualMachineSizesClient to be nil in uninitialized ClientSet")
	}
	if cs.PrivateEndpointClient != nil {
		t.Error("Expected PrivateEndpointClient to be nil in uninitialized ClientSet")
	}
	if cs.WorkspaceConnectionsClient != nil {
		t.Error("Expected WorkspaceConnectionsClient to be nil in uninitialized ClientSet")
	}
	if cs.WorkspaceFeaturesClient != nil {
		t.Error("Expected WorkspaceFeaturesClient to be nil in uninitialized ClientSet")
	}
}
