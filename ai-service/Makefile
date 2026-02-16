PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= us-central1
BUCKET_NAME ?= $(shell echo ${PROJECT_ID}-auto-claims)
_COMMIT_SHA ?= $(shell echo | find uv.lock main.py -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)

ASSESSOR_AGENT_URL ?= http://localhost:8081

.PHONY: local-ai-service

local-ai-service:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export BUCKET_NAME=${BUCKET_NAME} ;\
	export MOCK_MODE=false ;\
	export ASSESSOR_AGENT_URL=${ASSESSOR_AGENT_URL} ;\
	uv run main.py