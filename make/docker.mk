DOCKER_REPO ?= iamvkosarev/learning-cards
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +"%Y%m%d-%H%M%S")
VERSION ?= $(BUILD_TIME)-$(GIT_COMMIT)
IMAGE_NAME := $(DOCKER_REPO):$(VERSION)

docker_tag:
	@echo "Building image: $(IMAGE_NAME)"

docker_build: docker_tag
	docker build -t $(IMAGE_NAME) .

docker_push:
	docker push $(IMAGE_NAME)

docker_release: docker_build docker_push
	@echo "Docker released into $(IMAGE_NAME)"