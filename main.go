package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"reflect"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"alertmanager/models"
	"alertmanager/ruleengine"
	"alertmanager/utilities"

	"golang.org/x/exp/maps"
)

// const mongouri = "mongodb://localhost:27017/api"
// const mongodatabase = "spog_development"
// const mongocollection = "alerts"

var mongouri = os.Getenv("MONGO_URI")
var mongodatabase = os.Getenv("MONGO_DB")
var mongocollection = os.Getenv("MONGO_COLLECTION")
var NoderedEndpoint = os.Getenv("NODERED_ENDPOINT")


// Wrapper type around models.CustomTime
type CustomTimeWrapper struct {
    models.CustomTime
}

func (ct *CustomTimeWrapper) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), `"`) // Remove quotes from the string
	if len(b) < 3 {
		ct.Time = time.Time{}
		return nil
	}
    t, err := time.Parse("2006-01-02 15:04:05", s)
    if err != nil {
        return err
    }
    ct.Time = t
    return nil
}

func check_flap() (string,error){
	var rulesGroup ruleengine.RulesGroup

	rule := `{
		"condition": "AND",
		"rules": [
		  {
			"id": "name",
			"field": "name",
			"type": "string",
			"input": "text",
			"operator": "contains",
			"value": "John"
		  }
		]
	  }`

	err := json.Unmarshal([]byte(rule), &rulesGroup)
	if err != nil {
		fmt.Println("Error in rule evaluation ", err)
	}

	data := map[string]interface{}{
		"name":      "John",
		"age":       30,
		"birthdate": time.Date(1990, 6, 12, 0, 0, 0, 0, time.UTC),
	}

	res := ruleengine.EvaluateRulesGroup(data, rulesGroup)

	fmt.Println("The result is ", res)

	return "OK", nil
}

func main() {
	
	fmt.Println("\n\x1b[32mStarting EA API Server.....\x1b[0m\n")
	fmt.Println("\x1b[32mStarting mongo connection.....\x1b[0m\n")

 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongouri).SetServerAPIOptions(serverAPI)
	// Create a new mongoClient and connect to the server
	mongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := mongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("\x1b[32mPinged your deployment. You successfully connected to MongoDB!\x1b[0m\n ")
	fmt.Println("\x1b[32mWaiting for alerts.....\x1b[0m\n")

	http.HandleFunc("/api/v1/changes", func(w http.ResponseWriter, r *http.Request) {
		ChangeHandler(w, r, mongoClient)
	})

	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, mongoClient  )
	 })
	http.ListenAndServe(":8081", nil)
}

func Handler(w http.ResponseWriter, r *http.Request, mongoClient *mongo.Client ) {

	alertCollection := mongoClient.Database(mongodatabase).Collection(mongocollection)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var apiAlertData utilities.ApiAlertData

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &apiAlertData); err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// check if any of the required fields are empty.
	if key, err := apiAlertData.IsEmpty("entity", "alertTime", "alertSource", "serviceName", "alertSummary", "severity", "alertId"); err != nil {
		fmt.Printf("The key '%s' is missing or empty: %s\n", key, err)
		return
	} else {
		fmt.Println("All keys are present and not empty")
	}

	if apiAlertData["alertType"] == "CREATE"{
		fmt.Println("This is a create event")
		//status , err := check_flap()
		//fmt.Println("Here", status , err)
		fmt.Println("The API alertId is ", apiAlertData["alertId"])

		// De-duplication Starts
		filter := bson.M{
			"alertid":     apiAlertData["alertId"].(string) ,
            "alertstatus": "OPEN",
		}

		existingEvent := models.DbAlert{}
		
		opts := options.FindOne()
		err1 := alertCollection.FindOne(context.TODO(), filter, opts).Decode(&existingEvent)
	
		if err1 != nil {
			if err1 == mongo.ErrNoDocuments {
				fmt.Println("No matching event found. Creating Alert....")
				layout := "2006-01-02 15:04:05"
				parsedTime, err := time.Parse(layout, apiAlertData["alertTime"].(string))
				if err != nil {
					fmt.Println("Error parsing time:", err)
					return
				}
				// Create a alert in DB
				
				newAlert := models.DbAlert{
					ID: 				primitive.ObjectID{},
					Entity:				apiAlertData["entity"].(string),
					AlertFirstTime:		models.CustomTime{Time: parsedTime},
					AlertLastTime:		models.CustomTime{},
					AlertClearTime:		models.CustomTime{},
					AlertSource:		apiAlertData["alertSource"].(string),
					ServiceName: 		apiAlertData["serviceName"].(string),
					AlertSummary:		apiAlertData["alertSummary"].(string),
					AlertStatus:		"OPEN",
					AlertNotes:			apiAlertData["alertNotes"].(string),
					AlertAcked:			"NO",
					Severity:			apiAlertData["severity"].(string),
					AlertId:			apiAlertData["alertId"].(string),
					AlertPriority:		"NORMAL",
					IpAddress:			apiAlertData["ipAddress"].(string),
					AlertCount:			1,
					AdditionalDetails:	make(map[string]interface{}),
					Grouped: 			false ,	
					Parent:				false,
				}


				// Add additional Tags
				//fmt.Println("The object before addTags is " , newAlert )
				addTags(apiAlertData, &newAlert)
				fmt.Println("The object after addTags is " , newAlert )
				processAlertRules( &newAlert , mongoClient)
				processTagRules( &newAlert , mongoClient)

				insertResult , inserterr := alertCollection.InsertOne(context.TODO(), newAlert)

				if inserterr != nil {
					fmt.Println("Insert Error")
					log.Fatal(inserterr)
				}

				fmt.Println("The insert result is ", *insertResult)
				newAlert.ID = insertResult.InsertedID.(primitive.ObjectID)

				// Now do the Notification processing rules
				processGrouping(&newAlert , mongoClient)
				processNotifyRules( &newAlert , mongoClient)

				alertjsonData, err := json.Marshal(newAlert)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return
				}
				
				w.Header().Add("Content-Type" , "application/json")
				w.WriteHeader(201)
				w.Write([]byte(alertjsonData))
				fmt.Println("Inserted document successfully")
			} else {
				// Some other fatal error
				log.Fatal(err)
			}
		
		} else {
			// Duplicate Alert
			fmt.Printf("Found event: %+v\n", existingEvent)
			updatefilter := bson.M{"_id": existingEvent.ID }

			update := bson.M{
				"$set": bson.M{
					"alertcount": existingEvent.AlertCount + 1 ,
				},
			}
		
			updateResult , updateerr := alertCollection.UpdateOne(context.TODO(), updatefilter, update)
			if updateerr != nil {
				panic(err)
			}
			if updateResult.ModifiedCount > 0 {
				fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
			}
			alertjsonData, err := json.Marshal(apiAlertData)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(alertjsonData))
			
		}

	// De-duplication Ends
	}else{
		fmt.Println("This is a close event")
        // Find existing open alert first to check grouping info
		filter := bson.M{
			"alertid":     apiAlertData["alertId"].(string) ,
            "alertstatus": "OPEN",
		}
        
        var alertToClose models.DbAlert
        err := alertCollection.FindOne(context.TODO(), filter).Decode(&alertToClose)
        if err != nil {
             // Alert not found or error
             fmt.Println("Alert to close not found or error:", err)
             w.WriteHeader(200) // Idempotent success?
             return
        }
        
        parsedTime := time.Now()
        if val, ok := apiAlertData["alertTime"].(string); ok && val != "" {
             layout := "2006-01-02 15:04:05"
             if t, err := time.Parse(layout, val); err == nil {
                 parsedTime = t
             }
        }
        
        update := bson.M{
            "$set": bson.M{
                "alertstatus": "CLOSED",
                "alertcleartime": models.CustomTime{Time: parsedTime},
            },
        }
        
        updateResult, err := alertCollection.UpdateOne(context.TODO(), filter, update)
        if err != nil {
             log.Println("Error closing alert:", err)
             http.Error(w, "Error database update", http.StatusInternalServerError)
             return
        }
        
        fmt.Printf("Closed %v alerts matching %s\n", updateResult.ModifiedCount, apiAlertData["alertId"])
        
        // Check if Parent needs to be closed
        if alertToClose.Grouped && alertToClose.GroupIncidentId != "" {
             // It's a child. Check parent.
             
             // Safer: Query parent using _id converted from Hex string
             pID, err := primitive.ObjectIDFromHex(alertToClose.GroupIncidentId)
             if err == nil {
                 // Check if ANY child of this parent is still OPEN
                 
                 // Let's find Parent first to get GroupAlerts
                 var parent models.DbAlert
                 errP := alertCollection.FindOne(context.TODO(), bson.M{"_id": pID}).Decode(&parent)
                 if errP == nil {
                      // Count Checked Children
                      if len(parent.GroupAlerts) > 0 {
                          openChildrenCount, _ := alertCollection.CountDocuments(context.TODO(), bson.M{
                              "_id": bson.M{"$in": parent.GroupAlerts},
                              "alertstatus": "OPEN",
                          })
                          
                          if openChildrenCount == 0 {
                              fmt.Println("All children closed. Closing Parent Incident:", pID.Hex())
                              parentUpdate := bson.M{
                                  "$set": bson.M{
                                      "alertstatus": "CLOSED",
                                      "alertcleartime": models.CustomTime{Time: parsedTime},
                                  },
                              }
                              alertCollection.UpdateOne(context.TODO(), bson.M{"_id": pID}, parentUpdate)
                          } else {
                              fmt.Printf("Parent %s still has %d open children.\n", pID.Hex(), openChildrenCount)
                          }
                      }
                 }
             }
        }

        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Alert Closed", 
            "modifiedCount": updateResult.ModifiedCount,
        })
	}
}

func addTags(apiAlertData map[string]interface{}, newAlert *models.DbAlert) bool {
	Tags := make(map[string]interface{})
	exA := utilities.ExcludeAttributes{
		AlertTagExclude: []string{"entity","alertTime", "alertNotes", "severity","alertId","ipAddress","alertType","serviceName","alertSummary"},
	}

	for alertJsonKey, alertJsonValue  := range apiAlertData {

		if slices.Contains(exA.AlertTagExclude, alertJsonKey){
			// fmt.Printf(" %v is in exclude attribute\n",alertJsonKey )
		}else{
			// fmt.Printf(" %v is NOT in exclude attribute\n",alertJsonKey )
			// Add to the new map
			Tags[alertJsonKey] = alertJsonValue
		}
	}
	newAlert.AdditionalDetails = Tags 
	//fmt.Println("The object in addTags is " , newAlert )
	return true
}

func processAlertRules(newAlert *models.DbAlert, mongoClient *mongo.Client) bool {
	var rulesGroup ruleengine.RulesGroup
	alertRulesCollection := mongoClient.Database(mongodatabase).Collection("alertrules")

	cursor, err := alertRulesCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.TODO())

	var alertRules []models.DbAlertRule
	if err = cursor.All(context.TODO(), &alertRules); err != nil {
        log.Fatal(err)
    }

	for _, alertRule   := range alertRules {
		fmt.Println("Rule is ", alertRule.RuleObject)
		err := json.Unmarshal([]byte(alertRule.RuleObject), &rulesGroup)
		if err != nil {
			fmt.Println("Error in rule evaluation ", err)
		}
		var alertMap map[string]interface{}
		err1 := mapstructure.Decode(newAlert, &alertMap)
		if err1 != nil {
			fmt.Println("ERROR : Unable to convert struct to map")
		}
		fmt.Println("THE ALERT MAP IS ", alertMap)
		res := ruleengine.EvaluateRulesGroup(alertMap, rulesGroup)
		fmt.Printf("The Alert rule %v MATCH is %v \n", alertRule.RuleName , res)
		if res {
			if len(alertRule.SetField) == 0 {
				fmt.Println("The Set feild is empty. Skipping")
				continue
			}
			// Do the action specified in the rule
			v := reflect.ValueOf(newAlert).Elem()
			field := v.FieldByName(alertRule.SetField)
			if !field.IsValid() ||  !field.CanSet()  {
				fmt.Println("ERROR : The struct element is un settable")
				continue
			}
			fieldValue := reflect.ValueOf(alertRule.SetValue)
			field.Set(fieldValue)
		}
		fmt.Println("The MATCH is ", res)
	}
	return true
}

func processTagRules(newAlert *models.DbAlert, mongoClient *mongo.Client) bool {
	var rulesGroup ruleengine.RulesGroup
	tagRulesCollection := mongoClient.Database(mongodatabase).Collection("tagrules")

	cursor, err := tagRulesCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.TODO())

	var tagRules []models.DbTagRule
	if err = cursor.All(context.TODO(), &tagRules); err != nil {
        log.Fatal(err)
    }

	for _, tagRule   := range tagRules {
		fmt.Println("Rule is ", tagRule.RuleObject)
		err := json.Unmarshal([]byte(tagRule.RuleObject), &rulesGroup)
		if err != nil {
			fmt.Println("Error in rule evaluation ", err)
		}
		var alertMap map[string]interface{}
		err1 := mapstructure.Decode(newAlert, &alertMap)
		if err1 != nil {
			fmt.Println("ERROR : Unable to convert struct to map")
		}
		fmt.Println("THE ALERT MAP IS ", alertMap)
		res := ruleengine.EvaluateRulesGroup(alertMap, rulesGroup)
		fmt.Printf("The Tag rule %v MATCH is %v \n", tagRule.RuleName , res)
		if res {
			if len(tagRule.TagValue) != 0 {
				fmt.Println("The Tag Value is NOT empty. setting tag ")
				newAlert.AdditionalDetails[tagRule.TagName] = tagRule.TagValue
				continue
			}
			if len(tagRule.FieldExtraction) != 0 {
				pattern := tagRule.FieldExtraction
				re, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Println("Error compiling regex:", err)
					return false
				}
				// Do the action specified in the rule
				v := reflect.ValueOf(newAlert).Elem()
				field := v.FieldByName(tagRule.FieldName)
				if !field.IsValid() ||  !field.CanSet()  {
					fmt.Println("ERROR : The struct element is un settable")
					continue
				}

				if re.MatchString(field.String()) {
					// Extract the submatches
					matches := re.FindStringSubmatch(field.String())
					if len(matches) > 0 {

						fmt.Printf("Entire match: '%s'\n", matches[0])
						fmt.Printf("%s: '%s'\n",tagRule.TagName,  matches[1])
						// Add the extracted tag name to the additional details.
						newAlert.AdditionalDetails[tagRule.TagName] = matches[1]
					}
				}
			}
		}
		fmt.Println("The MATCH is ", res)
	}
	return true
}

func processGrouping(newAlert *models.DbAlert, mongoClient *mongo.Client) bool {
	fmt.Println("Processing group rules")
	alertGroupCollection := mongoClient.Database(mongodatabase).Collection("correlationrules")
	alertCollection := mongoClient.Database(mongodatabase).Collection("alerts")
	findOptions := options.Find()
    findOptions.SetSort(bson.D{{Key: "groupwindow", Value: 1}})
	cursor, err := alertGroupCollection.Find(context.TODO(), bson.D{},findOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.TODO())

	var alertGroupConfigs []models.DbAlertGroup

	if err = cursor.All(context.TODO(), &alertGroupConfigs); err != nil {
        log.Fatal(err)
    }

	for _, alertGroupConfig   := range alertGroupConfigs {

		fmt.Println("The alertConfig is ", alertGroupConfig)
        
        if alertGroupConfig.CorrelationMode == "SIMILARITY" {
             if processSimilarityRule(newAlert, alertGroupConfig, mongoClient) {
                 return true
             }
        } else {
		    // check if the alertPattern matches with the incomming event.
		    if (patternFound(alertGroupConfig.GroupTags , maps.Keys( newAlert.AdditionalDetails))){
			    // pattern is found in the incomming event
			    // construct the identifier
			    groupidentifier := ""
			    for _, tag := range alertGroupConfig.GroupTags {
				    groupidentifier = groupidentifier + "--" + newAlert.AdditionalDetails[tag].(string)
			    }
			    // if there is an event in open state with the same identifier
			    fmt.Println("THE IDENTIFIER IS ", groupidentifier)
			    findOptions := options.Find()
			    //findOptions.(bson.D{{Key: "groupidentifier", Value: groupidentifier}, {Key: "alertstatus" ,Value: "OPEN"}})
			    cursor, err := alertCollection.Find(context.TODO(), bson.M{"groupidentifier": groupidentifier , "alertstatus" : "OPEN"},findOptions)
			    if err != nil {
				    log.Fatal(err)
			    }
			    defer cursor.Close(context.TODO())

			    var identifiedalerts []models.DbAlert

			    if err = cursor.All(context.TODO(), &identifiedalerts); err != nil {
				    log.Fatal(err)
			    }
			    //
			    fmt.Println("THE TIME IS ", newAlert.AlertFirstTime.Unix() )
			    if len(identifiedalerts) != 0 {
				    fmt.Println("************************** FOUND OPEN EVENTS*********************************")
				    idn := identifiedalerts[len(identifiedalerts)-1] 
				    // There are open alerts with the same identifier
				    if newAlert.AlertFirstTime.Unix() - idn.AlertFirstTime.Unix() <= int64(alertGroupConfig.GroupWindow) {
					    // grpupEvent present and active
					    fmt.Printf("Identifier %s is within duration %v \n", groupidentifier , alertGroupConfig.GroupWindow )
					    // add the alert ID to the GroupAlerts 
					    updatefilter := bson.M{"_id": idn.ID }



					    update := bson.D{
						    {Key: "$push", Value: bson.D{
							    {Key: "groupalerts", Value: newAlert.ID},
						    }},
						    // {Key: "$set", Value: bson.D{
						    // 	{Key: "parent", Value: true},

						    // }},
					    }
				    
					    updateResult , updateerr := alertCollection.UpdateOne(context.TODO(), updatefilter, update)
					    if updateerr != nil {
						    panic(err)
					    }
					    if updateResult.ModifiedCount > 0 {
						    fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
					    }

					    // Update Grouped = true and GroupIncident ID for the alert,
					    updateOriginalAlertfilter := bson.M{"_id": newAlert.ID }

					    updateOriginal := bson.M{
						    "$set": bson.M{
							    "groupincidentid": idn.ID,
							    "grouped" : true,
						    },
					    }
			    
					    updateOriginalResult , updateerr := alertCollection.UpdateOne(context.TODO(), updateOriginalAlertfilter, updateOriginal)
					    if updateerr != nil {
						    panic(err)
					    }
					    if updateResult.ModifiedCount > 0 {
						    fmt.Printf("Matched %v documents and updated %v documents.\n", updateOriginalResult.MatchedCount, updateOriginalResult.ModifiedCount)
					    }
					    break
				    }else{ 
					    // create new spog event
					    copy := deepCopy(*newAlert)
					    copy.GroupIdentifier = groupidentifier
					    copy.AlertId = "grouped-"+groupidentifier
					    copy.GroupAlerts = append(copy.GroupAlerts, newAlert.ID)
					    copy.ID = primitive.ObjectID{}
					    // create a new parent alert
	    
					    insertResult , inserterr := alertCollection.InsertOne(context.TODO(), copy)
					    if inserterr != nil {
						    fmt.Println("Insert Error")
						    log.Fatal(inserterr)
					    }
					    fmt.Println("The insert result is ", *insertResult)
					    spogincidentId := insertResult.InsertedID.(primitive.ObjectID)
					    updatefilter := bson.M{"_id": newAlert.ID }
	    
					    update := bson.M{
						    "$set": bson.M{
							    "groupincidentid": spogincidentId,
							    "grouped" : true,
						    },
					    }
				    
					    updateResult , updateerr := alertCollection.UpdateOne(context.TODO(), updatefilter, update)
					    if updateerr != nil {
						    panic(err)
					    }
					    if updateResult.ModifiedCount > 0 {
						    fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
					    }
					    break
				    }

			    }else{
				    // create new spog event
				    copy := deepCopy(*newAlert)
				    copy.GroupIdentifier = groupidentifier
				    copy.AlertId = "grouped-"+groupidentifier
				    copy.GroupAlerts = append(copy.GroupAlerts, newAlert.ID)
				    copy.ID = primitive.ObjectID{}
				    copy.Parent = true
				    // create a new parent alert

				    insertResult , inserterr := alertCollection.InsertOne(context.TODO(), copy)
				    if inserterr != nil {
					    fmt.Println("Insert Error")
					    log.Fatal(inserterr)
				    }
				    fmt.Println("The insert result is ", *insertResult)
				    spogincidentId := insertResult.InsertedID.(primitive.ObjectID)
				    updatefilter := bson.M{"_id": newAlert.ID }

				    update := bson.M{
					    "$set": bson.M{
						    "groupincidentid": spogincidentId,
						    "grouped" : true,
					    },
				    }
			    
				    updateResult , updateerr := alertCollection.UpdateOne(context.TODO(), updatefilter, update)
				    if updateerr != nil {
					    panic(err)
				    }
				    if updateResult.ModifiedCount > 0 {
					    fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
				    }
				    break

			    }
		    }
        }
	}
	return true
}
type Item struct {
    ID primitive.ObjectID
}

func objectIdToString(id primitive.ObjectID) string {
    return id.Hex()
}

func processNotifyRules(newAlert *models.DbAlert, mongoClient *mongo.Client) bool {
	//const NoderedEndpoint = "http://192.168.1.201:1880/notifications"

	var rulesGroup ruleengine.RulesGroup
	notifyRulesCollection := mongoClient.Database(mongodatabase).Collection("notifyrules")

	cursor, err := notifyRulesCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.TODO())

	var notifyRules []models.DbNotifyRule
	if err = cursor.All(context.TODO(), &notifyRules); err != nil {
        log.Fatal(err)
    }

	for _, notifyRule   := range notifyRules {
		fmt.Println("Rule is ", notifyRule)
		err := json.Unmarshal([]byte(notifyRule.RuleObject), &rulesGroup)
		if err != nil {
			fmt.Println("Error in rule evaluation ", err)
		}
		var alertMap map[string]interface{}
		err1 := mapstructure.Decode(newAlert, &alertMap)
		if err1 != nil {
			fmt.Println("ERROR : Unable to convert struct to map")
		}
		fmt.Println("THE ALERT MAP IS ", alertMap)
		res := ruleengine.EvaluateRulesGroup(alertMap, rulesGroup)
		fmt.Printf("The Notify rule %v MATCH is %v \n", notifyRule.RuleName , res)

		if res {

			newAlert.AlertDestination = notifyRule.RuleName

			byteSlice, err := json.Marshal(newAlert)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println(byteSlice)
			response, err := http.Post(NoderedEndpoint, "application/json", bytes.NewBuffer(byteSlice))
			if err != nil {
				log.Fatalf("Error making POST request: %v", err)
			}
			defer response.Body.Close()
			
			fmt.Println("Here")
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

    		fmt.Println(string(body))
		}
		fmt.Println("The MATCH is ", res)
	}
	return true
}




func patternFound(array1, array2 []string) bool {
	// Create a map to store elements of array2
	elementMap := make(map[string]bool)
	
	// Populate the map with elements of array2
	for _, element := range array2 {
		elementMap[element] = true
	}
	
	// Check if all elements of array1 are present in the map
	for _, element := range array1 {
		if !elementMap[element] {
			return false
		}
	}
	
	return true
}

func deepCopy(original models.DbAlert) models.DbAlert {
    copy := original

    // Deep copy the map
    copy.AdditionalDetails = make(map[string]interface{})
    for k, v := range original.AdditionalDetails {
        copy.AdditionalDetails[k] = v
    }

    return copy
}



