# Auto Claims Application

[![CI](https://github.com/srinandan/auto-claims-demo/actions/workflows/ci.yaml/badge.svg)](https://github.com/srinandan/auto-claims-demo/actions/workflows/ci.yaml)
[![License](https://img.shields.io/github/license/srinandan/auto-claims-demo)](./LICENSE.txt)
[![Go Version](https://img.shields.io/github/go-mod/go-version/srinandan/auto-claims-demo?filename=backend/go.mod)](./backend/go.mod)
[![Python Version](https://img.shields.io/badge/python-3.11+-blue.svg)](./ai-service/pyproject.toml)
[![Node Version](https://img.shields.io/badge/node-20+-green.svg)](./frontend/package.json)
[![CodeQL](https://github.com/srinandan/auto-claims-demo/actions/workflows/codeql.yml/badge.svg)](https://github.com/srinandan/auto-claims-demo/actions/workflows/codeql.yml)

This project is an **Auto Claims Processing Application** that streamlines the claims process by allowing users to submit claims and automatically analyzing vehicle images for damage using an AI service.

## Architecture
![Architecture Diagram](./infra/infra.png)

This is a monorepo containing:
1.  **[Frontend](./frontend/README.md)**: A responsive web interface built with **Vue 3**, **Vite**, and **Tailwind CSS v4**.
2.  **[Backend](./backend/README.md)**: A robust REST API built with **Go (Gin)** and **SQLite**.
3.  **[AI Service](./ai-service/README.md)**: A specialized service built with **Python (FastAPI)** that uses **Vertex AI** and **YOLOv11** to analyze images and orchestrate the claims process.
4.  **[Assessor Agent](./assessor-agent/README.md)**: A remote A2A agent that assesses damage severity based on AI analysis.
5.  **[Processor Agent](./processor-agent/README.md)**: A remote A2A agent that generates repair estimates and makes final claim decisions.
6. **[Repair Shop Agent](./repair-shop-agent/README.md)**: A remote A2A agent that handles communication with repair shops to book appointments.
7. **[Agent Engine](./agent-engine/README.md)**: A specialized service that provides shared session and memory management across all agents.
8. **[Load Generator](./loadgen/README.md)**: A synthetic API load generation tool built with **Node.js** that injects realistic traffic patterns into the backend for OpenTelemetry analytics.

## Prerequisites

Ensure you have the following installed:
*   [Node.js](https://nodejs.org/) (v18+)
*   [Go](https://go.dev/) (v1.25.8+)
*   [Python](https://www.python.org/) (v3.11+) and [uv](https://docs.astral.sh/uv/)
*   [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (`gcloud`)
*   Google Cloud Project with Billing enabled

### 1. Automated Infrastructure Setup
We provide a script `infra/setup.py` to automate the one-time setup of GCP APIs, Service Accounts, and other fundamental cloud infrastructure required by the application.

```bash
python3 infra/setup.py
```

This will:
- Enable necessary APIs (Cloud Run, Vertex AI, Artifact Registry, Secret Manager, BigQuery, etc.).
- Create a service account `auto-claims-sa` with the required IAM roles.
- Create a GCS bucket, Artifact Registry repository, Secret Manager secrets, BigQuery dataset, and Cloud Build worker pool.

### 1. Installation
Install dependencies for all components according to their respective READMEs.

### 2. Execution
Start all microservices locally with a single command:
```bash
make local-all
```

To stop all locally running microservices, simply run:
```bash
make local-stop
```

### 3. Deployment
To deploy all components to Google Cloud Run natively from the root directory:
```bash
make deploy
```
This command leverages the root Makefile to build and submit deployments for the frontend, backend, AI service, and all agents.

### 4. Load Balancer Setup (Optional)
Once all services are deployed natively, you can optionally execute the load balancer script to tie the frontend and backend together under a single Global Application Load Balancer with a managed SSL certificate.
```bash
python3 infra/setup_lb.py
```

### 5. Individual Component Setup

For detailed setup, configuration, and individual execution of each component, refer to their respective READMEs:
- **[Frontend](./frontend/README.md)**
- **[Backend](./backend/README.md)**
- **[AI Service](./ai-service/README.md)**
- **[Assessor Agent](./assessor-agent/README.md)**
- **[Processor Agent](./processor-agent/README.md)**
- **[Repair Shop Agent](./repair-shop-agent/README.md)**
- **[Load Generator](./loadgen/README.md)**

## Built With

This application was built with the assistance of [Stitch](https://stitch.withgoogle.com/) and [Jules](https://jules.google.com).

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute to this project.

## Support

This demo is *NOT* endorsed by Google or Google Cloud. The repo is intended for educational/hobbyists use only.

## License

This project is licensed under the terms of the [LICENSE.txt](./LICENSE.txt) file. The AI generated car damage pictures are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit <http://creativecommons.org/licenses/by/4.0/>
