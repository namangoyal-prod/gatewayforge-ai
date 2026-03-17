package services

// Claude Desktop Integration via MCP
// This would allow using Claude Desktop app instead of API key
// Requires Claude Desktop with MCP enabled

type ClaudeDesktopClient struct {
    MCPEndpoint string
}

func NewClaudeDesktopClient() *ClaudeDesktopClient {
    return &ClaudeDesktopClient{
        MCPEndpoint: "http://localhost:3000", // Claude Desktop MCP server
    }
}

// TODO: Implement MCP protocol communication
// This would authenticate via your Claude Desktop session
// No API key needed - uses your existing Claude subscription
