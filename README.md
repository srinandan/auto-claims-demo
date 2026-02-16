# Auto Claims Application

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project is an **Auto Claims Processing Application** that streamlines the claims process by allowing users to submit claims and automatically analyzing vehicle images for damage using an AI service.

## Architecture Overview

The system consists of three main microservices:

1.  **[Frontend](./frontend/README.md)**: A responsive web interface built with **Vue 3**, **Vite**, and **Tailwind CSS v4** for users to submit claims and view results.
2.  **[Backend](./backend/README.md)**: A robust REST API built with **Go (Gin)** and **SQLite** that manages claim data, policies, and orchestration.
3.  **[AI Service](./ai-service/README.md)**: A specialized service built with **Python (FastAPI)** that leverages **Google Cloud Vertex AI / Gemini** (or a mock implementation) to detect vehicle damage from images.

## Prerequisites

Ensure you have the following installed:

*   [Node.js](https://nodejs.org/) (v18+)
*   [Go](https://go.dev/) (v1.23+)
*   [Python](https://www.python.org/) (v3.10+) and [uv](https://docs.astral.sh/uv/)
*   [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (`gcloud`)

## Quick Start

Detailed instructions for each service are available in their respective directories, but here is a quick summary to get everything running locally.

### 1. Start the Backend
```bash
cd backend
make local-backend
# Runs on http://localhost:8080
```

### 2. Start the AI Service
```bash
cd ai-service
make local-ai-service
# Runs on http://localhost:8000
```

### 3. Start the Frontend
```bash
cd frontend
make local-frontend
# Runs on http://localhost:5173
```

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute to this project.

## Support

This demo is *NOT* endorsed by Google or Google Cloud. The repo is intended for educational/hobbyists use only.

## License

This project is licensed under the terms of the [LICENSE.txt](./LICENSE.txt) file. The AI generated car damage pictures are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit <http://creativecommons.org/licenses/by/4.0/>