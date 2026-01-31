# PagerDuty Integration for Alert Grouping

## Overview
This implementation ensures that when alerts are grouped together, only the parent alert creates a PagerDuty incident. Child alerts that join the group will **update** the parent's existing PagerDuty incident with notes instead of creating new incidents.

## Changes Made

### 1. New Utility Module: `utilities/pagerduty.go`
Created a new utility module for PagerDuty operations with the following functionality:
- **Function**: `SendPagerDutyNote(incidentID string, content string) error`
- **Endpoint**: `http://192.168.1.201:5678/webhook/pagerduty/incident/update`
- **Payload Format**:
  ```json
  {
    "incident_id": "string",
    "content": "string"
  }
  ```
- Silently skips if no PagerDuty incident ID is present
- Returns error if the HTTP request fails

### 2. Updated `main.go` - `processNotifyRules()` Function
**Key Change**: Added logic to prevent grouped child alerts from creating new PagerDuty incidents.

**Behavior**:
- **Before processing**: Checks if `newAlert.Grouped == true` and `newAlert.GroupIncidentId != ""`
- **If grouped (child alert)**:
  1. Retrieves the parent alert from the database
  2. Checks if parent has a `pagerduty_incident_id`
  3. Sends an update note to the parent's incident: `<entity>:<alert summary> is OPENED`
  4. **Skips** creating a new PagerDuty incident for the child
  5. Continues to next rule
- **If not grouped (parent or standalone alert)**:
  1. Creates a new PagerDuty incident normally
  2. Stores the incident ID in the alert document

This ensures that:
- **Parent alerts** create PagerDuty incidents
- **Standalone alerts** create PagerDuty incidents
- **Child alerts** only update the parent's existing incident

### 3. Updated `main.go` - Alert Closure Logic
When a grouped alert is closed:
- Retrieves the parent alert
- Checks if the parent has a `pagerduty_incident_id`
- Sends a note with format: `<entity>:<alert summary> is CLOSED`
- Note is sent before checking if all children are closed

### 4. Updated `grouping_logic.go`
The `addToGroup()` function now only handles database updates:
- Adds child to parent's `groupalerts` array
- Sets child's `grouped=true` and `groupincidentid`
- **Does NOT** send PagerDuty notes (handled by `processNotifyRules`)

## Flow Diagram

```
New Alert Created
    ↓
processGrouping() - Determines if alert should be grouped
    ↓
    ├─→ If matches existing group:
    │       - Sets Grouped = true
    │       - Sets GroupIncidentId = parent ID
    │       - Adds to parent's GroupAlerts array
    │
    └─→ If no match:
            - Creates new parent (if needed)
            - Sets Parent = true
    ↓
processNotifyRules() - Handles PagerDuty incident creation/updates
    ↓
    ├─→ If Grouped = true (child alert):
    │       - Retrieves parent alert
    │       - Sends note to parent's PagerDuty incident
    │       - SKIPS creating new incident
    │
    └─→ If Grouped = false (parent/standalone):
            - Creates new PagerDuty incident
            - Stores incident ID in alert
```

## Note Format
All notes follow this format:
- **OPENED**: `<entity>:<alert summary> is OPENED`
- **CLOSED**: `<entity>:<alert summary> is CLOSED`

Example: `web-server-01:High CPU Usage is OPENED`

## Verification
The field `pagerduty_incident_id` is checked before sending any notes. This field is populated when:
- A PagerDuty incident is created for a parent or standalone alert
- The incident ID is stored in the alert document (see `processNotifyRules()` in main.go)

## Error Handling
- All PagerDuty note operations are non-blocking
- Failures are logged as warnings but don't interrupt alert processing
- Empty incident IDs are silently skipped

## Testing Recommendations
1. **Test 1: First Alert (Parent)**
   - Create an alert that triggers a PagerDuty incident
   - Verify the `pagerduty_incident_id` is populated
   - Verify a new PagerDuty incident is created

2. **Test 2: Second Alert (Child - Grouped)**
   - Create a second alert that groups with the first
   - Verify `grouped=true` and `groupincidentid` is set
   - Verify **NO new PagerDuty incident is created**
   - Check PagerDuty incident for the "OPENED" note

3. **Test 3: Close Child Alert**
   - Close the second alert
   - Check PagerDuty incident for the "CLOSED" note
   - Verify the parent incident remains open

4. **Test 4: Close All Alerts**
   - Close all child alerts
   - Verify parent incident is automatically closed

