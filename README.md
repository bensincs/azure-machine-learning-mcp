# Azure Machine Learning MCP Server

A comprehensive Model Context Protocol (MCP) server that provides access to Azure Machine Learning services through the Azure SDK for Go. This server enables you to manage Azure ML workspaces, compute resources, quotas, and more through standardized MCP tools.

## Features

This MCP server provides the following Azure ML management capabilities:

### Workspace Management
- **list_workspaces_by_subscription**: List all Azure ML workspaces in a subscription
- **get_workspace**: Get detailed information about a specific workspace
- **create_workspace**: Create a new Azure ML workspace

### Compute Resource Management
- **list_compute**: List all compute resources in a workspace
- **get_compute**: Get detailed information about a specific compute resource
- **start_compute**: Start a compute resource
- **stop_compute**: Stop a compute resource

### Resource Monitoring
- **list_quotas**: List resource quotas for a specific Azure region
- **list_usage**: List current resource usage for a specific Azure region
- **list_vm_sizes**: List available virtual machine sizes for compute

### Network & Security
- **list_private_endpoints**: List private endpoint connections for a workspace
- **list_workspace_connections**: List workspace connections
- **list_workspace_features**: List available features for a workspace

## Prerequisites

1. **Azure Subscription**: You need an active Azure subscription
2. **Azure Authentication**: The server will automatically prompt for authentication when you first use it. You can use any of these methods:
   - **Interactive Browser Login**: The server will open your browser for authentication (default)
   - **Device Code Flow**: Used automatically in headless environments
   - **Azure CLI**: If you're already logged in with `az login`
   - **Service Principal**: Set environment variables `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`
   - **Managed Identity**: When running on Azure resources

3. **Go 1.18+**: Required to build and run the server

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd aml-cmp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build ./cmd/mcp-server
```

## Configuration

### MCP Client Configuration

Add this server to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "servers": {
    "azure-ml": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server"],
      "cwd": "/path/to/aml-cmp",
      "env": {}
    }
  }
}
```

### Authentication Setup

The server will automatically prompt for Azure authentication when you first connect. No pre-authentication is required!

**Authentication Methods (in order of precedence):**

1. **Existing Credentials** (if available):
   - Azure CLI credentials (`az login`)
   - Managed Identity (when running on Azure)
   - Environment variables for Service Principal

2. **Interactive Authentication** (if no existing credentials):
   - **Browser Login**: Opens your default browser for Microsoft login
   - **Device Code**: Used in headless environments (displays a code to enter at https://microsoft.com/devicelogin)

**Manual Setup Options:**

If you prefer to set up authentication manually:

1. **Azure CLI** (Recommended for development):
```bash
az login
```

2. **Service Principal** (Recommended for production):
```bash
export AZURE_CLIENT_ID="your-client-id"
export AZURE_CLIENT_SECRET="your-client-secret"
export AZURE_TENANT_ID="your-tenant-id"
```

3. **Managed Identity** (When running on Azure):
   - No additional setup required when running on Azure VMs, App Service, etc.

## Usage Examples

### 1. List Workspaces
```
Use the list_workspaces_by_subscription tool with:
- subscription_id: "your-azure-subscription-id"
```

### 2. Create a Workspace
```
Use the create_workspace tool with:
- subscription_id: "your-azure-subscription-id"
- resource_group_name: "my-resource-group"
- workspace_name: "my-ml-workspace"
- location: "eastus"
- description: "My ML workspace for experiments"
- friendly_name: "ML Workspace"
```

### 3. Manage Compute Resources
```
# List compute resources
Use list_compute with workspace details

# Start a compute instance
Use start_compute with:
- subscription_id, resource_group_name, workspace_name, compute_name

# Stop a compute instance
Use stop_compute with:
- subscription_id, resource_group_name, workspace_name, compute_name
```

### 4. Monitor Resources
```
# Check quotas
Use list_quotas with:
- subscription_id: "your-subscription-id"
- location: "eastus"

# Check current usage
Use list_usage with same parameters
```

### 5. List Available VM Sizes
```
Use list_vm_sizes with:
- subscription_id: "your-subscription-id"
- location: "eastus"
```

## Tool Reference

### Workspace Tools

#### `list_workspaces_by_subscription`
Lists all Azure ML workspaces in a subscription.

**Parameters:**
- `subscription_id` (required): Azure subscription ID

**Returns:** List of workspaces with names, locations, and resource groups.

#### `get_workspace`
Gets detailed information about a specific workspace.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name

**Returns:** Detailed workspace information including properties, URLs, and configuration.

#### `create_workspace`
Creates a new Azure ML workspace.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name
- `location` (required): Azure region (e.g., "eastus", "westus2")
- `description` (optional): Workspace description
- `friendly_name` (optional): Friendly display name

**Returns:** Confirmation of workspace creation with workspace ID.

### Compute Tools

#### `list_compute`
Lists all compute resources in a workspace.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name

**Returns:** List of compute resources with types, states, and locations.

#### `get_compute`
Gets detailed information about a specific compute resource.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name
- `compute_name` (required): Compute resource name

**Returns:** Detailed compute information including state, creation time, and configuration.

#### `start_compute` / `stop_compute`
Start or stop a compute resource.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name
- `compute_name` (required): Compute resource name

**Returns:** Confirmation of operation completion.

### Monitoring Tools

#### `list_quotas`
Lists resource quotas for a specific region.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `location` (required): Azure region

**Returns:** List of quotas with limits, units, and resource types.

#### `list_usage`
Lists current resource usage for a specific region.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `location` (required): Azure region

**Returns:** Current usage with limits and available capacity.

#### `list_vm_sizes`
Lists available virtual machine sizes for compute.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `location` (required): Azure region

**Returns:** Available VM sizes with vCPU and memory specifications.

### Network & Security Tools

#### `list_private_endpoints`
Lists private endpoint connections for a workspace.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name

**Returns:** Private endpoint connections with status and configuration.

#### `list_workspace_connections`
Lists workspace connections (data stores, linked services).

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name

**Returns:** Workspace connections with types and authentication methods.

#### `list_workspace_features`
Lists available features for a workspace.

**Parameters:**
- `subscription_id` (required): Azure subscription ID
- `resource_group_name` (required): Resource group name
- `workspace_name` (required): Workspace name

**Returns:** Available workspace features and capabilities.

## Error Handling

The server provides detailed error messages for common scenarios:

- **Authentication Errors**: Issues with Azure credentials or permissions
- **Resource Not Found**: When specified resources don't exist
- **Permission Errors**: When the authenticated user lacks required permissions
- **API Errors**: Azure API-specific errors with detailed messages

## Development

### Project Structure
```
.
├── cmd/mcp-server/          # Main application entry point
│   └── main.go
├── internal/                # Internal packages (not importable externally)
│   ├── azure/              # Azure client management
│   │   ├── clients.go      # Azure service clients
│   │   └── tests/          # Tests for Azure client functionality
│   │       └── clients_test.go
│   ├── helpers/            # Utility functions
│   │   ├── azure.go        # Azure-specific helpers
│   │   └── tests/          # Tests for helper functions
│   │       └── azure_test.go
│   ├── server/             # MCP server implementation
│   │   ├── server.go       # Server setup and configuration
│   │   └── tests/          # Server tests
│   │       └── server_test.go
│   └── tools/              # MCP tool implementations
│       ├── workspace.go    # Workspace management tools
│       ├── compute.go      # Compute resource tools
│       ├── monitoring.go   # Monitoring tools (quotas, usage, VM sizes)
│       ├── network.go      # Network and security tools
│       └── tests/          # Tests for tool implementations
│           ├── workspace_test.go
│           ├── compute_test.go
│           ├── monitoring_test.go
│           └── network_test.go
├── .github/workflows/       # GitHub Actions CI/CD
├── Makefile                # Build and development commands
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # This file
```

### Building
```bash
# Build the server
make build

# Or manually
go build ./cmd/mcp-server
```

### Testing
```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run specific test suites
make test-helpers    # Test utility functions
make test-azure      # Test Azure client creation
make test-tools      # Test MCP tool implementations
make test-server     # Test server setup
```

### Development Commands
```bash
# Install dependencies
make deps

# Format code
make format

# Run linter
make lint

# Run all checks (format, lint, test)
make check

# Clean build artifacts
make clean

# Build for multiple platforms
make build-all

# Run the server locally
make run
```

### Running Locally
```bash
# Using make
make run

# Or directly
go run ./cmd/mcp-server
```

### Adding New Tools

To add new Azure ML functionality:

1. Create a new tool definition using `mcp.NewTool()`
2. Implement the handler function with Azure SDK calls
3. Add the tool to the server using `s.AddTool()`
4. Update this README with documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## License

This project is licensed under the MIT License. See LICENSE file for details.

## Support

For issues and questions:
- Check the Azure ML documentation: https://docs.microsoft.com/azure/machine-learning/
- Review Azure SDK for Go documentation: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go
- Open an issue in this repository

## Azure ML SDK Coverage

This MCP server covers the following Azure ML SDK clients:
- ✅ WorkspacesClient
- ✅ ComputeClient
- ✅ QuotasClient
- ✅ UsagesClient
- ✅ VirtualMachineSizesClient
- ✅ PrivateEndpointConnectionsClient
- ✅ WorkspaceConnectionsClient
- ✅ WorkspaceFeaturesClient
- ⏳ WorkspaceSKUsClient (planned)
- ⏳ PrivateLinkResourcesClient (planned)
- ⏳ OperationsClient (planned)

Future versions will expand coverage to include additional Azure ML services and capabilities.
