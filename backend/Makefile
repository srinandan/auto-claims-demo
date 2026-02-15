PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= us-central1
BUCKET_NAME ?= $(shell echo ${PROJECT_ID}-auto-claims)
# SERVICE_ACCOUNT_EMAIL ?= $(shell echo "movie-guru-chat-server-sa@${PROJECT_ID}.iam.gserviceaccount.com")
SERVICE_ACCOUNT_EMAIL ?= $(shell echo "auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com")
_COMMIT_SHA ?= $(shell echo | find go.mod go.sum database models handlers main.go -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum | cut -c 1-10)

.PHONY: local-backend
local-backend:
	export _COMMIT_SHA=${_COMMIT_SHA} ;\
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	export BUCKET_NAME=${BUCKET_NAME} ;\
	export SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT_EMAIL} ;\
	go mod tidy && go run main.go