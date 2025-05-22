local-docker-restart: local-docker-down local-docker-build local-docker-up

local-docker-build:
	docker-compose build

local-docker-up:
	@docker-compose up -d learning-cards-postgres learning-cards learning-cards-reviews

local-docker-down:
	docker-compose down

local-docker-logs:
	docker-compose logs -f