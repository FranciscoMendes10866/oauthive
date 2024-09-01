MIGRATIONS_DIR = db/migrations
DATABASE_URL = $(PWD)/database.db
SCRIPT = create_db.sh

GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) sqlite3 $(DATABASE_URL)

.PHONY: create up down

create:
	@echo "Creating new migration: $(NAME)"
	@$(GOOSE_CMD) create $(NAME) sql

up:
	@echo "Applying migrations..."
	@$(GOOSE_CMD) up

down:
	@echo "Rolling back migrations..."
	@$(GOOSE_CMD) down

reset:
	@echo "Resetting database..."
	@$(GOOSE_CMD) reset

status:
	@echo "Migration status:"
	@$(GOOSE_CMD) status

fix:
	@echo "Fixing migration (for testing only)..."
	@$(GOOSE_CMD) fix

create-database:
	@echo "Running script to create and optimize db..."
	@chmod +x $(SCRIPT)
	@./$(SCRIPT)
