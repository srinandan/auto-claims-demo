# AGENTS.md — Developer Guide for AI Coding Agents

This file is the primary reference for any AI coding agent working in this repository. It covers where things live, how to run and test each service, and how to extend the system. For the full product and architecture specification, see [`plan/spec.md`](plan/spec.md).

---

## Repo Layout

| Directory | Language | Role |
|---|---|---|
| `frontend/` | Vue.js / Vite | Web UI — policy login, claim submission, status dashboard |
| `backend/` | Go / Gin | API gateway, business logic, GCS uploads, DB persistence |
| `ai-service/` | Python / FastAPI | AI orchestrator — YOLO damage detection, Vertex AI analysis |
| `assessor-agent/` | Python / FastAPI + ADK | Remote A2A agent — damage severity assessment (port 8081) |
| `processor-agent/` | Python / FastAPI + ADK | Remote A2A agent — repair cost estimation (port 8082) |
| `repair-shop-agent/` | Python / FastAPI + ADK | Remote A2A agent — shop lookup & appointment booking (port 8083) |
| `agent-engine/` | Python | Shared session & memory store for all agents |
| `loadgen/` | Node.js | Synthetic load generator for perf testing (not production) |
| `infra/` | Python | Cloud infrastructure setup scripts |
| `plan/` | Markdown | Product & architecture specification |
| `.agents/skills/` | Markdown | Coding workflow skills for AI agents |

---

## Running Services Locally

Start everything at once from the repo root:

```bash
make local-all
```

Stop everything:

```bash
make local-stop
```

Or start services individually:

```bash
cd assessor-agent   && make local   # → http://localhost:8081
cd processor-agent  && make local   # → http://localhost:8082
cd repair-shop-agent && make local  # → http://localhost:8083
cd ai-service       && make local   # → http://localhost:8000
cd backend          && make local   # → http://localhost:8080
cd frontend         && make local   # → http://localhost:5173
```

**Prerequisites:** Go 1.25.8+, Python 3.11+, `uv`, Node.js 18+, a Google Cloud project with Vertex AI and GCS enabled.

---

## Linting

### Python services (`ai-service`, `assessor-agent`, `processor-agent`, `repair-shop-agent`)

Each service has a `lint` target that runs `codespell`, `ruff`, and `mypy`:

```bash
cd <service-dir> && make lint
```

To run only unused-import checks:

```bash
cd <service-dir> && uv run ruff check . --select F401
```

### Go backend

```bash
cd backend && go vet ./...
cd backend && go mod tidy   # removes unused dependencies
```

### Frontend

```bash
cd frontend && npm run build   # catches type errors and broken imports
```

---

## Running Tests

### Go backend

```bash
cd backend && go test ./...
```

### Python AI service

```bash
cd ai-service && uv run pytest tests/
```

No automated test suites exist yet for the A2A agents or frontend — add them under `<service>/tests/` following the `ai-service` pattern.

---

## Adding a New Agent

1. Create a new directory at the repo root (e.g., `./my-new-agent/`).
2. Implement a FastAPI service conforming to the **ADK A2A protocol** interface.
3. Add a `Makefile` with a `local` target that starts the service on an unused port.
4. Add a `lint` target using the same `codespell` / `ruff` / `mypy` stack as the other agents.
5. Register the agent's local URL as an env var in `ai-service` so the orchestrator can reach it.
6. Add entries to the root `Makefile` under `local-all` and `local-stop`.
7. Document the new agent in [`plan/spec.md`](plan/spec.md) under the System Architecture section.

---

## Coding Workflow Skill

Before opening a pull request, load and follow the coding workflow skill:

```
.agents/skills/auto-claims-reviewer/SKILL.md
```

It automates: linting all affected services → checking for unused packages → running tests → creating a branch → pushing → opening a PR.
