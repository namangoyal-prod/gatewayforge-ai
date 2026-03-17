package main

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// Integration represents a gateway integration through the pipeline
type Integration struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name            string         `json:"name" gorm:"not null"`
	PartnerName     string         `json:"partner_name" gorm:"not null"`
	IntegrationType string         `json:"integration_type" gorm:"not null"` // gateway, aggregator, direct
	PaymentMethods  pq.StringArray `json:"payment_methods" gorm:"type:text[]"`
	Geographies     pq.StringArray `json:"geographies" gorm:"type:text[]"`
	ExpectedGMV     float64        `json:"expected_gmv"`
	Status          string         `json:"status" gorm:"default:'brd_uploaded'"`
	Priority        string         `json:"priority" gorm:"default:'medium'"`
	CreatedBy       string         `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	CompletedAt     *time.Time     `json:"completed_at"`
	Metadata        JSONB          `json:"metadata" gorm:"type:jsonb;default:'{}'"`
}

// BRDDocument represents an uploaded BRD
type BRDDocument struct {
	ID                  string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID       string     `json:"integration_id" gorm:"not null"`
	FileName            string     `json:"file_name" gorm:"not null"`
	FilePath            string     `json:"file_path" gorm:"not null"`
	FileType            string     `json:"file_type"`
	UploadedBy          string     `json:"uploaded_by" gorm:"not null"`
	UploadedAt          time.Time  `json:"uploaded_at" gorm:"default:CURRENT_TIMESTAMP"`
	ValidationScore     *int       `json:"validation_score"`
	ValidationStatus    string     `json:"validation_status"`
	GapAnalysis         JSONB      `json:"gap_analysis" gorm:"type:jsonb"`
	AutoFixSuggestions  JSONB      `json:"auto_fix_suggestions" gorm:"type:jsonb"`
	ValidatedAt         *time.Time `json:"validated_at"`
	ValidatedBy         string     `json:"validated_by"`
}

// PRDDocument represents a generated PRD
type PRDDocument struct {
	ID            string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID string     `json:"integration_id" gorm:"not null"`
	BRDID         string     `json:"brd_id"`
	Content       string     `json:"content" gorm:"type:text;not null"`
	GeneratedAt   time.Time  `json:"generated_at" gorm:"default:CURRENT_TIMESTAMP"`
	ApprovedAt    *time.Time `json:"approved_at"`
	ApprovedBy    string     `json:"approved_by"`
	Status        string     `json:"status" gorm:"default:'generated'"`
	Modifications JSONB      `json:"modifications" gorm:"type:jsonb"`
	Diagrams      JSONB      `json:"diagrams" gorm:"type:jsonb"`
}

// CodeGeneration represents generated code
type CodeGeneration struct {
	ID                   string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID        string         `json:"integration_id" gorm:"not null"`
	PRDID                string         `json:"prd_id"`
	ReferenceIntegration string         `json:"reference_integration"`
	RepositoriesAffected pq.StringArray `json:"repositories_affected" gorm:"type:text[]"`
	GeneratedFiles       JSONB          `json:"generated_files" gorm:"type:jsonb"`
	PullRequestURLs      JSONB          `json:"pull_request_urls" gorm:"type:jsonb"`
	CodeReviewStatus     string         `json:"code_review_status" gorm:"default:'pending'"`
	SecurityScanResults  JSONB          `json:"security_scan_results" gorm:"type:jsonb"`
	GeneratedAt          time.Time      `json:"generated_at" gorm:"default:CURRENT_TIMESTAMP"`
	ReviewedAt           *time.Time     `json:"reviewed_at"`
	ReviewedBy           string         `json:"reviewed_by"`
	Approved             bool           `json:"approved" gorm:"default:false"`
}

// TestSuite represents generated tests
type TestSuite struct {
	ID                     string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID          string     `json:"integration_id" gorm:"not null"`
	CodeGenerationID       string     `json:"code_generation_id"`
	UnitTestsCount         int        `json:"unit_tests_count" gorm:"default:0"`
	IntegrationTestsCount  int        `json:"integration_tests_count" gorm:"default:0"`
	E2ETestsCount          int        `json:"e2e_tests_count" gorm:"default:0"`
	EdgeCaseTestsCount     int        `json:"edge_case_tests_count" gorm:"default:0"`
	PerformanceTestsCount  int        `json:"performance_tests_count" gorm:"default:0"`
	SecurityTestsCount     int        `json:"security_tests_count" gorm:"default:0"`
	TotalCoveragePercent   float64    `json:"total_coverage_percent"`
	LineCoveragePercent    float64    `json:"line_coverage_percent"`
	BranchCoveragePercent  float64    `json:"branch_coverage_percent"`
	TestResults            JSONB      `json:"test_results" gorm:"type:jsonb"`
	GeneratedAt            time.Time  `json:"generated_at" gorm:"default:CURRENT_TIMESTAMP"`
	ExecutedAt             *time.Time `json:"executed_at"`
	Passed                 *bool      `json:"passed"`
}

// Deployment represents a deployment to an environment
type Deployment struct {
	ID                 string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID      string         `json:"integration_id" gorm:"not null"`
	CodeGenerationID   string         `json:"code_generation_id"`
	TestSuiteID        string         `json:"test_suite_id"`
	Environment        string         `json:"environment" gorm:"not null"`
	ServicesDeployed   pq.StringArray `json:"services_deployed" gorm:"type:text[]"`
	HealthCheckStatus  string         `json:"health_check_status"`
	DeploymentLogs     string         `json:"deployment_logs" gorm:"type:text"`
	DeployedAt         time.Time      `json:"deployed_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeployedBy         string         `json:"deployed_by"`
	RollbackAt         *time.Time     `json:"rollback_at"`
	Status             string         `json:"status" gorm:"default:'deploying'"`
}

// AuditLog represents audit trail
type AuditLog struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID string    `json:"integration_id"`
	UserEmail     string    `json:"user_email" gorm:"not null"`
	Action        string    `json:"action" gorm:"not null"`
	EntityType    string    `json:"entity_type"`
	EntityID      string    `json:"entity_id"`
	Changes       JSONB     `json:"changes" gorm:"type:jsonb"`
	Timestamp     time.Time `json:"timestamp" gorm:"default:CURRENT_TIMESTAMP"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent" gorm:"type:text"`
}

// PipelineMetric represents pipeline stage metrics
type PipelineMetric struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID  string     `json:"integration_id" gorm:"not null"`
	Stage          string     `json:"stage" gorm:"not null"`
	StartedAt      time.Time  `json:"started_at" gorm:"not null"`
	CompletedAt    *time.Time `json:"completed_at"`
	DurationSeconds *int      `json:"duration_seconds"`
	Status         string     `json:"status"`
	AITokensUsed   *int       `json:"ai_tokens_used"`
	AICostUSD      *float64   `json:"ai_cost_usd"`
	Errors         JSONB      `json:"errors" gorm:"type:jsonb"`
}

// User represents a system user
type User struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Role      string     `json:"role" gorm:"not null"` // solutions, product, engineering, leadership
	Team      string     `json:"team"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	LastLogin *time.Time `json:"last_login"`
}

// Approval represents an approval workflow
type Approval struct {
	ID                string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID     string         `json:"integration_id" gorm:"not null"`
	Stage             string         `json:"stage" gorm:"not null"`
	EntityType        string         `json:"entity_type" gorm:"not null"`
	EntityID          string         `json:"entity_id" gorm:"not null"`
	ApproverRole      string         `json:"approver_role" gorm:"not null"`
	RequiredApprovers int            `json:"required_approvers" gorm:"default:1"`
	ApprovedBy        pq.StringArray `json:"approved_by" gorm:"type:text[]"`
	RejectedBy        pq.StringArray `json:"rejected_by" gorm:"type:text[]"`
	Comments          JSONB          `json:"comments" gorm:"type:jsonb"`
	Status            string         `json:"status" gorm:"default:'pending'"`
	RequestedAt       time.Time      `json:"requested_at" gorm:"default:CURRENT_TIMESTAMP"`
	ResolvedAt        *time.Time     `json:"resolved_at"`
}

// Comment represents a comment or annotation
type Comment struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IntegrationID   string     `json:"integration_id" gorm:"not null"`
	EntityType      string     `json:"entity_type" gorm:"not null"`
	EntityID        string     `json:"entity_id" gorm:"not null"`
	Section         string     `json:"section"`
	LineNumber      *int       `json:"line_number"`
	UserEmail       string     `json:"user_email" gorm:"not null"`
	CommentText     string     `json:"comment_text" gorm:"type:text;not null"`
	ParentCommentID *string    `json:"parent_comment_id"`
	Resolved        bool       `json:"resolved" gorm:"default:false"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ReferenceIntegration represents a reference integration in knowledge base
type ReferenceIntegration struct {
	ID               string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name             string         `json:"name" gorm:"unique;not null"`
	IntegrationType  string         `json:"integration_type"`
	Repositories     pq.StringArray `json:"repositories" gorm:"type:text[]"`
	FilePaths        JSONB          `json:"file_paths" gorm:"type:jsonb"`
	PatternsExtracted JSONB         `json:"patterns_extracted" gorm:"type:jsonb"`
	QualityScore     int            `json:"quality_score"`
	UsageCount       int            `json:"usage_count" gorm:"default:0"`
	LastUsed         *time.Time     `json:"last_used"`
	CreatedAt        time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal(value.([]byte), j)
	}

	return json.Unmarshal(bytes, j)
}

// Request/Response types
type CreateIntegrationRequest struct {
	Name            string   `json:"name" binding:"required"`
	PartnerName     string   `json:"partner_name" binding:"required"`
	IntegrationType string   `json:"integration_type" binding:"required"`
	PaymentMethods  []string `json:"payment_methods" binding:"required"`
	Geographies     []string `json:"geographies" binding:"required"`
	ExpectedGMV     float64  `json:"expected_gmv"`
	CreatedBy       string   `json:"created_by" binding:"required"`
}

type UpdateIntegrationRequest struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type ValidationResponse struct {
	QualityScore       int                      `json:"quality_score"`
	Status             string                   `json:"status"`
	ValidationReport   ValidationReport         `json:"validation_report"`
	GapAnalysis        []Gap                    `json:"gap_analysis"`
	AutoFixSuggestions []AutoFixSuggestion      `json:"auto_fix_suggestions"`
	ComparisonMatrix   ComparisonMatrix         `json:"comparison_matrix"`
}

type ValidationReport struct {
	Completeness        DimensionScore `json:"completeness"`
	TechnicalAccuracy   DimensionScore `json:"technical_accuracy"`
	Conformance         DimensionScore `json:"conformance"`
	Clarity             DimensionScore `json:"clarity"`
	RegulatoryCompliance DimensionScore `json:"regulatory_compliance"`
}

type DimensionScore struct {
	Score  int      `json:"score"`
	Issues []string `json:"issues"`
}

type Gap struct {
	Section      string `json:"section"`
	Gap          string `json:"gap"`
	Severity     string `json:"severity"`
	SuggestedFix string `json:"suggested_fix"`
}

type AutoFixSuggestion struct {
	Field     string      `json:"field"`
	Current   interface{} `json:"current"`
	Suggested interface{} `json:"suggested"`
	Rationale string      `json:"rationale"`
}

type ComparisonMatrix struct {
	ReferenceBRD    string   `json:"reference_brd"`
	SimilarityScore int      `json:"similarity_score"`
	KeyDifferences  []string `json:"key_differences"`
}
