# Skill: Coding Workflow

When invoked, execute the following steps in order. Do not skip steps or ask for confirmation between them unless a step fails.

---

## Step 0 — Plan Before Coding

Before writing any code, ensure you have a clear plan:
- **Think Before Coding**: Avoid hidden assumptions. Surface tradeoffs and plan implementation thoroughly.
- **Simplicity First**: Aim for the simplest solution that works. Avoid over-engineering.
- **Define Success Criteria**: Know how you will verify that your changes work.

---

## Step 1 — Inspect Uncommitted Changes

```bash
git status
git diff
```

Identify which service directories contain modified files. All subsequent steps operate only on affected services.

---

## Step 2 — Lint Affected Services

Run the linter for each service that has uncommitted changes.

**Python services** (`ai-service`, `assessor-agent`, `processor-agent`, `repair-shop-agent`):
```bash
cd <service-dir> && make lint
# runs: codespell, ruff check, ruff format, mypy
```

**Go backend**:
```bash
cd backend && go vet ./...
```

**Frontend**:
```bash
cd frontend && npm run build
```

Fix any errors before proceeding. Do not continue to Step 3 with lint failures outstanding.

---

## Step 3 — Check for Unused Packages

**Python services** — check for unused imports:
```bash
cd <service-dir> && uv run ruff check . --select F401
```

**Go backend** — remove unused module dependencies:
```bash
cd backend && go mod tidy
```

If `go mod tidy` modifies `go.mod` or `go.sum`, stage those changes.

---

## Step 4 — Run Tests

**Go backend**:
```bash
cd backend && go test ./...
```

**Python AI service**:
```bash
cd ai-service && uv run pytest tests/
```

All tests must pass before proceeding. If a test fails, fix it and re-run.

---

## Step 5 — Create a Feature Branch

Choose a branch name that reflects the change. Use the pattern `<type>/<short-description>`, e.g.:
- `feat/add-fraud-score-field`
- `fix/repair-shop-null-response`
- `refactor/claims-agent-retry-logic`

```bash
git checkout -b <branch-name>
```

---

## Step 6 — Commit and Push

Stage only the files relevant to the change (avoid `git add -A` blindly). Emphasize **Surgical Changes** — make precise, localized edits.

```bash
git add <specific files>
git commit -m "<type>: <concise description of what changed and why>"
git push -u origin <branch-name>
```

Commit message format: `<type>: <description>` where type is one of `feat`, `fix`, `refactor`, `test`, `docs`, `chore`.

---

## Step 7 — Open a Pull Request

Use the GitHub MCP tools (`mcp__github__create_pull_request`) to open a PR against `main`.

PR title: same format as the commit message — `<type>: <short description>`.

PR body must include:
- **Summary**: 2–4 bullet points describing what changed and why.
- **Services affected**: list of service directories touched.
- **Test plan**: what was run and what passed.

---

## Step 8 — Code Coverage and Unit Tests Review

The reviewer must examine code coverage reports and ensure that all new code is accompanied by corresponding unit tests. Verify that:
- Coverage metrics are collected and ideally meet the required threshold (80%).
- All unit tests run cleanly without network dependency issues (by utilizing mocks).

Ensure **Goal-Driven Execution** — verify that the success criteria defined in Step 0 are met.
