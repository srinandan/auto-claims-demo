PROJECT_ID ?= $(shell gcloud config get-value project)
PROJECT_NUMBER ?= $(shell gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")

.PHONY: infra lb

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

