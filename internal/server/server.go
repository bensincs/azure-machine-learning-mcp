package server

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"microsoft.com/aml-mcp/internal/tools"
)

// Config holds the configuration for the MCP server
type Config struct {
	Name    string
	Version string
}

// MCPServer wraps the underlying MCP server with our tools
type MCPServer struct {
	server *server.MCPServer
}

// New creates a new MCP server with all Azure ML tools registered
func New(config Config) *MCPServer {
	s := server.NewMCPServer(
		config.Name,
		config.Version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register all tool categories
	workspaceTools := tools.NewWorkspaceTools()
	workspaceTools.AddToServer(s)

	computeTools := tools.NewComputeTools()
	computeTools.AddToServer(s)

	monitoringTools := tools.NewMonitoringTools()
	monitoringTools.AddToServer(s)

	networkTools := tools.NewNetworkTools()
	networkTools.AddToServer(s)

	return &MCPServer{server: s}
}

// Serve starts the MCP server using stdio
func (ms *MCPServer) Serve() error {
	fmt.Println("Starting Azure Machine Learning MCP Server...")

	if err := server.ServeStdio(ms.server); err != nil {
		return fmt.Errorf("server error: %v", err)
	}

	fmt.Println("Azure ML MCP server is running. Press Ctrl+C to stop.")
	return nil
}
