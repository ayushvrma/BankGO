postgres:
	docker run --name some-postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:12-alpine

createdb:
	docker exec -it some-postgres12 createdb --username=root --owner=root BankGo

dropdb:
	docker exec -it some-postgres12 dropdb BankGo

migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/BankGo?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/BankGo?sslmode=disable" --verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc

# In terms of Make, a phony target is simply a target that is always out-of-date, so whenever you ask make <phony_target>, it will run, independent from the state of the file system. Some common make targets that are often phony are: all, install, clean, distclean, TAGS, info, check.