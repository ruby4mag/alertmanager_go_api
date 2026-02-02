# PagerDuty Environment Variables Configuration

## Overview
The AlertManager API now supports configurable PagerDuty endpoints via environment variables. This allows you to configure different endpoints for creating, updating, and closing PagerDuty incidents without modifying the code.

## Environment Variables

### 1. N8N_PD_CREATE_ENDPOINT
**Purpose**: Endpoint for creating new PagerDuty incidents  
**Example**: `http://192.168.1.201/webhook-test/pagerduty/incident/create`  
**Used for**: Creating new PagerDuty incidents when alerts are triggered

### 2. N8N_PD_UPDATE_ENDPOINT
**Purpose**: Endpoint for updating existing PagerDuty incidents  
**Example**: `http://192.168.1.201:5678/webhook/pagerduty/incident/update`  
**Used for**: 
- Sending notes when child alerts are opened
- Sending notes when child alerts are closed

### 3. N8N_PD_CLEAR_ENDPOINT
**Purpose**: Endpoint for closing PagerDuty incidents  
**Example**: `http://192.168.1.201:5678/webhook/pagerduty/incident/clear`  
**Used for**: Closing PagerDuty incidents when all child alerts are resolved

## Configuration

### Setting Environment Variables

You can set these environment variables in several ways:

#### Option 1: Export in Shell
```bash
export N8N_PD_CREATE_ENDPOINT="http://192.168.1.201/webhook-test/pagerduty/incident/create"
export N8N_PD_UPDATE_ENDPOINT="http://192.168.1.201:5678/webhook/pagerduty/incident/update"
export N8N_PD_CLEAR_ENDPOINT="http://192.168.1.201:5678/webhook/pagerduty/incident/clear"
```

#### Option 2: .env File
Create a `.env` file in your project root:
```
N8N_PD_CREATE_ENDPOINT=http://192.168.1.201/webhook-test/pagerduty/incident/create
N8N_PD_UPDATE_ENDPOINT=http://192.168.1.201:5678/webhook/pagerduty/incident/update
N8N_PD_CLEAR_ENDPOINT=http://192.168.1.201:5678/webhook/pagerduty/incident/clear
```

#### Option 3: Docker/Container Environment
If running in Docker, add to your docker-compose.yml or Dockerfile:
```yaml
environment:
  - N8N_PD_CREATE_ENDPOINT=http://192.168.1.201/webhook-test/pagerduty/incident/create
  - N8N_PD_UPDATE_ENDPOINT=http://192.168.1.201:5678/webhook/pagerduty/incident/update
  - N8N_PD_CLEAR_ENDPOINT=http://192.168.1.201:5678/webhook/pagerduty/incident/clear
```

## Behavior

### Graceful Degradation
If an endpoint is not configured (empty string), the system will:
- Log a warning message
- Skip the operation gracefully
- Continue processing without errors

Example log messages:
```
⚠️  No PagerDuty update endpoint configured, skipping note
⚠️  No PagerDuty clear endpoint configured, skipping close
```

### API Flow

1. **Alert Creation** → Uses `N8N_PD_CREATE_ENDPOINT` to create PagerDuty incident
2. **Child Alert Joins Group** → Uses `N8N_PD_UPDATE_ENDPOINT` to add note to parent incident
3. **Child Alert Closes** → Uses `N8N_PD_UPDATE_ENDPOINT` to add closure note to parent incident
4. **All Children Closed** → Uses `N8N_PD_CLEAR_ENDPOINT` to close parent incident

## Migration from Hardcoded Values

### Previous Implementation
The endpoints were hardcoded in `utilities/pagerduty.go`:
```go
const PagerDutyUpdateEndpoint = "http://192.168.1.201:5678/webhook/pagerduty/incident/update"
const PagerDutyClearEndpoint = "http://192.168.1.201:5678/webhook/pagerduty/incident/clear"
```

And the create endpoint used `NODERED_ENDPOINT`:
```go
var NoderedEndpoint = os.Getenv("NODERED_ENDPOINT")
```

### New Implementation
All endpoints are now read from environment variables in `main.go` with consistent naming:
```go
var PagerDutyCreateEndpoint = os.Getenv("N8N_PD_CREATE_ENDPOINT")
var PagerDutyUpdateEndpoint = os.Getenv("N8N_PD_UPDATE_ENDPOINT")
var PagerDutyClearEndpoint = os.Getenv("N8N_PD_CLEAR_ENDPOINT")
```

### Required Changes
- Replace `NODERED_ENDPOINT` with `N8N_PD_CREATE_ENDPOINT`
- Set `N8N_PD_UPDATE_ENDPOINT` for update operations
- Set `N8N_PD_CLEAR_ENDPOINT` for close operations

## Testing

To verify the configuration is working:

1. Set the environment variables
2. Start the application
3. Check the logs when alerts are processed - you should see messages like:
   ```
   📝 Sending PagerDuty UPDATE to: http://your-endpoint/update
   🔒 Closing PagerDuty incident via: http://your-endpoint/clear
   ```

## Troubleshooting

### Issue: No PagerDuty incidents being created/updated/closed
**Solution**: Verify all three environment variables are set correctly:
```bash
echo $N8N_PD_CREATE_ENDPOINT
echo $N8N_PD_UPDATE_ENDPOINT
echo $N8N_PD_CLEAR_ENDPOINT
```

### Issue: Warnings about missing endpoints
**Solution**: Ensure environment variables are exported before starting the application

### Issue: Wrong endpoint being used
**Solution**: Check for typos in environment variable names (they are case-sensitive)
