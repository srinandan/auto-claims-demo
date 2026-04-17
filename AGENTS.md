# AGENTS.md — Auto Claims Demo

This file describes the AI agents in this repository, their roles, how to run them, and how they interact with each other. It also documents the available skills that can be loaded by an AI coding agent working in this codebase.

---

## Table of Contents

1. [Agent Overview](#agent-overview)
2. [Agent Details](#agent-details)
   - [AI Service (Orchestrator)](#1-ai-service-orchestrator)
   - [Assessor Agent](#2-assessor-agent)
   - [Processor Agent](#3-processor-agent)
   - [Repair Shop Agent](#4-repair-shop-agent)
3. [Agent Communication (A2A Protocol)](#agent-communication-a2a-protocol)
4. [Running Agents Locally](#running-agents-locally)
5. [Skill: Auto Claims Reviewer](#skill-auto-claims-reviewer)
6. [Adding a New Agent](#adding-a-new-agent)

---

## Agent Overview

The system uses a **multi-agent architecture** built on the Google Cloud ADK A2A (Agent-to-Agent) protocol. A central AI service acts as an orchestrator, delegating specialized tasks to downstream remote agents.

```
Frontend (Vue)
     │
     ▼
Backend (Go/Gin)
     │
     ▼
AI Service / Orchestrator  ──────────────────────────┐
  ├── car_damage_detector.py  (YOLOv11 + Vertex AI)   │
  ├── claims_agent.py                                  │
  ├── appointment_agent.py                             │
  └── repair_shop_agent.py                             │
           │ A2A                                       │
     ┌─────┴──────┬─────────────────┐                 │
     ▼            ▼                 ▼                  │
assessor-agent  processor-agent  repair-shop-agent ◄──┘
  (port 8081)   (port 8082)        (port 8083)
     │            │                 │
     └────────────┼─────────────────┘
                  ▼
            agent-engine (Shared Session & Memory)
```

---

## Agent Details

### 1. AI Service (Orchestrator)

| Property | Value |
|---|---|
| **Location** | `./ai-service/` |
| **Runtime** | Python 3.11+ / FastAPI |
| **Default Port** | `8000` |
| **Role** | Primary entrypoint for all AI tasks. Orchestrates calls to downstream A2A agents. |

**Internal Modules:**

| Module | Responsibility |
|---|---|
| `claims_agent.py` | Core logic for processing the overall claim lifecycle. |
| `car_damage_detector.py` | Analyzes uploaded photos using **YOLOv11** and **Vertex AI** to detect damaged vehicle parts and classify severity. |
| `appointment_agent.py` | Manages the conversational flow for booking repair appointments. |
| `repair_shop_agent.py` | Finds and interacts with certified partner repair shops based on location and damage type. |

**Key Endpoints:**

| Endpoint | Description |
|---|---|
| `POST /process-claims` | Receives GCS photo URIs from the backend and runs the full damage analysis + estimation pipeline. |
| `POST /find-repair-shops` | Finds suitable repair shops based on claimant location and damage profile. |

---

### 2. Assessor Agent

| Property | Value |
|---|---|
| **Location** | `./assessor-agent/` |
| **Runtime** | Python / FastAPI + ADK A2A |
| **Default Port** | `8081` |
| **Role** | Remote A2A agent. Evaluates the **severity** of vehicle damage based on photo analysis results passed from the AI service. |

**Input:** Structured damage detection data (parts detected, quality scores, bounding boxes) from `car_damage_detector`.

**Output:** A severity classification (`Low`, `Medium`, `High`, `Critical`) used to inform the processor agent.

---

### 3. Processor Agent

| Property | Value |
|---|---|
| **Location** | `./processor-agent/` |
| **Runtime** | Python / FastAPI + ADK A2A |
| **Default Port** | `8082` |
| **Role** | Remote A2A agent. Generates **repair cost estimates** and makes a preliminary claim decision based on the assessor agent's output. |

**Input:** Severity assessment from the assessor agent and parts detection data.

**Output:** An `Estimate` object containing:
- `TotalAmount` — the overall estimated repair cost.
- `Items` — a JSON array of line-item repairs (e.g., `{ "part": "rear bumper", "action": "replace", "cost": 850 }`).
- A preliminary claim decision recommendation.

---

### 4. Repair Shop Agent

| Property | Value |
|---|---|
| **Location** | `./repair-shop-agent/` |
| **Runtime** | Python / FastAPI + ADK A2A |
| **Default Port** | `8083` |
| **Role** | Remote A2A agent. Standalone service for repair shop discovery and appointment booking. |

**Triggered by:** The `appointment_agent.py` inside the AI service.

**Capabilities:**
- Find certified repair shops near the claimant.
- Handle appointment scheduling requests.
- Relay booking confirmations back through the AI service to the backend.

---

### 5. Agent Engine

| Property | Value |
|---|---|
| **Location** | `./agent-engine/` |
| **Runtime** | Python / Agent Engine Instance |
| **Role** | Shared session and memory provider. |

**Key Responsibilities:**
- Provides a centralized service for session state management.
- Maintains shared long-term memory across all agents in the system.
- Ensures consistency in the claim review and processing workflow by allowing agents to access a unified context.

---

## Agent Communication (A2A Protocol)

All remote agents (`assessor-agent`, `processor-agent`, `repair-shop-agent`) communicate with the `ai-service` orchestrator using **Google Cloud ADK's A2A (Agent-to-Agent) protocol**.

Key conventions:
- Each agent exposes a standard A2A-compliant HTTP interface.
- The orchestrator (`ai-service`) is the only caller — agents do not call each other directly.
- Data is passed as structured JSON payloads between agents.
- The full conversation/state context is included in each request, as agents are stateless between calls.
- **Security:** Remote agents are configured to use an **Agent Gateway** for governed access and secure egress in cloud deployments.

---

## Running Agents Locally

Start all services at once from the repo root:

```bash
make local-all
```

Or start each agent individually:

```bash
# Remote A2A Agents
cd assessor-agent && make local-assessor-agent       # → http://localhost:8081
cd processor-agent && make local-processor-agent     # → http://localhost:8082
cd repair-shop-agent && make local-repair-shop-agent # → http://localhost:8083

# AI Service (Orchestrator)
cd ai-service && make local-ai-service               # → http://localhost:8000

# Backend (calls the AI service)
cd backend && make local-backend                     # → http://localhost:8080

# Frontend
cd frontend && make local-frontend                   # → http://localhost:5173
```

Stop all services:

```bash
make local-stop
```

**Prerequisites:** Python 3.11+, `uv`, Go 1.25.8+, Node.js 18+, and a Google Cloud project with Vertex AI enabled.

---

## Skill: Auto Claims Reviewer

The `.agents/skills/auto-claims-reviewer/SKILL.md` file defines an **AI agent persona and operating procedure** for automated claim review. An AI coding agent or LLM-powered workflow should load this skill when tasked with reviewing a claim end-to-end.

### What the Skill Does

It instructs an AI agent to act as a senior claims adjuster and execute a four-step workflow:

| Step | Action |
|---|---|
| **1. Ingest & Verify** | Retrieve claim and policy data from the backend API. Cross-reference vehicle info, customer name, and policy terms. |
| **2. Analyze Photos** | Review all photos, evaluate AI quality scores, correlate YOLO detections with the claimant's description. |
| **3. Assess Estimate** | Validate the processor agent's repair estimate line items against the detected parts and flag anomalous costs. |
| **4. Fraud Analysis** | Look for inconsistencies: mismatched damage, prior damage, suspicious timelines, or exaggerated claims. |

### Output Format

The skill produces a structured **Claim Review Summary** containing:
- Overall assessment and confidence score (`0.0` – `1.0`).
- Bulleted findings.
- Flags with severity (`Low` / `Medium` / `High` / `Critical`) and a recommended action for the human adjuster.

### When to Load This Skill

Load `.agents/skills/auto-claims-reviewer/SKILL.md` whenever an agent needs to:
- Perform an automated end-to-end review of a single claim.
- Produce a structured report for a human adjuster.
- Identify fraud indicators or data inconsistencies in a claim record.

---

## Adding a New Agent

1. Create a new directory at the repo root (e.g., `./my-new-agent/`).
2. Implement a FastAPI service that conforms to the **ADK A2A protocol** interface.
3. Add a `Makefile` with a `local-my-new-agent` target that starts the service on an unused port.
4. Register the agent's local URL as an environment variable in `ai-service` so the orchestrator can call it.
5. Add a corresponding entry in the root `Makefile` under `local-all` and `local-stop`.
6. Document the new agent in this file following the structure above.
7. If the agent introduces a new reviewable workflow, consider adding a corresponding skill under `.agents/skills/`.
