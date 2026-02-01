# AlertManager API - cURL Test Commands

## Base Configuration
```bash
BASE_URL="http://localhost:8081"
```

---

## 1. Basic Alert Creation

### Critical Alert - High CPU
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-01",
    "alertTime": "2026-02-01 08:30:00",
    "alertSource": "Prometheus",
    "serviceName": "WebService",
    "alertSummary": "High CPU Usage",
    "severity": "CRITICAL",
    "alertId": "alert-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.100",
    "alertNotes": "CPU usage exceeded 90%",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform"
  }'
```

### Warning Alert - Memory Usage
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-app-03",
    "alertTime": "2026-02-01 08:40:00",
    "alertSource": "Zabbix",
    "serviceName": "ApplicationService",
    "alertSummary": "Memory Usage Above Threshold",
    "severity": "WARNING",
    "alertId": "alert-003",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.150",
    "alertNotes": "Memory usage at 75%",
    "region": "eu-west-1",
    "environment": "staging",
    "team": "backend"
  }'
```

### Info Alert - Deployment
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-01",
    "alertTime": "2026-02-01 09:20:00",
    "alertSource": "Monitoring",
    "serviceName": "WebService",
    "alertSummary": "Deployment Completed",
    "severity": "INFO",
    "alertId": "alert-info-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.100",
    "alertNotes": "Version 2.5.0 deployed successfully",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "version": "2.5.0",
    "deployment_id": "deploy-12345"
  }'
```

---

## 2. Database Alerts

### Database Connection Pool
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "database-prod-01",
    "alertTime": "2026-02-01 08:35:00",
    "alertSource": "Nagios",
    "serviceName": "DatabaseService",
    "alertSummary": "Database Connection Pool Exhausted",
    "severity": "CRITICAL",
    "alertId": "alert-db-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.200",
    "alertNotes": "Connection pool at 100% capacity",
    "region": "us-west-2",
    "environment": "production",
    "team": "database",
    "application": "order-service",
    "cluster": "prod-cluster-01"
  }'
```

### Database Slow Query
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "database-prod-01",
    "alertTime": "2026-02-01 10:00:00",
    "alertSource": "MySQL",
    "serviceName": "DatabaseService",
    "alertSummary": "Slow Query Detected",
    "severity": "WARNING",
    "alertId": "alert-db-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.200",
    "alertNotes": "Query execution time > 5 seconds",
    "region": "us-west-2",
    "environment": "production",
    "team": "database",
    "query_time": "7.5s",
    "table": "orders"
  }'
```

---

## 3. Network Alerts

### Interface Down
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "router-core-01",
    "alertTime": "2026-02-01 09:05:00",
    "alertSource": "SNMP",
    "serviceName": "NetworkService",
    "alertSummary": "Interface Down",
    "severity": "CRITICAL",
    "alertId": "alert-network-001",
    "alertType": "CREATE",
    "ipAddress": "10.0.0.1",
    "alertNotes": "GigabitEthernet0/1 is down",
    "region": "us-east-1",
    "environment": "production",
    "team": "network",
    "device_type": "router",
    "interface": "GigabitEthernet0/1"
  }'
```

### High Packet Loss
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "switch-access-05",
    "alertTime": "2026-02-01 10:15:00",
    "alertSource": "SNMP",
    "serviceName": "NetworkService",
    "alertSummary": "High Packet Loss",
    "severity": "WARNING",
    "alertId": "alert-network-002",
    "alertType": "CREATE",
    "ipAddress": "10.0.1.5",
    "alertNotes": "Packet loss at 8%",
    "region": "us-east-1",
    "environment": "production",
    "team": "network",
    "device_type": "switch",
    "packet_loss": "8%"
  }'
```

---

## 4. Storage Alerts

### Disk Space Critical
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-storage-01",
    "alertTime": "2026-02-01 09:10:00",
    "alertSource": "Nagios",
    "serviceName": "StorageService",
    "alertSummary": "Disk Space Critical",
    "severity": "CRITICAL",
    "alertId": "alert-disk-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.2.50",
    "alertNotes": "/ partition at 95% capacity",
    "region": "us-west-2",
    "environment": "production",
    "team": "infrastructure",
    "mount_point": "/",
    "usage_percent": "95"
  }'
```

### Disk I/O High
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-storage-01",
    "alertTime": "2026-02-01 10:30:00",
    "alertSource": "Prometheus",
    "serviceName": "StorageService",
    "alertSummary": "High Disk I/O Wait",
    "severity": "WARNING",
    "alertId": "alert-disk-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.2.50",
    "alertNotes": "I/O wait time exceeds 40%",
    "region": "us-west-2",
    "environment": "production",
    "team": "infrastructure",
    "iowait": "42%"
  }'
```

---

## 5. Application Alerts

### Payment Service Failure
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "app-payment-service",
    "alertTime": "2026-02-01 09:15:00",
    "alertSource": "AppLogs",
    "serviceName": "PaymentService",
    "alertSummary": "Payment Processing Failures",
    "severity": "CRITICAL",
    "alertId": "alert-app-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.3.100",
    "alertNotes": "Multiple payment transaction failures detected",
    "region": "us-east-1",
    "environment": "production",
    "team": "payments",
    "application": "payment-gateway",
    "error_rate": "15%",
    "affected_users": "250"
  }'
```

### API Rate Limit Exceeded
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "api-gateway-prod",
    "alertTime": "2026-02-01 11:00:00",
    "alertSource": "APIGateway",
    "serviceName": "APIService",
    "alertSummary": "Rate Limit Exceeded",
    "severity": "WARNING",
    "alertId": "alert-app-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.3.200",
    "alertNotes": "Client exceeded API rate limits",
    "region": "us-east-1",
    "environment": "production",
    "team": "api",
    "client_id": "client-12345",
    "requests_per_minute": "1500"
  }'
```

---

## 6. Testing Alert Grouping

### First Alert in Group
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-02",
    "alertTime": "2026-02-01 08:45:00",
    "alertSource": "Prometheus",
    "serviceName": "WebService",
    "alertSummary": "Service Unavailable",
    "severity": "CRITICAL",
    "alertId": "alert-group-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.101",
    "alertNotes": "HTTP 503 errors",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "application": "web-frontend"
  }'
```

### Second Alert in Group (Similar Tags)
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-03",
    "alertTime": "2026-02-01 08:46:00",
    "alertSource": "Prometheus",
    "serviceName": "WebService",
    "alertSummary": "Service Degraded",
    "severity": "CRITICAL",
    "alertId": "alert-group-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.102",
    "alertNotes": "Slow response times",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "application": "web-frontend"
  }'
```

### Third Alert in Group
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-04",
    "alertTime": "2026-02-01 08:47:00",
    "alertSource": "Prometheus",
    "serviceName": "WebService",
    "alertSummary": "High Latency",
    "severity": "CRITICAL",
    "alertId": "alert-group-003",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.103",
    "alertNotes": "Response time > 5s",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "application": "web-frontend"
  }'
```

---

## 7. Testing Duplicate Alerts

### Original Alert
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-test-01",
    "alertTime": "2026-02-01 12:00:00",
    "alertSource": "Prometheus",
    "serviceName": "TestService",
    "alertSummary": "Test Alert for Deduplication",
    "severity": "WARNING",
    "alertId": "alert-dedup-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.50",
    "alertNotes": "First occurrence"
  }'
```

### Duplicate Alert (Should increment count)
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-test-01",
    "alertTime": "2026-02-01 12:05:00",
    "alertSource": "Prometheus",
    "serviceName": "TestService",
    "alertSummary": "Test Alert for Deduplication",
    "severity": "WARNING",
    "alertId": "alert-dedup-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.50",
    "alertNotes": "Second occurrence"
  }'
```

---

## 8. Closing Alerts

### Close Single Alert
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-003",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:00:00"
  }'
```

### Close Grouped Alert (Child)
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-group-001",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:30:00"
  }'
```

### Close Another Grouped Alert (Should close parent if all children closed)
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-group-002",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:35:00"
  }'
```

---

## 9. PagerDuty Integration Testing

### Alert with PagerDuty Fields
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-critical-01",
    "alertTime": "2026-02-01 09:25:00",
    "alertSource": "Prometheus",
    "serviceName": "CriticalService",
    "alertSummary": "Service Down - Immediate Action Required",
    "severity": "CRITICAL",
    "alertId": "alert-pd-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.250",
    "alertNotes": "Primary service is completely down",
    "region": "us-east-1",
    "environment": "production",
    "team": "sre",
    "pagerduty_service": "PXXXXXX",
    "pagerduty_escalation_policy": "PXXXXXX",
    "priority": "P1"
  }'
```

---

## 10. Minimal Alert (Required Fields Only)

```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "test-server",
    "alertTime": "2026-02-01 09:40:00",
    "alertSource": "TestSource",
    "serviceName": "TestService",
    "alertSummary": "Test Alert",
    "severity": "WARNING",
    "alertId": "alert-minimal-001",
    "alertType": "CREATE"
  }'
```

---

## 11. Kubernetes/Container Alerts

### Pod Crash Loop
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "pod-payment-service-7d8f9",
    "alertTime": "2026-02-01 13:00:00",
    "alertSource": "Kubernetes",
    "serviceName": "PaymentService",
    "alertSummary": "Pod CrashLoopBackOff",
    "severity": "CRITICAL",
    "alertId": "alert-k8s-001",
    "alertType": "CREATE",
    "ipAddress": "10.244.1.15",
    "alertNotes": "Pod restarting continuously",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "namespace": "production",
    "pod": "payment-service-7d8f9",
    "container": "payment-app",
    "restart_count": "15"
  }'
```

### Node Not Ready
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "k8s-node-worker-03",
    "alertTime": "2026-02-01 13:15:00",
    "alertSource": "Kubernetes",
    "serviceName": "KubernetesCluster",
    "alertSummary": "Node Not Ready",
    "severity": "CRITICAL",
    "alertId": "alert-k8s-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.10.103",
    "alertNotes": "Node status changed to NotReady",
    "region": "us-east-1",
    "environment": "production",
    "team": "platform",
    "cluster": "prod-cluster-01",
    "node": "worker-03",
    "condition": "NotReady"
  }'
```

---

## 12. Security Alerts

### Failed Login Attempts
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "auth-service-prod",
    "alertTime": "2026-02-01 14:00:00",
    "alertSource": "SecurityMonitor",
    "serviceName": "AuthenticationService",
    "alertSummary": "Multiple Failed Login Attempts",
    "severity": "WARNING",
    "alertId": "alert-sec-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.5.100",
    "alertNotes": "50+ failed login attempts from single IP",
    "region": "us-east-1",
    "environment": "production",
    "team": "security",
    "source_ip": "203.0.113.45",
    "failed_attempts": "52",
    "username": "admin"
  }'
```

### SSL Certificate Expiring
```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "www.example.com",
    "alertTime": "2026-02-01 14:30:00",
    "alertSource": "CertMonitor",
    "serviceName": "WebService",
    "alertSummary": "SSL Certificate Expiring Soon",
    "severity": "WARNING",
    "alertId": "alert-sec-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.100",
    "alertNotes": "SSL certificate expires in 7 days",
    "region": "us-east-1",
    "environment": "production",
    "team": "security",
    "domain": "www.example.com",
    "expiry_date": "2026-02-08",
    "days_remaining": "7"
  }'
```

---

## Testing Tips

1. **Check MongoDB** after each command:
   ```bash
   mongosh
   use spog_development
   db.alerts.find().pretty()
   ```

2. **Monitor Server Logs** to see processing:
   ```bash
   # Watch the server output for grouping, rules processing, etc.
   ```

3. **Test Grouping** by sending alerts with same tags within the grouping window

4. **Test Deduplication** by sending the same alertId multiple times

5. **Test Parent-Child Closure** by closing all child alerts to see parent auto-close

6. **Verify PagerDuty Integration** if configured with actual PagerDuty service IDs

---

## Required Fields

All CREATE alerts must include:
- `entity` - The entity/host/service affected
- `alertTime` - Timestamp in format "2006-01-02 15:04:05"
- `alertSource` - Source system generating the alert
- `serviceName` - Service name
- `alertSummary` - Brief description
- `severity` - CRITICAL, WARNING, INFO, etc.
- `alertId` - Unique identifier for the alert
- `alertType` - "CREATE" or "CLOSE"

## Optional Fields

- `ipAddress` - IP address of affected entity
- `alertNotes` - Additional notes
- Any custom tags (region, environment, team, etc.) for correlation/grouping
