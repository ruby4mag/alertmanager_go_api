package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbAlertGroup struct {
	ID 					primitive.ObjectID 	`bson:"_id,omitempty"`
	GroupName			string 				`bson:"groupname" json:"groupname"`
	GroupTags 			[]string 			`bson:"grouptags" json:"grouptags"`
	GroupWindow			int  				`bson:"groupwindow" json:"groupwindow"`
	ScopeTags           []string            `bson:"scope_tags" json:"scope_tags"`
	CorrelationMode     string              `bson:"correlation_mode" json:"correlation_mode"`
	Similarity          SimilarityConfig    `bson:"similarity" json:"similarity"`
}

type SimilarityConfig struct {
	Fields    []string `bson:"fields" json:"fields"`
	Threshold float64  `bson:"threshold" json:"threshold"`
}