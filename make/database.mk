MIGRATIONS_DIR := $(ROOT_DIR)/migrations


database-migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Need argument 'name=...'" && exit 1; \
	fi
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

DB_HOST ?= localhost
MIGRATE_CMD ?= up

define goose-cmd
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" $(MIGRATE_CMD)
endef

database-migrate:
	@$(goose-cmd)

database-migrate-up:
	$(MAKE) database-migrate DB_HOST=$(VPS_HOST) MIGRATE_CMD=status
	$(MAKE) database-migrate DB_HOST=$(VPS_HOST) MIGRATE_CMD=up


database-migrate-down:
	$(MAKE) database-migrate DB_HOST=$(VPS_HOST) MIGRATE_CMD=down

local-database-migrate-up:
	$(MAKE) database-migrate MIGRATE_CMD=status
	$(MAKE) database-migrate MIGRATE_CMD=up


local-database-migrate-down:
	$(MAKE) database-migrate MIGRATE_CMD=down


#local-database-migrate-reset:
#	$(MAKE) database-migrate MIGRATE_CMD=reset