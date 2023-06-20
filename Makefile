# Компилируем проект
build:
	go build -o myapp main.go

# Запускаем проект
run:
	go run main.go

# Запускаем тесты
test:
	go test ./...

# Устанавливаем зависимости
deps:
	go mod download

# Очищаем проект
clean:
	rm -f myapp