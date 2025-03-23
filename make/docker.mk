include make/common.mk

DOCKER_REPO ?= iamvkosarev/learning-cards

VERSION_FILE = $(ROOT_DIR)/VERSION
TAG := $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=new)

IMAGE_NAME := $(DOCKER_REPO):$(TAG)

docker_tag:
	@echo "Building image: $(IMAGE_NAME)"

docker_build: docker_tag
	docker build --platform=linux/amd64 -t $(IMAGE_NAME) .

docker_push:
	docker push $(IMAGE_NAME)

git_tag:
	@if git rev-parse $(TAG) >/dev/null 2>&1; then \
		echo "Git tag $(TAG) already exists. Skipping."; \
	else \
		git tag $(TAG) && \
		git push origin $(TAG) --quiet && \
		echo "Git tag $(TAG) created and pushed."; \
	fi

docker_release: git_tag docker_build docker_push
	@echo "Docker released into $(IMAGE_NAME)"