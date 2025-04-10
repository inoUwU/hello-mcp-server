package main

import (
	"flag"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	sseMode := flag.Bool("sse", false, "Run in SSE mode instead of stdio mode")
	flag.Parse()

	// Create MCP server with basic capabilities
	mcpServer := server.NewMCPServer(
		"python-executor",
		"1.0.0",
	)

	// Create and add the Python execution tool
	pythonTool := mcp.NewTool(
		"execute-python",
		mcp.WithDescription(
			"Execute Python code in an isolated environment. Playwright and headless browser are available for web scraping. Use this tool when you need real-time information, don't have the information internally and no other tools can provide this information. Only output printed to stdout or stderr is returned so ALWAYS use print statements! Please note all code is run in an ephemeral container so modules and code do NOT persist!",
		),
		mcp.WithString(
			"code",
			mcp.Description("The Python code to execute"),
			mcp.Required(),
		),
		mcp.WithString(
			"modules",
			mcp.Description(
				"Comma-separated list of Python modules your code requires. If your code requires external modules you MUST pass them here! These will installed automatically.",
			),
		),
	)

	mcpServer.AddTool(pythonTool, handlePythonExecution)

	// Run server in appropriate mode
	if *sseMode {
		// Create and start SSE server
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080"))
		log.Printf("Starting SSE server on localhost:8080")
		if err := sseServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		// Run as stdio server
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
