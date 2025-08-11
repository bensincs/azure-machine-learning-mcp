package helpers

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

// GetStringValue safely returns the value of a string pointer or "N/A" if nil
func GetStringValue(ptr *string) string {
	if ptr == nil {
		return "N/A"
	}
	return *ptr
}

// GetInt32Value safely returns the value of an int32 pointer or 0 if nil
func GetInt32Value(ptr *int32) int32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// GetFloat64Value safely returns the value of a float64 pointer or 0.0 if nil
func GetFloat64Value(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

// GetSKUString extracts SKU name from SKU object
func GetSKUString(sku *armmachinelearning.SKU) string {
	if sku == nil || sku.Name == nil {
		return "N/A"
	}
	return *sku.Name
}

// GetWorkspacePropertyString extracts specific property from workspace properties
func GetWorkspacePropertyString(props *armmachinelearning.WorkspaceProperties, field string) string {
	if props == nil {
		return "N/A"
	}

	switch field {
	case "Description":
		return GetStringValue(props.Description)
	case "FriendlyName":
		return GetStringValue(props.FriendlyName)
	case "DiscoveryUrl":
		return GetStringValue(props.DiscoveryURL)
	case "MlFlowTrackingUri":
		return GetStringValue(props.MlFlowTrackingURI)
	default:
		return "N/A"
	}
}

// ExtractResourceGroupFromID extracts resource group name from Azure resource ID
func ExtractResourceGroupFromID(id string) string {
	if id == "" {
		return "N/A"
	}
	parts := strings.Split(id, "/")
	for i, part := range parts {
		if part == "resourceGroups" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "N/A"
}

// GetComputeType extracts compute type from compute classification
func GetComputeType(compute armmachinelearning.ComputeClassification) string {
	if compute == nil {
		return "N/A"
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil || baseCompute.ComputeType == nil {
		return "N/A"
	}

	return string(*baseCompute.ComputeType)
}

// GetComputeDescription extracts description from compute classification
func GetComputeDescription(compute armmachinelearning.ComputeClassification) string {
	if compute == nil {
		return "N/A"
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil {
		return "N/A"
	}

	return GetStringValue(baseCompute.Description)
}

// GetComputeProvisioningState extracts provisioning state from compute classification
func GetComputeProvisioningState(compute armmachinelearning.ComputeClassification) string {
	if compute == nil {
		return "N/A"
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil || baseCompute.ProvisioningState == nil {
		return "N/A"
	}

	return string(*baseCompute.ProvisioningState)
}

// GetComputeCreatedOn extracts creation date from compute classification
func GetComputeCreatedOn(compute armmachinelearning.ComputeClassification) string {
	if compute == nil {
		return "N/A"
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil || baseCompute.CreatedOn == nil {
		return "N/A"
	}

	return baseCompute.CreatedOn.Format("2006-01-02 15:04:05")
}

// GetComputeModifiedOn extracts modification date from compute classification
func GetComputeModifiedOn(compute armmachinelearning.ComputeClassification) string {
	if compute == nil {
		return "N/A"
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil || baseCompute.ModifiedOn == nil {
		return "N/A"
	}

	return baseCompute.ModifiedOn.Format("2006-01-02 15:04:05")
}

// GetComputeIsAttached extracts attachment status from compute classification
func GetComputeIsAttached(compute armmachinelearning.ComputeClassification) bool {
	if compute == nil {
		return false
	}

	baseCompute := compute.GetCompute()
	if baseCompute == nil || baseCompute.IsAttachedCompute == nil {
		return false
	}

	return *baseCompute.IsAttachedCompute
}

// GetQuotaUnit extracts quota unit string
func GetQuotaUnit(unit *armmachinelearning.QuotaUnit) string {
	if unit == nil {
		return "N/A"
	}
	return string(*unit)
}

// GetUsageUnit extracts usage unit string
func GetUsageUnit(unit *armmachinelearning.UsageUnit) string {
	if unit == nil {
		return "N/A"
	}
	return string(*unit)
}

// GetPrivateEndpointStatus extracts private endpoint connection status
func GetPrivateEndpointStatus(status *armmachinelearning.PrivateEndpointServiceConnectionStatus) string {
	if status == nil {
		return "Unknown"
	}
	return string(*status)
}
