package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DashboardStats struct {
	OpenIncidents    int64   `json:"open_incidents"`
	CriticalActive   int64   `json:"critical_active"`
	AverageMTTR      float64 `json:"average_mttr_minutes"`
	SystemHealth     float64 `json:"system_health"`
	EventsProcessed  int64   `json:"events_processed_24h"`
}

type ServiceStat struct {
	Service  string `json:"service"`
	Critical int64  `json:"critical"`
	Warning  int64  `json:"warning"`
	Info     int64  `json:"info"`
	Total    int64  `json:"total"`
}

type TrendPoint struct {
	Timestamp string `json:"timestamp"`
	Count     int64  `json:"count"`
}

func DashboardStatsHandler(w http.ResponseWriter, r *http.Request, mongoClient *mongo.Client) {
	alertCollection := mongoClient.Database(mongodatabase).Collection(mongocollection)
	ctx := context.TODO()

	// 1. Open Incidents
	openFilter := bson.M{"alertstatus": "OPEN"}
	openCount, _ := alertCollection.CountDocuments(ctx, openFilter)

	// 2. Critical Active (P0 or P1)
	criticalFilter := bson.M{
		"alertstatus": "OPEN",
		"alertpriority": bson.M{"$in": bson.A{"P0", "P1"}},
	}
	criticalCount, _ := alertCollection.CountDocuments(ctx, criticalFilter)

	// 3. Average MTTR (Last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	mttrPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"alertstatus": "CLOSED",
			"alertcleartime.time": bson.M{"$gte": sevenDaysAgo},
		}}},
		{{Key: "$project", Value: bson.M{
			"duration": bson.M{"$divide": bson.A{
				bson.M{"$subtract": bson.A{"$alertcleartime.time", "$alertfirsttime.time"}},
				60000, // Convert to minutes
			}},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"avgMTTR": bson.M{"$avg": "$duration"},
		}}},
	}

	var mttrResult []bson.M
	var avgMTTR float64
	cursor, err := alertCollection.Aggregate(ctx, mttrPipeline)
	if err == nil {
		if cursor.All(ctx, &mttrResult); len(mttrResult) > 0 {
			if val, ok := mttrResult[0]["avgMTTR"].(float64); ok {
				avgMTTR = val
			}
		}
	}

	// 4. Events Processed (Last 24h)
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	eventCount, _ := alertCollection.CountDocuments(ctx, bson.M{
		"alertfirsttime.time": bson.M{"$gte": oneDayAgo},
	})

	stats := DashboardStats{
		OpenIncidents:   openCount,
		CriticalActive:  criticalCount,
		AverageMTTR:     avgMTTR,
		SystemHealth:    98.5, // Placeholder for logic
		EventsProcessed: eventCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int32:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	default:
		return 0
	}
}

func ServiceHeatmapHandler(w http.ResponseWriter, r *http.Request, mongoClient *mongo.Client) {
	alertCollection := mongoClient.Database(mongodatabase).Collection(mongocollection)
	ctx := context.TODO()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"alertstatus": "OPEN"}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$servicename",
			"critical": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$in": bson.A{"$alertpriority", bson.A{"P0", "P1"}}}, 1, 0}}},
			"warning":  bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$alertpriority", "P2"}}, 1, 0}}},
			"info":     bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$in": bson.A{"$alertpriority", bson.A{"P3", "P4"}}}, 1, 0}}},
			"total":    bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"total": -1}}},
		{{Key: "$limit", Value: 10}},
	}

	cursor, err := alertCollection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	heatmap := make([]ServiceStat, 0)
	for _, res := range results {
		heatmap = append(heatmap, ServiceStat{
			Service:  res["_id"].(string),
			Critical: toInt64(res["critical"]),
			Warning:  toInt64(res["warning"]),
			Info:     toInt64(res["info"]),
			Total:    toInt64(res["total"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heatmap)
}

func AlertTrendsHandler(w http.ResponseWriter, r *http.Request, mongoClient *mongo.Client) {
	alertCollection := mongoClient.Database(mongodatabase).Collection(mongocollection)
	ctx := context.TODO()

	// Last 24 hours, bucketed by hour
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"alertfirsttime.time": bson.M{"$gte": oneDayAgo}}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{"format": "%Y-%m-%d %H:00", "date": "$alertfirsttime.time"},
			},
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := alertCollection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trends := make([]TrendPoint, 0)
	for _, res := range results {
		trends = append(trends, TrendPoint{
			Timestamp: res["_id"].(string),
			Count:     toInt64(res["count"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trends)
}
