package main

import (
    "context"
    "fmt"
    "log"
    "strings"
    
    "alertmanager/models"
    "alertmanager/utilities"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mitchellh/mapstructure"
)

func processSimilarityRule(newAlert *models.DbAlert, rule models.DbAlertGroup, mongoClient *mongo.Client) bool {
    // 1. Check Scope
    var alertMap map[string]interface{}
    if err := mapstructure.Decode(newAlert, &alertMap); err != nil {
        fmt.Println("Error decoding alert to map:", err)
        return false
    }
    
    // Check Scope & Build Candidate Filter
    alertCollection := mongoClient.Database(mongodatabase).Collection("alerts")
    
    // Filter: Open alerts, not self
    filter := bson.M{
        "alertstatus": "OPEN", 
        "parent": true,
        "_id": bson.M{"$ne": newAlert.ID},
    }
    
    // For each Scope Tag, the candidate MUST have the SAME value as the incoming alert.
    scopeMatched := true
    for _, tagKey := range rule.ScopeTags {
        val, found := getMapValue(alertMap, tagKey)
        if !found {
            // Incoming alert doesn't have this tag, so it cannot be grouped based on this scope.
            // fmt.Printf("Alert missing scope tag '%s', skipping rule '%s'\n", tagKey, rule.GroupName)
            scopeMatched = false
            break
        }
        
        // Add exact match requirement to Mongo Query
        if isTopLevel(tagKey) {
            filter[strings.ToLower(tagKey)] = val
        } else {
            filter["additionaldetails." + tagKey] = val
        }
    }
    
    if !scopeMatched {
        return false
    }

    // Construct Group Identifier based on Scope Tags
    groupIdentifier := ""
    for _, tagKey := range rule.ScopeTags {
        val, _ := getMapValue(alertMap, tagKey) // We know it exists from loop above
        groupIdentifier = groupIdentifier + "--" + val
    }
    fmt.Printf(" constructed GroupIdentifier: %s\n", groupIdentifier)

    fmt.Printf("Similarity Rule '%s' Scope Valid. Searching candidates...\n", rule.GroupName)
    
    cursor, err := alertCollection.Find(context.TODO(), filter)
    if err != nil {
        fmt.Println("Error finding candidates:", err)
        return false
    }
    defer cursor.Close(context.TODO())

    var candidates []models.DbAlert
    if err = cursor.All(context.TODO(), &candidates); err != nil {
        fmt.Println("Error decoding candidates:", err)
        return false
    }
    
    // 3. Find Best Match
    var bestMatch *models.DbAlert
    var maxScore float64 = -1.0
    
    for i := range candidates {
        // Use pointer to avoid copying large struct
        candidate := &candidates[i]
        
        var candMap map[string]interface{}
        mapstructure.Decode(candidate, &candMap)
        
        score := calculateScore(alertMap, candMap, rule.Similarity.Fields)
        // fmt.Printf("Score with %s: %f\n", candidate.AlertId, score)
        
        if score > maxScore {
            maxScore = score
            bestMatch = candidate
        }
    }
    
    threshold := rule.Similarity.Threshold
    if threshold <= 0 { threshold = 0.8 } // Default
    
    if maxScore >= threshold && bestMatch != nil {
        fmt.Printf("Grouping with %s (Score: %f)\n", bestMatch.AlertId, maxScore)
        addToGroup(newAlert, bestMatch, alertCollection, groupIdentifier)
    } else {
        fmt.Printf("No match found (Max Score: %f). Creating new Incident.\n", maxScore)
        createSimilarityParent(newAlert, rule, alertCollection, groupIdentifier)
    }

    return true
}

func calculateScore(a, b map[string]interface{}, fields []string) float64 {
    if len(fields) == 0 { return 0.0 }
    
    total := 0.0
    validFields := 0
    for _, f := range fields {
        v1, ok1 := getMapValue(a, f)
        v2, ok2 := getMapValue(b, f)
        
        if ok1 && ok2 {
            total += utilities.ComputeSimilarity(v1, v2)
        }
        validFields++ 
    }
    if validFields == 0 { return 0.0 }
    return total / float64(validFields)
}

func getMapValue(m map[string]interface{}, key string) (string, bool) {
    target := strings.ToLower(key)
    
    // 1. Try Top Level (Case-insensitive search)
    for k, v := range m {
        if strings.ToLower(k) == target {
            return fmt.Sprintf("%v", v), true
        }
    }
    
    // 2. Try AdditionalDetails
    // We need to find the "AdditionalDetails" key in the map first (it might be "AdditionalDetails" or "additionaldetails" etc)
    var adMap map[string]interface{}
    foundAd := false
    
    for k, v := range m {
        if strings.ToLower(k) == "additionaldetails" {
            if val, ok := v.(map[string]interface{}); ok {
                adMap = val
                foundAd = true
                break
            }
        }
    }
    
    if foundAd {
        for k, v := range adMap {
            if strings.ToLower(k) == target {
                 return fmt.Sprintf("%v", v), true
            }
        }
    }
    
    return "", false
}

func isTopLevel(key string) bool {
    fields := []string{"entity", "alertsource", "servicename", "alertsummary", "severity", "alertid", "alertpriority", "alertstatus", "ipaddress"}
    k := strings.ToLower(key)
    for _, f := range fields {
        if k == f { return true }
    }
    return false
}

func addToGroup(child, parent *models.DbAlert, collection *mongo.Collection, groupIdentifier string) {
    // Add child to parent
    parentFilter := bson.M{"_id": parent.ID}
    parentUpdate := bson.D{
        {Key: "$push", Value: bson.D{{Key: "groupalerts", Value: child.ID}}},
        {Key: "$set", Value: bson.D{
            {Key: "parent", Value: true}, 
            {Key: "grouped", Value: false}, // Ensure parent remains visible
            {Key: "alertlasttime", Value: child.AlertFirstTime},
        }},
    }
    
    _, err := collection.UpdateOne(context.TODO(), parentFilter, parentUpdate)
    if err != nil {
        log.Println("Error updating parent group:", err)
    }
    
    // Update Child
    childFilter := bson.M{"_id": child.ID}
    childUpdate := bson.M{
        "$set": bson.M{
            "groupincidentid": parent.ID,
            "grouped": true,
            "groupidentifier": groupIdentifier,
        },
    }
    _, err = collection.UpdateOne(context.TODO(), childFilter, childUpdate)
    if err != nil {
        log.Println("Error updating child group info:", err)
    }

    // Update Parent Priority
    mongoClient := collection.Database().Client()
    UpdateParentPriority(parent.ID, mongoClient)
}

func createSimilarityParent(child *models.DbAlert, rule models.DbAlertGroup, collection *mongo.Collection, groupIdentifier string) {
    // Deep Copy using the helper in main.go
    // Note: main.go deepCopy is available since we are in package main
    copy := deepCopy(*child)
    copy.ID = primitive.ObjectID{}
    copy.Parent = true
    copy.Grouped = false // Parent incident is top-level, so Grouped must be false for UI to show it
    copy.GroupAlerts = []primitive.ObjectID{child.ID}
    copy.GroupIdentifier = groupIdentifier
    
    // Identifier
    // We add a random suffix or time to ensure uniqueness of AlertId if needed?
    // DbAlert.AlertId is string.
    copy.AlertId = fmt.Sprintf("group-%s-%s", rule.GroupName, child.AlertId)
    // Removed prefix to ensure better similarity matching for subsequent alerts
    // copy.AlertSummary = fmt.Sprintf("[%s Group] %s", rule.GroupName, child.AlertSummary) 
    copy.AlertSummary = child.AlertSummary
    
    insertResult, err := collection.InsertOne(context.TODO(), copy)
    if err != nil {
        log.Println("Error creating parent incident:", err)
        return
    }
    
    parentID := insertResult.InsertedID.(primitive.ObjectID)
    copy.ID = parentID
    
    // Update Child
    childFilter := bson.M{"_id": child.ID}
    childUpdate := bson.M{
        "$set": bson.M{
            "groupincidentid": parentID,
            "grouped": true,
            "groupidentifier": groupIdentifier,
        },
    }
    collection.UpdateOne(context.TODO(), childFilter, childUpdate)
    
    // CRITICAL: Process notify rules for the parent alert to create PagerDuty incident
    // Get mongoClient from the collection
    mongoClient := collection.Database().Client()
    fmt.Printf("🔔 Processing notify rules for newly created PARENT alert %s\n", copy.AlertId)
    processNotifyRules(&copy, mongoClient)
}

// UpdateParentPriority recalculates the priority of a parent alert based on its OPEN children
func UpdateParentPriority(parentID primitive.ObjectID, mongoClient *mongo.Client) {
	alertCollection := mongoClient.Database(mongodatabase).Collection("alerts")
	
	// Find the parent to get current GroupAlerts
	var parent models.DbAlert
	err := alertCollection.FindOne(context.TODO(), bson.M{"_id": parentID}).Decode(&parent)
	if err != nil {
		fmt.Println("Error finding parent for priority update:", err)
		return
	}

    // If no group alerts, nothing to do (or reset to default P5?)
    // But parent should have children.
    if len(parent.GroupAlerts) == 0 {
        return 
    }

	// Find all OPEN children
	filter := bson.M{
		"_id": bson.M{"$in": parent.GroupAlerts},
		"alertstatus": "OPEN",
	}

	cursor, err := alertCollection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error finding children for priority update:", err)
		return
	}
	defer cursor.Close(context.TODO())

	var children []models.DbAlert
	if err = cursor.All(context.TODO(), &children); err != nil {
		fmt.Println("Error decoding children:", err)
		return
	}

    if len(children) == 0 {
        // No open children? Parent might be closing soon or already closed.
        // We probably shouldn't receive this call if we are about to close it, 
        // but if we do, maybe leave it or set to lowest? 
        // Let's leave it to avoid race with closure logic.
        return
    }

	// Find highest priority (lowest integer value)
	highestPriority := 100
	for _, child := range children {
		p := utilities.PriorityToInt(child.AlertPriority)
		if p < highestPriority {
			highestPriority = p
		}
	}

	// Update Parent
    // Use P4 as default baseline if allowed, but here we strictly take from children.
    // If highestPriority is still 100 (unlikely if children > 0), assume P4.
    if highestPriority == 100 { highestPriority = 4 }

	newPriority := utilities.IntToPriority(highestPriority)
	
	if newPriority != parent.AlertPriority {
		fmt.Printf("Updating Parent %s Priority: %s -> %s\n", parent.AlertId, parent.AlertPriority, newPriority)
		update := bson.M{"$set": bson.M{"alertpriority": newPriority}}
		_, err := alertCollection.UpdateOne(context.TODO(), bson.M{"_id": parentID}, update)
		if err != nil {
			fmt.Println("Error updating parent priority:", err)
		}
	}
}

