CARDS_REMOTE_DIR := /root/$(CARDS_PROJECT_NAME)
REVIEWS_REMOTE_DIR := /root/$(REVIEWS_PROJECT_NAME)

vps-get-containers:
	@ssh $(VPS_USER)@$(VPS_HOST) "docker ps --format '{{ .ID}}\t{{.Names}}'"

vps-deploy-cards:
	@echo "Uploading .env, docker-compose.yml and config.yaml to VPS..."
	ssh $(VPS_USER)@$(VPS_HOST) "mkdir -p $(CARDS_REMOTE_DIR)/config"
	scp cmd/cards/.env  $(VPS_USER)@$(VPS_HOST):$(CARDS_REMOTE_DIR)/
	scp cmd/cards/docker-compose.yml $(VPS_USER)@$(VPS_HOST):$(CARDS_REMOTE_DIR)/
	scp cmd/cards/config.yaml $(VPS_USER)@$(VPS_HOST):$(CARDS_REMOTE_DIR)/
	@echo "Injecting IMAGE_NAME=$(CARDS_IMAGE_NAME) into .env..."
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(CARDS_REMOTE_DIR) && \
		grep -v '^IMAGE_NAME=' .env > .env.tmp || true && \
		echo 'IMAGE_NAME=$(CARDS_IMAGE_NAME)' >> .env.tmp && \
		mv .env.tmp .env"
	ssh $(VPS_USER)@$(VPS_HOST) "\
		docker network inspect iamvkosarev_network >/dev/null 2>&1 || \
		docker network create iamvkosarev_network"
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(CARDS_REMOTE_DIR) && \
		docker compose down && \
		docker compose pull && \
		docker compose up -d --remove-orphans"
	@echo "Deployed $(CARDS_IMAGE_NAME) to $(VPS_USER)@$(VPS_HOST)"


vps-deploy-reviews:
	@echo "Uploading .env, docker-compose.yml and config.yaml to VPS..."
	ssh $(VPS_USER)@$(VPS_HOST) "mkdir -p $(REVIEWS_REMOTE_DIR)/config"
	scp cmd/reviews/.env  $(VPS_USER)@$(VPS_HOST):$(REVIEWS_REMOTE_DIR)/
	scp cmd/reviews/docker-compose.yml $(VPS_USER)@$(VPS_HOST):$(REVIEWS_REMOTE_DIR)/
	scp cmd/reviews/config.yaml $(VPS_USER)@$(VPS_HOST):$(REVIEWS_REMOTE_DIR)/
	@echo "Injecting IMAGE_NAME=$(REVIEWS_IMAGE_NAME) into .env..."
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(REVIEWS_REMOTE_DIR) && \
		grep -v '^IMAGE_NAME=' .env > .env.tmp || true && \
		echo 'IMAGE_NAME=$(REVIEWS_IMAGE_NAME)' >> .env.tmp && \
		mv .env.tmp .env"
	ssh $(VPS_USER)@$(VPS_HOST) "\
		docker network inspect iamvkosarev_network >/dev/null 2>&1 || \
		docker network create iamvkosarev_network"
	ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(REVIEWS_REMOTE_DIR) && \
		docker compose down && \
		docker compose pull && \
        docker compose up -d --remove-orphans"
	@echo "Deployed $(REVIEWS_IMAGE_NAME) to $(VPS_USER)@$(VPS_HOST)"

docker-cards-logs:
	@ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(CARDS_REMOTE_DIR) && \
		docker compose logs -f"

docker-reviews-logs:
	@ssh $(VPS_USER)@$(VPS_HOST) "\
		cd $(REVIEWS_REMOTE_DIR) && \
		docker compose logs -f"