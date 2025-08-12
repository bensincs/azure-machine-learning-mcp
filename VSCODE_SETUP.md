# Getting Started with Azure ML MCP in VS Code

This guide will help you set up the Azure Machine Learning Model Context Protocol (MCP) server in Visual Studio Code, enabling you to manage Azure ML resources directly from your IDE.

## Prerequisites

1. **VS Code**: Make sure you have Visual Studio Code installed
2. **Azure Account**: An active Azure subscription
3. **Azure CLI**: Install Azure CLI and sign in with `az login`
4. **Docker**: Install Docker Desktop (for Docker-based setup)

## Setup Options

You can set up the Azure ML MCP server in two ways:

### Option 1: Docker Setup (Recommended)

This option uses a pre-built Docker image and is the easiest to get started with.

#### Step 1: Create MCP Configuration

1. Open VS Code
2. Create or open the `.vscode` folder in your workspace
3. Create a file named `mcp.json` in the `.vscode` folder
4. Add the following configuration:

```json
{
  "servers": {
    "azure-ml-mcp": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--pull=always",
        "-v", "~/.azure:/home/appuser/.azure",
        "ghcr.io/bensincs/azure-machine-learning-mcp:latest"
      ]
    }
  }
}
```

#### Step 2: Sign in to Azure

```bash
az login
```

This will open your browser and prompt you to sign in to Azure. Your credentials will be stored in `~/.azure` and shared with the Docker container.

#### Step 3: Install VS Code Extensions

Install the following VS Code extensions if you haven't already:

1. **GitHub Copilot** (required for MCP support)
2. **Azure Account** (optional, for additional Azure integration)

#### Step 4: Test the Setup

1. Open VS Code
2. Open the Command Palette (`Cmd+Shift+P` on macOS, `Ctrl+Shift+P` on Windows/Linux)
3. Look for Azure ML MCP tools in GitHub Copilot chat
4. Try asking: "How many Azure ML workspaces do I have?"

### Option 2: Local Binary Setup

This option builds and runs the MCP server locally from source.

#### Step 1: Clone and Build

```bash
git clone https://github.com/bensincs/azure-machine-learning-mcp.git
cd azure-machine-learning-mcp
make build
```

#### Step 2: Create MCP Configuration

Create `.vscode/mcp.json` in your workspace:

```json
{
  "servers": {
    "azure-ml-mcp": {
      "command": "/path/to/azure-machine-learning-mcp/bin/mcp-server",
      "args": []
    }
  }
}
```

Replace `/path/to/azure-machine-learning-mcp` with the actual path to your cloned repository.

#### Step 3: Sign in to Azure

```bash
az login
```

#### Step 4: Test the Setup

Same as Option 1, step 4.

## Using the Azure ML MCP

Once set up, you can use the Azure ML MCP through GitHub Copilot chat in VS Code. Here are some example prompts:

### Basic Usage Examples

1. **List Workspaces**:
   - "Show me all my Azure ML workspaces"
   - "How many ML workspaces do I have?"

2. **Workspace Details**:
   - "Get details for workspace 'my-workspace' in resource group 'my-rg'"
   - "What's the status of my ML workspace?"

3. **Compute Management**:
   - "List all compute resources in my workspace"
   - "Start the compute instance named 'my-compute'"
   - "Stop all running compute instances"

4. **Resource Monitoring**:
   - "Show my Azure ML quotas in East US region"
   - "What's my current resource usage?"
   - "List available VM sizes for ML compute"

5. **Create Resources**:
   - "Create a new ML workspace named 'test-workspace' in East US"
   - "Set up a workspace in resource group 'ml-rg'"

### Authentication

The MCP server automatically handles Azure authentication using this priority order:

1. **Azure CLI credentials** (from `az login`) - **Recommended**
2. **Environment variables** (Service Principal)
3. **Managed Identity** (when running on Azure)
4. **Interactive browser login** (fallback)

For the Docker setup, your local Azure CLI credentials are automatically shared with the container.

## Troubleshooting

### Common Issues

#### 1. "No MCP servers found"
- Ensure the `mcp.json` file is in the `.vscode` folder
- Check that the file path in the configuration is correct
- Restart VS Code after creating the configuration

#### 2. "Authentication failed"
- Run `az login` to refresh your credentials
- Check that you have the necessary permissions in your Azure subscription
- Verify your subscription ID is correct

#### 3. "Docker image not found"
- Ensure Docker is running
- Check your internet connection for pulling the image
- Try running the Docker command manually to test

#### 4. "Command not found" (Local setup)
- Verify the path to the `mcp-server` binary is correct
- Ensure the binary was built successfully with `make build`
- Check file permissions on the binary

### Debug Mode

To enable debug logging, modify your `mcp.json`:

```json
{
  "servers": {
    "azure-ml-mcp": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--pull=always",
        "-v", "~/.azure:/home/appuser/.azure:ro",
        "-e", "DEBUG=true",
        "ghcr.io/bensincs/azure-machine-learning-mcp:latest"
      ]
    }
  }
}
```

### Getting Help

1. **Check the logs**: Look at VS Code's output panel for error messages
2. **Verify Azure permissions**: Ensure your account has access to Azure ML resources
3. **Test Azure CLI**: Run `az ml workspace list` to verify basic connectivity
4. **Update the image**: For Docker setup, the `--pull=always` flag ensures you get the latest version

## Advanced Configuration

### Custom Subscription

If you have multiple Azure subscriptions, you can specify which one to use:

```bash
az account set --subscription "your-subscription-id"
```

### Service Principal Authentication

For automated scenarios, you can use a Service Principal:

```json
{
  "servers": {
    "azure-ml-mcp": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--pull=always",
        "-e", "AZURE_CLIENT_ID=your-client-id",
        "-e", "AZURE_CLIENT_SECRET=your-client-secret",
        "-e", "AZURE_TENANT_ID=your-tenant-id",
        "ghcr.io/bensincs/azure-machine-learning-mcp:latest"
      ]
    }
  }
}
```

### Multiple Environments

You can configure multiple MCP servers for different environments:

```json
{
  "servers": {
    "azure-ml-dev": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--pull=always",
        "-v", "~/.azure:/home/appuser/.azure:ro",
        "-e", "AZURE_SUBSCRIPTION_ID=dev-subscription-id",
        "ghcr.io/bensincs/azure-machine-learning-mcp:latest"
      ]
    },
    "azure-ml-prod": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--pull=always",
        "-v", "~/.azure:/home/appuser/.azure:ro",
        "-e", "AZURE_SUBSCRIPTION_ID=prod-subscription-id",
        "ghcr.io/bensincs/azure-machine-learning-mcp:latest"
      ]
    }
  }
}
```

## Next Steps

Once you have the Azure ML MCP set up:

1. **Explore the available tools**: Ask Copilot "What Azure ML tools are available?"
2. **Check out the main README**: See [README.md](README.md) for detailed tool documentation
3. **Try advanced scenarios**: Create workspaces, manage compute, monitor resources
4. **Integrate with your workflow**: Use the MCP in your ML development process

## Security Considerations

- The Docker setup mounts your Azure credentials as read-only
- Credentials are not stored in the container or transmitted over the network
- The MCP server runs with minimal privileges
- All Azure API calls use standard Azure authentication and authorization

For production use, consider using Service Principal authentication with minimal required permissions.
