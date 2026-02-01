#!/usr/bin/env python3
"""
Script to connect to Neo4j and analyze topology for alert generation
"""
import os
from neo4j import GraphDatabase
import json

# Get Neo4j credentials from environment
NEO4J_URI = os.getenv("NEO4J_URI", "neo4j://localhost:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD", "password")

class TopologyAnalyzer:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))
    
    def close(self):
        self.driver.close()
    
    def find_datacenter_hierarchy(self):
        """Find Datacenter -> Rack -> Host -> VM -> Application path"""
        with self.driver.session() as session:
            # Query to find the hierarchy
            query = """
            MATCH path = (dc)-[:CONTAINS*]->(node)
            WHERE (dc.type = 'Datacenter' OR dc.name CONTAINS 'DC' OR dc.name CONTAINS 'datacenter')
            RETURN dc.name as datacenter, 
                   dc.id as dc_id,
                   [n in nodes(path) | {name: n.name, type: n.type, id: n.id}] as path_nodes
            LIMIT 50
            """
            result = session.run(query)
            return [record.data() for record in result]
    
    def find_specific_path(self):
        """Find specific DC -> Rack -> Host -> VM -> Application path"""
        with self.driver.session() as session:
            query = """
            MATCH (dc)-[:CONTAINS*1..2]->(rack)-[:CONTAINS*1..2]->(host)-[:CONTAINS|RUNS*1..2]->(vm)-[:RUNS|HOSTS*1..2]->(app)
            WHERE (dc.type =~ '(?i)datacenter' OR dc.name =~ '(?i).*dc.*')
              AND (rack.type =~ '(?i)rack' OR rack.name =~ '(?i).*rack.*')
              AND (host.type =~ '(?i).*(host|esxi|hypervisor).*' OR host.name =~ '(?i).*(esxi|host).*')
              AND (vm.type =~ '(?i).*(vm|virtual).*' OR vm.name =~ '(?i).*(vm|virtual).*')
              AND (app.type =~ '(?i).*(app|service).*' OR app.name =~ '(?i).*(app|service).*')
            RETURN dc.name as datacenter, dc.type as dc_type,
                   rack.name as rack, rack.type as rack_type,
                   host.name as host, host.type as host_type,
                   vm.name as vm, vm.type as vm_type,
                   app.name as application, app.type as app_type
            LIMIT 10
            """
            result = session.run(query)
            return [record.data() for record in result]
    
    def find_all_node_types(self):
        """Get all unique node types in the graph"""
        with self.driver.session() as session:
            query = """
            MATCH (n)
            RETURN DISTINCT n.type as type, count(*) as count
            ORDER BY count DESC
            LIMIT 50
            """
            result = session.run(query)
            return [record.data() for record in result]
    
    def find_nodes_by_type(self, node_type):
        """Find nodes of a specific type"""
        with self.driver.session() as session:
            query = """
            MATCH (n)
            WHERE n.type = $type OR n.type CONTAINS $type
            RETURN n.name as name, n.type as type, n.id as id
            LIMIT 20
            """
            result = session.run(query, type=node_type)
            return [record.data() for record in result]
    
    def find_sample_hierarchy(self):
        """Find any hierarchical path in the graph"""
        with self.driver.session() as session:
            query = """
            MATCH path = (root)-[*1..4]->(leaf)
            WHERE NOT ()-[]->(root)
            WITH root, leaf, path, length(path) as depth
            WHERE depth >= 3
            RETURN root.name as root_name, root.type as root_type,
                   [n in nodes(path) | {name: n.name, type: n.type}] as hierarchy,
                   depth
            ORDER BY depth DESC
            LIMIT 5
            """
            result = session.run(query)
            return [record.data() for record in result]

def main():
    print("Connecting to Neo4j...")
    print(f"URI: {NEO4J_URI}")
    print(f"User: {NEO4J_USER}")
    
    analyzer = TopologyAnalyzer(NEO4J_URI, NEO4J_USER, NEO4J_PASSWORD)
    
    try:
        print("\n" + "="*80)
        print("1. Finding all node types in the graph...")
        print("="*80)
        node_types = analyzer.find_all_node_types()
        for nt in node_types:
            print(f"  - {nt['type']}: {nt['count']} nodes")
        
        print("\n" + "="*80)
        print("2. Searching for DC -> Rack -> Host -> VM -> App hierarchy...")
        print("="*80)
        specific_paths = analyzer.find_specific_path()
        if specific_paths:
            for i, path in enumerate(specific_paths, 1):
                print(f"\nPath {i}:")
                print(f"  Datacenter: {path['datacenter']} ({path['dc_type']})")
                print(f"  Rack: {path['rack']} ({path['rack_type']})")
                print(f"  Host: {path['host']} ({path['host_type']})")
                print(f"  VM: {path['vm']} ({path['vm_type']})")
                print(f"  Application: {path['application']} ({path['app_type']})")
        else:
            print("  No exact DC->Rack->Host->VM->App paths found")
        
        print("\n" + "="*80)
        print("3. Finding sample hierarchical paths...")
        print("="*80)
        sample_paths = analyzer.find_sample_hierarchy()
        for i, path in enumerate(sample_paths, 1):
            print(f"\nHierarchy {i} (depth={path['depth']}):")
            print(f"  Root: {path['root_name']} ({path['root_type']})")
            print("  Path:")
            for node in path['hierarchy']:
                print(f"    -> {node['name']} ({node['type']})")
        
        print("\n" + "="*80)
        print("4. Finding Datacenter nodes...")
        print("="*80)
        dc_hierarchy = analyzer.find_datacenter_hierarchy()
        if dc_hierarchy:
            print(f"Found {len(dc_hierarchy)} datacenter-related paths")
            for i, item in enumerate(dc_hierarchy[:3], 1):
                print(f"\nDatacenter {i}: {item['datacenter']}")
                print(f"  Path nodes: {len(item['path_nodes'])} nodes")
                for node in item['path_nodes'][:5]:
                    print(f"    - {node['name']} ({node.get('type', 'N/A')})")
        
        # Save results to JSON for curl command generation
        results = {
            'node_types': node_types,
            'specific_paths': specific_paths,
            'sample_paths': sample_paths,
            'dc_hierarchy': dc_hierarchy[:5] if dc_hierarchy else []
        }
        
        with open('/tmp/neo4j_topology.json', 'w') as f:
            json.dump(results, f, indent=2)
        
        print("\n" + "="*80)
        print("Results saved to /tmp/neo4j_topology.json")
        print("="*80)
        
    finally:
        analyzer.close()
        print("\nConnection closed.")

if __name__ == "__main__":
    main()
