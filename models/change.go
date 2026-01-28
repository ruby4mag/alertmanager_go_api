package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Enums for ChangeType
const (
	ChangeTypeDeployment = "deployment"
	ChangeTypeConfig     = "config"
	ChangeTypeInfra      = "infra"
	ChangeTypeDb         = "db"
	ChangeTypeNetwork    = "network"
	ChangeTypeOther      = "other"
)

// Enums for ChangeStatus
const (
	ChangeStatusScheduled  = "scheduled"
	ChangeStatusInProgress = "in_progress"
	ChangeStatusCompleted  = "completed"
	ChangeStatusFailed     = "failed"
	ChangeStatusRolledBack = "rolled_back"
)

type DbChange struct {
	ID               primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	ChangeID         string                 `json:"change_id" bson:"change_id"`
	Source           string                 `json:"source" bson:"source"`
	ChangeType       string                 `json:"change_type" bson:"change_type"`
	Name             string                 `json:"name" bson:"name"`
	Description      string                 `json:"description" bson:"description"`
	Status           string                 `json:"status" bson:"status"`
	StartTime        time.Time              `json:"start_time" bson:"start_time"`
	EndTime          *time.Time             `json:"end_time,omitempty" bson:"end_time,omitempty"`
	ImplementedBy    string                 `json:"implemented_by" bson:"implemented_by"`
	AffectedEntities []string               `json:"affected_entities" bson:"affected_entities"`
	RawPayload       map[string]interface{} `json:"raw_payload" bson:"raw_payload"`
	RiskScore        int                    `json:"risk_score" bson:"risk_score"`
	ChangeRiskAnalysis map[string]interface{} `json:"change_risk_analysis" bson:"change_risk_analysis"`
	CreatedAt        time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" bson:"updated_at"`
}

func (c *DbChange) Validate() error {
	if c.ChangeID == "" {
		return errors.New("change_id is required")
	}
	if c.Source == "" {
		return errors.New("source is required")
	}
	if c.StartTime.IsZero() {
		return errors.New("start_time is required")
	}

	validTypes := map[string]bool{
		ChangeTypeDeployment: true,
		ChangeTypeConfig:     true,
		ChangeTypeInfra:      true,
		ChangeTypeDb:         true,
		ChangeTypeNetwork:    true,
		ChangeTypeOther:      true,
	}
	if !validTypes[c.ChangeType] {
		return errors.New("invalid change_type")
	}

	validStatuses := map[string]bool{
		ChangeStatusScheduled:  true,
		ChangeStatusInProgress: true,
		ChangeStatusCompleted:  true,
		ChangeStatusFailed:     true,
		ChangeStatusRolledBack: true,
	}
	if !validStatuses[c.Status] {
		return errors.New("invalid status")
	}

	return nil
}
