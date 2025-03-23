local_docker_build:
	docker-compose build

local_docker_up:
	docker-compose up -d

local_docker_down:
	docker-compose down

local_docker_logs:
	docker-compose logs -f