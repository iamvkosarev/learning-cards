REPOSITORY_HOLDER := iamvkosarev
CARDS_PROJECT_NAME := learning-cards
REVIEWS_PROJECT_NAME := learning-cards-reviews

VERSION_FILE = $(ROOT_DIR)/VERSION
TAG := $(shell go run $(ROOT_DIR)/pkg/tools/versiongen -file $(VERSION_FILE) -mode=new)

CARDS_IMAGE_NAME := $(REPOSITORY_HOLDER)/$(CARDS_PROJECT_NAME):$(TAG)
REVIEWS_IMAGE_NAME := $(REPOSITORY_HOLDER)/$(REVIEWS_PROJECT_NAME):$(TAG)

.PHONY: docker-build docker-push git-tag docker-release

docker-build:
	@echo "Building: $(CARDS_IMAGE_NAME)"
	cd $(ROOT_DIR)/cmd/cards docker build --platform=linux/amd64 --target cards -t $(CARDS_IMAGE_NAME) .
	@echo "Building: $(REVIEWS_IMAGE_NAME)"
	cd $(ROOT_DIR)/cmd/cards  build --platform=linux/amd64 --target reviews -t $(REVIEWS_IMAGE_NAME) .

docker-push:
	docker push $(CARDS_IMAGE_NAME)
	docker push $(REVIEWS_IMAGE_NAME)

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
	@echo "Released:"
	@echo "  - $(CARDS_IMAGE_NAME)"
	@echo "  - $(REVIEWS_IMAGE_NAME)"