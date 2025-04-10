package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// handlePythonExecution handles the execute-python tool codes
func handlePythonExecution(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, ok := request.Params.Arguments["code"].(string)
	if !ok {
		return mcp.NewToolResultError("code is required"), nil
	}

	// Handle optionale modules argument
	var modules []string
	if modulesStr, ok := request.Params.Arguments["modules"].(string); ok && modulesStr != "" {
		modules = strings.Split(modulesStr, ",")
	}

	tmpDir, err := os.MkdirTemp("", "python_repl")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create temporary directory")), nil
	}

	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	err = os.WriteFile(path.Join(tmpDir, "script.py"), []byte(code), 0644)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write script to file")), nil
	}
}
