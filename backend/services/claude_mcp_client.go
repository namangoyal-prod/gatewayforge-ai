package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ClaudeMCPClient communicates with Claude Desktop via MCP protocol
// No API key needed - uses your Claude Desktop session
type ClaudeMCPClient struct {
	MCPEndpoint string
	HTTPClient  *http.Client
}

// MCPRequest represents a request to Claude via MCP
type MCPRequest struct {
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
}

// MCPResponse represents a response from Claude via MCP
type MCPResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *MCPError              `json:"error,omitempty"`
}

// MCPError represents an error in MCP response
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClaudeMCPClient creates a new MCP client
func NewClaudeMCPClient() *ClaudeMCPClient {
	mcpEndpoint := os.Getenv("CLAUDE_MCP_ENDPOINT")
	if mcpEndpoint == "" {
		// Default Claude Desktop MCP server endpoint
		mcpEndpoint = "http://localhost:52828"
	}

	return &ClaudeMCPClient{
		MCPEndpoint: mcpEndpoint,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// IsClaudeDesktopRunning checks if Claude Desktop is running and MCP is available
func (c *ClaudeMCPClient) IsClaudeDesktopRunning() bool {
	req := MCPRequest{
		Method:  "ping",
		Params:  map[string]interface{}{},
		JSONRPC: "2.0",
		ID:      1,
	}

	resp, err := c.sendMCPRequest(req)
	if err != nil {
		return false
	}

	return resp.Error == nil
}

// InvokeSkill invokes a Claude skill via MCP
func (c *ClaudeMCPClient) InvokeSkill(ctx context.Context, skillName string, input map[string]interface{}) (*MCPResponse, error) {
	req := MCPRequest{
		Method: "tools/call",
		Params: map[string]interface{}{
			"name":      skillName,
			"arguments": input,
		},
		JSONRPC: "2.0",
		ID:      int(time.Now().Unix()),
	}

	return c.sendMCPRequest(req)
}

// SendPrompt sends a prompt to Claude via MCP and returns the response
func (c *ClaudeMCPClient) SendPrompt(ctx context.Context, prompt string, systemPrompt string) (string, error) {
	req := MCPRequest{
		Method: "sampling/createMessage",
		Params: map[string]interface{}{
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"system":     systemPrompt,
			"max_tokens": 4096,
		},
		JSONRPC: "2.0",
		ID:      int(time.Now().Unix()),
	}

	resp, err := c.sendMCPRequest(req)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf("MCP error: %s", resp.Error.Message)
	}

	// Extract content from response
	if content, ok := resp.Result["content"].([]interface{}); ok {
		if len(content) > 0 {
			if textBlock, ok := content[0].(map[string]interface{}); ok {
				if text, ok := textBlock["text"].(string); ok {
					return text, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

// ValidateBRD uses Claude to validate a BRD document
func (c *ClaudeMCPClient) ValidateBRD(ctx context.Context, brdContent string) (*BRDValidationResult, error) {
	skillInput := map[string]interface{}{
		"brd_content": brdContent,
		"skill_name":  "brd-harmonizer",
	}

	resp, err := c.InvokeSkill(ctx, "brd-harmonizer", skillInput)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke BRD harmonizer: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("skill error: %s", resp.Error.Message)
	}

	// Parse result into BRDValidationResult
	result := &BRDValidationResult{}
	resultBytes, _ := json.Marshal(resp.Result)
	if err := json.Unmarshal(resultBytes, result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// GeneratePRD uses Claude to generate a PRD from a BRD
func (c *ClaudeMCPClient) GeneratePRD(ctx context.Context, brdContent map[string]interface{}) (*PRDGenerationResult, error) {
	skillInput := map[string]interface{}{
		"brd_content": brdContent,
		"skill_name":  "prd-generator",
	}

	resp, err := c.InvokeSkill(ctx, "prd-generator", skillInput)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke PRD generator: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("skill error: %s", resp.Error.Message)
	}

	// Parse result into PRDGenerationResult
	result := &PRDGenerationResult{}
	resultBytes, _ := json.Marshal(resp.Result)
	if err := json.Unmarshal(resultBytes, result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// GenerateCode uses Claude to generate code
func (c *ClaudeMCPClient) GenerateCode(ctx context.Context, prdContent map[string]interface{}, referenceIntegration string) (*CodeGenerationResult, error) {
	skillInput := map[string]interface{}{
		"prd_content":           prdContent,
		"reference_integration": referenceIntegration,
		"skill_name":            "coding-agent",
	}

	resp, err := c.InvokeSkill(ctx, "coding-agent", skillInput)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke coding agent: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("skill error: %s", resp.Error.Message)
	}

	// Parse result into CodeGenerationResult
	result := &CodeGenerationResult{}
	resultBytes, _ := json.Marshal(resp.Result)
	if err := json.Unmarshal(resultBytes, result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// sendMCPRequest sends an MCP request and returns the response
func (c *ClaudeMCPClient) sendMCPRequest(req MCPRequest) (*MCPResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.MCPEndpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("MCP request failed with status %d: %s", httpResp.StatusCode, string(body))
	}

	var mcpResp MCPResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &mcpResp, nil
}

// Result types for different operations
type BRDValidationResult struct {
	QualityScore       int                    `json:"quality_score"`
	Status             string                 `json:"status"`
	ValidationReport   map[string]interface{} `json:"validation_report"`
	GapAnalysis        []interface{}          `json:"gap_analysis"`
	AutoFixSuggestions []interface{}          `json:"auto_fix_suggestions"`
}

type PRDGenerationResult struct {
	PRDID      string                 `json:"prd_id"`
	Content    map[string]interface{} `json:"content"`
	Diagrams   map[string]interface{} `json:"diagrams"`
	Sections   []string               `json:"sections"`
	GeneratedAt time.Time             `json:"generated_at"`
}

type CodeGenerationResult struct {
	CodeGenerationID     string              `json:"code_generation_id"`
	GeneratedFiles       map[string][]string `json:"generated_files"`
	RepositoriesAffected []string            `json:"repositories_affected"`
	PullRequestURLs      map[string]string   `json:"pull_request_urls"`
	SecurityScanResults  map[string]interface{} `json:"security_scan_results"`
}
