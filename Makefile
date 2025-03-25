include make/common.mk
include make/database.mk
include make/grpc.mk
include make/docker.mk
include make/local-docker.mk
include make/vps.mk

release-and-deploy: docker-release vps-deploy
	@echo "Released and deployed $(IMAGE_NAME)"