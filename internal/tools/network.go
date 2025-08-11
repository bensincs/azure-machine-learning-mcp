package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/azure"
	"microsoft.com/aml-mcp/internal/helpers"
)

// NetworkTools contains all network and security-related MCP tools
type NetworkTools struct{}

// NewNetworkTools creates a new NetworkTools instance
func NewNetworkTools() *NetworkTools {
	return &NetworkTools{}
}

// AddToServer registers all network tools with the MCP server
func (nt *NetworkTools) AddToServer(s *server.MCPServer) {
	nt.addListPrivateEndpointsTool(s)
	nt.addListWorkspaceConnectionsTool(s)
	nt.addListWorkspaceFeaturesTool(s)
}

func (nt *NetworkTools) addListPrivateEndpointsTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_private_endpoints",
		mcp.WithDescription("List private endpoint connections for a workspace"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("resource_group_name",
			mcp.Required(),
			mcp.Description("Resource group name"),
		),
		mcp.WithString("workspace_name",
			mcp.Required(),
			mcp.Description("Workspace name"),
		),
	)

	s.AddTool(tool, nt.handleListPrivateEndpoints)
}

func (nt *NetworkTools) addListWorkspaceConnectionsTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_workspace_connections",
		mcp.WithDescription("List connections for a workspace"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("resource_group_name",
			mcp.Required(),
			mcp.Description("Resource group name"),
		),
		mcp.WithString("workspace_name",
			mcp.Required(),
			mcp.Description("Workspace name"),
		),
	)

	s.AddTool(tool, nt.handleListWorkspaceConnections)
}

func (nt *NetworkTools) addListWorkspaceFeaturesTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_workspace_features",
		mcp.WithDescription("List available features for a workspace"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("resource_group_name",
			mcp.Required(),
			mcp.Description("Resource group name"),
		),
		mcp.WithString("workspace_name",
			mcp.Required(),
			mcp.Description("Workspace name"),
		),
	)

	s.AddTool(tool, nt.handleListWorkspaceFeatures)
}

func (nt *NetworkTools) handleListPrivateEndpoints(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resourceGroupName, err := request.RequireString("resource_group_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspaceName, err := request.RequireString("workspace_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.PrivateEndpointClient.NewListPager(resourceGroupName, workspaceName, nil)
	var endpoints []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get private endpoints: %v", err)), nil
		}

		for _, endpoint := range page.Value {
			if endpoint.Name != nil {
				status := "Unknown"
				if endpoint.Properties != nil && endpoint.Properties.PrivateLinkServiceConnectionState != nil {
					status = helpers.GetPrivateEndpointStatus(endpoint.Properties.PrivateLinkServiceConnectionState.Status)
				}
				endpoints = append(endpoints, fmt.Sprintf("Name: %s, Status: %s, ID: %s",
					*endpoint.Name, status, helpers.GetStringValue(endpoint.ID)))
			}
		}
	}

	if len(endpoints) == 0 {
		return mcp.NewToolResultText("No private endpoint connections found."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d private endpoint connections:\n%s", len(endpoints), strings.Join(endpoints, "\n"))), nil
}

func (nt *NetworkTools) handleListWorkspaceConnections(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resourceGroupName, err := request.RequireString("resource_group_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspaceName, err := request.RequireString("workspace_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.WorkspaceConnectionsClient.NewListPager(resourceGroupName, workspaceName, nil)
	var connections []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get workspace connections: %v", err)), nil
		}

		for _, connection := range page.Value {
			if connection.Name != nil {
				authType := "Unknown"
				category := "Unknown"
				if connection.Properties != nil {
					authType = helpers.GetStringValue(connection.Properties.AuthType)
					category = helpers.GetStringValue(connection.Properties.Category)
				}
				connections = append(connections, fmt.Sprintf("Name: %s, Category: %s, Auth Type: %s, ID: %s",
					*connection.Name, category, authType, helpers.GetStringValue(connection.ID)))
			}
		}
	}

	if len(connections) == 0 {
		return mcp.NewToolResultText("No workspace connections found."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d workspace connections:\n%s", len(connections), strings.Join(connections, "\n"))), nil
}

func (nt *NetworkTools) handleListWorkspaceFeatures(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resourceGroupName, err := request.RequireString("resource_group_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspaceName, err := request.RequireString("workspace_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.WorkspaceFeaturesClient.NewListPager(resourceGroupName, workspaceName, nil)
	var features []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get workspace features: %v", err)), nil
		}

		for _, feature := range page.Value {
			if feature.ID != nil {
				displayName := helpers.GetStringValue(feature.DisplayName)
				description := helpers.GetStringValue(feature.Description)
				features = append(features, fmt.Sprintf("ID: %s, Name: %s, Description: %s",
					*feature.ID, displayName, description))
			}
		}
	}

	if len(features) == 0 {
		return mcp.NewToolResultText("No workspace features found."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d workspace features:\n%s", len(features), strings.Join(features, "\n"))), nil
}
