# GatewayForge AI - Deployment Guide

Complete guide for deploying GatewayForge AI to production.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development](#local-development)
3. [Staging Deployment](#staging-deployment)
4. [Production Deployment](#production-deployment)
5. [Infrastructure Setup](#infrastructure-setup)
6. [Monitoring & Alerting](#monitoring--alerting)
7. [Backup & Recovery](#backup--recovery)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

- **Docker** 24.0+
- **Kubernetes** 1.28+
- **Helm** 3.12+
- **kubectl** 1.28+
- **Terraform** 1.6+ (for infrastructure)
- **AWS CLI** / **GCP CLI** (based on cloud provider)

### Required Access

- AWS/GCP account with admin access
- GitHub repository access
- Claude API key (Anthropic)
- Slack webhook URL (for notifications)
- HashiCorp Vault access (for secrets)

---

## Local Development

### Quick Start

```bash
# Clone repository
git clone https://github.com/razorpay/gatewayforge-ai.git
cd gatewayforge-ai

# Setup environment
make setup

# Start all services
make docker-up

# View logs
make docker-logs
```

### Service URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Temporal UI**: http://localhost:8088
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **MinIO**: http://localhost:9001 (minioadmin/minioadmin)

### Development Workflow

```bash
# Backend development
make dev-backend

# Frontend development
make dev-frontend

# Run tests
make test

# Lint code
make lint

# Format code
make format
```

---

## Staging Deployment

### 1. Infrastructure Setup

```bash
cd terraform/staging

# Initialize Terraform
terraform init

# Plan infrastructure
terraform plan -out=staging.tfplan

# Apply infrastructure
terraform apply staging.tfplan
```

**Creates:**
- EKS/GKE Kubernetes cluster
- RDS PostgreSQL database
- ElastiCache Redis cluster
- S3/GCS bucket for artifacts
- VPC with private/public subnets
- Load balancers
- CloudWatch/Stackdriver logging

### 2. Configure kubectl

```bash
# AWS EKS
aws eks update-kubeconfig --name gatewayforge-staging --region ap-south-1

# GCP GKE
gcloud container clusters get-credentials gatewayforge-staging --region=asia-south1
```

### 3. Deploy with Helm

```bash
# Add Helm repository
helm repo add gatewayforge https://charts.gatewayforge.razorpay.com
helm repo update

# Create namespace
kubectl create namespace gatewayforge-staging

# Deploy
helm upgrade --install gatewayforge \
  gatewayforge/gatewayforge-ai \
  --namespace gatewayforge-staging \
  --values k8s/staging/values.yaml \
  --set image.tag=v1.0.0 \
  --set database.host=staging-db.c1a2b3c4d5e6.ap-south-1.rds.amazonaws.com \
  --set secrets.claudeApiKey=$CLAUDE_API_KEY
```

### 4. Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n gatewayforge-staging

# Check services
kubectl get svc -n gatewayforge-staging

# Check ingress
kubectl get ingress -n gatewayforge-staging

# View logs
kubectl logs -f deployment/gatewayforge-api -n gatewayforge-staging
```

### 5. Run Smoke Tests

```bash
# Health check
curl https://staging.gatewayforge.razorpay.com/api/v1/health

# Create test integration
curl -X POST https://staging.gatewayforge.razorpay.com/api/v1/integrations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Integration",
    "partner_name": "Test Partner",
    "integration_type": "gateway",
    "payment_methods": ["card"],
    "geographies": ["India"],
    "created_by": "test@razorpay.com"
  }'
```

---

## Production Deployment

### Pre-Deployment Checklist

- [ ] All staging tests passed
- [ ] Security scan completed (no high/critical vulnerabilities)
- [ ] Performance testing completed (load test with 1000 TPS)
- [ ] Database migrations reviewed
- [ ] Rollback plan prepared
- [ ] Change request approved
- [ ] Oncall engineer notified
- [ ] Monitoring dashboards configured
- [ ] Alerts configured
- [ ] Backup verified

### 1. Database Migration

```bash
# Create database backup
kubectl exec -it postgres-primary-0 -n gatewayforge-prod -- \
  pg_dump -U admin gatewayforge > backup-$(date +%Y%m%d-%H%M%S).sql

# Upload backup to S3
aws s3 cp backup-*.sql s3://gatewayforge-backups/prod/

# Run migration (dry-run first)
kubectl apply -f k8s/production/migrations/20260306-migration.yaml --dry-run=client

# Apply migration
kubectl apply -f k8s/production/migrations/20260306-migration.yaml

# Verify migration
kubectl logs -f job/migration-20260306 -n gatewayforge-prod
```

### 2. Deploy Backend

```bash
# Deploy API
helm upgrade --install gatewayforge \
  gatewayforge/gatewayforge-ai \
  --namespace gatewayforge-prod \
  --values k8s/production/values.yaml \
  --set image.tag=v1.0.0 \
  --set replicaCount=5 \
  --set resources.limits.cpu=2000m \
  --set resources.limits.memory=4Gi \
  --atomic \
  --timeout 10m

# Wait for rollout
kubectl rollout status deployment/gatewayforge-api -n gatewayforge-prod
```

### 3. Deploy Frontend

```bash
# Build and push Docker image
docker build -t gatewayforge-frontend:v1.0.0 -f frontend/Dockerfile .
docker push registry.razorpay.com/gatewayforge-frontend:v1.0.0

# Deploy
kubectl set image deployment/gatewayforge-frontend \
  frontend=registry.razorpay.com/gatewayforge-frontend:v1.0.0 \
  -n gatewayforge-prod

# Wait for rollout
kubectl rollout status deployment/gatewayforge-frontend -n gatewayforge-prod
```

### 4. Verify Production Deployment

```bash
# Health check
curl https://gatewayforge.razorpay.com/api/v1/health

# Check metrics
curl https://gatewayforge.razorpay.com/metrics

# Check all pods healthy
kubectl get pods -n gatewayforge-prod | grep -v Running

# Monitor logs
kubectl logs -f -l app=gatewayforge-api -n gatewayforge-prod --tail=100
```

### 5. Gradual Rollout (Canary Deployment)

```bash
# Deploy canary (10% traffic)
kubectl apply -f k8s/production/canary-10.yaml

# Monitor canary metrics
# - Error rate < 0.1%
# - P95 latency < 500ms
# - No customer complaints

# Increase to 50%
kubectl apply -f k8s/production/canary-50.yaml

# Full rollout
kubectl apply -f k8s/production/canary-100.yaml
```

---

## Infrastructure Setup

### AWS Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         AWS Cloud                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  VPC (10.0.0.0/16)                                   │  │
│  │                                                       │  │
│  │  ┌─────────────────┐    ┌─────────────────┐          │  │
│  │  │  Public Subnet  │    │  Public Subnet  │          │  │
│  │  │  (AZ-A)         │    │  (AZ-B)         │          │  │
│  │  │                 │    │                 │          │  │
│  │  │  ┌───────────┐  │    │  ┌───────────┐  │          │  │
│  │  │  │    ALB    │  │    │  │    NAT    │  │          │  │
│  │  │  └───────────┘  │    │  └───────────┘  │          │  │
│  │  └─────────────────┘    └─────────────────┘          │  │
│  │                                                       │  │
│  │  ┌─────────────────┐    ┌─────────────────┐          │  │
│  │  │ Private Subnet  │    │ Private Subnet  │          │  │
│  │  │  (AZ-A)         │    │  (AZ-B)         │          │  │
│  │  │                 │    │                 │          │  │
│  │  │  ┌───────────┐  │    │  ┌───────────┐  │          │  │
│  │  │  │    EKS    │  │    │  │    EKS    │  │          │  │
│  │  │  │   Nodes   │  │    │  │   Nodes   │  │          │  │
│  │  │  └───────────┘  │    │  └───────────┘  │          │  │
│  │  │                 │    │                 │          │  │
│  │  │  ┌───────────┐  │    │  ┌───────────┐  │          │  │
│  │  │  │    RDS    │  │    │  │   Redis   │  │          │  │
│  │  │  └───────────┘  │    │  └───────────┘  │          │  │
│  │  └─────────────────┘    └─────────────────┘          │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  External Services                                   │  │
│  │  • S3 (Artifact Storage)                             │  │
│  │  • CloudWatch (Logging & Monitoring)                 │  │
│  │  • Secrets Manager (Secrets)                         │  │
│  │  • CloudFront (CDN)                                  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Terraform Modules

```hcl
# terraform/production/main.tf
module "vpc" {
  source = "../modules/vpc"
  environment = "production"
  cidr_block = "10.0.0.0/16"
}

module "eks" {
  source = "../modules/eks"
  cluster_name = "gatewayforge-prod"
  node_groups = {
    general = {
      instance_types = ["m5.2xlarge"]
      min_size = 3
      max_size = 10
      desired_size = 5
    }
  }
}

module "rds" {
  source = "../modules/rds"
  instance_class = "db.r5.2xlarge"
  allocated_storage = 500
  multi_az = true
  backup_retention_period = 30
}

module "redis" {
  source = "../modules/elasticache"
  node_type = "cache.r5.xlarge"
  num_cache_nodes = 3
}

module "s3" {
  source = "../modules/s3"
  bucket_name = "gatewayforge-artifacts-prod"
  versioning = true
  lifecycle_rules = {
    archive_after_90_days = true
  }
}
```

---

## Monitoring & Alerting

### Grafana Dashboards

1. **Integration Pipeline Dashboard**
   - Integrations in progress
   - Average TAT per stage
   - Success/failure rates
   - AI cost per integration

2. **API Performance Dashboard**
   - Request rate (RPS)
   - Error rate
   - P50/P95/P99 latency
   - Top endpoints

3. **Infrastructure Dashboard**
   - CPU/Memory usage
   - Database connections
   - Redis cache hit rate
   - Disk usage

### Prometheus Alerts

```yaml
# k8s/production/prometheus-alerts.yaml
groups:
- name: gatewayforge
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"

  - alert: HighLatency
    expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "P95 latency > 1s"

  - alert: DatabaseConnectionsHigh
    expr: pg_stat_activity_count > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Database connections > 80"

  - alert: PodCrashLooping
    expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Pod is crash looping"
```

### PagerDuty Integration

```bash
# Set up PagerDuty webhook
kubectl create secret generic pagerduty-webhook \
  --from-literal=url=$PAGERDUTY_WEBHOOK_URL \
  -n gatewayforge-prod

# Configure Alertmanager
kubectl apply -f k8s/production/alertmanager-config.yaml
```

---

## Backup & Recovery

### Database Backup

```bash
# Automated daily backup (via CronJob)
kubectl apply -f k8s/production/cronjobs/db-backup.yaml

# Manual backup
kubectl exec -it postgres-primary-0 -n gatewayforge-prod -- \
  pg_dump -U admin -Fc gatewayforge > backup.dump

# Upload to S3
aws s3 cp backup.dump s3://gatewayforge-backups/prod/manual/$(date +%Y%m%d-%H%M%S).dump
```

### Database Restore

```bash
# Download backup from S3
aws s3 cp s3://gatewayforge-backups/prod/20260306-100000.dump backup.dump

# Restore
kubectl exec -i postgres-primary-0 -n gatewayforge-prod -- \
  pg_restore -U admin -d gatewayforge < backup.dump
```

### Disaster Recovery

**RTO**: 1 hour
**RPO**: 15 minutes

```bash
# Failover to secondary region (us-west-2)
terraform apply -var="active_region=us-west-2" -auto-approve

# Update DNS to point to new region
aws route53 change-resource-record-sets \
  --hosted-zone-id Z1234567890ABC \
  --change-batch file://dns-failover.json

# Restore from latest backup
./scripts/restore-from-backup.sh us-west-2
```

---

## Troubleshooting

### Common Issues

#### 1. Pod Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name> -n gatewayforge-prod

# Check logs
kubectl logs <pod-name> -n gatewayforge-prod

# Common causes:
# - Image pull error → Check registry credentials
# - OOMKilled → Increase memory limits
# - CrashLoopBackOff → Check application logs
```

#### 2. Database Connection Errors

```bash
# Check database status
kubectl exec -it postgres-primary-0 -n gatewayforge-prod -- \
  psql -U admin -d gatewayforge -c "SELECT 1"

# Check connection pool
kubectl exec -it gatewayforge-api-xxx -n gatewayforge-prod -- \
  curl localhost:8080/debug/pprof/heap

# Solutions:
# - Increase max_connections in PostgreSQL
# - Tune connection pool settings
# - Add read replicas
```

#### 3. High Latency

```bash
# Check APM traces
# Identify slow queries
kubectl exec -it postgres-primary-0 -n gatewayforge-prod -- \
  psql -U admin -d gatewayforge -c "
    SELECT query, mean_time, calls
    FROM pg_stat_statements
    ORDER BY mean_time DESC
    LIMIT 10;"

# Solutions:
# - Add database indexes
# - Optimize queries
# - Add caching layer
# - Scale horizontally
```

### Rollback Procedure

```bash
# Rollback deployment
helm rollback gatewayforge -n gatewayforge-prod

# Rollback database migration
kubectl apply -f k8s/production/migrations/rollback-20260306.yaml

# Verify rollback
kubectl rollout history deployment/gatewayforge-api -n gatewayforge-prod
```

---

## Security

### SSL/TLS Certificates

```bash
# Using cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Create certificate
kubectl apply -f k8s/production/certificates.yaml
```

### Network Policies

```bash
# Restrict inter-pod communication
kubectl apply -f k8s/production/network-policies.yaml
```

### Secrets Management

```bash
# Use Vault for sensitive secrets
vault kv put secret/gatewayforge/prod/api \
  database_password=xxx \
  claude_api_key=xxx \
  jwt_secret=xxx
```

---

## Scaling

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: gatewayforge-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: gatewayforge-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Database Scaling

```bash
# Vertical scaling (increase instance size)
terraform apply -var="db_instance_class=db.r5.4xlarge"

# Horizontal scaling (add read replicas)
terraform apply -var="read_replica_count=3"
```

---

## Cost Optimization

### Resource Optimization

```bash
# Analyze resource usage
kubectl top pods -n gatewayforge-prod
kubectl top nodes

# Right-size resources
# Reduce CPU/memory limits based on actual usage
```

### Spot Instances

```hcl
# Use spot instances for non-critical workloads
node_groups = {
  spot = {
    instance_types = ["m5.2xlarge", "m5a.2xlarge"]
    capacity_type = "SPOT"
    min_size = 2
    max_size = 10
  }
}
```

---

**For support, contact:** devops@razorpay.com
