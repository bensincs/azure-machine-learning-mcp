# Azure ML MCP Server - Development Guide

## Architecture Overview

The Azure ML MCP Server has been refactored into a clean, modular architecture following Go best practices:

### Package Structure

- **`cmd/mcp-server/`** - Main application entry point
- **`internal/azure/`** - Azure client management and credential handling
  - `tests/` - Unit tests for Azure client functionality
- **`internal/helpers/`** - Utility functions for Azure SDK data manipulation
  - `tests/` - Unit tests for helper functions
- **`internal/server/`** - MCP server setup and tool registration
  - `tests/` - Unit tests for server functionality
- **`internal/tools/`** - Individual MCP tool implementations organized by functionality
  - `tests/` - Unit tests for tool implementations

### Design Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Testability**: All packages are designed to be easily unit tested
3. **Modularity**: Tools are organized by functional areas (workspace, compute, monitoring, network)
4. **Error Handling**: Consistent error handling and user-friendly error messages
5. **Type Safety**: Extensive use of helper functions to safely handle pointer types from Azure SDK

### Testing Strategy

- **Unit Tests**: Each package has comprehensive unit tests
- **Integration Tests**: Tools are tested with mock requests
- **CI/CD**: Automated testing on multiple Go versions
- **Coverage**: Test coverage tracking and reporting

### Tool Categories

1. **Workspace Tools** (`tools/workspace.go`)
   - List workspaces by subscription
   - Get workspace details
   - Create new workspaces

2. **Compute Tools** (`tools/compute.go`)
   - List compute resources
   - Get compute details
   - Start/stop compute instances

3. **Monitoring Tools** (`tools/monitoring.go`)
   - List quotas and usage
   - VM size information
   - Resource monitoring

4. **Network Tools** (`tools/network.go`)
   - Private endpoint management
   - Workspace connections
   - Security features

### Adding New Tools

To add a new tool:

1. Create the tool in the appropriate category file or create a new category
2. Implement the tool handler function
3. Register the tool in the `AddToServer` method
4. Add unit tests for the new functionality
5. Update documentation

### Error Handling Patterns

The codebase uses consistent error handling patterns:

- Azure client creation errors are handled gracefully
- Missing parameters return user-friendly error messages
- Azure API errors are wrapped with context
- Helper functions return safe default values for nil pointers

### Testing Locally

Since the tools interact with real Azure resources, testing requires:

1. Azure CLI authentication (`az login`)
2. Valid Azure subscription
3. Appropriate permissions for the resources being tested

For CI/CD and development without Azure access, the tests are designed to skip gracefully when credentials are not available.
