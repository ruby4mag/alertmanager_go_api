package utilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CalculateRisk analyzes the risk for a given list of entities
func CalculateRisk(ctx context.Context, driver neo4j.DriverWithContext, entities []string) (int, map[string]interface{}, error) {
	if len(entities) == 0 {
		return 0, map[string]interface{}{"note": "No entities to analyze"}, nil
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	totalScore := 0
	analysisDetails := make(map[string]interface{})

	for _, entity := range entities {
		// Try to handle "type:name" format by extracting name, but checking both
		searchTerm := entity
		if strings.Contains(entity, ":") {
			parts := strings.SplitN(entity, ":", 2)
			if len(parts) == 2 {
				searchTerm = parts[1]
			}
		}

		// Query: Get Tier and Neighbor Count
		// We verify against 'name' or 'id' properties
		query := `
        MATCH (n)
        WHERE n.name = $name OR n.id = $name OR n.name = $full OR n.id = $full
        OPTIONAL MATCH (n)-[r]-(m)
        RETURN n.tier as tier, labels(n) as labels, count(m) as degree
        `

		result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, query, map[string]any{
				"name": searchTerm,
				"full": entity,
			})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				record := result.Record()

				// process tier
				tierRaw, _ := record.Get("tier")
				tierStr := ""
				if tierRaw != nil {
					tierStr = fmt.Sprintf("%v", tierRaw)
				}

				// process labels
				labelsRaw, _ := record.Get("labels")
				labels := []string{}
				if val, ok := labelsRaw.([]interface{}); ok {
					for _, v := range val {
						if s, ok := v.(string); ok {
							labels = append(labels, s)
						}
					}
				}

				// process degree
				degreeRaw, _ := record.Get("degree")
				degree := 0
				if val, ok := degreeRaw.(int64); ok {
					degree = int(val)
				}

				return map[string]interface{}{
					"tier":   tierStr,
					"labels": labels,
					"degree": degree,
				}, nil
			}

			return nil, nil // No match
		})

		if err != nil {
			fmt.Printf("Error querying Neo4j for %s: %v\n", entity, err)
			analysisDetails[entity] = map[string]string{"error": "Error querying topology"}
			continue
		}

		if result == nil {
			analysisDetails[entity] = map[string]string{"result": "Node not found in topology"}
			continue
		}

		data := result.(map[string]interface{})
		tier := data["tier"].(string)
		degree := data["degree"].(int)
		labels := data["labels"].([]string)

		// Scoring Logic
		// Tier 0/1: High Risk (50)
		// Tier 2: Medium (30)
		// Tier 3: Low (10)
		// Unknown: 10
		// Degree: +1 per dependency
		entityScore := 0

		switch strings.ToLower(tier) {
		case "tier 0", "tier 1", "1", "0":
			entityScore += 50
		case "tier 2", "2":
			entityScore += 30
		case "tier 3", "3":
			entityScore += 10
		default:
			entityScore += 10 // Default base for unknown tier
		}
		
		// Bonus for Critical Labels
		// If no Tier but label suggests importance?
		// For now just logging it.
		
		// Add weight for dependencies
		entityScore += degree * 1 

		totalScore += entityScore
		analysisDetails[entity] = map[string]interface{}{
			"tier":         tier,
			"labels":       labels,
			"dependencies": degree,
			"score":        entityScore,
		}
	}

	return totalScore, analysisDetails, nil
}
