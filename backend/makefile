URL_DB = postgres://postgres:1601@localhost:5432/whatsupp?sslmode=disable

run:
	go run main.go

migrate-create:
	migrate create -ext sql -dir db/migration -seq create_$(name)_table

migrate-up:
	migrate -database $(URL_DB) -path db/migration/ up
