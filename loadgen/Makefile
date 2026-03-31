PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= us-central1
_COMMIT_SHA ?= $(shell echo | find api-load-generator/package.json api-load-generator/package-lock.json api-load-generator/lib api-load-generator/config Dockerfile -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_TRACES_EXPORTER=otlp
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
OTEL_SERVICE_NAME=auto-claims-loadgen


.PHONY: loadgen local

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER} ;\
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;\
	export OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME} ;\
	cd api-load-generator && npm install && npm start

loadgen:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${REGION} --project=${PROJECT_ID} --quiet
