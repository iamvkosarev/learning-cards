LAST_TAG = $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=last)
REMOTE_IMAGE = $(DOCKER_REPO):$(LAST_TAG)
REMOTE_DIR := /root/$(PROJECT_NAME)

vps-deploy:
	@echo "Uploading .env, docker-compose.yml and config/ to VPS..."
	ssh $(VPS_USER)@$(VPS_HOST) "mkdir -p $(REMOTE_DIR)/config"
	scp .env docker-compose.yml $(VPS_USER)@$(VPS_HOST):$(REMOTE_DIR)/
	scp -r config/* $(VPS_USER)@$(VPS_HOST):$(REMOTE_DIR)/config/

	@echo "🛠️  Injecting IMAGE_NAME=$(REMOTE_IMAGE) into .env..."
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(REMOTE_DIR) && \
		grep -v '^IMAGE_NAME=' .env > .env.tmp || true && \
		echo 'IMAGE_NAME=$(REMOTE_IMAGE)' >> .env.tmp && \
		mv .env.tmp .env"

	@echo "Deploying image $(REMOTE_IMAGE) on VPS..."
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(REMOTE_DIR) && \
		docker compose pull && \
		docker compose up -d --remove-orphans"

	@echo "Deployed $(REMOTE_IMAGE) to $(VPS_USER)@$(VPS_HOST)"