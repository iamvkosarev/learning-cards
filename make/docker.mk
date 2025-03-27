DOCKER_REPO := $(REPOSITORY_HOLDER)/$(PROJECT_NAME)

VERSION_FILE = $(ROOT_DIR)/VERSION
TAG := $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=new)

IMAGE_NAME := $(DOCKER_REPO):$(TAG)

docker-tag:
	@echo "Building image: $(IMAGE_NAME)"

.PHONY: docker-build docker-tag
docker-build: docker-tag
	docker build --platform=linux/amd64 -t $(IMAGE_NAME) .

docker-push:
	docker push $(IMAGE_NAME)

git-tag:
	@if git rev-parse $(TAG) >/dev/null 2>&1; then \
		echo "Git tag $(TAG) already exists. Skipping."; \
	else \
		git tag $(TAG) && \
		git push origin $(TAG) --quiet && \
		echo "Git tag $(TAG) created and pushed."; \
	fi

.PHONY: docker-release git-tag docker-build docker-push
docker-release: git-tag docker-build docker-push
	@echo "Docker released into $(IMAGE_NAME)"

