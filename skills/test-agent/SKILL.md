# Test Agent Skill

Generates comprehensive test suites covering happy paths, edge cases, error scenarios, and performance benchmarks. Ensures 90%+ coverage and catches edge cases that have caused production incidents in past integrations.

## Purpose

Automatically generate exhaustive test coverage that would take engineers weeks to write manually. Learn from historical incidents to prioritize edge cases that are most likely to cause production issues.

## Input

```json
{
  "code_generation_id": "uuid",
  "generated_code": {...},
  "prd": {...},
  "reference_tests": "Codec"
}
```

## Output

```json
{
  "test_suite_id": "uuid",
  "test_counts": {
    "unit": 342,
    "integration": 87,
    "e2e": 45,
    "edge_case": 76,
    "performance": 15,
    "security": 28
  },
  "coverage": {
    "line": "92.4%",
    "branch": "89.7%",
    "function": "95.2%"
  },
  "execution_results": {
    "passed": 588,
    "failed": 5,
    "skipped": 0,
    "duration": "4m 32s"
  },
  "generated_files": [...]
}
```

## Test Coverage Matrix

| Category | Count | Purpose | Example |
|----------|-------|---------|---------|
| **Unit Tests** | 200-400 | Individual functions, transformers, validators | Test ISO 8583 field mapping logic |
| **Integration Tests** | 50-100 | Service-to-service, DB operations, API mocking | Test gateway adapter → database flow |
| **E2E Tests** | 30-50 | Complete payment lifecycle | Test onboarding → transaction → settlement |
| **Edge Case Tests** | 40-80 | Timeout, failures, malformed data, partitions | Test duplicate webhook delivery |
| **Performance Tests** | 10-20 | Load, latency, throughput | Test 1000 TPS sustained load |
| **Security Tests** | 20-30 | Injection, auth bypass, PCI compliance | Test SQL injection in merchant ID |

## Test Generation Strategy

### 1. Unit Tests

**Auto-generated for every function:**

```go
// Generated unit test for HDFC transformer
func TestHDFCTransformer_TransformAuthRequest(t *testing.T) {
    tests := []struct {
        name    string
        input   *Payment
        want    *HDFCAuthRequest
        wantErr bool
    }{
        {
            name: "valid card payment",
            input: &Payment{
                ID:         "pay_123",
                Amount:     10000,
                Currency:   "INR",
                Method:     "card",
                CardNumber: "4111111111111111",
                CardExpiry: "12/25",
                CVV:        "123",
            },
            want: &HDFCAuthRequest{
                MTI:         "0200",
                PAN:         "4111111111111111",
                Amount:      "000000010000",
                Currency:    "356",
                // ... full expected output
            },
            wantErr: false,
        },
        {
            name: "invalid card number",
            input: &Payment{
                CardNumber: "invalid",
            },
            want:    nil,
            wantErr: true,
        },
        {
            name: "amount exceeds limit",
            input: &Payment{
                Amount: 10000000, // 1 CR
            },
            want:    nil,
            wantErr: true,
        },
        // ... 15-20 more test cases per function
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            transformer := NewHDFCTransformer()
            got, err := transformer.TransformAuthRequest(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 2. Integration Tests

**Test service interactions:**

```go
func TestHDFCGateway_AuthorizeIntegration(t *testing.T) {
    // Setup: Mock HDFC server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request
        assert.Equal(t, "/payments/authorize", r.URL.Path)
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

        // Return mock response
        json.NewEncoder(w).Encode(HDFCAuthResponse{
            ResponseCode: "00",
            ApprovalCode: "123456",
            TransactionID: "txn_abc",
        })
    }))
    defer mockServer.Close()

    // Setup: Gateway with mock server
    gateway := NewHDFCGateway(&Config{
        BaseURL: mockServer.URL,
    })

    // Test: Authorize payment
    resp, err := gateway.Authorize(context.Background(), &Payment{
        Amount:   10000,
        Currency: "INR",
        Method:   "card",
    })

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "00", resp.Code)
    assert.True(t, resp.Approved)

    // Verify: Database entry created
    txn, err := db.GetTransaction(resp.TransactionID)
    assert.NoError(t, err)
    assert.Equal(t, "AUTHORIZED", txn.Status)
}
```

### 3. End-to-End Tests

**Test complete flows:**

```go
func TestHDFCIntegration_E2E_CardPayment(t *testing.T) {
    // Step 1: Onboard merchant
    merchant := onboardMerchant(t, &MerchantRequest{
        Name:    "Test Merchant",
        Gateway: "HDFC",
    })
    assert.NotEmpty(t, merchant.ID)

    // Step 2: Create order
    order := createOrder(t, merchant.ID, 10000)
    assert.Equal(t, "CREATED", order.Status)

    // Step 3: Initiate payment
    payment := initiatePayment(t, order.ID, &PaymentRequest{
        Method:     "card",
        CardNumber: "4111111111111111",
        CardExpiry: "12/25",
        CVV:        "123",
    })
    assert.Equal(t, "PROCESSING", payment.Status)

    // Step 4: Wait for authorization
    time.Sleep(2 * time.Second)
    payment = getPayment(t, payment.ID)
    assert.Equal(t, "AUTHORIZED", payment.Status)

    // Step 5: Capture payment
    capture := capturePayment(t, payment.ID, 10000)
    assert.Equal(t, "CAPTURED", capture.Status)

    // Step 6: Wait for settlement (simulate)
    simulateSettlement(t, payment.ID)

    // Step 7: Verify settlement
    settlement := getSettlement(t, payment.ID)
    assert.Equal(t, "SETTLED", settlement.Status)
    assert.Equal(t, 10000, settlement.Amount)

    // Step 8: Verify reconciliation
    recon := getReconciliation(t, payment.ID)
    assert.True(t, recon.Matched)
}
```

### 4. Edge Case Tests

**Learn from production incidents:**

```go
// Historical incident: Duplicate webhooks caused double refunds
func TestHDFCWebhook_DuplicateDelivery(t *testing.T) {
    handler := NewWebhookHandler()

    webhook := &Webhook{
        Type:          "payment.captured",
        PaymentID:     "pay_123",
        TransactionID: "txn_abc",
        Timestamp:     time.Now(),
    }

    // First delivery - should process
    err := handler.Process(webhook)
    assert.NoError(t, err)

    payment := getPayment(t, "pay_123")
    assert.Equal(t, "CAPTURED", payment.Status)

    // Duplicate delivery - should be idempotent
    err = handler.Process(webhook)
    assert.NoError(t, err)

    payment = getPayment(t, "pay_123")
    assert.Equal(t, "CAPTURED", payment.Status) // Still CAPTURED, not DOUBLE_CAPTURED

    // Verify: Only one capture in DB
    captures := db.GetCaptures("pay_123")
    assert.Len(t, captures, 1)
}

// Historical incident: Gateway timeout caused payment stuck in PROCESSING
func TestHDFCGateway_TimeoutHandling(t *testing.T) {
    // Mock: Gateway that times out
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(35 * time.Second) // Exceeds 30s timeout
    }))
    defer mockServer.Close()

    gateway := NewHDFCGateway(&Config{
        BaseURL: mockServer.URL,
        Timeout: 30 * time.Second,
    })

    // Test: Should timeout and handle gracefully
    _, err := gateway.Authorize(context.Background(), &Payment{
        ID:     "pay_timeout",
        Amount: 10000,
    })

    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "timeout")

    // Verify: Payment status is FAILED, not PROCESSING
    payment := getPayment(t, "pay_timeout")
    assert.Equal(t, "FAILED", payment.Status)
    assert.Equal(t, "GATEWAY_TIMEOUT", payment.ErrorCode)

    // Verify: Retry job is queued
    retryJob := getRetryJob(t, "pay_timeout")
    assert.NotNil(t, retryJob)
    assert.Equal(t, 1, retryJob.Attempt)
}

// Historical incident: Malformed ISO 8583 response crashed parser
func TestHDFCParser_MalformedResponse(t *testing.T) {
    parser := NewHDFCResponseParser()

    malformedResponses := []string{
        "",                              // Empty response
        "invalid",                       // Non-ISO 8583
        "0210abcd",                      // Incomplete message
        "0210" + strings.Repeat("0", 200), // Truncated fields
    }

    for _, resp := range malformedResponses {
        t.Run(resp, func(t *testing.T) {
            // Should not panic, should return error
            _, err := parser.Parse(resp)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), "malformed")
        })
    }
}
```

### 5. Performance Tests

**Load testing with k6:**

```javascript
// k6 performance test
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 100 },   // Ramp up to 100 TPS
        { duration: '5m', target: 100 },   // Stay at 100 TPS
        { duration: '2m', target: 500 },   // Spike to 500 TPS
        { duration: '5m', target: 500 },   // Stay at 500 TPS
        { duration: '2m', target: 0 },     // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
        http_req_failed: ['rate<0.01'], // < 1% errors
    },
};

export default function() {
    let payload = JSON.stringify({
        amount: 10000,
        currency: 'INR',
        method: 'card',
        card_number: '4111111111111111',
    });

    let params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test_token',
        },
    };

    let res = http.post('http://localhost:8080/payments', payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
        'payment authorized': (r) => JSON.parse(r.body).status === 'authorized',
    });

    sleep(1);
}
```

### 6. Security Tests

```go
// SQL Injection test
func TestHDFCAPI_SQLInjection(t *testing.T) {
    maliciousInputs := []string{
        "1' OR '1'='1",
        "1; DROP TABLE payments; --",
        "1' UNION SELECT * FROM users --",
    }

    for _, input := range maliciousInputs {
        resp := makeRequest(t, &PaymentRequest{
            OrderID: input, // Inject into order_id
        })

        // Should reject, not execute SQL
        assert.Equal(t, 400, resp.StatusCode)
        assert.Contains(t, resp.Body, "invalid")

        // Verify: No data leaked
        payments := db.GetAllPayments()
        assert.Empty(t, payments) // Table should be empty
    }
}

// XSS test
func TestHDFCAPI_XSS(t *testing.T) {
    resp := makeRequest(t, &PaymentRequest{
        MerchantName: "<script>alert('xss')</script>",
    })

    assert.Equal(t, 200, resp.StatusCode)

    payment := getPayment(t, resp.PaymentID)
    // Should be sanitized
    assert.NotContains(t, payment.MerchantName, "<script>")
    assert.NotContains(t, payment.MerchantName, "alert")
}

// PCI-DSS compliance test
func TestHDFC_PCIDSSCompliance(t *testing.T) {
    // Test: Card number should never be logged
    logs := captureLogs(func() {
        gateway.Authorize(context.Background(), &Payment{
            CardNumber: "4111111111111111",
        })
    })

    for _, log := range logs {
        assert.NotContains(t, log, "4111111111111111")
        assert.NotContains(t, log, "CardNumber")
    }

    // Test: Card number should not be stored in DB
    payment := getPayment(t, "pay_123")
    assert.Empty(t, payment.CardNumber) // Should be empty or masked
    assert.Equal(t, "1111", payment.CardLast4) // Only last 4 digits
}
```

## Test Data Generation

**Auto-generate realistic test data:**

```go
// Test data generator
type TestDataGenerator struct {
    faker *gofakeit.Faker
}

func (g *TestDataGenerator) GeneratePayment() *Payment {
    return &Payment{
        ID:          "pay_" + g.faker.UUID(),
        Amount:      g.faker.IntRange(100, 100000),
        Currency:    "INR",
        Method:      g.faker.RandomString([]string{"card", "upi", "netbanking"}),
        CardNumber:  g.generateValidCardNumber(),
        CardExpiry:  g.generateFutureExpiry(),
        CVV:         g.faker.Numerify("###"),
        Email:       g.faker.Email(),
        Phone:       "+91" + g.faker.Numerify("##########"),
    }
}

func (g *TestDataGenerator) generateValidCardNumber() string {
    // Generate Luhn-valid card numbers
    bins := []string{"411111", "555555", "378282"} // Visa, MC, Amex
    bin := bins[g.faker.IntRange(0, len(bins)-1)]
    suffix := g.faker.Numerify("##########")
    return luhn.Generate(bin + suffix)
}
```

## Mock Server Generation

**Auto-generate mock gateway servers:**

```go
// Mock HDFC server for testing
func NewMockHDFCServer(behavior string) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch behavior {
        case "success":
            json.NewEncoder(w).Encode(HDFCResponse{
                ResponseCode: "00",
                Message:      "Approved",
            })

        case "failure":
            json.NewEncoder(w).Encode(HDFCResponse{
                ResponseCode: "05",
                Message:      "Do not honor",
            })

        case "timeout":
            time.Sleep(35 * time.Second)

        case "malformed":
            w.Write([]byte("invalid response"))

        case "intermittent":
            if rand.Float64() < 0.5 {
                json.NewEncoder(w).Encode(HDFCResponse{ResponseCode: "00"})
            } else {
                w.WriteHeader(500)
            }
        }
    }))
}
```

## Coverage Analysis

**Auto-generate coverage reports:**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Coverage breakdown by package
go tool cover -func=coverage.out | grep -v "100.0%"

# Uncovered lines
go-cover-treemap -coverprofile coverage.out > coverage.svg
```

## Metrics

- **Test Generation Time**: 15-20 minutes
- **Total Tests Generated**: 500-600 per integration
- **Target Coverage**: 90%+ line, 85%+ branch
- **Execution Time**: Full suite in < 5 minutes
