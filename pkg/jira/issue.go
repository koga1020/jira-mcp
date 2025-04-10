package jira

import (
	"context"
	"encoding/json"

	"github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_issue",
			mcp.WithDescription("Get issue details from Jira"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The issue key to retrieve"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey := request.Params.Arguments["issue_key"].(string)
			issue, _, _ := client.Issue.Get(issueKey, nil)
			issueJson, err := json.Marshal((issue))
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(issueJson)), nil

		}
}
