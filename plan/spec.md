# Auto-Claims Demo - Project Specification

## 1. Project Overview
The **Auto-Claims Demo** is an end-to-end, event-driven system designed to automate and streamline the process of auto insurance claims. The system leverages Generative AI to handle various aspects of the claims process, from initial photo analysis to damage estimation and repair shop interaction. The goal is to reduce manual effort, accelerate claim resolution times, and improve the customer experience.

### Goals:
- Provide a modern, user-friendly web interface for policyholders to submit and track claims.
- Automate the analysis of vehicle damage from uploaded photos using AI.
- Generate preliminary repair cost estimates using AI agents.
- Facilitate finding and booking appointments with approved repair shops.
- Create a scalable and observable microservices-based architecture.
- Demonstrate the power of a multi-agent system in a real-world business process.

---

## 2. System Architecture

The project is built on a microservices architecture, with distinct services for the frontend, backend, and various AI agent functionalities.

### A. Frontend (Vue.js)
- **Framework:** Vue.js with Vite.
- **API Communication:** Standard HTTP client (e.g., fetch or axios).
- **Key Features:**
    - User authentication (based on policy number).
    - Dashboard to view existing claims.
    - A detailed view for each claim, showing status, photos, and estimates.
    - A wizard-like form for submitting new claims, including photo uploads.
    - Interaction with AI agents for finding repair shops and booking appointments.

### B. Backend (Go)
- **Framework:** Gin.
- **Database:** GORM with a relational database (e.g., PostgreSQL, MySQL, or SQLite for local dev).
- **Responsibilities:**
    - Serves as the main API gateway for the frontend.
    - Handles core business logic for claims and policyholders.
    - Manages data persistence for all primary models.
    - Interacts with Google Cloud Storage (GCS) for photo uploads and retrievals.
    - Orchestrates calls to the various downstream AI services.
    - Provides endpoints for creating, retrieving, and updating claims.

### C. AI Service & Agents (Python)
A collection of specialized AI services, built with FastAPI, uv, and the Google Cloud ADK (A2A protocol).

- **`ai-service`**: The primary AI service (port 8000) that acts as the entrypoint for AI tasks. It leverages YOLOv11 for specific damage detection, and Vertex AI for generative analysis. It also acts as the orchestrator to call other agents.
    - **`claims_agent.py`**: A core agent for processing the overall claim logic.
    - **`car_damage_detector.py`**: An agent specialized in analyzing photos to detect and classify vehicle damage using YOLOv11 and Vertex AI.
    - **`appointment_agent.py`**: An agent that handles the conversational flow for booking appointments with repair shops.
    - **`repair_shop_agent.py`**: An agent responsible for finding and interacting with repair shops.
- **`assessor-agent`**: A remote A2A agent (port 8081) for performing assessment tasks, evaluating the severity of vehicle damage based on visual analysis and claim details.
- **`processor-agent`**: A remote A2A agent (port 8082) focused on generating repair estimates and making final claim decisions based on the initial assessment.
- **`repair-shop-agent`**: A remote A2A agent (port 8083) standalone service for repair shop interactions, finding suitable repair shops and handling appointment bookings.

### D. Load Generator
- **`loadgen`**: A Node.js-based utility for generating synthetic load against the backend API. This is crucial for performance testing, OpenTelemetry traces, and ensuring the system's scalability and reliability under pressure. It is not part of the core production application.

---

## 3. Personas

- **Claimant/Policyholder**: The primary user of the system. They submit claims for their vehicles and use the frontend to track the status and interact with the automated process.
- **Insurance Adjuster**: An employee of the insurance company. They might use the system to review complex cases flagged by the AI, override automated decisions, or monitor the overall claims pipeline.
- **Repair Shop Technician/Manager**: An external user at a partner repair shop. They might interact with the system to receive appointment requests, provide detailed estimates, and update repair statuses.
- **Developer/System Administrator**: Responsible for maintaining, deploying, and monitoring the health of the entire application stack.

---

## 4. User Journeys

### A. Claim Submission
1.  The **Claimant** logs into the web portal using their policy number.
2.  They navigate to "File a New Claim" and fill out a form with details about the incident (date, description, etc.).
3.  They upload multiple photos of the damaged vehicle from their computer or mobile device.
4.  The frontend sends this data to the **Backend**, which creates a new `Claim` record and uploads the photos to GCS.

### B. Automated Analysis & Estimation
1.  Upon successful claim creation, the **Backend** triggers the `AnalyzeClaim` process (`/api/claims/:id/analyze`).
2.  It calls the **`ai-service`** (`/process-claims` endpoint), providing the GCS URIs of the uploaded photos.
3.  The **`car_damage_detector`** within `ai-service` analyzes each photo using YOLOv11/VertexAI, identifying damaged parts and assessing severity.
4.  The `ai-service` orchestrates calls to the **`assessor-agent`** and **`processor-agent`** via the A2A protocol to fully assess the claim and generate a preliminary repair estimate.
5.  The AI service returns the comprehensive analysis results and estimate to the Backend, which persists this information in the `AnalysisResult` and `Estimate` tables.
6.  The claim status is updated to "Assessed", and the **Claimant** can now view the results.

### C. Repair Shop Selection
1.  The **Claimant** views the AI-generated estimate and decides to proceed.
2.  They click "Find Repair Shops". The **Backend** calls the **`ai-service`** (`/find-repair-shops` endpoint).
3.  The **`repair_shop_agent`** finds suitable, certified shops based on the claimant's location and the type of damage.
4.  The **Claimant** can then initiate a conversation with the **`appointment_agent`** to book an appointment.

---

## 5. Data Model

Based on `backend/models/models.go`.

### `Claim`
| Field | Type | Description |
|---|---|---|
| `ID` | uint | Primary Key |
| `PolicyNumber` | string | Foreign key to PolicyHolder |
| `CustomerName` | string | Name of the customer filing the claim |
| `Status` | string | The current stage of the claim (e.g., "New", "Analyzing", "Assessed") |
| `Description` | string | User-provided description of the incident |
| `AccidentDate` | time.Time | Date of the accident |
| `IncidentCity` | string | City where the incident occurred |
| `IncidentState`| string | State where the incident occurred |
| `IncidentType` | string | Type of incident (e.g., "Collision", "Theft") |
| `CollisionType`| string | More specific collision details |
| `Severity` | string | Initial assessment of severity |
| `Photos` | []Photo | Associated photos (one-to-many) |
| `Estimates` | []Estimate | Associated estimates (one-to-many) |

### `Photo`
| Field | Type | Description |
|---|---|---|
| `ID` | uint | Primary Key |
| `ClaimID` | uint | Foreign key to Claim |
| `URL` | string | GCS object path or signed URL for the photo |
| `AnalysisResult` | AnalysisResult | The result of the AI analysis for this photo |

### `AnalysisResult`
| Field | Type | Description |
|---|---|---|
| `ID` | uint | Primary Key |
| `PhotoID` | uint | Foreign key to Photo |
| `QualityScore` | string | Quality of the image (e.g., "Good", "Blurry") |
| `Detections` | string | JSON string of detected damaged parts with bounding boxes |
| `PartsDetected`| string | Comma-separated list of unique parts detected |
| `Severity` | string | Assessed severity from the photo (e.g., "Low", "Medium", "High") |

### `Estimate`
| Field | Type | Description |
|---|---|---|
| `ID` | uint | Primary Key |
| `ClaimID` | uint | Foreign key to Claim |
| `TotalAmount` | float64 | The total estimated cost |
| `Items` | string | JSON string of line items for the repair |
| `Source` | string | Who generated the estimate ("AI", "Shop") |

### `PolicyHolder`
| Field | Type | Description |
|---|---|---|
| `ID` | uint | Primary Key |
| `PolicyNumber` | string | Unique policy number |
| `FirstName` | string | Policyholder's first name |
| `LastName` | string | Policyholder's last name |
| `AutoMake` | string | Make of the insured vehicle |
| `AutoModel` | string | Model of the insured vehicle |
| `AutoYear` | int | Year of the insured vehicle |
| `...` | ... | Other policy and demographic details |

---

## 6. API Specification (Backend)

Key endpoints exposed by the Go backend service.

| Endpoint | Method | Description | Payload/Params |
|---|---|---|---|
| `/api/claims` | GET | List all claims, optionally filtered by `policy_number`. | `?policy_number=...` |
| `/api/claims` | POST | Create a new claim. Expects `multipart/form-data`. | Form data with claim fields and `files`. |
| `/api/claims/:id` | GET | Get details for a single claim. | `id` (path param) |
| `/api/claims/:id` | PATCH | Update a claim's status. | `id` (path param), `{"status": "..."}` |
| `/api/claims/:id` | DELETE | Delete a claim and all its associations. | `id` (path param) |
| `/api/claims/:id/analyze` | POST | Trigger the AI analysis pipeline for a claim. | `id` (path param) |
| `/api/claims/:id/find-shops` | POST | Find repair shops near the policyholder. | `id` (path param) |
| `/api/claims/:id/book-appointment` | POST | Interact with the appointment booking agent. | `id` (path param), Chat-like payload. |
| `/api/policies/:number` | GET | Get details for a policyholder. | `number` (path param) |

---

## 7. Deployment & Execution

### Prerequisites
- Go (v1.25.8+)
- Node.js (v18+) & npm
- Python (v3.11+) and `uv`
- Docker (for containerized deployment)
- Access to a Google Cloud project with GCS enabled.

### Local Development
Each service can be run locally using commands specified in its `Makefile`.

1.  **Agents (A2A):**
    ```bash
    cd assessor-agent && make local-assessor-agent
    cd processor-agent && make local-processor-agent
    cd repair-shop-agent && make local-repair-shop-agent
    ```

2.  **AI Service:**
    ```bash
    cd ai-service
    make local-ai-service
    ```
    *This runs on port 8000.*

3.  **Backend:**
    ```bash
    cd backend
    make local-backend
    ```
    *This runs the Go server on port 8080 by default.*

4.  **Frontend:**
    ```bash
    cd frontend
    make local-frontend
    ```
    *This starts the Vite dev server, typically on port 5173.*

### Cloud Deployment
Deploying to Google Cloud involves three phases:
1. `python3 infra/setup.py` (Foundation Setup)
2. Deploying each service to Cloud Run / Vertex AI (`gcloud builds submit --config .cloudbuild/deploy.yaml .`)
3. `python3 infra/setup_lb.py` (Setup Load Balancer)
