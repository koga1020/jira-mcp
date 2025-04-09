package main

import (
	"fmt"
	mcpjira "mcp-sample/pkg/jira"
	"os"

	jira "github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Jira MCP",
		"1.0.0",
		server.WithResourceCapabilities(false, false),
		server.WithLogging(),
	)

	// Read JIRA credentials from environment variables
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_API_TOKEN"),
	}
	client, _ := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))

	s.AddTool(mcpjira.GetIssue(client))

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
