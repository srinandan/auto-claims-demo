PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= us-central1

.PHONY: infra

infra:
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export GOOGLE_CLOUD_REGION=${REGION} ;\
	gcloud storage buckets create gs://${PROJECT_ID}-auto-claims --location=${REGION} --uniform-bucket-level-access --public-access-prevention --default-storage-class=STANDARD --retention-period=30d ;\
	gcloud iam service-accounts create auto-claims-sa --display-name="auto-claims-sa" --project=${PROJECT_ID} ;\
	gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/apphub.editor" ;\
	gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/artifactregistry.admin" ;\
	gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/run.admin" ;\
	gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:auto-claims-sa@${PROJECT_ID}.iam.gserviceaccount.com" --role="roles/storage.admin"

