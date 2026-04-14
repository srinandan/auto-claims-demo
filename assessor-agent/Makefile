PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)
_COMMIT_SHA ?= $(shell echo | find app pyproject.toml uv.lock -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
PORT ?= 8081
REASONING_ENGINE_ID ?= 3024167149779288064
SHARED_AGENT_ENGINE_ID ?= 3064462601669836800

.PHONY: local deploy

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_LOCATION=${REGION} ;\
	export BQ_ANALYTICS_GCS_BUCKET=${PROJECT_ID} ;\
	export BQ_ANALYTICS_CONNECTION_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export LOGS_BUCKET_NAME=${PROJECT_ID} ;\
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;\
	export PORT=${PORT} ;\
	uv run uvicorn app.fast_api_app:app --host localhost --port ${PORT} --reload

lint:
	uv run codespell
	uv run ruff check . --diff
	uv run ruff format .
	uv run mypy .

deploy:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_LOCATION=${REGION} ;\
	export BQ_ANALYTICS_GCS_BUCKET=${PROJECT_ID} ;\
	export BQ_ANALYTICS_CONNECTION_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export LOGS_BUCKET_NAME=${PROJECT_ID} ;\
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;\
	export REASONING_ENGINE_ID=${REASONING_ENGINE_ID} ;\
	export SHARED_AGENT_ENGINE_ID=${SHARED_AGENT_ENGINE_ID} ;\
	(uv export --no-hashes --no-header --no-dev --no-emit-project --no-annotate > app/app_utils/.requirements.txt 2>/dev/null || \
	uv export --no-hashes --no-header --no-dev --no-emit-project > app/app_utils/.requirements.txt) && \
	uv run -m app.app_utils.deploy \
		--source-packages=./app \
		--entrypoint-module=app.agent_engine_app \
		--entrypoint-object=agent_engine \
		--location=${REGION} \
		--agent-identity \
		--set-env-vars=OTEL_SERVICE_NAME=auto-claims-assessor-agent,OTEL_PYTHON_LOGGING_AUTO_INSTRUMENTATION_ENABLED=true,REASONING_ENGINE_ID=${REASONING_ENGINE_ID},SHARED_AGENT_ENGINE_ID=${SHARED_AGENT_ENGINE_ID},GOOGLE_CLOUD_AGENT_ENGINE_ENABLE_TELEMETRY=true \
		--requirements-file=app/app_utils/.requirements.txt

