include .envrc
MIGRATIONS_DIR := ./cmd/migrate/migrations


migrate-create:
	goose create $(name) sql -dir $(MIGRATIONS_DIR) -s
# make migrate-create name=create_users
# make migrate-create name=create_posts

#docker exec -it shotseek-db-1 psql -U admin shotseek

#goose create create_users sql -s -table "users" 

#direnv allow .
