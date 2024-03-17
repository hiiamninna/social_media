db_up:
	migrate -path ./db/migrations -database postgres://postgres:password@localhost:5432/social_media?sslmode=disable -verbose up

db_down:
	migrate -path ./db/migrations -database postgres://postgres:password@localhost:5432/social_media?sslmode=disable -verbose down

go-run:
	go run main.go

go-build:
	go build main.go