# Auto Claims Backend

This is the backend service for the Auto Claims Demo, built with **Go** and **Gin**.

## Tech Stack

-   **Language**: [Go](https://go.dev/) (v1.23+)
-   **Framework**: [Gin](https://gin-gonic.com/)
-   **Database**: [SQLite](https://www.sqlite.org/index.html) (via [GORM](https://gorm.io/))
-   **ORM**: GORM

## Prerequisites

-   Go (v1.23 or later)
-   Google Cloud SDK (`gcloud`) - for authentication if using GCP features

## Setup & Run

### Environment Variables

The application expects the following environment variables (handled by `Makefile` for local dev):

-   `PROJECT_ID`: Google Cloud Project ID
-   `PROJECT_NUMBER`: Google Cloud Project Number
-   `REGION`: Google Cloud Region (e.g., `us-central1`)
-   `BUCKET_NAME`: GCS Bucket for storing images
-   `SERVICE_ACCOUNT_EMAIL`: Service Account Email

### Run Locally

```bash
make local-backend
```

This command will:
1.  Set default environment variables.
2.  Install Go dependencies (`go mod tidy`).
3.  Start the server on `http://localhost:8080`.

## API Endpoints

### Claims
-   `GET /api/claims`: List all claims
-   `GET /api/claims/:id`: Get specific claim details
-   `POST /api/claims`: Create a new claim
-   `PUT /api/claims/:id`: Update a claim
-   `POST /api/claims/:id/analyze`: Trigger AI analysis for a claim
-   `POST /api/claims/:id/repair-shops`: Find repair shops for a claim

### Policies
-   `GET /api/policies/:number`: Get a policy by its number

### Health
-   `GET /ping`: Health check endpoint (returns `pong`)
