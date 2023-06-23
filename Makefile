deps:
	go mod download
build:
	go build -o myapp main.go
run:
	go run main.go
clean:
	rm -f myapp