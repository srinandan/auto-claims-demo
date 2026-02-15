# Auto Claims Application

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project is an Auto Claims processing application that allows users to submit claims and automatically analyzes images for damage using an AI service.

## Project Structure

The project is divided into three main components:

*   **Frontend (`frontend/`)**: A Vue 3 + Vite application for the user interface, styled with Tailwind CSS v4.
*   **Backend (`backend/`)**: A Go (Gin) REST API that manages claims and interacts with the database.
*   **AI Service (`ai-service/`)**: A Python (FastAPI) service that analyzes images for damage. It can run in mock mode or use real inference (mocked or using Vertex AI/Vision API).

## Prerequisites

Ensure you have the following installed:

*   [Node.js](https://nodejs.org/) (v18 or later)
*   [Go](https://go.dev/) (v1.23 or later)
*   [Python](https://www.python.org/) (v3.10 or later) and [uv](https://docs.astral.sh/uv/)
*   [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (`gcloud`)

## Getting Started

### 1. Backend

The backend handles API requests and database interactions.

**Setup & Run:**

```bash
cd backend
make local-backend
```

This command will:
*   Set necessary environment variables (`PROJECT_ID`, `PROJECT_NUMBER`, `REGION`, `_COMMIT_SHA`).
*   Download Go dependencies (`go mod tidy`).
*   Run the server on `http://localhost:8080`.

### 2. AI Service

The AI service analyzes images uploaded for claims.

**Setup & Run:**

```bash
cd ai-service
make local-ai-service
```

This command will:
*   Set necessary environment variables.
*   Run the FastAPI server using `uv` on `http://localhost:8000`.
*   Note: By default, it runs in `MOCK_MODE=false`. To run in mock mode, set `MOCK_MODE=true` in the Makefile or environment.

### 3. Frontend

The frontend provides the user interface for submitting claims.

**Setup & Run:**

```bash
cd frontend
make local-frontend
```

This command will:
*   Set necessary environment variables.
*   Install Node.js dependencies (`npm install`).
*   Start the Vite development server.

## API Endpoints

### Backend
*   `GET /ping`: Health check.
*   `GET /api/claims`: List all claims.
*   `GET /api/claims/:id`: Get a specific claim.
*   `POST /api/claims`: Create a new claim.
*   `POST /api/claims/:id/analyze`: Trigger analysis for a claim.
*   `GET /api/policies/:number`: Get a policy by number.

### AI Service
*   `POST /analyze`: Analyzes an uploaded image for damage.
*   `POST /process-claims`: Processes a batch of claims images.

## License

This project is licensed under the terms of the [Apache License 2.0](./LICENSE).