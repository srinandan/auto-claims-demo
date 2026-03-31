.PHONY: all deploy frontend backend ai-service assessor-agent processor-agent repair-shop-agent

all: deploy

deploy:
	$(MAKE) -j6 frontend backend ai-service assessor-agent processor-agent repair-shop-agent

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
