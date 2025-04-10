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

func CreateIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_issue",
			mcp.WithDescription("Create a new issue in Jira"),
			mcp.WithString("project_key",
				mcp.Required(),
				mcp.Description("The project key where the issue will be created"),
			),
			mcp.WithString("summary",
				mcp.Required(),
				mcp.Description("The summary/title of the issue"),
			),
			mcp.WithString("description",
				mcp.Required(),
				mcp.Description("The description of the issue"),
			),
			mcp.WithString("issue_type",
				mcp.Required(),
				mcp.Description("The type of the issue (e.g., 'Bug', 'Task', 'Story')"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectKey := request.Params.Arguments["project_key"].(string)
			summary := request.Params.Arguments["summary"].(string)
			description := request.Params.Arguments["description"].(string)
			issueType := request.Params.Arguments["issue_type"].(string)

			issueFields := jira.IssueFields{
				Project: jira.Project{
					Key: projectKey,
				},
				Summary:     summary,
				Description: description,
				Type: jira.IssueType{
					Name: issueType,
				},
			}

			issue := jira.Issue{
				Fields: &issueFields,
			}

			createdIssue, _, err := client.Issue.Create(&issue)
			if err != nil {
				return nil, err
			}

			issueJSON, err := json.Marshal(createdIssue)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(issueJSON)), nil
		}
}
