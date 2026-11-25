# Лабораторная работа №3 — Go

Структуры данных, сериализация, модульные тесты и бенчмарки.

Исходный проект: https://github.com/VorobevAB/Lab_3 (Максим Воробьев). В этой версии исправлены проверки односвязного списка в тестах.

## Запуск

```bash
go mod download
go build -o code .
./code
```

## Линтер, тесты и покрытие

```bash
go vet ./...
go test -v ./...
go test ./structs -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Бенчмарки

```bash
go test -bench=. -benchmem ./...
```

Покрытие из приложенного исходного отчёта — 42,3% операторов; для достижения 85% нужны дополнительные тесты.
