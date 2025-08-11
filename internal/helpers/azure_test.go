package helpers_test

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"microsoft.com/aml-mcp/internal/helpers"
)

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "N/A",
		},
		{
			name:     "valid string",
			input:    to.Ptr("test-value"),
			expected: "test-value",
		},
		{
			name:     "empty string",
			input:    to.Ptr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetStringValue(tt.input)
			if result != tt.expected {
				t.Errorf("GetStringValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetInt32Value(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected int32
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: 0,
		},
		{
			name:     "valid int32",
			input:    to.Ptr(int32(42)),
			expected: 42,
		},
		{
			name:     "zero value",
			input:    to.Ptr(int32(0)),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetInt32Value(tt.input)
			if result != tt.expected {
				t.Errorf("GetInt32Value() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetFloat64Value(t *testing.T) {
	tests := []struct {
		name     string
		input    *float64
		expected float64
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: 0.0,
		},
		{
			name:     "valid float64",
			input:    to.Ptr(3.14),
			expected: 3.14,
		},
		{
			name:     "zero value",
			input:    to.Ptr(0.0),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetFloat64Value(tt.input)
			if result != tt.expected {
				t.Errorf("GetFloat64Value() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSKUString(t *testing.T) {
	tests := []struct {
		name     string
		input    *armmachinelearning.SKU
		expected string
	}{
		{
			name:     "nil SKU",
			input:    nil,
			expected: "N/A",
		},
		{
			name: "SKU with nil name",
			input: &armmachinelearning.SKU{
				Name: nil,
			},
			expected: "N/A",
		},
		{
			name: "valid SKU",
			input: &armmachinelearning.SKU{
				Name: to.Ptr("Basic"),
			},
			expected: "Basic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetSKUString(tt.input)
			if result != tt.expected {
				t.Errorf("GetSKUString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractResourceGroupFromID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "N/A",
		},
		{
			name:     "valid resource ID",
			input:    "/subscriptions/sub-id/resourceGroups/my-rg/providers/Microsoft.MachineLearningServices/workspaces/my-workspace",
			expected: "my-rg",
		},
		{
			name:     "malformed resource ID",
			input:    "/subscriptions/sub-id/providers/Microsoft.MachineLearningServices/workspaces/my-workspace",
			expected: "N/A",
		},
		{
			name:     "resource ID without resource group name",
			input:    "/subscriptions/sub-id/resourceGroups/",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.ExtractResourceGroupFromID(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractResourceGroupFromID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetWorkspacePropertyString(t *testing.T) {
	props := &armmachinelearning.WorkspaceProperties{
		Description:       to.Ptr("Test workspace"),
		FriendlyName:      to.Ptr("My Workspace"),
		DiscoveryURL:      to.Ptr("https://discovery.azureml.ms/"),
		MlFlowTrackingURI: to.Ptr("https://mlflow.azureml.ms/"),
	}

	tests := []struct {
		name     string
		props    *armmachinelearning.WorkspaceProperties
		field    string
		expected string
	}{
		{
			name:     "nil properties",
			props:    nil,
			field:    "Description",
			expected: "N/A",
		},
		{
			name:     "valid description",
			props:    props,
			field:    "Description",
			expected: "Test workspace",
		},
		{
			name:     "valid friendly name",
			props:    props,
			field:    "FriendlyName",
			expected: "My Workspace",
		},
		{
			name:     "valid discovery URL",
			props:    props,
			field:    "DiscoveryUrl",
			expected: "https://discovery.azureml.ms/",
		},
		{
			name:     "valid MLFlow tracking URI",
			props:    props,
			field:    "MlFlowTrackingUri",
			expected: "https://mlflow.azureml.ms/",
		},
		{
			name:     "unknown field",
			props:    props,
			field:    "UnknownField",
			expected: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetWorkspacePropertyString(tt.props, tt.field)
			if result != tt.expected {
				t.Errorf("GetWorkspacePropertyString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Mock compute implementation for testing
type mockCompute struct {
	compute *armmachinelearning.Compute
}

func (m *mockCompute) GetCompute() *armmachinelearning.Compute {
	return m.compute
}

func TestGetComputeType(t *testing.T) {
	tests := []struct {
		name     string
		input    armmachinelearning.ComputeClassification
		expected string
	}{
		{
			name:     "nil compute",
			input:    nil,
			expected: "N/A",
		},
		{
			name: "compute with nil base",
			input: &mockCompute{
				compute: nil,
			},
			expected: "N/A",
		},
		{
			name: "compute with nil type",
			input: &mockCompute{
				compute: &armmachinelearning.Compute{
					ComputeType: nil,
				},
			},
			expected: "N/A",
		},
		{
			name: "valid compute type",
			input: &mockCompute{
				compute: &armmachinelearning.Compute{
					ComputeType: to.Ptr(armmachinelearning.ComputeTypeComputeInstance),
				},
			},
			expected: "ComputeInstance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetComputeType(tt.input)
			if result != tt.expected {
				t.Errorf("GetComputeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetComputeCreatedOn(t *testing.T) {
	testTime := time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    armmachinelearning.ComputeClassification
		expected string
	}{
		{
			name:     "nil compute",
			input:    nil,
			expected: "N/A",
		},
		{
			name: "compute with nil created on",
			input: &mockCompute{
				compute: &armmachinelearning.Compute{
					CreatedOn: nil,
				},
			},
			expected: "N/A",
		},
		{
			name: "valid created on time",
			input: &mockCompute{
				compute: &armmachinelearning.Compute{
					CreatedOn: &testTime,
				},
			},
			expected: "2023-10-15 14:30:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetComputeCreatedOn(tt.input)
			if result != tt.expected {
				t.Errorf("GetComputeCreatedOn() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetQuotaUnit(t *testing.T) {
	tests := []struct {
		name     string
		input    *armmachinelearning.QuotaUnit
		expected string
	}{
		{
			name:     "nil quota unit",
			input:    nil,
			expected: "N/A",
		},
		{
			name:     "valid quota unit",
			input:    to.Ptr(armmachinelearning.QuotaUnitCount),
			expected: "Count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetQuotaUnit(tt.input)
			if result != tt.expected {
				t.Errorf("GetQuotaUnit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUsageUnit(t *testing.T) {
	tests := []struct {
		name     string
		input    *armmachinelearning.UsageUnit
		expected string
	}{
		{
			name:     "nil usage unit",
			input:    nil,
			expected: "N/A",
		},
		{
			name:     "valid usage unit",
			input:    to.Ptr(armmachinelearning.UsageUnitCount),
			expected: "Count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetUsageUnit(tt.input)
			if result != tt.expected {
				t.Errorf("GetUsageUnit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPrivateEndpointStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    *armmachinelearning.PrivateEndpointServiceConnectionStatus
		expected string
	}{
		{
			name:     "nil status",
			input:    nil,
			expected: "Unknown",
		},
		{
			name:     "approved status",
			input:    to.Ptr(armmachinelearning.PrivateEndpointServiceConnectionStatusApproved),
			expected: "Approved",
		},
		{
			name:     "pending status",
			input:    to.Ptr(armmachinelearning.PrivateEndpointServiceConnectionStatusPending),
			expected: "Pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetPrivateEndpointStatus(tt.input)
			if result != tt.expected {
				t.Errorf("GetPrivateEndpointStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}
