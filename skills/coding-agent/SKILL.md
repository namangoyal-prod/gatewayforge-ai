# Coding Agent Skill

Generates production-ready code across all required repositories by analyzing reference integrations and adapting patterns to the new gateway's requirements. Powered by SWE Agent for deep codebase understanding and pattern extraction.

## Purpose

Transform approved PRDs into complete, tested, deployable code that follows Razorpay's:
- Architectural patterns
- Coding standards
- Security best practices
- Observability conventions
- Configuration patterns

## Input

```json
{
  "prd_id": "uuid",
  "prd_content": {...},
  "reference_integration": "Codec",
  "target_repositories": [
    "gateway-adapters",
    "routing-engine",
    "merchant-onboarding",
    "payment-processing-core",
    "settlement-service",
    "api-gateway",
    "monitoring-config"
  ]
}
```

## Output

```json
{
  "code_generation_id": "uuid",
  "repositories_affected": [...],
  "generated_files": {
    "gateway-adapters": [
      "adapters/hdfc_gateway.go",
      "adapters/hdfc_gateway_test.go",
      "transformers/hdfc_transformer.go"
    ],
    "routing-engine": [
      "rules/hdfc_routing.yaml",
      "config/hdfc_priority.json"
    ]
  },
  "pull_requests": {
    "gateway-adapters": "https://github.com/razorpay/gateway-adapters/pull/1234"
  },
  "security_scan": {
    "status": "PASSED",
    "vulnerabilities": []
  },
  "line_count": 15420,
  "test_coverage": "92%"
}
```

## Repository Coverage & Code Generation

### 1. Gateway Adapter Service
**Purpose:** Core integration logic

**Generated Files:**
```
gateway-adapters/
├── adapters/
│   ├── hdfc_gateway.go              # Main adapter implementing Gateway interface
│   ├── hdfc_gateway_test.go         # Unit tests
│   └── hdfc_iso8583_builder.go      # ISO 8583 message construction
├── transformers/
│   ├── hdfc_request_transformer.go   # Razorpay → Gateway format
│   └── hdfc_response_transformer.go  # Gateway → Razorpay format
└── config/
    └── hdfc_config.yaml              # Gateway-specific configuration
```

**Key Pattern:**
```go
// Learn from Codec implementation
type HDFCGateway struct {
    config     *Config
    httpClient *http.Client
    encryptor  crypto.Encryptor
    logger     *Logger
    metrics    *Metrics
}

func (g *HDFCGateway) Authorize(ctx context.Context, payment *Payment) (*AuthResponse, error) {
    // Pattern extracted from Codec:
    // 1. Build ISO 8583 message
    // 2. Encrypt sensitive fields
    // 3. Send to gateway
    // 4. Parse response
    // 5. Transform to internal format
    // 6. Emit metrics
    // 7. Log (PCI-compliant)

    span := g.startSpan(ctx, "HDFCGateway.Authorize")
    defer span.End()

    msg, err := g.buildAuthMessage(payment)
    if err != nil {
        g.metrics.IncrementError("build_message_failed")
        return nil, errors.Wrap(err, "failed to build auth message")
    }

    encrypted, err := g.encryptor.Encrypt(msg)
    if err != nil {
        return nil, errors.Wrap(err, "encryption failed")
    }

    resp, err := g.sendWithRetry(ctx, encrypted)
    if err != nil {
        return nil, err
    }

    return g.parseAuthResponse(resp)
}
```

### 2. Routing Engine
**Purpose:** Decide when to route payments to this gateway

**Generated Files:**
```
routing-engine/
├── rules/
│   └── hdfc_routing.yaml          # Routing rules
├── config/
│   └── hdfc_priority.json         # Priority configuration
└── handlers/
    └── hdfc_route_handler.go      # Custom routing logic
```

**Routing Rule Example:**
```yaml
name: hdfc_upi_routing
gateway: HDFC
enabled: true
priority: 10
conditions:
  payment_method: UPI
  merchant_category: [e-commerce, retail]
  amount_range:
    min: 100
    max: 100000
  geography: India
  time_window:
    start: "00:00"
    end: "23:59"
success_rate_threshold: 95.0
fallback_gateway: PhonePe
```

### 3. Merchant Onboarding Service
**Purpose:** Onboard merchants to use this gateway

**Generated Files:**
```
merchant-onboarding/
├── handlers/
│   └── hdfc_onboarding_handler.go
├── validators/
│   └── hdfc_kyc_validator.go
├── workflows/
│   └── hdfc_activation_workflow.go
└── templates/
    └── hdfc_credential_template.json
```

### 4. Payment Processing Core
**Purpose:** Transaction lifecycle management

**Generated Files:**
```
payment-processing-core/
├── handlers/
│   ├── hdfc_auth_handler.go
│   ├── hdfc_capture_handler.go
│   └── hdfc_void_handler.go
├── state_machines/
│   └── hdfc_transaction_sm.go
└── idempotency/
    └── hdfc_idempotency.go
```

### 5. Settlement Service
**Purpose:** Settlement & reconciliation

**Generated Files:**
```
settlement-service/
├── parsers/
│   └── hdfc_settlement_parser.go    # Parse settlement files
├── reconcilers/
│   └── hdfc_reconciler.go           # Match transactions
└── payouts/
    └── hdfc_payout_trigger.go       # Trigger merchant payouts
```

### 6. Callback/Webhook Handler
**Purpose:** Receive async notifications from gateway

**Generated Files:**
```
webhook-handler/
├── handlers/
│   └── hdfc_webhook_handler.go
├── validators/
│   └── hdfc_signature_validator.go
└── transformers/
    └── hdfc_webhook_transformer.go
```

### 7. API Gateway Config
**Purpose:** Expose merchant-facing APIs

**Generated Files:**
```
api-gateway/
├── routes/
│   └── hdfc_routes.yaml
├── middlewares/
│   └── hdfc_auth_middleware.go
└── validators/
    └── hdfc_request_validator.go
```

### 8. Monitoring & Alerts
**Purpose:** Observability

**Generated Files:**
```
monitoring-config/
├── grafana/
│   └── hdfc_dashboard.json
├── prometheus/
│   └── hdfc_alerts.yaml
└── metrics/
    └── hdfc_metrics.go
```

### 9. Config Service
**Purpose:** Centralized configuration

**Generated Files:**
```
config-service/
├── configs/
│   ├── hdfc_dev.yaml
│   ├── hdfc_staging.yaml
│   └── hdfc_prod.yaml
└── schemas/
    └── hdfc_config_schema.json
```

### 10. Database Migrations
**Purpose:** Schema changes

**Generated Files:**
```
migrations/
└── 20260306_hdfc_integration.sql
```

## Code Generation Intelligence

### Pattern Recognition

The Coding Agent uses SWE Agent to:

1. **Analyze Reference Integration:**
```bash
# SWE Agent explores Codec integration
swe-agent analyze \
  --repos gateway-adapters,routing-engine,payment-processing-core \
  --integration Codec \
  --extract-patterns
```

2. **Identify Common Patterns:**
- Interface implementations
- Error handling conventions
- Retry logic patterns
- Logging format
- Metric naming
- Configuration structure
- Test patterns

3. **Extract Conventions:**
```go
// Convention discovered:
// - All gateway adapters implement Gateway interface
// - All errors are wrapped with context
// - All external calls have timeout
// - All sensitive data is encrypted before logging
// - All metrics follow pattern: gateway_<name>_<operation>_<status>
```

### Security-First Generation

**Auto-Implemented Security:**

1. **PCI-DSS Compliant Logging:**
```go
// NEVER log sensitive fields
logger.Info("processing payment",
    zap.String("payment_id", payment.ID),
    zap.String("amount", payment.Amount),
    zap.String("currency", payment.Currency),
    // zap.String("card_number", payment.CardNumber), // ❌ BLOCKED
    zap.String("card_last4", payment.CardLast4),      // ✅ OK
)
```

2. **Encryption:**
```go
// Auto-generated encryption for sensitive fields
type SecurePayload struct {
    CardNumber   string `encrypt:"true"`   // Auto-encrypted
    CVV          string `encrypt:"true"`
    AccountNumber string `encrypt:"true"`
    PlainField   string `encrypt:"false"`
}
```

3. **Input Validation:**
```go
// Auto-generated validators
func ValidateAuthRequest(req *AuthRequest) error {
    if err := validateAmount(req.Amount); err != nil {
        return err
    }
    if err := validateCardNumber(req.CardNumber); err != nil {
        return err
    }
    // XSS prevention
    req.MerchantName = sanitize(req.MerchantName)
    // SQL injection prevention
    req.OrderID = validateAlphanumeric(req.OrderID)
    return nil
}
```

### Observability Built-In

**Every Generated Function:**

```go
func (g *HDFCGateway) ProcessPayment(ctx context.Context, p *Payment) (*Result, error) {
    // 1. Distributed tracing
    span := trace.SpanFromContext(ctx)
    span.SetAttributes(
        attribute.String("gateway", "HDFC"),
        attribute.String("payment_id", p.ID),
    )

    // 2. Structured logging
    logger := log.WithContext(ctx).With(
        zap.String("gateway", "HDFC"),
        zap.String("payment_id", p.ID),
    )
    logger.Info("processing payment")

    // 3. Metrics
    start := time.Now()
    defer func() {
        g.metrics.RecordLatency("hdfc.process_payment", time.Since(start))
    }()

    // 4. Error tracking
    result, err := g.doProcessPayment(ctx, p)
    if err != nil {
        g.metrics.IncrementError("hdfc.process_payment.error")
        logger.Error("payment processing failed", zap.Error(err))
        return nil, err
    }

    g.metrics.IncrementSuccess("hdfc.process_payment.success")
    return result, nil
}
```

### Configuration-Driven Design

**Gateway-specific behavior externalized:**

```yaml
# hdfc_config.yaml
gateway:
  name: HDFC
  type: acquiring_bank
  enabled: true

  endpoints:
    base_url: https://api.hdfc.com
    auth: /oauth/token
    authorize: /payments/authorize
    capture: /payments/capture

  timeouts:
    connect: 10s
    read: 30s
    write: 30s

  retry:
    max_attempts: 3
    backoff: exponential
    initial_interval: 1s

  encryption:
    algorithm: AES256
    key_rotation: 90d

  rate_limits:
    per_second: 100
    burst: 200
```

### Backward Compatibility

**Generated code follows existing contracts:**

```go
// Implements existing Gateway interface
type Gateway interface {
    Authorize(context.Context, *Payment) (*AuthResponse, error)
    Capture(context.Context, *CaptureRequest) (*CaptureResponse, error)
    Void(context.Context, *VoidRequest) (*VoidResponse, error)
    Refund(context.Context, *RefundRequest) (*RefundResponse, error)
    GetStatus(context.Context, string) (*StatusResponse, error)
}

// No changes to interface = no breaking changes upstream
```

## Integration with SWE Agent

```python
def generate_code_with_swe_agent(prd, reference_integration):
    """Use SWE Agent to generate code"""

    # Step 1: Analyze reference codebase
    swe_agent.analyze_codebase(
        repos=["gateway-adapters", "routing-engine", "payment-processing-core"],
        focus=reference_integration
    )

    # Step 2: Extract patterns
    patterns = swe_agent.extract_patterns(
        integration=reference_integration,
        aspects=["architecture", "error_handling", "testing", "config"]
    )

    # Step 3: Generate code for new integration
    generated_code = swe_agent.generate_integration(
        prd=prd,
        patterns=patterns,
        target_integration="HDFC",
        follow_conventions=True
    )

    # Step 4: Run security scan
    security_results = swe_agent.security_scan(generated_code)

    # Step 5: Generate tests
    tests = swe_agent.generate_tests(generated_code, coverage_target=90)

    return {
        "code": generated_code,
        "tests": tests,
        "security": security_results,
        "patterns_used": patterns
    }
```

## Quality Gates

Before code is submitted:

1. ✅ **Compiles**: All generated code compiles without errors
2. ✅ **Tests Pass**: All generated tests pass
3. ✅ **Coverage**: ≥ 85% line coverage, ≥ 90% branch coverage
4. ✅ **Security Scan**: Zero high/critical vulnerabilities
5. ✅ **Lint**: Passes golangci-lint
6. ✅ **Conventions**: Follows Razorpay code standards
7. ✅ **Dependencies**: No new vulnerable dependencies

## Metrics

- **Lines of Code Generated**: Avg 15,000-20,000 per integration
- **Files Generated**: Avg 50-80 files across 10 repos
- **Generation Time**: 30-45 minutes
- **First-Pass Success Rate**: % of generated code that passes all gates
- **Code Review Acceptance**: % of generated code accepted by engineers
