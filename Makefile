postgres:
	docker run -d -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=sandbox -e POSTGRES_USER=postgres  --name postgres --rm  postgres

redis:
	docker run --name redis -d -p 6379:6379  --rm redis

run:
	docker-compose up -d

stop:
	docker-compose down