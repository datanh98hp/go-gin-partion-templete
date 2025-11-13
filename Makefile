include .env
export

MIRGRATION_DIR = internal/db/migrations
# MIRGRATION_DIR = .//nternal//db//migrations
CONN_STRING = postgresql://${DB_USER}:${DB_PASS}@$(DB_HOST):$(DB_PORT)/${DB_NAME}?sslmode=${DB_SSLMODE}
PROD_COMPOSE=docker-compose.prod.yml
NOAPP_COMPOSE=docker-compose.noapp.yml
DEV_COMPOSE=docker-compose.dev.yml
ENV_FILE=.env
# Import and export databases
importdb:
	docker exec -i db psql -U root -d ${DB_NAME} < ./backup_data-db.sql
exportdb:
	docker exec -i db pg_dump -U root -d ${DB_NAME} > ./backup_data-db.sql
start_container:
	docker compose up -d
remove_container:
	docker compose down
# Run the Go server
server:
	go run ./cmd/api/main.go
# 	Buid binary file
build:
	go build -o ./bin/myapp ./cmd/api/main.go
# 	Run binary file
run_binary:
	./bin/myapp
# Generate sqlc
sqlc:
	sqlc generate

# Migration database example : make mirgrate_create name="profiles"
migrate_create:
	migrate create -ext sql -dir $(MIRGRATION_DIR) -seq $(name)

# Run all pending migration
migrate_up:
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" up
# Rollback last migration
migrate_down:
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" down 1
# Rollback n migration
migrate_rollback:
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" down $(step)

migrate_force: # example : make mirgrate_create version=1 to set the version to 1.useful when you manually change the database
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" force $(version)
# Drop everything from database
migrate_drop:
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" drop

# 	Migrate to a specific version
# example : make mirgrate_goto version=1 --- to migrate to version 1
migrate_goto:
	migrate -path $(MIRGRATION_DIR) -database "$(CONN_STRING)" goto $(version)

# Build container for noapp
noapp:
	docker compose -f $(NOAPP_COMPOSE) down
	docker compose -f $(NOAPP_COMPOSE) --env-file $(ENV_FILE) up -d --build
# stop all container noapp
stop-noapp:
	docker compose -f $(NOAPP_COMPOSE) down
# --------------------------------
# Build container for dev
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(ENV_FILE) up --build
# stop all container dev
stop-dev:
	docker compose -f $(DEV_COMPOSE) down

# -----------------------------------
# Build container for production
prod:
	docker compose -f $(PROD_COMPOSE) down
	docker compose -f $(PROD_COMPOSE) --env-file $(ENV_FILE) up -d --build
# stop all container production
stop-prod:
	docker compose -f $(PROD_COMPOSE) down
# view logs product
log-prod:
	docker compose -f $(PROD_COMPOSE) logs -f --tail=100
# Go to container api
bash:
	docker exec -it go-api /bin/sh


.PHONY: importdb exportdb build server run_binary start_container remove_container migrate_create migrate_up migrate_down migrate_force migrate_goto migrate_drop migrate_rollback sqlc prod stop-prod log-prod noapp stop-noapp dev stop-dev bash 