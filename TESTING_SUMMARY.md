# Summary: AlertManager Testing Resources

## 🎯 What Was Created

I've created a comprehensive testing suite for your AlertManager API with **real topology data from Neo4j**.

---

## 📦 Files Created

### 1. **TESTING_GUIDE.md** (Main Documentation)
Complete guide explaining all test files and how to use them.

### 2. **test_curl_commands.sh** (15 Generic Tests)
Automated script with predefined test scenarios:
- Basic alerts (CRITICAL, WARNING, INFO)
- Database, network, storage alerts
- Alert grouping and deduplication
- Alert closure and parent-child logic

**Usage:**
```bash
./test_curl_commands.sh
```

### 3. **topology_alerts.sh** ⭐ (17 Topology-Based Tests)
**THIS IS THE SPECIAL ONE!** Generated from your actual Neo4j topology.

Creates alerts following your real infrastructure:
```
DC-1 (Datacenter)
  └── RACK-1-14 (Rack)
      └── ESX-1-14-7 (ESXi Host)
          └── vm-marriage-9029, vm-serve-2968, etc. (VMs)
              └── Applications (Various services)
```

**Alert Types:**
- 🏢 Datacenter: Power redundancy issues
- 🗄️ Rack: Temperature warnings
- 💻 ESXi Host: Memory pressure
- 🖥️ VM: High CPU/Memory usage
- 📱 Application: High response times

**Usage:**
```bash
./topology_alerts.sh
```

### 4. **CURL_TEST_COMMANDS.md** (Documentation)
Markdown file with 50+ copy-paste curl commands organized by category:
- Basic alerts
- Database alerts
- Network alerts
- Storage alerts
- Application alerts
- Kubernetes alerts
- Security alerts
- Testing scenarios

### 5. **generate_topology_alerts.py** (Generator Script)
Python script that:
- Connects to your Neo4j database
- Finds DC → Rack → Host → VM → App paths
- Generates `topology_alerts.sh` with real entity names

**Regenerate when topology changes:**
```bash
export NEO4J_URI=neo4j://192.168.1.201:7687
export NEO4J_USER=neo4j
export NEO4J_PASSWORD=kl8j2300
python3 generate_topology_alerts.py
```

### 6. **analyze_topology.py** (Analysis Tool)
Explores your Neo4j topology structure and saves results to `/tmp/neo4j_topology.json`.

---

## 🐛 Bug Fix Applied

**Fixed panic when optional fields are missing:**
- Added `getStringOrEmpty()` helper function
- Fixed handling of `ipAddress` and `alertNotes` fields
- Application rebuilt successfully

**Before:** Panic when `ipAddress` was nil
**After:** Gracefully handles missing optional fields

---

## 🚀 Quick Start

### Option 1: Generic Testing
```bash
cd /opt/alertninja/alertmanager_go_api
./test_curl_commands.sh
```

### Option 2: Topology-Based Testing (Recommended!)
```bash
cd /opt/alertninja/alertmanager_go_api
./topology_alerts.sh
```

### Option 3: Individual Commands
```bash
# Copy commands from CURL_TEST_COMMANDS.md
cat CURL_TEST_COMMANDS.md
```

---

## 📊 What Makes Topology Testing Special

The `topology_alerts.sh` script uses **actual entity names** from your Neo4j database:

**Real Entities Found:**
- Datacenter: `DC-1`
- Rack: `RACK-1-14`
- ESXi Host: `ESX-1-14-7`
- VMs: `vm-marriage-9029`, `vm-serve-2968`, `vm-character-4807`, etc.
- Applications: `Incubate Ubiquitous E-Services`, `Seize Out-Of-The-Box E-Services`, etc.

**Benefits:**
✅ Tests correlation rules with real topology relationships
✅ Validates parent-child alert grouping
✅ Tests impact analysis across infrastructure layers
✅ Each alert includes full topology context (datacenter, rack, host, vm, app)

---

## 📝 Example Alert from Topology Script

```bash
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "vm-marriage-9029",
  "alertTime": "2026-02-01 11:43:09",
  "alertSource": "vCenter",
  "serviceName": "VirtualizationService",
  "alertSummary": "High CPU Usage on Virtual Machine",
  "severity": "CRITICAL",
  "alertId": "alert-vm-001",
  "alertType": "CREATE",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-marriage-9029",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'
```

Notice how it includes the full topology path as tags!

---

## 🔍 Verifying Results

```bash
# Check MongoDB
mongosh
use spog_development
db.alerts.find().pretty()

# Find grouped alerts
db.alerts.find({"grouped": true}).pretty()

# Find parent alerts
db.alerts.find({"parent": true}).pretty()
```

---

## 📚 All Files Summary

| File | Size | Purpose |
|------|------|---------|
| `TESTING_GUIDE.md` | 8.9K | Complete documentation |
| `test_curl_commands.sh` | 11K | 15 generic test scenarios |
| `topology_alerts.sh` | 16K | **17 topology-based tests** ⭐ |
| `CURL_TEST_COMMANDS.md` | 16K | 50+ copy-paste commands |
| `generate_topology_alerts.py` | 11K | Topology alert generator |
| `analyze_topology.py` | 6.8K | Topology analysis tool |

---

## 🎉 You're All Set!

You now have:
1. ✅ Generic test commands for basic API testing
2. ✅ **Topology-aware test commands using your real Neo4j data**
3. ✅ Tools to regenerate tests when topology changes
4. ✅ Complete documentation
5. ✅ Bug fix for optional fields

**Start testing with:**
```bash
./topology_alerts.sh
```

This will create 17 alerts across all layers of your infrastructure using real entity names from your Neo4j topology!
