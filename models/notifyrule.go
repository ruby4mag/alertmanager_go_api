package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbNotifyRule struct {
	ID 					primitive.ObjectID `bson:"_id,omitempty"`
	RuleName			string 				`bson:"rulename" json:"rulename"`
	RuleDescription 	string 				`bson:"ruledescription" json:"ruledescription"`
	RuleObject			string  			`bson:"ruleobject" json:"ruleobject"`
	Order				int  				`bson:"order" json:"order"`
	PayLoad				string				`bson:"payload" json:"payload"`
	EndPoint			string 				`bson:"endpoint" json:"endpoint"`
	PagerDutyService		string				`bson:"pagerduty_service,omitempty" json:"pagerduty_service,omitempty"`
	PagerDutyEscalationPolicy	string			`bson:"pagerduty_escalation_policy,omitempty" json:"pagerduty_escalation_policy,omitempty"`
}

