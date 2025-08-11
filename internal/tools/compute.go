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

// ComputeTools contains all compute-related MCP tools
type ComputeTools struct{}

// NewComputeTools creates a new ComputeTools instance
func NewComputeTools() *ComputeTools {
	return &ComputeTools{}
}

// AddToServer registers all compute tools with the MCP server
func (ct *ComputeTools) AddToServer(s *server.MCPServer) {
	ct.addListComputeTool(s)
	ct.addGetComputeTool(s)
	ct.addStartComputeTool(s)
	ct.addStopComputeTool(s)
}

func (ct *ComputeTools) addListComputeTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_compute",
		mcp.WithDescription("List all compute resources in an Azure ML workspace"),
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

	s.AddTool(tool, ct.handleListCompute)
}

func (ct *ComputeTools) addGetComputeTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_compute",
		mcp.WithDescription("Get details of a specific compute resource"),
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
		mcp.WithString("compute_name",
			mcp.Required(),
			mcp.Description("Compute resource name"),
		),
	)

	s.AddTool(tool, ct.handleGetCompute)
}

func (ct *ComputeTools) addStartComputeTool(s *server.MCPServer) {
	tool := mcp.NewTool("start_compute",
		mcp.WithDescription("Start a compute resource"),
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
		mcp.WithString("compute_name",
			mcp.Required(),
			mcp.Description("Compute resource name"),
		),
	)

	s.AddTool(tool, ct.handleStartCompute)
}

func (ct *ComputeTools) addStopComputeTool(s *server.MCPServer) {
	tool := mcp.NewTool("stop_compute",
		mcp.WithDescription("Stop a compute resource"),
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
		mcp.WithString("compute_name",
			mcp.Required(),
			mcp.Description("Compute resource name"),
		),
	)

	s.AddTool(tool, ct.handleStopCompute)
}

func (ct *ComputeTools) handleListCompute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	pager := clients.ComputeClient.NewListPager(resourceGroupName, workspaceName, nil)
	var computes []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get compute resources: %v", err)), nil
		}

		for _, compute := range page.Value {
			if compute.Name != nil && compute.Properties != nil {
				computes = append(computes, fmt.Sprintf("Name: %s, Type: %s, Location: %s, State: %s",
					*compute.Name,
					helpers.GetComputeType(compute.Properties),
					helpers.GetStringValue(compute.Location),
					helpers.GetComputeProvisioningState(compute.Properties)))
			}
		}
	}

	if len(computes) == 0 {
		return mcp.NewToolResultText("No compute resources found in the workspace."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d compute resources:\n%s", len(computes), strings.Join(computes, "\n"))), nil
}

func (ct *ComputeTools) handleGetCompute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	computeName, err := request.RequireString("compute_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := clients.ComputeClient.Get(ctx, resourceGroupName, workspaceName, computeName, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get compute resource: %v", err)), nil
	}

	compute := resp.ComputeResource
	details := fmt.Sprintf(`Compute Resource Details:
Name: %s
Type: %s
Location: %s
Description: %s
Resource ID: %s
Provisioning State: %s
Created On: %s
Modified On: %s
Is Attached: %t`,
		helpers.GetStringValue(compute.Name),
		helpers.GetComputeType(compute.Properties),
		helpers.GetStringValue(compute.Location),
		helpers.GetComputeDescription(compute.Properties),
		helpers.GetStringValue(compute.ID),
		helpers.GetComputeProvisioningState(compute.Properties),
		helpers.GetComputeCreatedOn(compute.Properties),
		helpers.GetComputeModifiedOn(compute.Properties),
		helpers.GetComputeIsAttached(compute.Properties))

	return mcp.NewToolResultText(details), nil
}

func (ct *ComputeTools) handleStartCompute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	computeName, err := request.RequireString("compute_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	poller, err := clients.ComputeClient.BeginStart(ctx, resourceGroupName, workspaceName, computeName, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start compute: %v", err)), nil
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start compute: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully started compute resource '%s'", computeName)), nil
}

func (ct *ComputeTools) handleStopCompute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	computeName, err := request.RequireString("compute_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	poller, err := clients.ComputeClient.BeginStop(ctx, resourceGroupName, workspaceName, computeName, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop compute: %v", err)), nil
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop compute: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully stopped compute resource '%s'", computeName)), nil
}
