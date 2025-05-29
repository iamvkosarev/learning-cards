.PHONY: ensure-network up-cards up-reviews down-cards down-reviews

ensure-network:
	@if ! docker network ls --format '{{.Name}}' | grep -q '^iamvkosarev_network$$'; then \
		echo "Creating external network iamvkosarev_network..."; \
		docker network create iamvkosarev_network; \
	else \
		echo "Network iamvkosarev_network already exists."; \
	fi
.PHONY: up-cards up-reviews


local-up-cards: ensure-network
	cd cmd/cards && docker compose -f docker-compose.yml --env-file .env up -d --build

local-up-reviews: ensure-network
	cd cmd/reviews && docker compose -f docker-compose.yml --env-file .env up -d --build

local-down-cards:
	cd cmd/cards && docker compose -f docker-compose.yml --env-file .env down

local-down-reviews:
	cd cmd/reviews && docker compose -f docker-compose.yml --env-file .env down

logs-cards:
	cd cmd/cards && docker compose -f docker-compose.yml logs -f

logs-reviews:
	cd cmd/reviews && docker compose -f docker-compose.yml logs -f