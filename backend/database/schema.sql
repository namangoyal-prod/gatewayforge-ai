-- GatewayForge AI Database Schema
-- PostgreSQL 14+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Integrations Table: Tracks each gateway integration through the pipeline
CREATE TABLE integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    partner_name VARCHAR(255) NOT NULL,
    integration_type VARCHAR(50) NOT NULL, -- gateway, aggregator, direct
    payment_methods TEXT[], -- Array of payment methods
    geographies TEXT[], -- Array of countries/regions
    expected_gmv DECIMAL(15, 2),
    status VARCHAR(50) NOT NULL DEFAULT 'brd_uploaded', -- Pipeline stage
    priority VARCHAR(20) DEFAULT 'medium', -- low, medium, high, critical
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb
);

-- BRD Documents Table
CREATE TABLE brd_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(20), -- pdf, docx, gdoc
    uploaded_by VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validation_score INTEGER, -- 0-100
    validation_status VARCHAR(20), -- pending, approved, rejected
    gap_analysis JSONB, -- Structured gap analysis report
    auto_fix_suggestions JSONB,
    validated_at TIMESTAMP,
    validated_by VARCHAR(255)
);

-- PRD Documents Table
CREATE TABLE prd_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    brd_id UUID REFERENCES brd_documents(id),
    content TEXT NOT NULL,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP,
    approved_by VARCHAR(255),
    status VARCHAR(20) DEFAULT 'generated', -- generated, under_review, approved, rejected
    modifications JSONB, -- Track PM edits
    diagrams JSONB -- Sequence diagrams, swim lanes
);

-- Code Generation Table
CREATE TABLE code_generations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    prd_id UUID REFERENCES prd_documents(id),
    reference_integration VARCHAR(255), -- e.g., "Codec/JustPay"
    repositories_affected TEXT[], -- List of repos
    generated_files JSONB, -- {repo: [files]}
    pull_request_urls JSONB, -- {repo: pr_url}
    code_review_status VARCHAR(20) DEFAULT 'pending',
    security_scan_results JSONB,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by VARCHAR(255),
    approved BOOLEAN DEFAULT FALSE
);

-- Test Suites Table
CREATE TABLE test_suites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    code_generation_id UUID REFERENCES code_generations(id),
    unit_tests_count INTEGER DEFAULT 0,
    integration_tests_count INTEGER DEFAULT 0,
    e2e_tests_count INTEGER DEFAULT 0,
    edge_case_tests_count INTEGER DEFAULT 0,
    performance_tests_count INTEGER DEFAULT 0,
    security_tests_count INTEGER DEFAULT 0,
    total_coverage_percent DECIMAL(5, 2),
    line_coverage_percent DECIMAL(5, 2),
    branch_coverage_percent DECIMAL(5, 2),
    test_results JSONB, -- Execution results
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    executed_at TIMESTAMP,
    passed BOOLEAN
);

-- Deployments Table
CREATE TABLE deployments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    code_generation_id UUID REFERENCES code_generations(id),
    test_suite_id UUID REFERENCES test_suites(id),
    environment VARCHAR(50) NOT NULL, -- dev, staging, prod
    services_deployed TEXT[], -- List of service names
    health_check_status VARCHAR(20), -- passed, failed, partial
    deployment_logs TEXT,
    deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deployed_by VARCHAR(255),
    rollback_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'deploying' -- deploying, deployed, failed, rolled_back
);

-- Audit Log Table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    user_email VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50), -- brd, prd, code, test, deployment
    entity_id UUID,
    changes JSONB, -- Before/after state
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT
);

-- Pipeline Metrics Table
CREATE TABLE pipeline_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    stage VARCHAR(50) NOT NULL, -- brd_validation, prd_generation, etc.
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    duration_seconds INTEGER,
    status VARCHAR(20), -- success, failed, skipped
    ai_tokens_used INTEGER,
    ai_cost_usd DECIMAL(10, 4),
    errors JSONB
);

-- User Roles and Permissions
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- solutions, product, engineering, leadership
    team VARCHAR(100),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Approval Workflows
CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    stage VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    approver_role VARCHAR(50) NOT NULL,
    required_approvers INTEGER DEFAULT 1,
    approved_by TEXT[], -- Array of user emails
    rejected_by TEXT[], -- Array of user emails
    comments JSONB, -- {user: comment}
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP
);

-- Comments and Annotations
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    section VARCHAR(255), -- For PRD section-specific comments
    line_number INTEGER, -- For code comments
    user_email VARCHAR(255) NOT NULL,
    comment_text TEXT NOT NULL,
    parent_comment_id UUID REFERENCES comments(id), -- For threaded comments
    resolved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reference Integrations (Knowledge Base)
CREATE TABLE reference_integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    integration_type VARCHAR(50),
    repositories TEXT[], -- List of repo names
    file_paths JSONB, -- {repo: [critical_files]}
    patterns_extracted JSONB, -- Coding patterns, conventions
    quality_score INTEGER, -- How good is this reference
    usage_count INTEGER DEFAULT 0,
    last_used TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_integrations_status ON integrations(status);
CREATE INDEX idx_integrations_created_at ON integrations(created_at DESC);
CREATE INDEX idx_integrations_partner ON integrations(partner_name);
CREATE INDEX idx_brd_integration ON brd_documents(integration_id);
CREATE INDEX idx_prd_integration ON prd_documents(integration_id);
CREATE INDEX idx_code_integration ON code_generations(integration_id);
CREATE INDEX idx_audit_integration ON audit_logs(integration_id);
CREATE INDEX idx_audit_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX idx_metrics_integration ON pipeline_metrics(integration_id);
CREATE INDEX idx_approvals_integration ON approvals(integration_id);
CREATE INDEX idx_approvals_status ON approvals(status);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_integrations_updated_at BEFORE UPDATE ON integrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
