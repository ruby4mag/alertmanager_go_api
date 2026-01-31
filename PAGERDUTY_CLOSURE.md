# PagerDuty Incident Closure Integration

## Overview
When all child alerts in a group are closed, the parent alert is automatically closed. This implementation ensures that the PagerDuty incident is also closed when the parent alert is closed.

## Implementation

### New Function: `ClosePagerDutyIncident()`
**Location**: `utilities/pagerduty.go`

**Endpoint**: `http://192.168.1.201/webhook-test/pagerduty/incident/clear`

**Payload**:
```json
{
  "incident_id": "string"
}
```

**Behavior**:
- Sends HTTP POST request to the clear endpoint
- Logs the request and response for debugging
- Returns error if the request fails (non-blocking)

### Integration Point
**Location**: `main.go` - Alert closure logic (line ~395)

**When triggered**:
- A grouped child alert is closed
- System checks if all children are closed
- If all children are closed:
  1. Parent alert status is updated to "CLOSED"
  2. **PagerDuty incident is closed** (NEW)

## Flow Diagram

```
Child Alert Closed
  ↓
Check parent's children
  ↓
All children closed?
  ├─→ NO: Parent remains OPEN
  │       PagerDuty incident remains OPEN
  │
  └─→ YES: Close parent alert in DB
           ↓
           Check if parent has PagerDuty incident ID
           ↓
           🔒 Call ClosePagerDutyIncident()
           ↓
           POST to /webhook-test/pagerduty/incident/clear
           ↓
           ✅ PagerDuty incident CLOSED
```

## Log Output

When all children are closed and parent is closed:

```
All children closed. Closing Parent Incident: 697d9146ca1f8955660bfe9a
🔒 Closing PagerDuty incident for parent alert 697d9146ca1f8955660bfe9a
🔒 Closing PagerDuty incident via: http://192.168.1.201/webhook-test/pagerduty/incident/clear
   Incident ID: Q0R898N63GAL9A
   Payload: {"incident_id":"Q0R898N63GAL9A"}
   Response Status: 200
   Response Body: {...}
✅ Successfully closed PagerDuty incident Q0R898N63GAL9A
```

## Complete Alert Lifecycle with PagerDuty

### 1. First Alert (Parent Created)
```
Alert 1 arrives
  ↓
Creates parent alert
  ↓
🆕 Creates PagerDuty incident
  ↓
Stores incident_id in parent
```

### 2. Second Alert (Child Added)
```
Alert 2 arrives
  ↓
Matches parent
  ↓
Marked as child
  ↓
📝 Sends "OPENED" note to parent's incident
```

### 3. Child Alert Closed
```
Alert 2 closed
  ↓
📝 Sends "CLOSED" note to parent's incident
  ↓
Check: All children closed? NO
  ↓
Parent remains OPEN
```

### 4. All Children Closed
```
Alert 1 closed
  ↓
📝 Sends "CLOSED" note to parent's incident
  ↓
Check: All children closed? YES
  ↓
Close parent alert in DB
  ↓
🔒 Close PagerDuty incident
```

## Error Handling
- All PagerDuty operations are non-blocking
- Failures are logged as warnings
- Empty incident IDs are silently skipped
- Parent alert is still closed even if PagerDuty close fails

## Testing Checklist
- [ ] Create first alert → Verify PagerDuty incident created
- [ ] Create second alert (grouped) → Verify note added to incident
- [ ] Close second alert → Verify "CLOSED" note added
- [ ] Close first alert → Verify parent closed AND PagerDuty incident closed
- [ ] Verify PagerDuty incident status is "resolved" or "closed"
