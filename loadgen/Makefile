GOOGLE_CLOUD_PROJECT ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${GOOGLE_CLOUD_PROJECT} --format="value(projectNumber)")
GOOGLE_CLOUD_LOCATION ?= us-central1
_COMMIT_SHA ?= $(shell echo | find api-load-generator/package.json api-load-generator/package-lock.json api-load-generator/lib api-load-generator/config Dockerfile -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_TRACES_EXPORTER=otlp
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
OTEL_SERVICE_NAME=auto-claims-loadgen


.PHONY: loadgen local local-loadgen

local-loadgen: local

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${GOOGLE_CLOUD_LOCATION} ;\
	export OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER} ;\
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;\
	export OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME} ;\
	cd api-load-generator && npm install && npm start

loadgen:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${GOOGLE_CLOUD_LOCATION} ;\
	export PROJECT_ID=${GOOGLE_CLOUD_PROJECT} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${GOOGLE_CLOUD_LOCATION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${GOOGLE_CLOUD_LOCATION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${GOOGLE_CLOUD_LOCATION} --project=${GOOGLE_CLOUD_PROJECT} --quiet
