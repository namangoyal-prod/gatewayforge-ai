package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configuration
type Config struct {
	Port         string
	DatabaseURL  string
	ClaudeAPIKey string
	TemporalURL  string
	Environment  string
}

// Global instances
var (
	db     *gorm.DB
	config *Config
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize configuration
	config = &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://admin:password@localhost:5432/gatewayforge?sslmode=disable"),
		ClaudeAPIKey: getEnv("CLAUDE_API_KEY", ""),
		TemporalURL:  getEnv("TEMPORAL_URL", "localhost:7233"),
		Environment:  getEnv("ENV", "development"),
	}

	// Initialize database
	initDB()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(CORSMiddleware())

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Integration routes
		integrations := v1.Group("/integrations")
		{
			integrations.GET("", listIntegrations)
			integrations.POST("", createIntegration)
			integrations.GET("/:id", getIntegration)
			integrations.PUT("/:id", updateIntegration)
			integrations.DELETE("/:id", deleteIntegration)
			integrations.GET("/:id/status", getIntegrationStatus)
			integrations.GET("/:id/timeline", getIntegrationTimeline)
		}

		// BRD routes
		brds := v1.Group("/brds")
		{
			brds.POST("", uploadBRD)
			brds.GET("/:id", getBRD)
			brds.POST("/:id/validate", validateBRD)
			brds.POST("/:id/approve", approveBRD)
			brds.POST("/:id/reject", rejectBRD)
			brds.GET("/:id/gap-analysis", getBRDGapAnalysis)
		}

		// PRD routes
		prds := v1.Group("/prds")
		{
			prds.GET("/:id", getPRD)
			prds.PUT("/:id", updatePRD)
			prds.POST("/:id/approve", approvePRD)
			prds.POST("/:id/reject", rejectPRD)
		}

		// Code generation routes
		code := v1.Group("/code")
		{
			code.GET("/:id", getCodeGeneration)
			code.GET("/:id/files", getGeneratedFiles)
			code.POST("/:id/approve", approveCode)
			code.GET("/:id/security-scan", getSecurityScanResults)
		}

		// Test routes
		tests := v1.Group("/tests")
		{
			tests.GET("/:id", getTestSuite)
			tests.POST("/:id/execute", executeTests)
			tests.GET("/:id/coverage", getTestCoverage)
			tests.GET("/:id/results", getTestResults)
		}

		// Deployment routes
		deployments := v1.Group("/deployments")
		{
			deployments.GET("/:id", getDeployment)
			deployments.POST("/:id/health-check", runHealthCheck)
			deployments.POST("/:id/rollback", rollbackDeployment)
		}

		// Approval routes
		approvals := v1.Group("/approvals")
		{
			approvals.GET("", listPendingApprovals)
			approvals.POST("/:id/approve", approveEntity)
			approvals.POST("/:id/reject", rejectEntity)
		}

		// Comment routes
		comments := v1.Group("/comments")
		{
			comments.GET("", listComments)
			comments.POST("", createComment)
			comments.PUT("/:id", updateComment)
			comments.DELETE("/:id", deleteComment)
			comments.POST("/:id/resolve", resolveComment)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/pipeline-metrics", getPipelineMetrics)
			analytics.GET("/quality-trends", getQualityTrends)
			analytics.GET("/cost-analysis", getCostAnalysis)
			analytics.GET("/velocity", getVelocityMetrics)
		}

		// Reference integrations
		references := v1.Group("/references")
		{
			references.GET("", listReferenceIntegrations)
			references.GET("/:name", getReferenceIntegration)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/me", getCurrentUser)
			users.GET("", listUsers)
		}

		// Estimation engine — AI-assisted vs traditional timeline estimator
		estimation := v1.Group("/estimation")
		{
			estimation.POST("/estimate", estimateIntegration)
		}

		// Health check
		v1.GET("/health", healthCheck)
	}

	// Start server
	log.Printf("Starting GatewayForge AI API on port %s", config.Port)
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Database initialization
func initDB() {
	var err error
	db, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Health check handler
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "gatewayforge-ai-api",
		"version":   "1.0.0",
	})
}

// Integration handlers
func listIntegrations(c *gin.Context) {
	var integrations []Integration
	result := db.Order("created_at DESC").Find(&integrations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"integrations": integrations,
		"total":        len(integrations),
	})
}

func createIntegration(c *gin.Context) {
	var req CreateIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	integration := Integration{
		Name:          req.Name,
		PartnerName:   req.PartnerName,
		IntegrationType: req.IntegrationType,
		PaymentMethods: req.PaymentMethods,
		Geographies:   req.Geographies,
		ExpectedGMV:   req.ExpectedGMV,
		Status:        "brd_uploaded",
		CreatedBy:     req.CreatedBy,
	}

	if result := db.Create(&integration); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, integration)
}

func getIntegration(c *gin.Context) {
	id := c.Param("id")
	var integration Integration

	if result := db.First(&integration, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func updateIntegration(c *gin.Context) {
	id := c.Param("id")
	var req UpdateIntegrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var integration Integration
	if result := db.First(&integration, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	// Update fields
	if req.Status != "" {
		integration.Status = req.Status
	}
	if req.Priority != "" {
		integration.Priority = req.Priority
	}

	if result := db.Save(&integration); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func deleteIntegration(c *gin.Context) {
	id := c.Param("id")
	if result := db.Delete(&Integration{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Integration deleted successfully"})
}

func getIntegrationStatus(c *gin.Context) {
	id := c.Param("id")

	var integration Integration
	if result := db.First(&integration, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	// Fetch related entities
	var brd BRDDocument
	db.Where("integration_id = ?", id).First(&brd)

	var prd PRDDocument
	db.Where("integration_id = ?", id).First(&prd)

	var code CodeGeneration
	db.Where("integration_id = ?", id).First(&code)

	var test TestSuite
	db.Where("integration_id = ?", id).First(&test)

	var deployment Deployment
	db.Where("integration_id = ?", id).First(&deployment)

	c.JSON(http.StatusOK, gin.H{
		"integration": integration,
		"brd":         brd,
		"prd":         prd,
		"code":        code,
		"test":        test,
		"deployment":  deployment,
	})
}

func getIntegrationTimeline(c *gin.Context) {
	id := c.Param("id")

	var metrics []PipelineMetric
	result := db.Where("integration_id = ?", id).Order("started_at ASC").Find(&metrics)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeline": metrics,
	})
}

// BRD handlers
func uploadBRD(c *gin.Context) {
	// Handle file upload
	file, err := c.FormFile("brd")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	integrationID := c.PostForm("integration_id")
	uploadedBy := c.PostForm("uploaded_by")

	// Save file to storage
	filepath := fmt.Sprintf("./uploads/brds/%s/%s", integrationID, file.Filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create BRD record
	brd := BRDDocument{
		IntegrationID:    integrationID,
		FileName:         file.Filename,
		FilePath:         filepath,
		FileType:         getFileType(file.Filename),
		UploadedBy:       uploadedBy,
		ValidationStatus: "pending",
	}

	if result := db.Create(&brd); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Trigger BRD Harmonizer workflow
	go triggerBRDValidation(brd.ID)

	c.JSON(http.StatusCreated, brd)
}

func getBRD(c *gin.Context) {
	id := c.Param("id")
	var brd BRDDocument

	if result := db.First(&brd, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BRD not found"})
		return
	}

	c.JSON(http.StatusOK, brd)
}

func validateBRD(c *gin.Context) {
	id := c.Param("id")

	// Trigger validation workflow
	go triggerBRDValidation(id)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Validation started",
		"brd_id":  id,
	})
}

func approveBRD(c *gin.Context) {
	id := c.Param("id")

	var brd BRDDocument
	if result := db.First(&brd, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BRD not found"})
		return
	}

	brd.ValidationStatus = "approved"
	brd.ValidatedAt = timePtr(time.Now())
	brd.ValidatedBy = c.PostForm("validated_by")

	db.Save(&brd)

	// Update integration status
	db.Model(&Integration{}).Where("id = ?", brd.IntegrationID).Update("status", "prd_generation")

	// Trigger PRD generation
	go triggerPRDGeneration(brd.IntegrationID, id)

	c.JSON(http.StatusOK, brd)
}

func rejectBRD(c *gin.Context) {
	id := c.Param("id")

	var brd BRDDocument
	if result := db.First(&brd, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BRD not found"})
		return
	}

	brd.ValidationStatus = "rejected"
	brd.ValidatedAt = timePtr(time.Now())
	brd.ValidatedBy = c.PostForm("validated_by")

	db.Save(&brd)

	// Update integration status
	db.Model(&Integration{}).Where("id = ?", brd.IntegrationID).Update("status", "brd_rejected")

	c.JSON(http.StatusOK, brd)
}

func getBRDGapAnalysis(c *gin.Context) {
	id := c.Param("id")

	var brd BRDDocument
	if result := db.First(&brd, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BRD not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gap_analysis":        brd.GapAnalysis,
		"auto_fix_suggestions": brd.AutoFixSuggestions,
		"validation_score":    brd.ValidationScore,
	})
}

// Placeholder functions for async workflow triggers
func triggerBRDValidation(brdID string) {
	log.Printf("Triggering BRD validation for ID: %s", brdID)
	// TODO: Implement Temporal workflow trigger
}

func triggerPRDGeneration(integrationID, brdID string) {
	log.Printf("Triggering PRD generation for integration: %s", integrationID)
	// TODO: Implement Temporal workflow trigger
}

// Helper functions
func getFileType(filename string) string {
	if len(filename) < 4 {
		return "unknown"
	}
	ext := filename[len(filename)-4:]
	switch ext {
	case ".pdf":
		return "pdf"
	case "docx":
		return "docx"
	default:
		return "unknown"
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Stub implementations for remaining handlers
func getPRD(c *gin.Context)                  { c.JSON(http.StatusOK, gin.H{}) }
func updatePRD(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{}) }
func approvePRD(c *gin.Context)              { c.JSON(http.StatusOK, gin.H{}) }
func rejectPRD(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{}) }
func getCodeGeneration(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{}) }
func getGeneratedFiles(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{}) }
func approveCode(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{}) }
func getSecurityScanResults(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{}) }
func getTestSuite(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{}) }
func executeTests(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{}) }
func getTestCoverage(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{}) }
func getTestResults(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{}) }
func getDeployment(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{}) }
func runHealthCheck(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{}) }
func rollbackDeployment(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{}) }
func listPendingApprovals(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{}) }
func approveEntity(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{}) }
func rejectEntity(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{}) }
func listComments(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{}) }
func createComment(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{}) }
func updateComment(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{}) }
func deleteComment(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{}) }
func resolveComment(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{}) }
func getPipelineMetrics(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{}) }
func getQualityTrends(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{}) }
func getCostAnalysis(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{}) }
func getVelocityMetrics(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{}) }
func listReferenceIntegrations(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func getReferenceIntegration(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func getCurrentUser(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{}) }
func listUsers(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{}) }
