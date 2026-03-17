package orchestration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/razorpay/gatewayforge-ai/services"
)

// MCPActivities implements Temporal activities using Claude Desktop MCP
type MCPActivities struct {
	claudeClient *services.ClaudeMCPClient
}

// NewMCPActivities creates a new instance of MCP-based activities
func NewMCPActivities() *MCPActivities {
	return &MCPActivities{
		claudeClient: services.NewClaudeMCPClient(),
	}
}

// ValidateBRDActivityMCP validates BRD using Claude Desktop via MCP
func (a *MCPActivities) ValidateBRDActivityMCP(ctx context.Context, brdID string) (*BRDValidationResult, error) {
	log.Printf("Validating BRD via MCP: %s", brdID)

	// Check if Claude Desktop is running
	if !a.claudeClient.IsClaudeDesktopRunning() {
		return nil, fmt.Errorf("Claude Desktop is not running. Please start Claude Desktop and ensure MCP is enabled")
	}

	// TODO: Fetch BRD content from database
	brdContent := "Sample BRD content for " + brdID

	// Call Claude via MCP
	result, err := a.claudeClient.ValidateBRD(ctx, brdContent)
	if err != nil {
		return nil, fmt.Errorf("BRD validation failed: %w", err)
	}

	// Convert to workflow result type
	workflowResult := &BRDValidationResult{
		QualityScore: result.QualityScore,
		Status:       result.Status,
		BRDContent: map[string]interface{}{
			"validated": true,
			"timestamp": time.Now(),
		},
		GapAnalysis: []Gap{},
	}

	log.Printf("BRD Validation completed via MCP: Score=%d, Status=%s", result.QualityScore, result.Status)
	return workflowResult, nil
}

// GeneratePRDActivityMCP generates PRD using Claude Desktop via MCP
func (a *MCPActivities) GeneratePRDActivityMCP(ctx context.Context, input GeneratePRDInput) (*PRDGenerationResult, error) {
	log.Printf("Generating PRD via MCP for Integration: %s", input.IntegrationID)

	// Check if Claude Desktop is running
	if !a.claudeClient.IsClaudeDesktopRunning() {
		return nil, fmt.Errorf("Claude Desktop is not running. Please start Claude Desktop and ensure MCP is enabled")
	}

	// Call Claude via MCP
	mcpResult, err := a.claudeClient.GeneratePRD(ctx, input.BRDContent)
	if err != nil {
		return nil, fmt.Errorf("PRD generation failed: %w", err)
	}

	// Convert to workflow result type
	result := &PRDGenerationResult{
		PRDID:   mcpResult.PRDID,
		Content: mcpResult.Content,
	}

	log.Printf("PRD Generation completed via MCP: PRDID=%s", result.PRDID)
	return result, nil
}

// GenerateCodeActivityMCP generates code using Claude Desktop via MCP
func (a *MCPActivities) GenerateCodeActivityMCP(ctx context.Context, input GenerateCodeInput) (*CodeGenerationResult, error) {
	log.Printf("Generating code via MCP for Integration: %s, Reference: %s",
		input.IntegrationID, input.ReferenceIntegration)

	// Check if Claude Desktop is running
	if !a.claudeClient.IsClaudeDesktopRunning() {
		return nil, fmt.Errorf("Claude Desktop is not running. Please start Claude Desktop and ensure MCP is enabled")
	}

	// Call Claude via MCP
	mcpResult, err := a.claudeClient.GenerateCode(ctx, input.PRDContent, input.ReferenceIntegration)
	if err != nil {
		return nil, fmt.Errorf("Code generation failed: %w", err)
	}

	// Convert to workflow result type
	result := &CodeGenerationResult{
		CodeGenerationID:     mcpResult.CodeGenerationID,
		GeneratedFiles:       mcpResult.GeneratedFiles,
		RepositoriesAffected: mcpResult.RepositoriesAffected,
		PullRequestURLs:      mcpResult.PullRequestURLs,
	}

	log.Printf("Code Generation completed via MCP: Files=%d, Repos=%d",
		len(result.GeneratedFiles), len(result.RepositoriesAffected))
	return result, nil
}

// HealthCheckActivityMCP checks if Claude Desktop MCP is available
func (a *MCPActivities) HealthCheckActivityMCP(ctx context.Context) error {
	if !a.claudeClient.IsClaudeDesktopRunning() {
		return fmt.Errorf("Claude Desktop MCP is not available. Please ensure:\n" +
			"1. Claude Desktop is running\n" +
			"2. MCP server is enabled in Claude Desktop settings\n" +
			"3. MCP endpoint is accessible at %s", a.claudeClient.MCPEndpoint)
	}

	log.Println("✅ Claude Desktop MCP health check passed")
	return nil
}
