# Auto Claims AI Service

This service handles image analysis for detecting vehicle damage using **Google Cloud Vertex AI / Gemini**.

## Tech Stack

-   **Language**: [Python](https://www.python.org/) (v3.10+)
-   **Framework**: [FastAPI](https://fastapi.tiangolo.com/)
-   **AI/ML**: [Google Cloud Vertex AI](https://cloud.google.com/vertex-ai) / Gemini Pro Vision
-   **Package Manager**: [uv](https://docs.astral.sh/uv/)

## Prerequisites

-   Python 3.10+
-   `uv` package manager
-   Google Cloud SDK (`gcloud`)

## Configuration

### Mock Mode
You can run the service in **Mock Mode** to avoid calling actual Google Cloud APIs. This is useful for local development and testing without incurring costs or needing full credentials.

Set `MOCK_MODE=true` in your environment or `Makefile`.

## Setup & Run

### Run Locally

```bash
make local-ai-service
```

This command will:
1.  Set default environment variables.
2.  Start the FastAPI server using `uv` on `http://localhost:8000`.

To run in mock mode:
```bash
MOCK_MODE=true make local-ai-service
```

## API Endpoints

-   `POST /analyze`: Analyzes a single uploaded image for damage.
    -   Inputs: `file` (UploadFile), `photo_id` (Form)
-   `POST /process-claims`: Batch processes images for a claim using a sequential agent.
    -   Inputs: `file_uris` (List[str])
-   `POST /find-repair-shops`: Finds repair shops based on location and damage.
    -   Inputs: `zip_code`, `state`, `make`, `model`, `damage_type`
