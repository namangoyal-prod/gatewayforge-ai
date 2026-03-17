# GatewayForge AI

**Autonomous Gateway Integration Platform**
*From BRD to Production-Ready Code in Hours, Not Weeks*

---

## Overview

GatewayForge AI transforms how Razorpay integrates with new payment gateways, banking partners, and payment networks. The platform reduces integration timelines from 6-12 weeks to 2-3 days by leveraging AI-powered automation at every stage of the lifecycle.

### Key Features

- **BRD Harmonizer**: Validates BRDs against Razorpay's integration standards (completeness, technical accuracy, compliance)
- **PRD Generator**: Auto-generates production-grade PRDs from approved BRDs
- **Coding Agent**: Generates code across all repositories using SWE Agent and reference patterns
- **Test Agent**: Creates comprehensive test suites with 90%+ coverage
- **Deploy Agent**: Automates dev stack deployment with health checks

### Impact

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Integration TAT | 6-12 weeks | 2-3 days | **95% reduction** |
| BRD Rejection Rate | 40% | <5% | **88% improvement** |
| Code Coverage | 60-70% | 90%+ | **30% increase** |
| PM Capacity | 8-10/year | 40-50/year | **5x increase** |
| Cost per Integration | ₹35-50L | ₹5-10L | **80% reduction** |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    GatewayForge AI Platform                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │  Frontend    │    │   Backend    │    │  AI Skills   │  │
│  │  (React +    │───▶│   (Go API)   │───▶│  (Claude)    │  │
│  │   Blade)     │    │              │    │              │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│                              │                              │
│                              ▼                              │
│                    ┌──────────────────┐                     │
│                    │  Orchestration   │                     │
│                    │  (Temporal.io)   │                     │
│                    └──────────────────┘                     │
│                              │                              │
│                              ▼                              │
│                    ┌──────────────────┐                     │
│                    │    Database      │                     │
│                    │   (PostgreSQL)   │                     │
│                    └──────────────────┘                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Quick Start

### Prerequisites

- **Go** 1.21+
- **Node.js** 18+
- **PostgreSQL** 14+
- **Docker** & Docker Compose
- **Claude Authentication** (choose one):
  - **Option A**: Claude Desktop (recommended for development) - [Quick Start](QUICK_START_CLAUDE_DESKTOP.md)
  - **Option B**: Claude API Key (for production) - Get from https://console.anthropic.com/

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/razorpay/gatewayforge-ai.git
   cd gatewayforge-ai
   ```

2. **Set up the database**
   ```bash
   createdb gatewayforge
   psql gatewayforge < backend/database/schema.sql
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration:
   # - DATABASE_URL
   # - CLAUDE_API_KEY
   # - TEMPORAL_URL
   ```

4. **Start backend API**
   ```bash
   cd backend/api
   go mod download
   go run main.go
   # API will be available at http://localhost:8080
   ```

5. **Start frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   # Frontend will be available at http://localhost:5173
   ```

6. **Start Temporal (Optional - for full workflow orchestration)**
   ```bash
   docker-compose up -d temporal
   cd backend/orchestration
   go run worker/main.go
   ```

---

## Usage

### 1. Upload BRD

Navigate to the dashboard and click "New Integration". Fill in partner details and upload the BRD document (PDF or DOCX).

### 2. BRD Validation

The BRD Harmonizer skill automatically validates the document:
- **Completeness**: All mandatory sections present
- **Technical Accuracy**: Valid ISO 8583, API specs, encryption
- **Conformance**: Follows Razorpay's template
- **Clarity**: No ambiguous language
- **Compliance**: Regulatory requirements met

**Score ≥70 = GREEN** → Proceeds to PRD generation
**Score 50-69 = AMBER** → Minor fixes needed
**Score <50 = RED** → Major rework required

### 3. PRD Review & Approval

Auto-generated PRD covers:
- Integration overview
- Merchant onboarding journey
- Payment processing flows
- Settlement & reconciliation
- Error handling
- API specifications
- Non-functional requirements

PMs can review, edit inline, and approve the PRD.

### 4. Code Generation

The Coding Agent analyzes reference integrations (e.g., Codec/JustPay) and generates code across:
- Gateway adapter service
- Routing engine
- Payment processing core
- Settlement service
- Webhook handlers
- API gateway config
- Monitoring dashboards

Engineers review generated code, approve pull requests.

### 5. Test Generation & Execution

Test Agent generates:
- Unit tests (200-400)
- Integration tests (50-100)
- E2E tests (30-50)
- Edge case tests (40-80)
- Performance tests (10-20)
- Security tests (20-30)

Target: 90%+ coverage. All tests are executed automatically.

### 6. Deployment

Deploy Agent provisions services on dev stack:
- Builds Docker images
- Deploys to Kubernetes
- Runs health checks
- Executes smoke tests
- Generates deployment report

---

## Project Structure

```
gatewayforge-ai/
├── backend/
│   ├── api/                    # REST API (Go + Gin)
│   │   ├── main.go
│   │   └── handlers/
│   ├── orchestration/          # Temporal workflows
│   │   ├── workflows.go
│   │   ├── activities.go
│   │   └── worker/
│   ├── database/
│   │   └── schema.sql          # PostgreSQL schema
│   ├── models/
│   │   └── models.go           # Data models
│   └── services/               # Business logic
├── frontend/
│   ├── src/
│   │   ├── pages/              # React pages
│   │   │   ├── Dashboard.tsx
│   │   │   ├── BRDUpload.tsx
│   │   │   ├── BRDValidation.tsx
│   │   │   ├── PRDReview.tsx
│   │   │   ├── CodeReview.tsx
│   │   │   ├── TestExecution.tsx
│   │   │   └── Deployment.tsx
│   │   ├── components/         # Reusable components
│   │   ├── api/                # API client
│   │   └── App.tsx
│   └── package.json
├── skills/                     # AI Skills (Claude)
│   ├── brd-harmonizer/
│   │   └── SKILL.md
│   ├── prd-generator/
│   │   └── SKILL.md
│   ├── coding-agent/
│   │   └── SKILL.md
│   ├── test-agent/
│   │   └── SKILL.md
│   └── deploy-agent/
│       └── SKILL.md
├── config/
│   ├── temporal/
│   ├── database/
│   └── nginx/
├── docs/
│   ├── api/
│   ├── architecture/
│   └── deployment/
└── README.md
```

---

## AI Skills

Each skill is a specialized Claude agent with domain expertise:

### BRD Harmonizer
**Input**: BRD document
**Output**: Quality score, gap analysis, auto-fix suggestions
**Knowledge Base**: Historical BRDs, ISO 8583 specs, NPCI guidelines, PCI-DSS requirements

### PRD Generator
**Input**: Approved BRD
**Output**: Complete PRD with 8 sections, diagrams, test scenarios
**Intelligence**: Cross-references similar integrations, fills implicit requirements

### Coding Agent
**Input**: Approved PRD
**Output**: Production-ready code across 10+ repositories
**Powered By**: SWE Agent for codebase analysis and pattern extraction

### Test Agent
**Input**: Generated code
**Output**: 500-600 tests with 90%+ coverage
**Intelligence**: Learns from production incident history to prioritize edge cases

### Deploy Agent
**Input**: Code + tests
**Output**: Deployed services on dev stack with health checks
**Capabilities**: Docker builds, Kubernetes deployment, monitoring setup

---

## API Reference

### Integrations

- `GET /api/v1/integrations` - List all integrations
- `POST /api/v1/integrations` - Create new integration
- `GET /api/v1/integrations/:id` - Get integration details
- `GET /api/v1/integrations/:id/status` - Get full pipeline status
- `GET /api/v1/integrations/:id/timeline` - Get stage timeline

### BRDs

- `POST /api/v1/brds` - Upload BRD document
- `GET /api/v1/brds/:id` - Get BRD details
- `POST /api/v1/brds/:id/validate` - Trigger validation
- `POST /api/v1/brds/:id/approve` - Approve BRD
- `GET /api/v1/brds/:id/gap-analysis` - Get gap analysis report

### PRDs

- `GET /api/v1/prds/:id` - Get PRD
- `PUT /api/v1/prds/:id` - Update PRD
- `POST /api/v1/prds/:id/approve` - Approve PRD

### Code

- `GET /api/v1/code/:id` - Get generated code
- `GET /api/v1/code/:id/files` - List generated files
- `POST /api/v1/code/:id/approve` - Approve code

### Tests

- `GET /api/v1/tests/:id` - Get test suite
- `POST /api/v1/tests/:id/execute` - Execute tests
- `GET /api/v1/tests/:id/coverage` - Get coverage report

### Deployments

- `GET /api/v1/deployments/:id` - Get deployment details
- `POST /api/v1/deployments/:id/health-check` - Run health checks
- `POST /api/v1/deployments/:id/rollback` - Rollback deployment

---

## Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/gatewayforge

# API
PORT=8080
ENV=development

# AI
CLAUDE_API_KEY=your_claude_api_key

# Orchestration
TEMPORAL_URL=localhost:7233

# Storage
S3_BUCKET=gatewayforge-artifacts
S3_REGION=ap-south-1

# Monitoring
PROMETHEUS_URL=http://localhost:9090
GRAFANA_URL=http://localhost:3000
```

---

## Deployment

### Development

```bash
docker-compose up
```

### Production

```bash
# Build Docker images
docker build -t gatewayforge-api:latest -f backend/Dockerfile .
docker build -t gatewayforge-frontend:latest -f frontend/Dockerfile .

# Deploy to Kubernetes
kubectl apply -f k8s/
```

---

## Monitoring

- **Grafana Dashboards**: `http://localhost:3000`
  - Integration Pipeline Metrics
  - AI Skill Performance
  - API Latency & Throughput
  - Cost Analysis

- **Prometheus Metrics**: `http://localhost:9090`
  - `gatewayforge_integrations_total`
  - `gatewayforge_stage_duration_seconds`
  - `gatewayforge_ai_tokens_used`
  - `gatewayforge_test_coverage_percent`

---

## Security

- **Authentication**: OAuth 2.0 + JWT
- **Authorization**: Role-based access control (Solutions, Product, Engineering, Leadership)
- **Encryption**: AES-256 at rest, TLS 1.3 in transit
- **Secrets Management**: HashiCorp Vault
- **PCI-DSS Compliance**: All generated code follows PCI-DSS guidelines
- **Audit Trail**: Complete audit log of all actions

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

Copyright © 2026 Razorpay. All rights reserved.

---

## Support

- **Documentation**: `docs/`
- **Issues**: GitHub Issues
- **Email**: gatewayforge-support@razorpay.com
- **Slack**: #gatewayforge-ai

---

## Roadmap

### Phase 1: Foundation (Q1 FY27) ✅
- BRD Harmonizer skill
- Frontend MVP
- Knowledge base seeding

### Phase 2: PRD Automation (Q2 FY27)
- PRD Generator skill
- PRD review workflow
- Pilot with 3-4 integrations

### Phase 3: Code Generation (Q3 FY27)
- Coding Agent skill
- Test Agent skill
- Engineering approval workflow

### Phase 4: Full Automation (Q4 FY27)
- Deploy Agent skill
- Analytics dashboard
- Self-improvement loop
- Scale to all integrations

---

## Team

- **Product**: Naman (PM2, Payment Processing POD)
- **Engineering**: Platform Engineering Team
- **AI/ML**: AI Engineering Team
- **Stakeholders**: CPO, CTO, VP Engineering, Solutions Lead

---

**Built with ❤️ by Razorpay's Payment Processing POD**
