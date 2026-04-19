GOOGLE_CLOUD_PROJECT ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${GOOGLE_CLOUD_PROJECT} --format="value(projectNumber)")
GOOGLE_CLOUD_LOCATION ?= $(shell gcloud config get-value compute/region)
BUCKET_NAME ?= $(shell echo ${GOOGLE_CLOUD_PROJECT}-auto-claims)
SERVICE_ACCOUNT_EMAIL ?= $(shell echo "auto-claims-sa@${GOOGLE_CLOUD_PROJECT}.iam.gserviceaccount.com")
_COMMIT_SHA ?= $(shell echo | find go.mod go.sum database models handlers main.go -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
AI_SERVICE_URL ?= http://localhost:8000

.PHONY: local deploy local-backend

local-backend: local

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${GOOGLE_CLOUD_LOCATION} ;\
	export BUCKET_NAME=${BUCKET_NAME} ;\
	export SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT_EMAIL} ;\
	export AI_SERVICE_URL=${AI_SERVICE_URL} ;\
	go mod tidy && go run main.go

deploy:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${GOOGLE_CLOUD_LOCATION} ;\
	export PROJECT_ID=${GOOGLE_CLOUD_PROJECT} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${GOOGLE_CLOUD_LOCATION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${GOOGLE_CLOUD_LOCATION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${GOOGLE_CLOUD_LOCATION} --project=${GOOGLE_CLOUD_PROJECT} --quiet
