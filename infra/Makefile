PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
REGION ?= $(shell gcloud config get-value compute/region)

.PHONY: infra lb base

infra:
	@echo "Setting up foundation..."
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_LOCATION=${REGION} ;\
	python3 setup.py

lb:
	@echo "Setting up Global Application Load Balancer..."
	export GOOGLE_CLOUD_PROJECT=${PROJECT_ID} ;\
	export GOOGLE_CLOUD_LOCATION=${REGION} ;\
	python3 setup_lb.py

base:
	@echo "Building auto-claims-ai-service-base image..."
	export REGION=${REGION} ;\
	export PROJECT_ID=${PROJECT_ID} ;\
	export PROJECT_NUMBER=${PROJECT_NUMBER} ;\
	export LOCATION=${REGION} ;\
	gcloud builds submit base/ --config=base/base.yaml --region=${REGION} --substitutions=_REGION=${REGION}
	@echo "\n=================================================================================="
	@echo "IMPORTANT: Don't forget to update the base image path in ai-service/Dockerfile"
	@echo "to use the correct project ID before running 'make ai-service'!"
	@echo "==================================================================================\n"

