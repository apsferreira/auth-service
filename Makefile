.PHONY: up down logs ps clean test migrate-up migrate-down seed help

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

# Shared postgres container name
PG=shared-postgres
DB=auth_db

help: ## Show this help message
	@echo "$(GREEN)Auth Service - Makefile Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

up: ## Start all containers (requires shared-infra running)
	@echo "$(GREEN)Starting Auth Service...$(NC)"
	docker compose up -d

down: ## Stop all containers
	@echo "$(YELLOW)Stopping Auth Service...$(NC)"
	docker compose down

logs: ## Show logs
	docker compose logs -f

ps: ## Show running containers
	docker compose ps

clean: ## Remove all containers (CAUTION! Data is in shared-infra)
	@echo "$(RED)WARNING: Stopping containers only. Data lives in shared-infra.$(NC)"
	docker compose down

migrate-up: ## Run database migrations (up)
	@echo "$(GREEN)Running migrations...$(NC)"
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/001_create_tenants.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/002_create_roles_and_permissions.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/003_create_users.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/004_create_otp_codes.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/005_create_refresh_tokens.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/007_create_services.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/008_add_service_to_permissions.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/010_add_performance_indexes.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/011_add_password_auth.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/013_create_auth_events.up.sql
	@echo "$(GREEN)Migrations completed!$(NC)"

migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/013_create_auth_events.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/011_add_password_auth.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/010_add_performance_indexes.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/009_seed_services.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/008_add_service_to_permissions.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/007_create_services.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/005_create_refresh_tokens.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/004_create_otp_codes.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/003_create_users.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/002_create_roles_and_permissions.down.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/001_create_tenants.down.sql
	@echo "$(GREEN)Migrations rolled back!$(NC)"

seed: ## Seed mock data
	@echo "$(GREEN)Seeding mock data...$(NC)"
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/006_seed_data.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/009_seed_services.up.sql
	docker exec -i $(PG) psql -U postgres -d $(DB) < backend/migrations/012_seed_admin_user.up.sql
	@echo "$(GREEN)Data seeded!$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	cd backend && go test ./...

restart: ## Restart all services
	@echo "$(YELLOW)Restarting services...$(NC)"
	docker compose restart

build: ## Build images
	@echo "$(GREEN)Building images...$(NC)"
	docker compose build
