# Jira MCP Server

Jira用のModel Context Protocol (MCP) サーバー実装です。modelcontextprotocol/go-sdkを使用しています。

## Installation

### Using go install (推奨)

```bash
go install github.com/koga1020/jira-mcp@latest
```

### From source

```bash
git clone https://github.com/koga1020/jira-mcp.git
cd jira-mcp
go build
```

## Usage

### Setting up in Claude Code

`.claude/mcp.json`:

```json
{
  "mcpServers": {
    "jira": {
      "command": "jira-mcp",
      "env": {
        "JIRA_USERNAME": "your-email@example.com",
        "JIRA_API_TOKEN": "your-api-token",
        "JIRA_URL": "https://your-domain.atlassian.net"
      }
    }
  }
}
```

### Setting up in VS Code

```json
"mcp": {
    "inputs": [
        {
            "type": "promptString",
            "id": "jira_user_name",
            "description": "Jira User Email Address",
            "password": false
        },
        {
            "type": "promptString",
            "id": "jira_api_token",
            "description": "Jira API Token",
            "password": true
        },
        {
            "type": "promptString",
            "id": "jira_url",
            "description": "Jira URL (e.g. https://example.atlassian.net)",
            "password": false
        }
    ],
    "servers": {
        "jira-mcp-server": {
            "type": "stdio",
            "command": "jira-mcp",
            "args": [],
            "env": {
                "JIRA_USERNAME": "${input:jira_user_name}",
                "JIRA_API_TOKEN": "${input:jira_api_token}",
                "JIRA_URL": "${input:jira_url}"
            }
        }
    }
}
```

## Available Tools

- `get_jira_issue`: Get issue details from Jira
- `create_jira_issue`: Create a new issue in Jira
- `search_jira_issue`: Search for issues using JQL
- `edit_jira_issue`: Edit an existing issue

## Obtaining Jira API Token

1. Log in to Jira
2. Go to https://id.atlassian.com/manage-profile/security/api-tokens
3. Click "Create API token"
4. Give it a label and create
5. Copy the generated token
