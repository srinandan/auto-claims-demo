.PHONY: all deploy frontend backend ai-service assessor-agent processor-agent repair-shop-agent loadgen local-all local-stop

all: deploy

deploy:
	$(MAKE) -j6 frontend backend ai-service assessor-agent processor-agent repair-shop-agent

local-all:
	@echo "Starting all services locally..."
	$(MAKE) -C assessor-agent local & \
	$(MAKE) -C processor-agent local & \
	$(MAKE) -C repair-shop-agent local & \
	$(MAKE) -C ai-service local & \
	$(MAKE) -C backend local & \
	$(MAKE) -C frontend local & \
	$(MAKE) -C loadgen local & \
	wait

local-stop:
	@echo "Stopping all local services..."
	-pkill -f "uvicorn"
	-pkill -f "uv run main.py"
	-pkill -f "go run main.go"
	-pkill -f "vite"
	-pkill -f "node"


frontend:
	$(MAKE) -C frontend deploy

backend:
	$(MAKE) -C backend deploy

ai-service:
	$(MAKE) -C ai-service deploy

assessor-agent:
	$(MAKE) -C assessor-agent deploy

processor-agent:
	$(MAKE) -C processor-agent deploy

repair-shop-agent:
	$(MAKE) -C repair-shop-agent deploy
