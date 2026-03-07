[🔙 Back to Load Generator README](../README.md)

# API Load Generator Engine

This component runs a set of REST requests from Node.js, as specified in a job definition file (`config.json`), to generate synthetic load for APIs.

## Prerequisites

Before running this script directly, install the required npm modules:

```bash
npm install
```

## Management API

The server exposes a management API (default port `5950`) to control and monitor the continuous load generation process dynamically.

### Check Status

**Endpoint:** `GET /status`

Returns a carefully tracked JSON payload providing the current status of the service loop. 
*(If deployed with multiple instances, this returns the status of the specific node answering your request.)*

**Response Example:**
```json
{
  "version": "20180111-0814",
  "times": {
    "start": "Fri Feb 13 2015 02:58:10 GMT-0000 (UTC)",
    "lastRun": "Fri Feb 13 2015 03:07:28 GMT-0000 (UTC)",
    "wake": "Fri Feb 13 2015 03:08:38 GMT-0000 (UTC)",
    "current": "Fri Feb 13 2015 03:08:06 GMT-0000 (UTC)"
  },
  "loglevel": 2,
  "nRequests": 42,
  "jobId": "auto-claims-demo-template",
  "description": "Drive synthetic API traffic",
  "status": "waiting",
  "responseCounts": {
    "total": 42,
    "200": 41,
    "401": 1
  },
  "durationOfLastRunInMs": 1632,
  "currentRunsPerHour": 51
}
```

### Control Generator

**Endpoint:** `POST /control`
**Headers:** `Content-Type: application/x-www-form-urlencoded`

Use this endpoint to pause/resume load tests or inspect detailed logging without needing to push a brand new deployment or kill the service.

#### Option 1: Start/Stop Traffic

Pass the `action` parameter to halt or resume the traffic generator node.
*   `action=start`
*   `action=stop`

#### Option 2: Adjust Log Level

Pass the `action` parameter alongside the `loglevel` to dramatically alter console output for debugging purposes.

*   `action=setlog&loglevel=N`

*(Where N is between 0 and 10)*
*   **0:** Almost no logging
*   **2:** Very minimal logging (only wake/sleep and errors)
*   **3:** See every single API call out
*   **10:** Maximum debug logging
