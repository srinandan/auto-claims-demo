[🔙 Back to Main Project README](../README.md)

# Auto Claims Load Generator

This is the load generator for the Auto Claims Demo, intended to drive synthetic live traffic to the backend APIs so that the analytics engine and OpenTelemetry traces have interesting data.

## Tech Stack

-   **Runtime**: [Node.js](https://nodejs.org/) (v22+)
-   **Templating**: [Handlebars](https://handlebarsjs.com/) (for request variation)
-   **HTTP Framework**: [Express](https://expressjs.com/) (for control server)
-   **Telemetry**: [OpenTelemetry](https://opentelemetry.io/)

## Prerequisites

-   Node.js (v22.16.0 or later)
-   npm
-   Google Cloud SDK (`gcloud`)

## Setup & Run

The project includes a `Makefile` for easy local execution and Cloud Build deployment.

### Run Locally

1. Set your environment variables in `.env` (or `ENV`).
2. Run the load generator locally via Make:

```bash
make local-loadgen
```

This will automatically install dependencies and start the generator.

### Deploy to Cloud Run

You can deploy the load generator to run continuously on Cloud Run.

Make sure you are signed in (`gcloud auth login`) and that your user has permissions to deploy to Cloud Run (`roles/run.admin`) in your Google Cloud project.

**Using Cloud Build (Recommended):**
```bash
make loadgen
```

## Project Structure

*   `api-load-generator.js`: The main entrypoint that bootstraps everything.
*   `lib/`: Core modules.
    *   `load-generator.js`: The engine that sends the actual API traffic.
    *   `ScheduleSelector.js` & `WeightedRandomSelector.js`: Handlebars helpers for varying requests by time or weighted probabilities.
    *   `data/`: Files supporting simulation of different originating IP Addresses/Geographies.
    *   `server.js`: Express server to handle status/control requests (listening on port 5950).
*   `config/`: Contains `config.json` and `defaults.json` defining target APIs, API keys, and traffic rates. `config.json` heavily uses Handlebars templates for dynamic request generation.
*   `tracing.js`: OpenTelemetry auto-instrumentation setup.