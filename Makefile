include .env
export

database-local:
	docker exec -it chirpydb psql -U postgres

up-migrate:
	goose -dir ./sql/schema postgres "$(DB_URL)" up

down-migrate:
	goose -dir ./sql/schema postgres "$(DB_URL)" down

	
