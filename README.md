# Auto Claims Application

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

![Architecture Diagram](./infra/infra.png)

This project is an **Auto Claims Processing Application** that streamlines the claims process by allowing users to submit claims and automatically analyzing vehicle images for damage using an AI service.

The system consists of several microservices:

1.  **[Frontend](./frontend/README.md)**: A responsive web interface built with **Vue 3**, **Vite**, and **Tailwind CSS v4**.
2.  **[Backend](./backend/README.md)**: A robust REST API built with **Go (Gin)** and **SQLite**.
3.  **[AI Service](./ai-service/README.md)**: A specialized service built with **Python (FastAPI)** that uses **Vertex AI** and **YOLOv11** to analyze images and orchestrate the claims process.
4.  **[Assessor Agent](./assessor-agent/README.md)**: A remote A2A agent that assesses damage severity based on AI analysis.
5.  **[Processor Agent](./processor-agent/README.md)**: A remote A2A agent that generates repair estimates and makes final claim decisions.
6.  **[Repair Shop Agent](./repair-shop-agent/README.md)**: A remote A2A agent that handles communication with repair shops to book appointments.
7.  **[Load Generator](./loadgen/README.md)**: A synthetic API load generation tool built with **Node.js** that injects realistic traffic patterns into the backend for OpenTelemetry analytics.

## Prerequisites

Ensure you have the following installed:

*   [Node.js](https://nodejs.org/) (v18+)
*   [Go](https://go.dev/) (v1.23+)
*   [Python](https://www.python.org/) (v3.10+) and [uv](https://docs.astral.sh/uv/)
*   [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (`gcloud`)

## Quick Start

Detailed instructions for each service are available in their respective directories, but here is a quick summary to get everything running locally.

### 1. Start the Agents

In separate terminals:
```bash
# Terminal 1: Assessor Agent
cd assessor-agent
make local-assessor-agent
# Runs on http://localhost:8081

# Terminal 2: Processor Agent
cd processor-agent
make local-processor-agent
# Runs on http://localhost:8082

# Terminal 3: Repair Shop Agent
cd repair-shop-agent
make local-repair-shop-agent
# Runs on http://localhost:8083
```

### 2. Start the AI Service
```bash
cd ai-service
make local-ai-service
# Runs on http://localhost:8000
```

### 3. Start the Backend
```bash
cd backend
make local-backend
# Runs on http://localhost:8080
```

### 4. Start the Frontend
```bash
cd frontend
make local-frontend
# Runs on http://localhost:5173
```

### 5. Start the Load Generator (Optional)
If you want to generate synthetic traffic to your APIs:
```bash
cd loadgen
make local-loadgen
# Begins generating traffic immediately
```

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute to this project.

## Support

This demo is *NOT* endorsed by Google or Google Cloud. The repo is intended for educational/hobbyists use only.

## License

This project is licensed under the terms of the [LICENSE.txt](./LICENSE.txt) file. The AI generated car damage pictures are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit <http://creativecommons.org/licenses/by/4.0/>
