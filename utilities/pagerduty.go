package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PagerDutyUpdatePayload represents the payload to send to PagerDuty update endpoint
type PagerDutyUpdatePayload struct {
	IncidentID string `json:"incident_id"`
	Content    string `json:"content"`
}

// PagerDutyClearPayload represents the payload to send to PagerDuty clear endpoint
type PagerDutyClearPayload struct {
	IncidentID string `json:"incident_id"`
}

// SendPagerDutyNote sends a note to a PagerDuty incident
// Returns error if the request fails
func SendPagerDutyNote(endpoint string, incidentID string, content string) error {
	if incidentID == "" {
		// No PagerDuty incident associated, skip silently
		fmt.Println("⚠️  No PagerDuty incident ID provided, skipping note")
		return nil
	}

	if endpoint == "" {
		fmt.Println("⚠️  No PagerDuty update endpoint configured, skipping note")
		return nil
	}

	payload := PagerDutyUpdatePayload{
		IncidentID: incidentID,
		Content:    content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling PagerDuty payload: %v", err)
	}

	fmt.Printf("📝 Sending PagerDuty UPDATE to: %s\n", endpoint)
	fmt.Printf("   Incident ID: %s\n", incidentID)
	fmt.Printf("   Content: %s\n", content)
	fmt.Printf("   Payload: %s\n", string(jsonData))

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending PagerDuty update: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("   Response Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response Body: %s\n", string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PagerDuty update returned status: %d", resp.StatusCode)
	}

	fmt.Printf("✅ Successfully sent PagerDuty note to incident %s\n", incidentID)
	return nil
}

// ClosePagerDutyIncident closes a PagerDuty incident
// Returns error if the request fails
func ClosePagerDutyIncident(endpoint string, incidentID string) error {
	if incidentID == "" {
		// No PagerDuty incident associated, skip silently
		fmt.Println("⚠️  No PagerDuty incident ID provided, skipping close")
		return nil
	}

	if endpoint == "" {
		fmt.Println("⚠️  No PagerDuty clear endpoint configured, skipping close")
		return nil
	}

	payload := PagerDutyClearPayload{
		IncidentID: incidentID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling PagerDuty clear payload: %v", err)
	}

	fmt.Printf("🔒 Closing PagerDuty incident via: %s\n", endpoint)
	fmt.Printf("   Incident ID: %s\n", incidentID)
	fmt.Printf("   Payload: %s\n", string(jsonData))

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error closing PagerDuty incident: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("   Response Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response Body: %s\n", string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PagerDuty clear returned status: %d", resp.StatusCode)
	}

	fmt.Printf("✅ Successfully closed PagerDuty incident %s\n", incidentID)
	return nil
}
