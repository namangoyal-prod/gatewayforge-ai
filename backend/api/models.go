package main

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Integration represents a gateway integration pipeline
type Integration struct {
	ID              string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string         `json:"name" gorm:"not null"`
	PartnerName     string         `json:"partner_name" gorm:"not null"`
	IntegrationType string         `json:"integration_type" gorm:"not null"`
	PaymentMethods  pq.StringArray `json:"payment_methods" gorm:"type:text[]"`
	Geographies     pq.StringArray `json:"geographies" gorm:"type:text[]"`
	ExpectedGMV     float64        `json:"expected_gmv" gorm:"type:decimal(15,2)"`
	Status          string         `json:"status" gorm:"default:brd_uploaded"`
	Priority        string         `json:"priority" gorm:"default:medium"`
	CreatedBy       string         `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	Metadata        datatypes.JSON `json:"metadata" gorm:"type:jsonb;default:'{}'"`
}

// BRDDocument represents an uploaded BRD file and its validation state
type BRDDocument struct {
	ID                 string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID      string         `json:"integration_id" gorm:"type:uuid"`
	FileName           string         `json:"file_name" gorm:"not null"`
	FilePath           string         `json:"file_path" gorm:"not null"`
	FileType           string         `json:"file_type"`
	UploadedBy         string         `json:"uploaded_by" gorm:"not null"`
	UploadedAt         time.Time      `json:"uploaded_at"`
	ValidationScore    *int           `json:"validation_score,omitempty"`
	ValidationStatus   string         `json:"validation_status"`
	GapAnalysis        datatypes.JSON `json:"gap_analysis" gorm:"type:jsonb"`
	AutoFixSuggestions datatypes.JSON `json:"auto_fix_suggestions" gorm:"type:jsonb"`
	ValidatedAt        *time.Time     `json:"validated_at,omitempty"`
	ValidatedBy        string         `json:"validated_by,omitempty"`
}

func (BRDDocument) TableName() string { return "brd_documents" }

// PRDDocument represents a generated PRD for an integration
type PRDDocument struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID string         `json:"integration_id" gorm:"type:uuid"`
	BRDID         string         `json:"brd_id" gorm:"type:uuid"`
	Content       string         `json:"content" gorm:"not null"`
	GeneratedAt   time.Time      `json:"generated_at"`
	ApprovedAt    *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy    string         `json:"approved_by,omitempty"`
	Status        string         `json:"status" gorm:"default:generated"`
	Modifications datatypes.JSON `json:"modifications" gorm:"type:jsonb"`
	Diagrams      datatypes.JSON `json:"diagrams" gorm:"type:jsonb"`
}

func (PRDDocument) TableName() string { return "prd_documents" }

// CodeGeneration represents AI-generated code for an integration
type CodeGeneration struct {
	ID                    string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID         string         `json:"integration_id" gorm:"type:uuid"`
	PRDID                 string         `json:"prd_id" gorm:"type:uuid"`
	ReferenceIntegration  string         `json:"reference_integration"`
	RepositoriesAffected  pq.StringArray `json:"repositories_affected" gorm:"type:text[]"`
	GeneratedFiles        datatypes.JSON `json:"generated_files" gorm:"type:jsonb"`
	PullRequestURLs       datatypes.JSON `json:"pull_request_urls" gorm:"type:jsonb"`
	CodeReviewStatus      string         `json:"code_review_status" gorm:"default:pending"`
	SecurityScanResults   datatypes.JSON `json:"security_scan_results" gorm:"type:jsonb"`
	GeneratedAt           time.Time      `json:"generated_at"`
	ReviewedAt            *time.Time     `json:"reviewed_at,omitempty"`
	ReviewedBy            string         `json:"reviewed_by,omitempty"`
	Approved              bool           `json:"approved" gorm:"default:false"`
}

func (CodeGeneration) TableName() string { return "code_generations" }

// TestSuite represents generated and executed tests for an integration
type TestSuite struct {
	ID                      string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID           string         `json:"integration_id" gorm:"type:uuid"`
	CodeGenerationID        string         `json:"code_generation_id" gorm:"type:uuid"`
	UnitTestsCount          int            `json:"unit_tests_count" gorm:"default:0"`
	IntegrationTestsCount   int            `json:"integration_tests_count" gorm:"default:0"`
	E2ETestsCount           int            `json:"e2e_tests_count" gorm:"default:0"`
	EdgeCaseTestsCount      int            `json:"edge_case_tests_count" gorm:"default:0"`
	PerformanceTestsCount   int            `json:"performance_tests_count" gorm:"default:0"`
	SecurityTestsCount      int            `json:"security_tests_count" gorm:"default:0"`
	TotalCoveragePercent    float64        `json:"total_coverage_percent" gorm:"type:decimal(5,2)"`
	LineCoveragePercent     float64        `json:"line_coverage_percent" gorm:"type:decimal(5,2)"`
	BranchCoveragePercent   float64        `json:"branch_coverage_percent" gorm:"type:decimal(5,2)"`
	TestResults             datatypes.JSON `json:"test_results" gorm:"type:jsonb"`
	GeneratedAt             time.Time      `json:"generated_at"`
	ExecutedAt              *time.Time     `json:"executed_at,omitempty"`
	Passed                  *bool          `json:"passed,omitempty"`
}

func (TestSuite) TableName() string { return "test_suites" }

// Deployment represents a deployment of generated code
type Deployment struct {
	ID                string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID     string         `json:"integration_id" gorm:"type:uuid"`
	CodeGenerationID  string         `json:"code_generation_id" gorm:"type:uuid"`
	TestSuiteID       string         `json:"test_suite_id" gorm:"type:uuid"`
	Environment       string         `json:"environment" gorm:"not null"`
	ServicesDeployed  pq.StringArray `json:"services_deployed" gorm:"type:text[]"`
	HealthCheckStatus string         `json:"health_check_status"`
	DeploymentLogs    string         `json:"deployment_logs"`
	DeployedAt        time.Time      `json:"deployed_at"`
	DeployedBy        string         `json:"deployed_by"`
	RollbackAt        *time.Time     `json:"rollback_at,omitempty"`
	Status            string         `json:"status" gorm:"default:deploying"`
}

// AuditLog records all actions taken in the pipeline
type AuditLog struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID string         `json:"integration_id" gorm:"type:uuid"`
	UserEmail     string         `json:"user_email" gorm:"not null"`
	Action        string         `json:"action" gorm:"not null"`
	EntityType    string         `json:"entity_type"`
	EntityID      string         `json:"entity_id" gorm:"type:uuid"`
	Changes       datatypes.JSON `json:"changes" gorm:"type:jsonb"`
	Timestamp     time.Time      `json:"timestamp"`
	IPAddress     string         `json:"ip_address"`
	UserAgent     string         `json:"user_agent"`
}

func (AuditLog) TableName() string { return "audit_logs" }

// PipelineMetric tracks timing and cost for each pipeline stage
type PipelineMetric struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID string         `json:"integration_id" gorm:"type:uuid"`
	Stage         string         `json:"stage" gorm:"not null"`
	StartedAt     time.Time      `json:"started_at" gorm:"not null"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
	DurationSecs  *int           `json:"duration_seconds,omitempty"`
	Status        string         `json:"status"`
	AITokensUsed  *int           `json:"ai_tokens_used,omitempty"`
	AICostUSD     *float64       `json:"ai_cost_usd,omitempty" gorm:"type:decimal(10,4)"`
	Errors        datatypes.JSON `json:"errors" gorm:"type:jsonb"`
}

func (PipelineMetric) TableName() string { return "pipeline_metrics" }

// User represents a platform user with a role
type User struct {
	ID        string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Role      string     `json:"role" gorm:"not null"`
	Team      string     `json:"team"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

// Approval represents an approval workflow record
type Approval struct {
	ID                string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID     string         `json:"integration_id" gorm:"type:uuid"`
	Stage             string         `json:"stage" gorm:"not null"`
	EntityType        string         `json:"entity_type" gorm:"not null"`
	EntityID          string         `json:"entity_id" gorm:"type:uuid;not null"`
	ApproverRole      string         `json:"approver_role" gorm:"not null"`
	RequiredApprovers int            `json:"required_approvers" gorm:"default:1"`
	ApprovedBy        pq.StringArray `json:"approved_by" gorm:"type:text[]"`
	RejectedBy        pq.StringArray `json:"rejected_by" gorm:"type:text[]"`
	Comments          datatypes.JSON `json:"comments" gorm:"type:jsonb"`
	Status            string         `json:"status" gorm:"default:pending"`
	RequestedAt       time.Time      `json:"requested_at"`
	ResolvedAt        *time.Time     `json:"resolved_at,omitempty"`
}

// Comment represents a comment or annotation on any entity
type Comment struct {
	ID              string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IntegrationID   string     `json:"integration_id" gorm:"type:uuid"`
	EntityType      string     `json:"entity_type" gorm:"not null"`
	EntityID        string     `json:"entity_id" gorm:"type:uuid;not null"`
	Section         string     `json:"section,omitempty"`
	LineNumber      *int       `json:"line_number,omitempty"`
	UserEmail       string     `json:"user_email" gorm:"not null"`
	CommentText     string     `json:"comment_text" gorm:"not null"`
	ParentCommentID *string    `json:"parent_comment_id,omitempty" gorm:"type:uuid"`
	Resolved        bool       `json:"resolved" gorm:"default:false"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ReferenceIntegration is a known-good integration used as a coding reference
type ReferenceIntegration struct {
	ID                 string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name               string         `json:"name" gorm:"uniqueIndex;not null"`
	IntegrationType    string         `json:"integration_type"`
	Repositories       pq.StringArray `json:"repositories" gorm:"type:text[]"`
	FilePaths          datatypes.JSON `json:"file_paths" gorm:"type:jsonb"`
	PatternsExtracted  datatypes.JSON `json:"patterns_extracted" gorm:"type:jsonb"`
	QualityScore       *int           `json:"quality_score,omitempty"`
	UsageCount         int            `json:"usage_count" gorm:"default:0"`
	LastUsed           *time.Time     `json:"last_used,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
}

func (ReferenceIntegration) TableName() string { return "reference_integrations" }

// Request types

// CreateIntegrationRequest is the payload for creating a new integration
type CreateIntegrationRequest struct {
	Name            string   `json:"name" binding:"required"`
	PartnerName     string   `json:"partner_name" binding:"required"`
	IntegrationType string   `json:"integration_type" binding:"required"`
	PaymentMethods  []string `json:"payment_methods"`
	Geographies     []string `json:"geographies"`
	ExpectedGMV     float64  `json:"expected_gmv"`
	CreatedBy       string   `json:"created_by" binding:"required"`
}

// UpdateIntegrationRequest is the payload for updating an integration
type UpdateIntegrationRequest struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
}
