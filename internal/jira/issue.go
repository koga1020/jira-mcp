package jira

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddTools registers all JIRA tools to the MCP server
func AddTools(server *mcp.Server, client *jira.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_jira_issue",
		Description: "Get issue details from Jira",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetIssueArgs) (*mcp.CallToolResult, any, error) {
		return getIssue(ctx, client, args)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_jira_issue",
		Description: "Create a new issue in Jira",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args CreateIssueArgs) (*mcp.CallToolResult, any, error) {
		return createIssue(ctx, client, args)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_jira_issue",
		Description: "Search for issues in Jira using JQL",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchIssueArgs) (*mcp.CallToolResult, any, error) {
		return searchIssue(ctx, client, args)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_jira_issue",
		Description: "Edit issue in Jira",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args EditIssueArgs) (*mcp.CallToolResult, any, error) {
		return editIssue(ctx, client, args)
	})
}

// GetIssueArgs defines the parameters for getting a JIRA issue
type GetIssueArgs struct {
	IssueKey string `json:"issue_key" jsonschema:"The issue key to retrieve"`
}

func getIssue(ctx context.Context, client *jira.Client, args GetIssueArgs) (*mcp.CallToolResult, any, error) {
	// 必要なフィールドのみ取得してトークン消費を削減
	issue, _, err := client.Issue.Get(args.IssueKey, &jira.GetQueryOptions{
		Fields: "key,summary,description,status,assignee,reporter,priority,created,updated,comment,parent,project,issuetype,labels",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get issue: %w", err)
	}

	issueJSON, err := json.Marshal(issue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal issue: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(issueJSON)},
		},
	}, nil, nil
}

// CreateIssueArgs defines the parameters for creating a JIRA issue
type CreateIssueArgs struct {
	ProjectKey  string  `json:"project_key" jsonschema:"The project key where the issue will be created"`
	Summary     string  `json:"summary" jsonschema:"The summary/title of the issue"`
	Description string  `json:"description" jsonschema:"The description of the issue"`
	IssueType   string  `json:"issue_type" jsonschema:"The type of the issue (e.g. 'Bug' 'Task' 'Story')"`
	Parent      *string `json:"parent,omitempty" jsonschema:"The parent issue key (optional)"`
}

func createIssue(ctx context.Context, client *jira.Client, args CreateIssueArgs) (*mcp.CallToolResult, any, error) {
	issueFields := jira.IssueFields{
		Project: jira.Project{
			Key: args.ProjectKey,
		},
		Summary:     args.Summary,
		Description: args.Description,
		Type: jira.IssueType{
			Name: args.IssueType,
		},
	}

	// Add parent if specified
	if args.Parent != nil && *args.Parent != "" {
		issueFields.Parent = &jira.Parent{
			Key: *args.Parent,
		}
	}

	issue := jira.Issue{
		Fields: &issueFields,
	}

	createdIssue, _, err := client.Issue.Create(&issue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create issue: %w", err)
	}

	issueJSON, err := json.Marshal(createdIssue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal created issue: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(issueJSON)},
		},
	}, nil, nil
}

// SearchIssueArgs defines the parameters for searching JIRA issues
type SearchIssueArgs struct {
	JQL        string `json:"jql" jsonschema:"JQL query string to search issues"`
	MaxResults *int   `json:"max_results,omitempty" jsonschema:"Maximum number of results to return (default: 50)"`
}

func searchIssue(ctx context.Context, client *jira.Client, args SearchIssueArgs) (*mcp.CallToolResult, any, error) {
	// 必要なフィールドのみ取得してトークン消費を削減
	searchOptions := &jira.SearchOptions{
		MaxResults: 50, // デフォルト値
		Fields:     []string{"key", "summary", "description", "status", "assignee", "reporter", "priority", "created", "updated", "parent", "project", "issuetype"},
	}

	// max_resultsが指定されていれば上書き
	if args.MaxResults != nil {
		searchOptions.MaxResults = *args.MaxResults
	}

	// JQLを使ってイシューを検索
	issues, _, err := client.Issue.Search(args.JQL, searchOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search issues: %w", err)
	}

	// 検索結果をJSONに変換
	issuesJSON, err := json.Marshal(issues)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal search results: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(issuesJSON)},
		},
	}, nil, nil
}

// EditIssueArgs defines the parameters for editing a JIRA issue
type EditIssueArgs struct {
	IssueKey    string `json:"issue_key" jsonschema:"The issue key"`
	Description string `json:"description" jsonschema:"The description of the issue"`
}

func editIssue(ctx context.Context, client *jira.Client, args EditIssueArgs) (*mcp.CallToolResult, any, error) {
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Description: args.Description,
		},
	}

	req, err := client.NewRequest("PUT", "rest/api/2/issue/"+args.IssueKey, issue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	updatedIssue := new(jira.Issue)
	_, err = client.Do(req, updatedIssue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update issue: %w", err)
	}

	issueJSON, err := json.Marshal(updatedIssue)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal updated issue: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(issueJSON)},
		},
	}, nil, nil
}
