package azure

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)// ClientSet holds all Azure ML service clients
type ClientSet struct {
	WorkspacesClient           *armmachinelearning.WorkspacesClient
	ComputeClient              *armmachinelearning.ComputeClient
	QuotasClient               *armmachinelearning.QuotasClient
	UsagesClient               *armmachinelearning.UsagesClient
	VirtualMachineSizesClient  *armmachinelearning.VirtualMachineSizesClient
	PrivateEndpointClient      *armmachinelearning.PrivateEndpointConnectionsClient
	WorkspaceConnectionsClient *armmachinelearning.WorkspaceConnectionsClient
	WorkspaceFeaturesClient    *armmachinelearning.WorkspaceFeaturesClient
}

// getAzureCredential attempts to get Azure credentials using multiple methods
func getAzureCredential() (azcore.TokenCredential, error) {
	// First, try DefaultAzureCredential (includes Azure CLI, managed identity, etc.)
	if cred, err := azidentity.NewDefaultAzureCredential(nil); err == nil {
		// Test the credential by trying to get a token
		ctx := context.Background()
		_, err = cred.GetToken(ctx, policy.TokenRequestOptions{
			Scopes: []string{"https://management.azure.com/.default"},
		})
		if err == nil {
			log.Println("Using existing Azure credentials (Azure CLI, Managed Identity, or Environment)")
			return cred, nil
		}
		log.Printf("Default credentials failed: %v", err)
	}

	// If no existing credentials, try interactive browser authentication
	log.Println("No existing Azure credentials found. Opening browser for interactive login...")
	
	// Check if we're in a headless environment
if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" && os.Getenv("XDG_SESSION_TYPE") == "" {
// Headless environment, use device code flow
log.Println("Headless environment detected. Using device code authentication...")
return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
fmt.Printf("\n%s\n", message.Message)
return nil
},
})
}

// Use interactive browser authentication
return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
RedirectURL: "http://localhost:8080",
})
}

// NewClientSet creates a new set of Azure ML clients
func NewClientSet(subscriptionID string) (*ClientSet, error) {
// Get Azure credential with interactive fallback
cred, err := getAzureCredential()
if err != nil {
return nil, fmt.Errorf("failed to obtain Azure credential: %v", err)
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
WorkspacesClient:           workspacesClient,
ComputeClient:              computeClient,
QuotasClient:               quotasClient,
UsagesClient:               usagesClient,
VirtualMachineSizesClient:  vmSizesClient,
PrivateEndpointClient:      privateEndpointClient,
WorkspaceConnectionsClient: workspaceConnectionsClient,
WorkspaceFeaturesClient:    workspaceFeaturesClient,
}, nil
}
