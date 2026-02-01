# AlertManager API Testing Guide

This directory contains comprehensive curl command examples and scripts for testing the AlertManager API.

## 📁 Available Test Files

### 1. **CURL_TEST_COMMANDS.md**
Comprehensive markdown documentation with categorized curl commands for various alert scenarios.

**Categories include:**
- Basic alerts (Critical, Warning, Info)
- Database alerts
- Network alerts
- Storage alerts
- Application alerts
- Kubernetes/Container alerts
- Security alerts
- Alert grouping tests
- Duplicate alert tests
- PagerDuty integration tests

**Usage:**
```bash
# Copy any command from the markdown file and run it directly
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{ ... }'
```

---

### 2. **test_curl_commands.sh**
Automated bash script with 15 pre-defined test scenarios.

**Features:**
- Creates alerts with various severities
- Tests deduplication
- Tests alert grouping
- Tests alert closure
- Color-coded output

**Usage:**
```bash
# Run all tests
./test_curl_commands.sh

# Or run individual commands by copying from the script
```

---

### 3. **topology_alerts.sh** ⭐ (Generated from Neo4j)
**Topology-aware alerts based on your actual Neo4j graph data!**

This script creates alerts following the infrastructure hierarchy:
```
Datacenter (DC-1)
    └── Rack (RACK-1-14)
        └── ESXi Host (ESX-1-14-7)
            └── Virtual Machine (vm-xxx-xxxx)
                └── Application (Service Name)
```

**Alert Types:**
1. **Application Layer** - High response times, errors
2. **VM Layer** - CPU/Memory usage
3. **ESXi Host Layer** - Memory pressure, resource contention
4. **Rack Layer** - Temperature, power issues
5. **Datacenter Layer** - Power redundancy, cooling

**Usage:**
```bash
# Run topology-based alerts
./topology_alerts.sh
```

**Why use this?**
- Tests correlation rules with real topology relationships
- Validates parent-child alert grouping
- Tests impact analysis across infrastructure layers
- Uses actual entity names from your Neo4j database

---

## 🔧 Generating Custom Topology Alerts

### Prerequisites
```bash
# Set Neo4j credentials
export NEO4J_URI=neo4j://192.168.1.201:7687
export NEO4J_USER=neo4j
export NEO4J_PASSWORD=kl8j2300
```

### Generate New Topology-Based Alerts
```bash
# Analyze topology and generate curl commands
python3 generate_topology_alerts.py

# This will:
# 1. Connect to Neo4j
# 2. Find DC -> Rack -> Host -> VM -> App paths
# 3. Generate topology_alerts.sh with real entity names
```

### Analyze Topology Structure
```bash
# Explore your Neo4j topology
python3 analyze_topology.py

# Output saved to: /tmp/neo4j_topology.json
```

---

## 📊 Alert Structure

### Required Fields
All CREATE alerts must include:
```json
{
  "entity": "server-name",              // Entity/host/service affected
  "alertTime": "2026-02-01 11:30:00",  // Format: YYYY-MM-DD HH:MM:SS
  "alertSource": "Prometheus",          // Source system
  "serviceName": "WebService",          // Service name
  "alertSummary": "High CPU Usage",     // Brief description
  "severity": "CRITICAL",               // CRITICAL, WARNING, INFO
  "alertId": "alert-001",               // Unique identifier
  "alertType": "CREATE"                 // CREATE or CLOSE
}
```

### Optional Fields
```json
{
  "ipAddress": "192.168.1.100",
  "alertNotes": "Additional context",
  "region": "us-east-1",
  "environment": "production",
  "team": "platform",
  // ... any custom tags for correlation/grouping
}
```

---

## 🎯 Testing Scenarios

### 1. Test Deduplication
```bash
# Send same alertId twice - should increment count
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"entity":"test","alertTime":"2026-02-01 12:00:00","alertSource":"Test",
       "serviceName":"Test","alertSummary":"Test","severity":"WARNING",
       "alertId":"dedup-test-001","alertType":"CREATE"}'

# Send again (same alertId)
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"entity":"test","alertTime":"2026-02-01 12:05:00","alertSource":"Test",
       "serviceName":"Test","alertSummary":"Test","severity":"WARNING",
       "alertId":"dedup-test-001","alertType":"CREATE"}'
```

### 2. Test Alert Grouping
Create alerts with matching tags within the grouping window:
```bash
# Alert 1
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"entity":"web-01","alertTime":"2026-02-01 12:00:00","alertSource":"Prometheus",
       "serviceName":"Web","alertSummary":"Issue 1","severity":"CRITICAL",
       "alertId":"group-001","alertType":"CREATE",
       "region":"us-east-1","environment":"production","team":"platform"}'

# Alert 2 (same tags, within time window)
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"entity":"web-02","alertTime":"2026-02-01 12:01:00","alertSource":"Prometheus",
       "serviceName":"Web","alertSummary":"Issue 2","severity":"CRITICAL",
       "alertId":"group-002","alertType":"CREATE",
       "region":"us-east-1","environment":"production","team":"platform"}'
```

### 3. Test Alert Closure
```bash
# Close an alert
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"alertId":"alert-001","alertType":"CLOSE","alertTime":"2026-02-01 13:00:00"}'
```

### 4. Test Parent-Child Auto-Closure
```bash
# Create grouped alerts, then close all children
# The parent should auto-close when all children are closed

# Close child 1
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"alertId":"group-001","alertType":"CLOSE","alertTime":"2026-02-01 13:00:00"}'

# Close child 2 (parent should auto-close)
curl -X POST "http://localhost:8081/" -H "Content-Type: application/json" \
  -d '{"alertId":"group-002","alertType":"CLOSE","alertTime":"2026-02-01 13:01:00"}'
```

---

## 🔍 Verifying Results

### Check MongoDB
```bash
mongosh
use spog_development
db.alerts.find().pretty()

# Find specific alert
db.alerts.find({"alertid": "alert-001"}).pretty()

# Find grouped alerts
db.alerts.find({"grouped": true}).pretty()

# Find parent alerts
db.alerts.find({"parent": true}).pretty()
```

### Monitor Server Logs
Watch the alertmanager console output to see:
- Alert processing
- Rule evaluation
- Grouping logic
- PagerDuty integration calls

---

## 🌐 Topology-Based Testing Benefits

Using `topology_alerts.sh` provides several advantages:

1. **Real Entity Names** - Uses actual DC, Rack, Host, VM, and App names from Neo4j
2. **Hierarchical Relationships** - Tests alerts across infrastructure layers
3. **Correlation Testing** - Validates that related alerts are properly correlated
4. **Impact Analysis** - Tests how alerts propagate through the topology
5. **Realistic Scenarios** - Mimics real-world alert patterns

### Example Topology Path
```
DC-1 (Datacenter)
  └── RACK-1-14 (Rack)
      └── ESX-1-14-7 (ESXi Host)
          └── vm-marriage-9029 (Virtual Machine)
              └── Incubate Ubiquitous E-Services (Application)
```

Alerts are created at each level with appropriate context tags, allowing you to test:
- How an application alert relates to its VM
- How a VM alert relates to its host
- How a host alert relates to its rack
- How rack issues impact the datacenter

---

## 📝 Quick Reference

| File | Purpose | When to Use |
|------|---------|-------------|
| `CURL_TEST_COMMANDS.md` | Copy-paste individual commands | Quick testing, specific scenarios |
| `test_curl_commands.sh` | Run predefined test suite | General API testing, CI/CD |
| `topology_alerts.sh` | Test with real topology data | Correlation testing, topology validation |
| `generate_topology_alerts.py` | Generate new topology tests | Update tests when topology changes |
| `analyze_topology.py` | Explore Neo4j structure | Understanding your topology |

---

## 🚀 Quick Start

```bash
# 1. Basic testing
./test_curl_commands.sh

# 2. Topology-aware testing (recommended)
export NEO4J_URI=neo4j://192.168.1.201:7687
export NEO4J_USER=neo4j
export NEO4J_PASSWORD=kl8j2300
./topology_alerts.sh

# 3. Verify in MongoDB
mongosh
use spog_development
db.alerts.find().count()
db.alerts.find({"parent": true}).pretty()
```

---

## 💡 Tips

1. **Start Simple** - Use `test_curl_commands.sh` first to verify basic functionality
2. **Use Topology Tests** - Run `topology_alerts.sh` to test correlation and grouping
3. **Monitor Logs** - Watch the server output to understand alert processing
4. **Check MongoDB** - Verify alerts are stored correctly with proper grouping
5. **Test Closure** - Always test closing alerts to verify parent-child logic
6. **Regenerate** - Run `generate_topology_alerts.py` when your Neo4j topology changes

---

## 🔗 Related Documentation

- `PAGERDUTY_ENDPOINTS.md` - PagerDuty integration details
- `PAGERDUTY_INTEGRATION.md` - PagerDuty setup guide
- `PAGERDUTY_CLOSURE.md` - Alert closure with PagerDuty

---

**Happy Testing! 🎉**
