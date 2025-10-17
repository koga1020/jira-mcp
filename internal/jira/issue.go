package jira

import (
	"context"
	"encoding/json"

	"github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_jira_issue",
			mcp.WithDescription("Get issue details from Jira"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The issue key to retrieve"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey := request.Params.Arguments["issue_key"].(string)
			issue, _, err := client.Issue.Get(issueKey, nil)
			if err != nil {
				return nil, err
			}
			issueJson, err := json.Marshal((issue))
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(issueJson)), nil

		}
}

func CreateIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_jira_issue",
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
			mcp.WithString("parent",
				mcp.Description("The parent issue key"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectKey := request.Params.Arguments["project_key"].(string)
			summary := request.Params.Arguments["summary"].(string)
			description := request.Params.Arguments["description"].(string)
			issueType := request.Params.Arguments["issue_type"].(string)
			parentRaw, parentExists := request.Params.Arguments["parent"]

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

			if parentExists {
				parent, ok := parentRaw.(string)
				if ok && parent != "" {
					issueFields.Parent = &jira.Parent{
						Key: parent,
					}
				}
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

func SearchIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_jira_issue",
			mcp.WithDescription("Search for issues in Jira using JQL"),
			mcp.WithString("jql",
				mcp.Required(),
				mcp.Description("JQL query string to search issues"),
			),
			mcp.WithNumber("max_results",
				mcp.Description("Maximum number of results to return (default: 50)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			jql := request.Params.Arguments["jql"].(string)

			// 検索オプションの設定
			searchOptions := &jira.SearchOptions{
				MaxResults: 50, // デフォルト値
			}

			// max_resultsが指定されていれば上書き
			if maxResults, ok := request.Params.Arguments["max_results"]; ok {
				searchOptions.MaxResults = int(maxResults.(float64))
			}

			// JQLを使ってイシューを検索
			issues, _, err := client.Issue.Search(jql, searchOptions)
			if err != nil {
				return nil, err
			}

			// 検索結果をJSONに変換
			issuesJSON, err := json.Marshal(issues)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(issuesJSON)), nil
		}
}

func EditIssue(client *jira.Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("edit_jira_issue",
			mcp.WithDescription("Edit issue in Jira"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The issue key"),
			),
			mcp.WithString("description",
				mcp.Description("The description of the issue"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey := request.Params.Arguments["issue_key"].(string)
			description := request.Params.Arguments["description"].(string)

			issue := &jira.Issue{
				Fields: &jira.IssueFields{
					Description: description,
				},
			}

			req, err := client.NewRequest("PUT", "rest/api/2/issue/"+issueKey, issue)
			if err != nil {
				return nil, err
			}

			updatedIssue := new(jira.Issue)
			_, err = client.Do(req, updatedIssue)
			if err != nil {
				return nil, err
			}

			issueJSON, err := json.Marshal(updatedIssue)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(issueJSON)), nil
		}
}
