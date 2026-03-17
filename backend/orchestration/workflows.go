package orchestration

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// GatewayIntegrationWorkflowInput defines the input for the workflow
type GatewayIntegrationWorkflowInput struct {
	IntegrationID string
	BRDID         string
	ReferenceIntegration string
}

// Stage represents a pipeline stage
type Stage struct {
	Name      string
	Status    string
	StartedAt time.Time
	CompletedAt *time.Time
	Duration  *time.Duration
	Output    interface{}
	Error     error
}

// GatewayIntegrationWorkflow orchestrates the 5-stage integration pipeline
func GatewayIntegrationWorkflow(ctx workflow.Context, input GatewayIntegrationWorkflowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Gateway Integration Workflow", "integrationID", input.IntegrationID)

	var stages []Stage

	// Configure activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// ============================================================
	// STAGE 1: BRD Validation (BRD Harmonizer Skill)
	// ============================================================
	stage1 := Stage{Name: "BRD Validation", StartedAt: workflow.Now(ctx), Status: "running"}
	logger.Info("Stage 1: BRD Validation - Starting")

	var brdValidationResult BRDValidationResult
	err := workflow.ExecuteActivity(ctx, ValidateBRDActivity, input.BRDID).Get(ctx, &brdValidationResult)
	if err != nil {
		stage1.Status = "failed"
		stage1.Error = err
		logger.Error("Stage 1: BRD Validation - Failed", "error", err)
		return err
	}

	stage1.Status = "completed"
	completedAt := workflow.Now(ctx)
	stage1.CompletedAt = &completedAt
	duration := completedAt.Sub(stage1.StartedAt)
	stage1.Duration = &duration
	stage1.Output = brdValidationResult
	stages = append(stages, stage1)
	logger.Info("Stage 1: BRD Validation - Completed", "score", brdValidationResult.QualityScore)

	// Check if BRD passed validation
	if brdValidationResult.Status != "GREEN" {
		logger.Warn("BRD validation failed - workflow paused for human review")
		// Signal to wait for human approval
		var approved bool
		signalChan := workflow.GetSignalChannel(ctx, "brd-approval")
		signalChan.Receive(ctx, &approved)

		if !approved {
			return fmt.Errorf("BRD rejected by reviewer")
		}
	}

	// ============================================================
	// STAGE 2: PRD Generation (PRD Generator Skill)
	// ============================================================
	stage2 := Stage{Name: "PRD Generation", StartedAt: workflow.Now(ctx), Status: "running"}
	logger.Info("Stage 2: PRD Generation - Starting")

	var prdResult PRDGenerationResult
	err = workflow.ExecuteActivity(ctx, GeneratePRDActivity, GeneratePRDInput{
		IntegrationID: input.IntegrationID,
		BRDID:         input.BRDID,
		BRDContent:    brdValidationResult.BRDContent,
	}).Get(ctx, &prdResult)
	if err != nil {
		stage2.Status = "failed"
		stage2.Error = err
		logger.Error("Stage 2: PRD Generation - Failed", "error", err)
		return err
	}

	stage2.Status = "completed"
	completedAt = workflow.Now(ctx)
	stage2.CompletedAt = &completedAt
	duration = completedAt.Sub(stage2.StartedAt)
	stage2.Duration = &duration
	stage2.Output = prdResult
	stages = append(stages, stage2)
	logger.Info("Stage 2: PRD Generation - Completed", "prdID", prdResult.PRDID)

	// Wait for PM approval
	logger.Info("Waiting for PRD approval")
	var prdApproved bool
	prdSignalChan := workflow.GetSignalChannel(ctx, "prd-approval")
	prdSignalChan.Receive(ctx, &prdApproved)

	if !prdApproved {
		return fmt.Errorf("PRD rejected by PM")
	}

	// ============================================================
	// STAGE 3: Code Generation (Coding Agent Skill)
	// ============================================================
	stage3 := Stage{Name: "Code Generation", StartedAt: workflow.Now(ctx), Status: "running"}
	logger.Info("Stage 3: Code Generation - Starting")

	var codeResult CodeGenerationResult
	err = workflow.ExecuteActivity(ctx, GenerateCodeActivity, GenerateCodeInput{
		IntegrationID:        input.IntegrationID,
		PRDID:                prdResult.PRDID,
		PRDContent:           prdResult.Content,
		ReferenceIntegration: input.ReferenceIntegration,
	}).Get(ctx, &codeResult)
	if err != nil {
		stage3.Status = "failed"
		stage3.Error = err
		logger.Error("Stage 3: Code Generation - Failed", "error", err)
		return err
	}

	stage3.Status = "completed"
	completedAt = workflow.Now(ctx)
	stage3.CompletedAt = &completedAt
	duration = completedAt.Sub(stage3.StartedAt)
	stage3.Duration = &duration
	stage3.Output = codeResult
	stages = append(stages, stage3)
	logger.Info("Stage 3: Code Generation - Completed",
		"filesGenerated", len(codeResult.GeneratedFiles),
		"repositories", len(codeResult.RepositoriesAffected))

	// Wait for engineering approval
	logger.Info("Waiting for code approval")
	var codeApproved bool
	codeSignalChan := workflow.GetSignalChannel(ctx, "code-approval")
	codeSignalChan.Receive(ctx, &codeApproved)

	if !codeApproved {
		return fmt.Errorf("Code rejected by engineer")
	}

	// ============================================================
	// STAGE 4: Test Generation & Execution (Test Agent Skill)
	// ============================================================
	stage4 := Stage{Name: "Test Generation & Execution", StartedAt: workflow.Now(ctx), Status: "running"}
	logger.Info("Stage 4: Test Generation - Starting")

	var testResult TestGenerationResult
	err = workflow.ExecuteActivity(ctx, GenerateTestsActivity, GenerateTestsInput{
		IntegrationID:    input.IntegrationID,
		CodeGenerationID: codeResult.CodeGenerationID,
		GeneratedCode:    codeResult.GeneratedFiles,
		PRDContent:       prdResult.Content,
	}).Get(ctx, &testResult)
	if err != nil {
		stage4.Status = "failed"
		stage4.Error = err
		logger.Error("Stage 4: Test Generation - Failed", "error", err)
		return err
	}

	logger.Info("Stage 4: Tests Generated", "totalTests", testResult.TotalTests)

	// Execute tests
	logger.Info("Stage 4: Executing Tests")
	var testExecutionResult TestExecutionResult
	err = workflow.ExecuteActivity(ctx, ExecuteTestsActivity, ExecuteTestsInput{
		TestSuiteID: testResult.TestSuiteID,
	}).Get(ctx, &testExecutionResult)
	if err != nil {
		stage4.Status = "failed"
		stage4.Error = err
		logger.Error("Stage 4: Test Execution - Failed", "error", err)
		return err
	}

	stage4.Status = "completed"
	completedAt = workflow.Now(ctx)
	stage4.CompletedAt = &completedAt
	duration = completedAt.Sub(stage4.StartedAt)
	stage4.Duration = &duration
	stage4.Output = testExecutionResult
	stages = append(stages, stage4)
	logger.Info("Stage 4: Test Execution - Completed",
		"passed", testExecutionResult.PassedTests,
		"failed", testExecutionResult.FailedTests,
		"coverage", fmt.Sprintf("%.2f%%", testExecutionResult.Coverage))

	// Check if tests passed
	if testExecutionResult.FailedTests > 0 {
		logger.Warn("Some tests failed - manual review required")
		var testsApproved bool
		testsSignalChan := workflow.GetSignalChannel(ctx, "tests-approval")
		testsSignalChan.Receive(ctx, &testsApproved)

		if !testsApproved {
			return fmt.Errorf("Tests not approved")
		}
	}

	// ============================================================
	// STAGE 5: Deployment (Deploy Agent Skill)
	// ============================================================
	stage5 := Stage{Name: "Deployment", StartedAt: workflow.Now(ctx), Status: "running"}
	logger.Info("Stage 5: Deployment - Starting")

	var deployResult DeploymentResult
	err = workflow.ExecuteActivity(ctx, DeployToDevStackActivity, DeployToDevStackInput{
		IntegrationID:    input.IntegrationID,
		CodeGenerationID: codeResult.CodeGenerationID,
		TestSuiteID:      testResult.TestSuiteID,
		Environment:      "dev",
	}).Get(ctx, &deployResult)
	if err != nil {
		stage5.Status = "failed"
		stage5.Error = err
		logger.Error("Stage 5: Deployment - Failed", "error", err)
		return err
	}

	stage5.Status = "completed"
	completedAt = workflow.Now(ctx)
	stage5.CompletedAt = &completedAt
	duration = completedAt.Sub(stage5.StartedAt)
	stage5.Duration = &duration
	stage5.Output = deployResult
	stages = append(stages, stage5)
	logger.Info("Stage 5: Deployment - Completed",
		"servicesDeployed", len(deployResult.ServicesDeployed),
		"healthStatus", deployResult.HealthCheckStatus)

	// ============================================================
	// Workflow Complete
	// ============================================================
	logger.Info("Gateway Integration Workflow - Completed Successfully",
		"integrationID", input.IntegrationID,
		"totalStages", len(stages))

	// Update integration status to deployed
	err = workflow.ExecuteActivity(ctx, UpdateIntegrationStatusActivity, UpdateIntegrationStatusInput{
		IntegrationID: input.IntegrationID,
		Status:        "deployed",
	}).Get(ctx, nil)
	if err != nil {
		logger.Warn("Failed to update integration status", "error", err)
	}

	// Send completion notification
	workflow.ExecuteActivity(ctx, SendNotificationActivity, SendNotificationInput{
		IntegrationID: input.IntegrationID,
		Message:       "Integration successfully deployed to dev stack",
		Recipients:    []string{"product@razorpay.com", "engineering@razorpay.com"},
	})

	return nil
}

// Activity input/output types

type BRDValidationResult struct {
	QualityScore int
	Status       string // GREEN, AMBER, RED
	BRDContent   map[string]interface{}
	GapAnalysis  []Gap
}

type Gap struct {
	Section      string
	Gap          string
	Severity     string
	SuggestedFix string
}

type GeneratePRDInput struct {
	IntegrationID string
	BRDID         string
	BRDContent    map[string]interface{}
}

type PRDGenerationResult struct {
	PRDID   string
	Content map[string]interface{}
}

type GenerateCodeInput struct {
	IntegrationID        string
	PRDID                string
	PRDContent           map[string]interface{}
	ReferenceIntegration string
}

type CodeGenerationResult struct {
	CodeGenerationID     string
	GeneratedFiles       map[string][]string
	RepositoriesAffected []string
	PullRequestURLs      map[string]string
}

type GenerateTestsInput struct {
	IntegrationID    string
	CodeGenerationID string
	GeneratedCode    map[string][]string
	PRDContent       map[string]interface{}
}

type TestGenerationResult struct {
	TestSuiteID string
	TotalTests  int
}

type ExecuteTestsInput struct {
	TestSuiteID string
}

type TestExecutionResult struct {
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Coverage     float64
}

type DeployToDevStackInput struct {
	IntegrationID    string
	CodeGenerationID string
	TestSuiteID      string
	Environment      string
}

type DeploymentResult struct {
	DeploymentID      string
	ServicesDeployed  []string
	HealthCheckStatus string
}

type UpdateIntegrationStatusInput struct {
	IntegrationID string
	Status        string
}

type SendNotificationInput struct {
	IntegrationID string
	Message       string
	Recipients    []string
}
