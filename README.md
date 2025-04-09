# Jira MCP Server

## Setup

```bash
git clone https://github.com/koga1020/jira-mcp.git
cd jira-mcp
```

```bash
go mod tidy
```

```bash
go build
```

## Usage

### Running the server

```bash
JIRA_USERNAME=your-email@example.com JIRA_API_TOKEN=your-api-token ./jira-mcp
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
            "description": "Jira URL(e.g. https://example.atlassian.net)",
            "password": false
        }
    ],
    "servers": {
        "jira-mcp-server": {
            "type": "stdio",
            "command": "/path/to/your/jira-mcp",
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
