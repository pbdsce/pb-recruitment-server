# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

ifeq ($(DB_ADDR),)
    $(error DB_ADDR environment variable is not set. Please add DB_ADDR to your .env file)
endif

MIGRATIONS_PATH = ./internal/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@, $(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@, $(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@, $(MAKECMDGOALS))
