deps:
	go mod download
down_migrate:
	migrate -path ./migrations -database 'postgres://postgres:5218521111@localhost:5436/postgres?sslmode=disable' down	
up_migrate:
	migrate -path ./migrations -database 'postgres://postgres:5218521111@localhost:5436/postgres?sslmode=disable' up
build:
	go build -o myapp main.go
run:
	go run main.go
clean:
	rm -f myapp