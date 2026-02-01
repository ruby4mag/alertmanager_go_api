#!/usr/bin/env python3
"""
Generate realistic curl commands based on Neo4j topology
"""
import os
from neo4j import GraphDatabase
import json
from datetime import datetime, timedelta

NEO4J_URI = os.getenv("NEO4J_URI", "neo4j://192.168.1.201:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD", "kl8j2300")

class CurlCommandGenerator:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))
        self.base_url = "http://localhost:8081"
        self.base_time = datetime.now()
    
    def close(self):
        self.driver.close()
    
    def get_sample_paths(self, limit=10):
        """Get sample DC -> Rack -> ESX -> VM -> App paths"""
        with self.driver.session() as session:
            query = """
            MATCH path = (dc)-[*4]->(app)
            WHERE dc.name STARTS WITH 'DC-'
              AND ANY(n IN nodes(path) WHERE n.name STARTS WITH 'RACK-')
              AND ANY(n IN nodes(path) WHERE n.name STARTS WITH 'ESX-')
              AND ANY(n IN nodes(path) WHERE n.name STARTS WITH 'vm-')
            WITH dc, path, nodes(path) as path_nodes
            WHERE size(path_nodes) = 5
            RETURN DISTINCT
                   path_nodes[0].name as datacenter,
                   path_nodes[1].name as rack,
                   path_nodes[2].name as host,
                   path_nodes[3].name as vm,
                   path_nodes[4].name as application
            LIMIT $limit
            """
            result = session.run(query, limit=limit)
            return [record.data() for record in result]
    
    def generate_alert_curl(self, entity, alert_summary, severity, alert_id, 
                           alert_source="Monitoring", service_name="Infrastructure",
                           additional_tags=None, time_offset_minutes=0):
        """Generate a curl command for creating an alert"""
        alert_time = (self.base_time + timedelta(minutes=time_offset_minutes)).strftime("%Y-%m-%d %H:%M:%S")
        
        payload = {
            "entity": entity,
            "alertTime": alert_time,
            "alertSource": alert_source,
            "serviceName": service_name,
            "alertSummary": alert_summary,
            "severity": severity,
            "alertId": alert_id,
            "alertType": "CREATE",
            "alertNotes": f"Alert generated for {entity}",
            "environment": "production",
            "region": "datacenter-1"
        }
        
        if additional_tags:
            payload.update(additional_tags)
        
        json_payload = json.dumps(payload, indent=2)
        
        curl_cmd = f'''curl -X POST "{self.base_url}/" \\
  -H "Content-Type: application/json" \\
  -d '{json_payload}'
'''
        return curl_cmd
    
    def generate_topology_alerts(self):
        """Generate alerts for different levels of the topology"""
        paths = self.get_sample_paths(5)
        
        if not paths:
            print("No paths found in topology!")
            return []
        
        commands = []
        alert_scenarios = []
        
        for idx, path in enumerate(paths):
            dc = path['datacenter']
            rack = path['rack']
            host = path['host']
            vm = path['vm']
            app = path['application']
            
            # Scenario 1: Application-level alert
            alert_scenarios.append({
                'title': f'Application Alert - {app}',
                'entity': app,
                'summary': 'High Response Time Detected',
                'severity': 'WARNING',
                'alert_id': f'alert-app-{idx+1:03d}',
                'source': 'APM',
                'service': 'ApplicationService',
                'tags': {
                    'datacenter': dc,
                    'rack': rack,
                    'host': host,
                    'vm': vm,
                    'application': app,
                    'response_time_ms': '2500',
                    'layer': 'application'
                },
                'offset': idx * 5
            })
            
            # Scenario 2: VM-level alert
            alert_scenarios.append({
                'title': f'VM Alert - {vm}',
                'entity': vm,
                'summary': 'High CPU Usage on Virtual Machine',
                'severity': 'CRITICAL',
                'alert_id': f'alert-vm-{idx+1:03d}',
                'source': 'vCenter',
                'service': 'VirtualizationService',
                'tags': {
                    'datacenter': dc,
                    'rack': rack,
                    'host': host,
                    'vm': vm,
                    'cpu_usage': '92%',
                    'memory_usage': '85%',
                    'layer': 'virtualization'
                },
                'offset': idx * 5 + 1
            })
            
            # Scenario 3: ESXi Host alert
            alert_scenarios.append({
                'title': f'ESXi Host Alert - {host}',
                'entity': host,
                'summary': 'ESXi Host Memory Pressure',
                'severity': 'WARNING',
                'alert_id': f'alert-esx-{idx+1:03d}',
                'source': 'vCenter',
                'service': 'HypervisorService',
                'tags': {
                    'datacenter': dc,
                    'rack': rack,
                    'host': host,
                    'memory_pressure': 'high',
                    'vm_count': '45',
                    'layer': 'hypervisor'
                },
                'offset': idx * 5 + 2
            })
            
            # Scenario 4: Rack-level alert (power/cooling)
            if idx == 0:  # Only for first path
                alert_scenarios.append({
                    'title': f'Rack Alert - {rack}',
                    'entity': rack,
                    'summary': 'Rack Temperature Warning',
                    'severity': 'WARNING',
                    'alert_id': f'alert-rack-{idx+1:03d}',
                    'source': 'DCIM',
                    'service': 'DataCenterService',
                    'tags': {
                        'datacenter': dc,
                        'rack': rack,
                        'temperature_celsius': '32',
                        'threshold_celsius': '30',
                        'layer': 'infrastructure'
                    },
                    'offset': idx * 5 + 3
                })
            
            # Scenario 5: Datacenter-level alert
            if idx == 0:  # Only for first path
                alert_scenarios.append({
                    'title': f'Datacenter Alert - {dc}',
                    'entity': dc,
                    'summary': 'Datacenter Power Redundancy Lost',
                    'severity': 'CRITICAL',
                    'alert_id': f'alert-dc-{idx+1:03d}',
                    'source': 'DCIM',
                    'service': 'DataCenterService',
                    'tags': {
                        'datacenter': dc,
                        'power_feed': 'A',
                        'status': 'offline',
                        'redundancy': 'lost',
                        'layer': 'datacenter'
                    },
                    'offset': idx * 5 + 4
                })
        
        return alert_scenarios
    
    def generate_script(self, output_file):
        """Generate a complete bash script with all curl commands"""
        scenarios = self.generate_topology_alerts()
        
        script_content = '''#!/bin/bash

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

'''
        
        for idx, scenario in enumerate(scenarios, 1):
            script_content += f'''
# ==========================================
# {idx}. {scenario['title']}
# ==========================================
echo "Creating alert: {scenario['title']}"
'''
            
            curl_cmd = self.generate_alert_curl(
                entity=scenario['entity'],
                alert_summary=scenario['summary'],
                severity=scenario['severity'],
                alert_id=scenario['alert_id'],
                alert_source=scenario['source'],
                service_name=scenario['service'],
                additional_tags=scenario['tags'],
                time_offset_minutes=scenario['offset']
            )
            
            script_content += curl_cmd
            script_content += '\necho ""\n'
        
        script_content += '''
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
'''
        
        with open(output_file, 'w') as f:
            f.write(script_content)
        
        print(f"Generated {len(scenarios)} alert scenarios")
        print(f"Script saved to: {output_file}")
        
        return scenarios

def main():
    print("Connecting to Neo4j to analyze topology...")
    generator = CurlCommandGenerator(NEO4J_URI, NEO4J_USER, NEO4J_PASSWORD)
    
    try:
        scenarios = generator.generate_script('/opt/alertninja/alertmanager_go_api/topology_alerts.sh')
        
        print("\n" + "="*80)
        print("Generated Alert Scenarios:")
        print("="*80)
        for scenario in scenarios:
            print(f"  - {scenario['title']}")
            print(f"    Entity: {scenario['entity']}")
            print(f"    Severity: {scenario['severity']}")
            print(f"    Layer: {scenario['tags'].get('layer', 'N/A')}")
            print()
        
    finally:
        generator.close()

if __name__ == "__main__":
    main()
