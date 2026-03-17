# Deploy Agent Skill

Automates dev stack deployment, service provisioning, health checks, and smoke tests. Ensures all services are running correctly before marking integration as deployment-ready.

## Purpose

- Provision services on dev stack
- Resolve dependencies (databases, caches, queues)
- Run health checks across all services
- Execute smoke tests to validate end-to-end functionality
- Generate deployment artifacts for staging promotion
- Provide comprehensive deployment reports

## Input

```json
{
  "integration_id": "uuid",
  "code_generation_id": "uuid",
  "test_suite_id": "uuid",
  "services_to_deploy": [
    "gateway-adapter-hdfc",
    "routing-engine",
    "payment-processor",
    "settlement-service"
  ],
  "environment": "dev"
}
```

## Output

```json
{
  "deployment_id": "uuid",
  "environment": "dev",
  "status": "SUCCESS",
  "services_deployed": {
    "gateway-adapter-hdfc": {
      "version": "1.0.0",
      "replicas": 2,
      "health": "HEALTHY",
      "endpoints": ["http://dev-gateway-adapter-hdfc:8080"]
    },
    "routing-engine": {...},
    "payment-processor": {...},
    "settlement-service": {...}
  },
  "health_checks": {
    "liveness": "PASSED",
    "readiness": "PASSED",
    "integration": "PASSED"
  },
  "smoke_test_results": {
    "total": 15,
    "passed": 15,
    "failed": 0,
    "duration": "2m 14s"
  },
  "deployment_artifacts": {
    "docker_images": [...],
    "helm_charts": [...],
    "terraform_configs": [...]
  },
  "deployed_at": "2026-03-06T10:30:00Z"
}
```

## Deployment Workflow

### Phase 1: Pre-Deployment Validation

**1.1 Dependency Check**
```yaml
dependencies:
  databases:
    - postgres:14
    - redis:7
  message_queues:
    - kafka:3.4
  config_services:
    - consul:1.15
  observability:
    - prometheus:2.45
    - grafana:9.5
```

```bash
# Verify all dependencies are available
check_dependency postgres "SELECT 1"
check_dependency redis "PING"
check_dependency kafka "kafka-topics --list"
```

**1.2 Configuration Validation**
```bash
# Validate all config files
validate_config config/hdfc_dev.yaml
validate_config config/routing_dev.yaml

# Check secrets exist in Vault
vault kv get secret/hdfc/api_key
vault kv get secret/hdfc/encryption_key
```

**1.3 Resource Availability**
```bash
# Check cluster has sufficient resources
kubectl top nodes
kubectl describe quota dev-namespace
```

### Phase 2: Service Provisioning

**2.1 Build Docker Images**
```bash
# Build images for all modified services
docker build -t gateway-adapter-hdfc:1.0.0 ./gateway-adapters
docker build -t routing-engine:1.0.0 ./routing-engine
docker build -t payment-processor:1.0.0 ./payment-processing-core

# Push to registry
docker push registry.razorpay.com/gateway-adapter-hdfc:1.0.0
```

**2.2 Generate Kubernetes Manifests**
```yaml
# gateway-adapter-hdfc-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-adapter-hdfc
  namespace: dev
  labels:
    app: gateway-adapter-hdfc
    integration: hdfc
    generated-by: gatewayforge-ai
spec:
  replicas: 2
  selector:
    matchLabels:
      app: gateway-adapter-hdfc
  template:
    metadata:
      labels:
        app: gateway-adapter-hdfc
    spec:
      containers:
      - name: gateway-adapter
        image: registry.razorpay.com/gateway-adapter-hdfc:1.0.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: GATEWAY_NAME
          value: "HDFC"
        - name: ENV
          value: "dev"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: db-config
              key: host
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: hdfc-secrets
              key: api_key
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: gateway-adapter-hdfc
  namespace: dev
spec:
  selector:
    app: gateway-adapter-hdfc
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
```

**2.3 Deploy Services**
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/gateway-adapter-hdfc-deployment.yaml
kubectl apply -f k8s/routing-engine-deployment.yaml
kubectl apply -f k8s/payment-processor-deployment.yaml

# Wait for rollout
kubectl rollout status deployment/gateway-adapter-hdfc -n dev
kubectl rollout status deployment/routing-engine -n dev
```

**2.4 Database Migrations**
```bash
# Run migrations
psql -h dev-db -U admin -d razorpay < migrations/20260306_hdfc_integration.sql

# Verify migrations
psql -h dev-db -U admin -d razorpay -c "SELECT * FROM integrations WHERE name='HDFC';"
```

### Phase 3: Health Checks

**3.1 Liveness Checks**
```bash
# Check all pods are running
kubectl get pods -n dev -l integration=hdfc

# Expected output:
# NAME                                   READY   STATUS    RESTARTS   AGE
# gateway-adapter-hdfc-7d8f9c-abc12     1/1     Running   0          2m
# gateway-adapter-hdfc-7d8f9c-def34     1/1     Running   0          2m
```

**3.2 Readiness Checks**
```bash
# Test /health/ready endpoints
for service in gateway-adapter-hdfc routing-engine payment-processor; do
    curl -f http://$service:8080/health/ready || exit 1
done
```

**3.3 Integration Health Checks**
```go
// Integration health check
func CheckIntegrationHealth(integration string) error {
    checks := []HealthCheck{
        CheckDatabase(integration),
        CheckGatewayConnectivity(integration),
        CheckRoutingRules(integration),
        CheckConfigLoaded(integration),
        CheckSecretsAccessible(integration),
    }

    for _, check := range checks {
        if err := check.Run(); err != nil {
            return fmt.Errorf("health check failed: %s - %w", check.Name, err)
        }
    }

    return nil
}
```

**3.4 Dependency Health**
```bash
# Verify service can connect to dependencies
kubectl exec -it gateway-adapter-hdfc-7d8f9c-abc12 -- sh -c "
  psql -h postgres -U admin -c 'SELECT 1' &&
  redis-cli -h redis PING &&
  curl -f http://consul:8500/v1/health/state/passing
"
```

### Phase 4: Smoke Tests

**4.1 Basic Connectivity**
```bash
# Test: Can service respond to requests?
curl -X POST http://gateway-adapter-hdfc:8080/healthz
# Expected: 200 OK

# Test: Can service access configuration?
curl http://gateway-adapter-hdfc:8080/config
# Expected: {"gateway": "HDFC", "enabled": true}
```

**4.2 End-to-End Smoke Tests**
```go
func RunSmokeTests(integration string) error {
    tests := []SmokeTest{
        {
            Name: "Create test payment",
            Test: func() error {
                payment, err := createPayment(&PaymentRequest{
                    Amount:   100,
                    Currency: "INR",
                    Method:   "card",
                    Gateway:  integration,
                })
                if err != nil {
                    return err
                }
                if payment.Status != "CREATED" {
                    return fmt.Errorf("unexpected status: %s", payment.Status)
                }
                return nil
            },
        },
        {
            Name: "Test gateway selection",
            Test: func() error {
                gateway, err := selectGateway(&RoutingRequest{
                    Amount:        100,
                    PaymentMethod: "card",
                    Merchant:      "test_merchant",
                })
                if err != nil {
                    return err
                }
                if gateway != integration {
                    return fmt.Errorf("wrong gateway selected: %s", gateway)
                }
                return nil
            },
        },
        {
            Name: "Test authorization (mock)",
            Test: func() error {
                // Use mock mode for smoke test
                resp, err := authorizePayment("pay_test_123", true)
                if err != nil {
                    return err
                }
                if !resp.Approved {
                    return fmt.Errorf("auth not approved")
                }
                return nil
            },
        },
        {
            Name: "Test metrics emission",
            Test: func() error {
                // Verify metrics are being published
                metrics, err := scrapeMetrics("http://gateway-adapter-hdfc:9090/metrics")
                if err != nil {
                    return err
                }
                if !strings.Contains(metrics, "hdfc_") {
                    return fmt.Errorf("HDFC metrics not found")
                }
                return nil
            },
        },
    }

    for _, test := range tests {
        if err := test.Test(); err != nil {
            return fmt.Errorf("smoke test '%s' failed: %w", test.Name, err)
        }
    }

    return nil
}
```

### Phase 5: Monitoring Setup

**5.1 Prometheus Alerts**
```yaml
# hdfc-alerts.yaml
groups:
- name: hdfc_gateway
  interval: 30s
  rules:
  - alert: HDFCHighErrorRate
    expr: rate(hdfc_requests_failed_total[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "HDFC gateway error rate > 5%"

  - alert: HDFCHighLatency
    expr: histogram_quantile(0.95, hdfc_request_duration_seconds) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "HDFC P95 latency > 1s"

  - alert: HDFCServiceDown
    expr: up{job="gateway-adapter-hdfc"} == 0
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "HDFC gateway service is down"
```

**5.2 Grafana Dashboard**
```json
{
  "dashboard": {
    "title": "HDFC Gateway - Dev",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {"expr": "rate(hdfc_requests_total[5m])"}
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {"expr": "rate(hdfc_requests_failed_total[5m])"}
        ]
      },
      {
        "title": "P95 Latency",
        "targets": [
          {"expr": "histogram_quantile(0.95, hdfc_request_duration_seconds)"}
        ]
      },
      {
        "title": "Authorization Success Rate",
        "targets": [
          {"expr": "rate(hdfc_auth_success_total[5m]) / rate(hdfc_auth_total[5m])"}
        ]
      }
    ]
  }
}
```

**5.3 Log Aggregation**
```bash
# Configure Fluentd to collect logs
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-hdfc-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/gateway-adapter-hdfc*.log
      tag hdfc.gateway
      format json
    </source>

    <filter hdfc.**>
      @type record_transformer
      <record>
        integration hdfc
        environment dev
      </record>
    </filter>

    <match hdfc.**>
      @type elasticsearch
      host elasticsearch
      index_name hdfc-logs
    </match>
EOF
```

### Phase 6: Generate Deployment Artifacts

**6.1 Helm Chart**
```yaml
# Chart.yaml
apiVersion: v2
name: hdfc-gateway
version: 1.0.0
description: HDFC Gateway Integration
type: application

# values.yaml
replicaCount: 2
image:
  repository: registry.razorpay.com/gateway-adapter-hdfc
  tag: "1.0.0"
service:
  type: ClusterIP
  port: 8080
resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 500m
    memory: 512Mi
```

**6.2 Terraform Config**
```hcl
# terraform/hdfc-integration.tf
resource "kubernetes_deployment" "gateway_adapter_hdfc" {
  metadata {
    name      = "gateway-adapter-hdfc"
    namespace = var.namespace
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "gateway-adapter-hdfc"
      }
    }

    template {
      metadata {
        labels = {
          app         = "gateway-adapter-hdfc"
          integration = "hdfc"
        }
      }

      spec {
        container {
          name  = "gateway-adapter"
          image = "registry.razorpay.com/gateway-adapter-hdfc:${var.version}"

          port {
            container_port = 8080
          }

          resources {
            limits = {
              cpu    = "1000m"
              memory = "1Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }
      }
    }
  }
}
```

## Deployment Report

```markdown
# HDFC Integration - Dev Deployment Report

**Deployment ID:** dep_abc123
**Integration:** HDFC
**Environment:** dev
**Status:** ✅ SUCCESS
**Deployed At:** 2026-03-06 10:30:00 UTC
**Duration:** 8m 42s

---

## Services Deployed

| Service | Version | Replicas | Status | Endpoint |
|---------|---------|----------|--------|----------|
| gateway-adapter-hdfc | 1.0.0 | 2/2 | ✅ HEALTHY | http://gateway-adapter-hdfc:8080 |
| routing-engine | 1.0.0 | 2/2 | ✅ HEALTHY | http://routing-engine:8080 |
| payment-processor | 1.0.0 | 3/3 | ✅ HEALTHY | http://payment-processor:8080 |
| settlement-service | 1.0.0 | 2/2 | ✅ HEALTHY | http://settlement-service:8080 |

---

## Health Checks

| Check | Status | Details |
|-------|--------|---------|
| Liveness | ✅ PASSED | All pods running |
| Readiness | ✅ PASSED | All endpoints responding |
| Database | ✅ PASSED | Migrations applied successfully |
| Redis | ✅ PASSED | Cache connectivity OK |
| Kafka | ✅ PASSED | Message queue connectivity OK |

---

## Smoke Test Results

**Total:** 15 tests | **Passed:** 15 | **Failed:** 0 | **Duration:** 2m 14s

✅ Create test payment
✅ Test gateway selection
✅ Test authorization (mock)
✅ Test metrics emission
✅ Test logging
✅ Test configuration loading
✅ Test secret access
✅ Test database operations
✅ Test cache operations
✅ Test message publishing
... [5 more tests]

---

## Artifacts Generated

- Docker Images: 4
- Helm Charts: 1
- Terraform Configs: 1
- Kubernetes Manifests: 12

---

## Next Steps

1. **QA Validation:** Run full E2E test suite on dev stack
2. **Performance Testing:** Load test with 1000 TPS
3. **Security Scan:** Run penetration tests
4. **Staging Promotion:** Deploy to staging after QA sign-off

---

## Monitoring

- **Grafana Dashboard:** http://grafana.dev/d/hdfc-gateway
- **Prometheus Alerts:** Configured for error rate, latency, downtime
- **Logs:** Available in Elasticsearch index `hdfc-logs`

---

## Rollback Plan

If issues are found:
```bash
kubectl rollout undo deployment/gateway-adapter-hdfc -n dev
helm rollback hdfc-gateway 0
```
```

## Metrics

- **Avg Deployment Time:** 8-10 minutes
- **Success Rate:** > 95%
- **Rollback Time:** < 2 minutes
- **Health Check Coverage:** 100% of deployed services
