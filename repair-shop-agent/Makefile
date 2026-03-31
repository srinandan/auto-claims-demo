PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)
_COMMIT_SHA ?= $(shell echo | find app pyproject.toml uv.lock -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
PORT ?= 8083
REASONING_ENGINE_ID ?= 147311468377800704

.PHONY: local deploy

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;
	export REGION=${REGION} ;
	export PROJECT_ID=${PROJECT_ID} ;
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;
	export GOOGLE_CLOUD_LOCATION=${REGION} ;
	export BQ_ANALYTICS_GCS_BUCKET=${PROJECT_ID} ;
	export BQ_ANALYTICS_CONNECTION_ID=${PROJECT_ID} ;
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;
	export LOGS_BUCKET_NAME=${PROJECT_ID} ;
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;
	(uv export --no-hashes --no-header --no-dev --no-emit-project --no-annotate > app/app_utils/.requirements.txt 2>/dev/null || \
	uv export --no-hashes --no-header --no-dev --no-emit-project > app/app_utils/.requirements.txt) && \
	uv run -m app.app_utils.deploy \
		--source-packages=./app \
		--entrypoint-module=app.agent_engine_app \
		--entrypoint-object=agent_engine \
		--agent-identity \
		--location=${REGION} \
		--set-env-vars=REASONING_ENGINE_ID=${REASONING_ENGINE_ID},GOOGLE_CLOUD_AGENT_ENGINE_ENABLE_TELEMETRY=true \
		--requirements-file=app/app_utils/.requirements.txt

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
	(uv export --no-hashes --no-header --no-dev --no-emit-project --no-annotate > app/app_utils/.requirements.txt 2>/dev/null || \
	uv export --no-hashes --no-header --no-dev --no-emit-project > app/app_utils/.requirements.txt) && \
	uv run -m app.app_utils.deploy \
		--source-packages=./app \
		--entrypoint-module=app.agent_engine_app \
		--entrypoint-object=agent_engine \
		--agent-identity \
		--location=${REGION} \
		--set-env-vars=REASONING_ENGINE_ID=${REASONING_ENGINE_ID},GOOGLE_CLOUD_AGENT_ENGINE_ENABLE_TELEMETRY=true \
		--requirements-file=app/app_utils/.requirements.txt
