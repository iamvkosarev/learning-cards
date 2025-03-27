# learning-cards

**Learning Cards** is a microservice built with Go that allows users to create, manage, and test memory flashcards grouped into collections.

The project supports both gRPC and REST interfaces.

## Quick Start

You must have [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/) and [Goose](https://github.com/pressly/goose) installed.

```bash
make quick-start
```

This command prepares your local environment and launches the service using Docker Compose.
It performs the following steps:

1. `setup-config` - copies default .env and config files from templates (if not already present)
2. `local-docker-restart` - builds and starts the service and database containers
3. `local-database-migrate-up` - applies pending database migrations via Goose

Once completed, your service will be available at:
- REST:  localhost:8080/api/learning-cards/v1
- gRPC:  localhost:50051

## Local API Testing
You can quickly test the REST API using the predefined Makefile targets:

### Create a new cards group
```bash
make local-test-add-group name="Test"
```

### Add a card to a group
```bash
make local-test-add-card group_id=1 front_text=こんいちは 
```

### Get all cards from a group
```bash
make local-test-add-card group_id=1
```

### View service logs:

```bash
make local-docker-logs
```
