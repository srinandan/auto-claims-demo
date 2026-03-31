PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)
_COMMIT_SHA ?= $(shell echo | find public src package.json package-lock.json -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
OTEL_TRACES_EXPORTER=otlp
OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
VITE_OTLP_EXPORTER_URL=https://telemetry.googleapis.com/v1/traces


.PHONY: local deploy

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER} ;\
	export OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT} ;\
	export VITE_OTLP_EXPORTER_URL=${VITE_OTLP_EXPORTER_URL} ;\
	npm install && npm run dev

deploy:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${REGION} --project=${PROJECT_ID} --quiet
