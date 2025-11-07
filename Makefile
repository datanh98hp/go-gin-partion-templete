include .env
export

MIRGRATION_DIR = internal/db/migrations
# MIRGRATION_DIR = .//nternal//db//migrations
CONN_STRING = postgresql://${DB_USER}:${DB_PASS}@$(DB_HOST):$(DB_PORT)/${DB_NAME}?sslmode=${DB_SSLMODE}

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
	
.PHONY: importdb exportdb server start_container remove_container migrate_create migrate_up migrate_down migrate_force migrate_goto migrate_drop migrate_rollback sqlc