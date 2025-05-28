include make/common.mk
include make/quick-start.mk
include make/database.mk
include make/grpc.mk
include make/docker.mk
include make/local-docker.mk
include make/vps.mk
include make/test-cmd.mk
include make/test.mk

release-and-deploy: docker-release vps-deploy-cards vps-deploy-reviews
	@echo "Released and deployed $(IMAGE_NAME)"

quick-start: setup-config local-docker-restart local-database-migrate-up