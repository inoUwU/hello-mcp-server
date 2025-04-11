package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"os/exec"
	"path"
	"strings"
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

	cmdArgs := []string{
		"run",
		"--rm",
		"-v",
		fmt.Sprintf("%s:/app", tmpDir),
		"mcr.microsoft.com/playwright/python:1.49.1-noble",
	}

	shArgs := []string{}

	// modules is given by llm
	if len(modules) > 0 {
		shArgs = append(shArgs, "python", "-m", "pip", "install", "--quiet")
		shArgs = append(shArgs, modules...)
		shArgs = append(shArgs, "&&")
	}

	shArgs = append(shArgs, "python", path.Join("app", "script.py"))
	cmdArgs = append(cmdArgs, "sh", "-c", strings.Join(shArgs, " "))

	cmd := exec.Command("docker", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return mcp.NewToolResultError(
				fmt.Sprintf("Python exited with code %d: %s",
					exitError.ExitCode(),
					string(exitError.Stderr),
				),
			), nil
		}
		return mcp.NewToolResultError(
			fmt.Sprintf("Execution failed: %v")), nil
	}
	return mcp.NewToolResultText(string(out)), nil

}
