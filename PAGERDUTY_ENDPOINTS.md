# PagerDuty Endpoint Usage Summary

## Overview
The application uses **TWO different endpoints** for PagerDuty operations:

### 1. CREATE Endpoint (for new incidents)
- **Endpoint**: `NoderedEndpoint` environment variable (e.g., `http://192.168.1.201/webhook-test/pagerduty/incident/create`)
- **Used for**: Creating NEW PagerDuty incidents
- **When**: 
  - First alert in a group (parent alert)
  - Standalone alerts (not grouped)
- **Location**: `main.go` - `processNotifyRules()` function (line ~804)
- **Payload**: Full alert object with PagerDuty service and escalation policy

### 2. UPDATE Endpoint (for adding notes)
- **Endpoint**: `http://192.168.1.201:5678/webhook/pagerduty/incident/update` (hardcoded in `utilities/pagerduty.go`)
- **Used for**: Adding notes to EXISTING PagerDuty incidents
- **When**:
  - Child alert is added to a group (OPENED note)
  - Grouped alert is closed (CLOSED note)
- **Location**: `utilities/pagerduty.go` - `SendPagerDutyNote()` function
- **Payload**: 
  ```json
  {
    "incident_id": "string",
    "content": "entity:summary is OPENED/CLOSED"
  }
  ```

## Flow Example

```
Alert 1 arrives (matches notify rule)
  ↓
processGrouping() → No match → Creates parent (Grouped=false, Parent=true)
  ↓
processNotifyRules() → Grouped=false
  ↓
🆕 Creates NEW incident via CREATE endpoint
  ↓ 
Stores pagerduty_incident_id in Alert 1

---

Alert 2 arrives (matches notify rule, groups with Alert 1)
  ↓
processGrouping() → Matches Alert 1 → Sets Grouped=true, GroupIncidentId=Alert1.ID
  ↓
processNotifyRules() → Grouped=true
  ↓
🔗 Retrieves parent (Alert 1)
  ↓
📝 Sends UPDATE note via UPDATE endpoint
  ↓
Uses Alert 1's pagerduty_incident_id
  ↓
✅ Skips creating new incident

---

Alert 2 is closed
  ↓
Retrieves parent (Alert 1)
  ↓
📝 Sends UPDATE note via UPDATE endpoint
  ↓
Content: "entity:summary is CLOSED"
```

## Logging Output

With the enhanced logging, you should see:

### For Parent/Standalone Alerts:
```
🆕 Alert alert-123 is PARENT or STANDALONE. Creating NEW PagerDuty incident.
   Endpoint: http://192.168.1.201/webhook-test/pagerduty/incident/create
   Payload size: 1234 bytes
   Received response from PagerDuty create endpoint
```

### For Grouped Child Alerts:
```
🔗 Alert alert-456 is a GROUPED CHILD. Will update parent's PagerDuty incident instead of creating new one.
📝 Sending PagerDuty UPDATE to: http://192.168.1.201:5678/webhook/pagerduty/incident/update
   Incident ID: PXXXXXX
   Content: server-01:High CPU is OPENED
   Payload: {"incident_id":"PXXXXXX","content":"server-01:High CPU is OPENED"}
   Response Status: 200
   Response Body: {...}
✅ Successfully sent PagerDuty note to incident PXXXXXX
✅ Skipped creating new PagerDuty incident for grouped child alert
```

### For Alert Closure:
```
📝 Sending PagerDuty UPDATE to: http://192.168.1.201:5678/webhook/pagerduty/incident/update
   Incident ID: PXXXXXX
   Content: server-01:High CPU is CLOSED
   Payload: {"incident_id":"PXXXXXX","content":"server-01:High CPU is CLOSED"}
   Response Status: 200
   Response Body: {...}
✅ Successfully sent PagerDuty note to incident PXXXXXX
```

## Verification

To verify the correct endpoints are being used:

1. **Check logs** for the emoji indicators:
   - 🆕 = Creating new incident (CREATE endpoint)
   - 🔗 = Grouped child (will use UPDATE endpoint)
   - 📝 = Sending update note (UPDATE endpoint)

2. **Monitor network traffic** to see which endpoints are called

3. **Check PagerDuty** to ensure:
   - Only ONE incident is created per group
   - Notes are added to the incident when alerts join/leave

## Environment Variables

Make sure `NODERED_ENDPOINT` is set correctly:
```bash
export NODERED_ENDPOINT="http://192.168.1.201/webhook-test/pagerduty/incident/create"
```

The UPDATE endpoint is hardcoded in `utilities/pagerduty.go` as:
```go
const PagerDutyUpdateEndpoint = "http://192.168.1.201:5678/webhook/pagerduty/incident/update"
```
