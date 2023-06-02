SHELL := /bin/bash

VERSION := 0.0.1
BUILD_INFO := Manual build

ENV_FILE := .env
ifeq ($(filter $(MAKECMDGOALS),config clean),)
	ifneq ($(strip $(wildcard $(ENV_FILE))),)
		ifneq ($(MAKECMDGOALS),config)
			include $(ENV_FILE)
			export
		endif
	endif
endif

.PHONY: help lint deploy run-api-image pushimages images runapiops all 
.DEFAULT_GOAL := help

ACR_FQDN := ${ACR_NAME}.azurecr.io

all: 
	@make deploy
	@make buildimage
	@make pushimage

help: ## 💬 This help message
	@grep -E '[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deploy: ## 🝡 Deploy azure resources
	@bash ./deploy.sh

buildimage: ## 🐳 Build docker image
	@docker build -t $(ACR_FQDN)/${ACR_REPO_NAME}:$(IMAGE_TAG) -f src/identityApp/Dockerfile .

runimage: ## 🐳 Run docker image
	@docker run -it --rm -p 8080:80 $(ACR_FQDN)/${ACR_REPO_NAME}:$(IMAGE_TAG)

pushimage: ## 🐳 Push docker image
	az acr login --name ${ACR_NAME} && docker push $(ACR_FQDN)/${ACR_REPO_NAME}:$(IMAGE_TAG)
