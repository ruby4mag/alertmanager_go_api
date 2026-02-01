#!/bin/bash

# AlertManager API Test Commands
# Base URL - adjust if needed
BASE_URL="http://localhost:8081"

echo "=================================="
echo "AlertManager API Test Commands"
echo "=================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ========================================
# 1. CREATE ALERT - Basic Alert
# ========================================
echo -e "${GREEN}1. Creating a basic alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 2. CREATE ALERT - With Additional Tags
# ========================================
echo -e "${GREEN}2. Creating alert with additional tags${NC}"
curl -X POST "${BASE_URL}/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "database-prod-01",
    "alertTime": "2026-02-01 08:35:00",
    "alertSource": "Nagios",
    "serviceName": "DatabaseService",
    "alertSummary": "Database Connection Pool Exhausted",
    "severity": "CRITICAL",
    "alertId": "alert-002",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.200",
    "alertNotes": "Connection pool at 100% capacity",
    "region": "us-west-2",
    "environment": "production",
    "team": "database",
    "application": "order-service",
    "cluster": "prod-cluster-01"
  }'
echo -e "\n"

# ========================================
# 3. CREATE ALERT - Different Severity
# ========================================
echo -e "${GREEN}3. Creating WARNING severity alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 4. CREATE ALERT - For Grouping Test (Same Tags)
# ========================================
echo -e "${GREEN}4. Creating first alert for grouping test${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 5. CREATE ALERT - Second Alert for Grouping (Similar Tags)
# ========================================
echo -e "${GREEN}5. Creating second alert for grouping test (should group with #4)${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 6. DUPLICATE ALERT - Should increment count
# ========================================
echo -e "${GREEN}6. Sending duplicate alert (should increment count)${NC}"
curl -X POST "${BASE_URL}/" \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "server-web-01",
    "alertTime": "2026-02-01 08:50:00",
    "alertSource": "Prometheus",
    "serviceName": "WebService",
    "alertSummary": "High CPU Usage",
    "severity": "CRITICAL",
    "alertId": "alert-001",
    "alertType": "CREATE",
    "ipAddress": "192.168.1.100",
    "alertNotes": "CPU usage still high"
  }'
echo -e "\n"

# ========================================
# 7. CLOSE ALERT
# ========================================
echo -e "${GREEN}7. Closing an alert${NC}"
curl -X POST "${BASE_URL}/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-003",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:00:00"
  }'
echo -e "\n"

# ========================================
# 8. CREATE ALERT - Network Issue
# ========================================
echo -e "${GREEN}8. Creating network-related alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 9. CREATE ALERT - Disk Space
# ========================================
echo -e "${GREEN}9. Creating disk space alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 10. CREATE ALERT - Application Error
# ========================================
echo -e "${GREEN}10. Creating application error alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 11. CREATE ALERT - INFO Severity
# ========================================
echo -e "${GREEN}11. Creating INFO severity alert${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 12. CREATE ALERT - With PagerDuty Fields
# ========================================
echo -e "${GREEN}12. Creating alert with PagerDuty fields${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

# ========================================
# 13. CLOSE ALERT - Close Grouped Alert Child
# ========================================
echo -e "${GREEN}13. Closing first grouped alert (child)${NC}"
curl -X POST "${BASE_URL}/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-group-001",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:30:00"
  }'
echo -e "\n"

# ========================================
# 14. CLOSE ALERT - Close Second Grouped Alert
# ========================================
echo -e "${GREEN}14. Closing second grouped alert (should close parent too)${NC}"
curl -X POST "${BASE_URL}/" \
  -H "Content-Type: application/json" \
  -d '{
    "alertId": "alert-group-002",
    "alertType": "CLOSE",
    "alertTime": "2026-02-01 09:35:00"
  }'
echo -e "\n"

# ========================================
# 15. CREATE ALERT - Minimal Required Fields
# ========================================
echo -e "${GREEN}15. Creating alert with minimal required fields${NC}"
curl -X POST "${BASE_URL}/" \
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
echo -e "\n"

echo -e "${BLUE}=================================="
echo "Test commands completed!"
echo "==================================${NC}"
echo ""
echo "Tips:"
echo "  - Check MongoDB to verify alerts were created"
echo "  - Monitor console output for grouping behavior"
echo "  - Test correlation rules by creating alerts with matching tags"
echo "  - Verify PagerDuty integration if configured"
echo ""
