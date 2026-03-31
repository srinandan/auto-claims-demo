PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)
BUCKET_NAME ?= $(shell echo ${PROJECT_ID}-auto-claims)
SERVICE_ACCOUNT_EMAIL ?= $(shell echo "auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com")
_COMMIT_SHA ?= $(shell echo | find go.mod go.sum database models handlers main.go -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)
AI_SERVICE_URL ?= http://localhost:8000

.PHONY: local deploy

local:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export BUCKET_NAME=${BUCKET_NAME} ;\
	export SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT_EMAIL} ;\
	export AI_SERVICE_URL=${AI_SERVICE_URL} ;\
	go mod tidy && go run main.go

deploy:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	gcloud builds submit --config=.cloudbuild/deploy.yaml --region=${REGION} --substitutions=_COMMIT_SHA=${_COMMIT_SHA},_REGION=${REGION}

clean:
	gcloud run services delete auto-claims-$$(basename $$PWD) --region=${REGION} --project=${PROJECT_ID} --quiet
