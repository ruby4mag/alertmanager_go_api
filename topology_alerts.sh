#!/bin/bash

# AlertManager API - Topology-Based Test Commands
# Generated from actual Neo4j topology data
# 
# Topology Structure: Datacenter -> Rack -> ESXi Host -> VM -> Application

BASE_URL="http://localhost:8081"

echo "=========================================="
echo "Topology-Based Alert Testing"
echo "=========================================="
echo ""
echo "This script creates alerts based on your actual Neo4j topology:"
echo "  - Datacenter level alerts"
echo "  - Rack level alerts"
echo "  - ESXi Host alerts"
echo "  - Virtual Machine alerts"
echo "  - Application alerts"
echo ""


# ==========================================
# 1. Application Alert - Incubate Ubiquitous E-Services
# ==========================================
echo "Creating alert: Application Alert - Incubate Ubiquitous E-Services"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "Incubate Ubiquitous E-Services",
  "alertTime": "2026-02-01 11:42:09",
  "alertSource": "APM",
  "serviceName": "ApplicationService",
  "alertSummary": "High Response Time Detected",
  "severity": "WARNING",
  "alertId": "alert-app-001",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for Incubate Ubiquitous E-Services",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-marriage-9029",
  "application": "Incubate Ubiquitous E-Services",
  "response_time_ms": "2500",
  "layer": "application"
}'

echo ""

# ==========================================
# 2. VM Alert - vm-marriage-9029
# ==========================================
echo "Creating alert: VM Alert - vm-marriage-9029"
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
  "alertNotes": "Alert generated for vm-marriage-9029",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-marriage-9029",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'

echo ""

# ==========================================
# 3. ESXi Host Alert - ESX-1-14-7
# ==========================================
echo "Creating alert: ESXi Host Alert - ESX-1-14-7"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "ESX-1-14-7",
  "alertTime": "2026-02-01 11:44:09",
  "alertSource": "vCenter",
  "serviceName": "HypervisorService",
  "alertSummary": "ESXi Host Memory Pressure",
  "severity": "WARNING",
  "alertId": "alert-esx-001",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for ESX-1-14-7",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "memory_pressure": "high",
  "vm_count": "45",
  "layer": "hypervisor"
}'

echo ""

# ==========================================
# 4. Rack Alert - RACK-1-14
# ==========================================
echo "Creating alert: Rack Alert - RACK-1-14"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "RACK-1-14",
  "alertTime": "2026-02-01 11:45:09",
  "alertSource": "DCIM",
  "serviceName": "DataCenterService",
  "alertSummary": "Rack Temperature Warning",
  "severity": "WARNING",
  "alertId": "alert-rack-001",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for RACK-1-14",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "temperature_celsius": "32",
  "threshold_celsius": "30",
  "layer": "infrastructure"
}'

echo ""

# ==========================================
# 5. Datacenter Alert - DC-1
# ==========================================
echo "Creating alert: Datacenter Alert - DC-1"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "DC-1",
  "alertTime": "2026-02-01 11:46:09",
  "alertSource": "DCIM",
  "serviceName": "DataCenterService",
  "alertSummary": "Datacenter Power Redundancy Lost",
  "severity": "CRITICAL",
  "alertId": "alert-dc-001",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for DC-1",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "power_feed": "A",
  "status": "offline",
  "redundancy": "lost",
  "layer": "datacenter"
}'

echo ""

# ==========================================
# 6. Application Alert - Seize Out-Of-The-Box E-Services
# ==========================================
echo "Creating alert: Application Alert - Seize Out-Of-The-Box E-Services"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "Seize Out-Of-The-Box E-Services",
  "alertTime": "2026-02-01 11:47:09",
  "alertSource": "APM",
  "serviceName": "ApplicationService",
  "alertSummary": "High Response Time Detected",
  "severity": "WARNING",
  "alertId": "alert-app-002",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for Seize Out-Of-The-Box E-Services",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-serve-2968",
  "application": "Seize Out-Of-The-Box E-Services",
  "response_time_ms": "2500",
  "layer": "application"
}'

echo ""

# ==========================================
# 7. VM Alert - vm-serve-2968
# ==========================================
echo "Creating alert: VM Alert - vm-serve-2968"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "vm-serve-2968",
  "alertTime": "2026-02-01 11:48:09",
  "alertSource": "vCenter",
  "serviceName": "VirtualizationService",
  "alertSummary": "High CPU Usage on Virtual Machine",
  "severity": "CRITICAL",
  "alertId": "alert-vm-002",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for vm-serve-2968",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-serve-2968",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'

echo ""

# ==========================================
# 8. ESXi Host Alert - ESX-1-14-7
# ==========================================
echo "Creating alert: ESXi Host Alert - ESX-1-14-7"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "ESX-1-14-7",
  "alertTime": "2026-02-01 11:49:09",
  "alertSource": "vCenter",
  "serviceName": "HypervisorService",
  "alertSummary": "ESXi Host Memory Pressure",
  "severity": "WARNING",
  "alertId": "alert-esx-002",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for ESX-1-14-7",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "memory_pressure": "high",
  "vm_count": "45",
  "layer": "hypervisor"
}'

echo ""

# ==========================================
# 9. Application Alert - Re-Intermediate Wireless E-Tailers
# ==========================================
echo "Creating alert: Application Alert - Re-Intermediate Wireless E-Tailers"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "Re-Intermediate Wireless E-Tailers",
  "alertTime": "2026-02-01 11:52:09",
  "alertSource": "APM",
  "serviceName": "ApplicationService",
  "alertSummary": "High Response Time Detected",
  "severity": "WARNING",
  "alertId": "alert-app-003",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for Re-Intermediate Wireless E-Tailers",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-character-4807",
  "application": "Re-Intermediate Wireless E-Tailers",
  "response_time_ms": "2500",
  "layer": "application"
}'

echo ""

# ==========================================
# 10. VM Alert - vm-character-4807
# ==========================================
echo "Creating alert: VM Alert - vm-character-4807"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "vm-character-4807",
  "alertTime": "2026-02-01 11:53:09",
  "alertSource": "vCenter",
  "serviceName": "VirtualizationService",
  "alertSummary": "High CPU Usage on Virtual Machine",
  "severity": "CRITICAL",
  "alertId": "alert-vm-003",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for vm-character-4807",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-character-4807",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'

echo ""

# ==========================================
# 11. ESXi Host Alert - ESX-1-14-7
# ==========================================
echo "Creating alert: ESXi Host Alert - ESX-1-14-7"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "ESX-1-14-7",
  "alertTime": "2026-02-01 11:54:09",
  "alertSource": "vCenter",
  "serviceName": "HypervisorService",
  "alertSummary": "ESXi Host Memory Pressure",
  "severity": "WARNING",
  "alertId": "alert-esx-003",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for ESX-1-14-7",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "memory_pressure": "high",
  "vm_count": "45",
  "layer": "hypervisor"
}'

echo ""

# ==========================================
# 12. Application Alert - Re-Contextualize Transparent E-Business
# ==========================================
echo "Creating alert: Application Alert - Re-Contextualize Transparent E-Business"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "Re-Contextualize Transparent E-Business",
  "alertTime": "2026-02-01 11:57:09",
  "alertSource": "APM",
  "serviceName": "ApplicationService",
  "alertSummary": "High Response Time Detected",
  "severity": "WARNING",
  "alertId": "alert-app-004",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for Re-Contextualize Transparent E-Business",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-surface-9665",
  "application": "Re-Contextualize Transparent E-Business",
  "response_time_ms": "2500",
  "layer": "application"
}'

echo ""

# ==========================================
# 13. VM Alert - vm-surface-9665
# ==========================================
echo "Creating alert: VM Alert - vm-surface-9665"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "vm-surface-9665",
  "alertTime": "2026-02-01 11:58:09",
  "alertSource": "vCenter",
  "serviceName": "VirtualizationService",
  "alertSummary": "High CPU Usage on Virtual Machine",
  "severity": "CRITICAL",
  "alertId": "alert-vm-004",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for vm-surface-9665",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-surface-9665",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'

echo ""

# ==========================================
# 14. ESXi Host Alert - ESX-1-14-7
# ==========================================
echo "Creating alert: ESXi Host Alert - ESX-1-14-7"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "ESX-1-14-7",
  "alertTime": "2026-02-01 11:59:09",
  "alertSource": "vCenter",
  "serviceName": "HypervisorService",
  "alertSummary": "ESXi Host Memory Pressure",
  "severity": "WARNING",
  "alertId": "alert-esx-004",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for ESX-1-14-7",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "memory_pressure": "high",
  "vm_count": "45",
  "layer": "hypervisor"
}'

echo ""

# ==========================================
# 15. Application Alert - Aggregate Integrated Portals
# ==========================================
echo "Creating alert: Application Alert - Aggregate Integrated Portals"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "Aggregate Integrated Portals",
  "alertTime": "2026-02-01 12:02:09",
  "alertSource": "APM",
  "serviceName": "ApplicationService",
  "alertSummary": "High Response Time Detected",
  "severity": "WARNING",
  "alertId": "alert-app-005",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for Aggregate Integrated Portals",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-ball-3412",
  "application": "Aggregate Integrated Portals",
  "response_time_ms": "2500",
  "layer": "application"
}'

echo ""

# ==========================================
# 16. VM Alert - vm-ball-3412
# ==========================================
echo "Creating alert: VM Alert - vm-ball-3412"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "vm-ball-3412",
  "alertTime": "2026-02-01 12:03:09",
  "alertSource": "vCenter",
  "serviceName": "VirtualizationService",
  "alertSummary": "High CPU Usage on Virtual Machine",
  "severity": "CRITICAL",
  "alertId": "alert-vm-005",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for vm-ball-3412",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "vm": "vm-ball-3412",
  "cpu_usage": "92%",
  "memory_usage": "85%",
  "layer": "virtualization"
}'

echo ""

# ==========================================
# 17. ESXi Host Alert - ESX-1-14-7
# ==========================================
echo "Creating alert: ESXi Host Alert - ESX-1-14-7"
curl -X POST "http://localhost:8081/" \
  -H "Content-Type: application/json" \
  -d '{
  "entity": "ESX-1-14-7",
  "alertTime": "2026-02-01 12:04:09",
  "alertSource": "vCenter",
  "serviceName": "HypervisorService",
  "alertSummary": "ESXi Host Memory Pressure",
  "severity": "WARNING",
  "alertId": "alert-esx-005",
  "alertType": "CREATE",
  "alertNotes": "Alert generated for ESX-1-14-7",
  "environment": "production",
  "region": "datacenter-1",
  "datacenter": "DC-1",
  "rack": "RACK-1-14",
  "host": "ESX-1-14-7",
  "memory_pressure": "high",
  "vm_count": "45",
  "layer": "hypervisor"
}'

echo ""

echo "=========================================="
echo "All topology-based alerts created!"
echo "=========================================="
echo ""
echo "Alert Summary:"
echo "  - Application layer alerts: Check APM metrics"
echo "  - VM layer alerts: Check vCenter for resource usage"
echo "  - ESXi Host alerts: Check hypervisor health"
echo "  - Rack alerts: Check DCIM for temperature/power"
echo "  - Datacenter alerts: Check facility infrastructure"
echo ""
echo "These alerts should trigger correlation rules based on topology relationships."
echo ""
