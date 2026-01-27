package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"alertmanager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChangeInfo struct {
	ChangeID         string     `json:"change_id"`
	Source           string     `json:"source"`
	ChangeType       string     `json:"change_type"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	StartTime        string     `json:"start_time"` // String for parsing
	EndTime          string     `json:"end_time"`   // String for parsing
	ImplementedBy    string     `json:"implemented_by"`
	AffectedEntities []string   `json:"affected_entities"`
	RawPayload       map[string]interface{} `json:"raw_payload"`
}

func ChangeHandler(w http.ResponseWriter, r *http.Request, mongoClient *mongo.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	changeCollection := mongoClient.Database(mongodatabase).Collection("changes")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var input ChangeInfo
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	
	// Create DbChange object from input
	dbChange := models.DbChange{
		ChangeID:         input.ChangeID,
		Source:           input.Source,
		ChangeType:       input.ChangeType,
		Name:             input.Name,
		Description:      input.Description,
		Status:           input.Status,
		ImplementedBy:    input.ImplementedBy,
		AffectedEntities: input.AffectedEntities,
		RawPayload:       input.RawPayload,
		UpdatedAt:        time.Now(),
	}

	// Parse StartTime
	if input.StartTime != "" {
		t, err := time.Parse(time.RFC3339, input.StartTime)
		if err != nil {
			// Try other formats if needed, or stick to RFC3339
			http.Error(w, "Invalid start_time format (use RFC3339)", http.StatusBadRequest)
			return
		}
		dbChange.StartTime = t
	}

	// Parse EndTime
	if input.EndTime != "" {
		t, err := time.Parse(time.RFC3339, input.EndTime)
		if err != nil {
			http.Error(w, "Invalid end_time format (use RFC3339)", http.StatusBadRequest)
			return
		}
		dbChange.EndTime = &t
	}

	// Validate
	if err := dbChange.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Idempotency Check
	filter := bson.M{"change_id": dbChange.ChangeID}
	
	// Check if exists
	var existingChange models.DbChange
	err = changeCollection.FindOne(context.TODO(), filter).Decode(&existingChange)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create new
			dbChange.ID = primitive.NewObjectID()
			dbChange.CreatedAt = time.Now()
			
			_, err := changeCollection.InsertOne(context.TODO(), dbChange)
			if err != nil {
				http.Error(w, "Error inserting change", http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(dbChange)
			return
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		// Update existing
		// We only update fields that are provided/modifyable. The prompt implies full update or status update.
		// "Allow status updates for existing change records"
		
		update := bson.M{
			"$set": bson.M{
				"source":            dbChange.Source,
				"change_type":       dbChange.ChangeType,
				"name":              dbChange.Name,
				"description":       dbChange.Description,
				"status":            dbChange.Status,
				"start_time":        dbChange.StartTime,
				"end_time":          dbChange.EndTime,
				"implemented_by":    dbChange.ImplementedBy,
				"affected_entities": dbChange.AffectedEntities,
				"raw_payload":       dbChange.RawPayload,
				"updated_at":        time.Now(),
			},
		}

		_, err := changeCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "Error updating change", http.StatusInternalServerError)
			return
		}

		// Retrieve updated document
		changeCollection.FindOne(context.TODO(), filter).Decode(&existingChange)
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingChange)
		return
	}
}
