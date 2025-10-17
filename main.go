package main

import (
	"context"
	"log"
	"os"

	jira "github.com/andygrunwald/go-jira"
	mcpjira "github.com/koga1020/jira-mcp/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Read JIRA credentials from environment variables
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_API_TOKEN"),
	}
	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		log.Fatalf("Failed to create JIRA client: %v", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "jira-mcp",
		Version: "1.0.0",
	}, nil)

	// Add tools
	mcpjira.AddTools(server, client)

	// Start the server
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
