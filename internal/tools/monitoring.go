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

// MonitoringTools contains all monitoring-related MCP tools (quotas, usage, VM sizes)
type MonitoringTools struct{}

// NewMonitoringTools creates a new MonitoringTools instance
func NewMonitoringTools() *MonitoringTools {
	return &MonitoringTools{}
}

// AddToServer registers all monitoring tools with the MCP server
func (mt *MonitoringTools) AddToServer(s *server.MCPServer) {
	mt.addListQuotasTool(s)
	mt.addListUsageTool(s)
	mt.addListVMSizesTool(s)
}

func (mt *MonitoringTools) addListQuotasTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_quotas",
		mcp.WithDescription("List quotas for Azure ML resources in a location"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("Azure region location (e.g., eastus, westus2)"),
		),
	)

	s.AddTool(tool, mt.handleListQuotas)
}

func (mt *MonitoringTools) addListUsageTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_usage",
		mcp.WithDescription("List current usage for Azure ML resources in a location"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("Azure region location (e.g., eastus, westus2)"),
		),
	)

	s.AddTool(tool, mt.handleListUsage)
}

func (mt *MonitoringTools) addListVMSizesTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_vm_sizes",
		mcp.WithDescription("List available virtual machine sizes for Azure ML compute"),
		mcp.WithString("subscription_id",
			mcp.Required(),
			mcp.Description("Azure subscription ID"),
		),
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("Azure region location (e.g., eastus, westus2)"),
		),
	)

	s.AddTool(tool, mt.handleListVMSizes)
}

func (mt *MonitoringTools) handleListQuotas(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	location, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.QuotasClient.NewListPager(location, nil)
	var quotas []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get quotas: %v", err)), nil
		}

		for _, quota := range page.Value {
			if quota.Name != nil && quota.Limit != nil {
				quotas = append(quotas, fmt.Sprintf("Resource: %s, Limit: %d, Unit: %s, Type: %s",
					helpers.GetStringValue(quota.Name.Value),
					*quota.Limit,
					helpers.GetQuotaUnit(quota.Unit),
					helpers.GetStringValue(quota.Type)))
			}
		}
	}

	if len(quotas) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No quotas found for location '%s'.", location)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d quotas for location '%s':\n%s", len(quotas), location, strings.Join(quotas, "\n"))), nil
}

func (mt *MonitoringTools) handleListUsage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	location, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pager := clients.UsagesClient.NewListPager(location, nil)
	var usages []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get usage: %v", err)), nil
		}

		for _, usage := range page.Value {
			if usage.Name != nil && usage.CurrentValue != nil && usage.Limit != nil {
				usages = append(usages, fmt.Sprintf("Resource: %s, Current: %d, Limit: %d, Unit: %s",
					helpers.GetStringValue(usage.Name.Value),
					*usage.CurrentValue,
					*usage.Limit,
					helpers.GetUsageUnit(usage.Unit)))
			}
		}
	}

	if len(usages) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No usage data found for location '%s'.", location)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d usage entries for location '%s':\n%s", len(usages), location, strings.Join(usages, "\n"))), nil
}

func (mt *MonitoringTools) handleListVMSizes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subscriptionID, err := request.RequireString("subscription_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	location, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clients, err := azure.NewClientSet(subscriptionID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := clients.VirtualMachineSizesClient.List(ctx, location, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get VM sizes: %v", err)), nil
	}

	var vmSizes []string
	for _, vmSize := range resp.Value {
		if vmSize.Name != nil {
			vmSizes = append(vmSizes, fmt.Sprintf("Name: %s, vCPUs: %d, Memory: %.1f GB",
				*vmSize.Name,
				helpers.GetInt32Value(vmSize.VCPUs),
				helpers.GetFloat64Value(vmSize.MemoryGB)))
		}
	}

	if len(vmSizes) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No VM sizes found for location '%s'.", location)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d VM sizes for location '%s':\n%s", len(vmSizes), location, strings.Join(vmSizes, "\n"))), nil
}
