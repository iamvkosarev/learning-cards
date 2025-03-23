VPS_USER ?= root
VPS_HOST ?= kosarev.app
VPS_NAME ?= learning-cards

LAST_TAG := $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=last)
REMOTE_IMAGE ?= $(DOCKER_REPO):$(LAST_TAG)

deploy_vps:
	@ssh $(VPS_USER)@$(VPS_HOST) "\
		docker pull $(REMOTE_IMAGE) && \
		(docker stop $(VPS_NAME) || true) && \
		(docker rm $(VPS_NAME) || true) && \
		docker run -d --name $(VPS_NAME) \
			-p 8080:8080 -p 50051:50051 \
			$(REMOTE_IMAGE)"
	@echo "Deployed $(REMOTE_IMAGE) to $(VPS_USER)@$(VPS_HOST)"