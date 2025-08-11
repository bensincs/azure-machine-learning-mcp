package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

// ClientSet holds all Azure ML service clients
type ClientSet struct {
	WorkspacesClient             *armmachinelearning.WorkspacesClient
	ComputeClient                *armmachinelearning.ComputeClient
	QuotasClient                 *armmachinelearning.QuotasClient
	UsagesClient                 *armmachinelearning.UsagesClient
	VirtualMachineSizesClient    *armmachinelearning.VirtualMachineSizesClient
	PrivateEndpointClient        *armmachinelearning.PrivateEndpointConnectionsClient
	WorkspaceConnectionsClient   *armmachinelearning.WorkspaceConnectionsClient
	WorkspaceFeaturesClient      *armmachinelearning.WorkspaceFeaturesClient
}

// NewClientSet creates a new set of Azure ML clients
func NewClientSet(subscriptionID string) (*ClientSet, error) {
	// Use DefaultAzureCredential which will automatically use Azure CLI credentials when available
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain Azure credential (make sure you're logged in with 'az login'): %v", err)
	}

	workspacesClient, err := armmachinelearning.NewWorkspacesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspaces client: %v", err)
	}

	computeClient, err := armmachinelearning.NewComputeClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %v", err)
	}

	quotasClient, err := armmachinelearning.NewQuotasClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create quotas client: %v", err)
	}

	usagesClient, err := armmachinelearning.NewUsagesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create usages client: %v", err)
	}

	vmSizesClient, err := armmachinelearning.NewVirtualMachineSizesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create VM sizes client: %v", err)
	}

	privateEndpointClient, err := armmachinelearning.NewPrivateEndpointConnectionsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create private endpoint client: %v", err)
	}

	workspaceConnectionsClient, err := armmachinelearning.NewWorkspaceConnectionsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace connections client: %v", err)
	}

	workspaceFeaturesClient, err := armmachinelearning.NewWorkspaceFeaturesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace features client: %v", err)
	}

	return &ClientSet{
		WorkspacesClient:             workspacesClient,
		ComputeClient:                computeClient,
		QuotasClient:                 quotasClient,
		UsagesClient:                 usagesClient,
		VirtualMachineSizesClient:    vmSizesClient,
		PrivateEndpointClient:        privateEndpointClient,
		WorkspaceConnectionsClient:   workspaceConnectionsClient,
		WorkspaceFeaturesClient:      workspaceFeaturesClient,
	}, nil
}
