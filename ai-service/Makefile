PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)
BUCKET_NAME ?= $(shell echo ${PROJECT_ID}-auto-claims)
_COMMIT_SHA ?= $(shell echo | find uv.lock *.py -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
ASSESSOR_AGENT_URL ?= http://localhost:8081
PROCESSOR_AGENT_URL ?= http://localhost:8082
REPAIR_SHOP_AGENT_URL ?= http://localhost:8083
PORT ?= 8000
AGENT_ENGINE_ENABLED ?= false

.PHONY: local base deploy

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export BUCKET_NAME=${BUCKET_NAME} ;\
	export MOCK_MODE=false ;\
	export ASSESSOR_AGENT_URL=${ASSESSOR_AGENT_URL} ;\
	export PROCESSOR_AGENT_URL=${PROCESSOR_AGENT_URL} ;\
	export REPAIR_SHOP_AGENT_URL=${REPAIR_SHOP_AGENT_URL} ;\
	export OTEL_INSTRUMENTATION_GENAI_CAPTURE_MESSAGE_CONTENT=NO_CONTENT ;\
	export PORT=${PORT} ;\
	export AGENT_ENGINE_ENABLED=${AGENT_ENGINE_ENABLED} ;\
	uv run main.py

deploy:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}

base:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit base/ --config=base/base.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${REGION} --project=${PROJECT_ID} --quiet
