package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/azure"
	"microsoft.com/aml-mcp/internal/helpers"
)

// WorkspaceTools contains all workspace-related MCP tools
type WorkspaceTools struct{}

// NewWorkspaceTools creates a new WorkspaceTools instance
func NewWorkspaceTools() *WorkspaceTools {
	return &WorkspaceTools{}
}

// AddToServer registers all workspace tools with the MCP server
func (wt *WorkspaceTools) AddToServer(s *server.MCPServer) {
	wt.addListWorkspacesBySubscriptionTool(s)
	wt.addGetWorkspaceTool(s)
	wt.addCreateWorkspaceTool(s)
}

func (wt *WorkspaceTools) addListWorkspacesBySubscriptionTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_workspaces_by_subscription",
		mcp.WithDescription("List all Azure ML workspaces in a subscription"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
	)

	s.AddTool(tool, wt.handleListWorkspacesBySubscription)
}

func (wt *WorkspaceTools) addGetWorkspaceTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_workspace",
		mcp.WithDescription("Get details of a specific Azure ML workspace"),
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

	s.AddTool(tool, wt.handleGetWorkspace)
}

func (wt *WorkspaceTools) addCreateWorkspaceTool(s *server.MCPServer) {
	tool := mcp.NewTool("create_workspace",
		mcp.WithDescription("Create a new Azure ML workspace"),
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
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("Azure region location (e.g., eastus, westus2)"),
		),
		mcp.WithString("description",
			mcp.Description("Workspace description"),
		),
		mcp.WithString("friendly_name",
			mcp.Description("Friendly name for the workspace"),
		),
	)

	s.AddTool(tool, wt.handleCreateWorkspace)
}

func (wt *WorkspaceTools) handleListWorkspacesBySubscription(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.WorkspacesClient.NewListBySubscriptionPager(nil)
	var workspaces []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get workspaces: %v", err)), nil
		}

		for _, workspace := range page.Value {
			if workspace.Name != nil {
				workspaces = append(workspaces, fmt.Sprintf("Name: %s, Location: %s, Resource Group: %s",
					*workspace.Name,
					helpers.GetStringValue(workspace.Location),
					helpers.ExtractResourceGroupFromID(helpers.GetStringValue(workspace.ID))))
			}
		}
	}

	if len(workspaces) == 0 {
		return mcp.NewToolResultText("No Azure ML workspaces found in the subscription."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d Azure ML workspaces:\n%s", len(workspaces), strings.Join(workspaces, "\n"))), nil
}

func (wt *WorkspaceTools) handleGetWorkspace(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	resp, err := clients.WorkspacesClient.Get(ctx, resourceGroupName, workspaceName, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get workspace: %v", err)), nil
	}

	workspace := resp.Workspace
	details := fmt.Sprintf(`Workspace Details:
Name: %s
Location: %s
Resource Group: %s
ID: %s
Type: %s
SKU: %s
Description: %s
Friendly Name: %s
Discovery URL: %s
ML Flow Tracking URI: %s`,
		helpers.GetStringValue(workspace.Name),
		helpers.GetStringValue(workspace.Location),
		resourceGroupName,
		helpers.GetStringValue(workspace.ID),
		helpers.GetStringValue(workspace.Type),
		helpers.GetSKUString(workspace.SKU),
		helpers.GetWorkspacePropertyString(workspace.Properties, "Description"),
		helpers.GetWorkspacePropertyString(workspace.Properties, "FriendlyName"),
		helpers.GetWorkspacePropertyString(workspace.Properties, "DiscoveryUrl"),
		helpers.GetWorkspacePropertyString(workspace.Properties, "MlFlowTrackingUri"))

	return mcp.NewToolResultText(details), nil
}

func (wt *WorkspaceTools) handleCreateWorkspace(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	location, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	description := request.GetString("description", "")
	friendlyName := request.GetString("friendly_name", "")

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspace := armmachinelearning.Workspace{
		Location: to.Ptr(location),
		Properties: &armmachinelearning.WorkspaceProperties{
			Description:  to.Ptr(description),
			FriendlyName: to.Ptr(friendlyName),
		},
	}

	poller, err := clients.WorkspacesClient.BeginCreateOrUpdate(ctx, resourceGroupName, workspaceName, workspace, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start workspace creation: %v", err)), nil
	}

	result, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create workspace: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully created workspace '%s' in resource group '%s' at location '%s'. Workspace ID: %s",
		workspaceName, resourceGroupName, location, helpers.GetStringValue(result.Workspace.ID))), nil
}
