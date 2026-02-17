PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= us-central1
_COMMIT_SHA ?= $(shell echo | find app pyproject.toml uv.lock -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
PORT ?= 8082

.PHONY: processor-agent local-processor-agent

processor-agent:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}


local-processor-agent:
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