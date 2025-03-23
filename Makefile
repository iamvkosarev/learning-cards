include make/grpc.mk
include make/local-docker.mk
include make/docker.mk
include make/vps.mk

release_and_deploy: docker_release deploy_vps
	@echo "Released and deployed $(IMAGE_NAME)"