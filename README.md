# go-musthave-group-diploma-tpl

Шаблон репозитория для группового дипломного проекта курса "Go-разработчик"

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-group-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.


## Миграция базы данных
Установка необходимых инструментов
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
Создание таблички. Будут созданы два пустых файла в ./internal/storage/migrations. Их необходимо заполнить нужными данными
```
migrate create -ext sql -dir internal/storage/migrations -seq create_users_table
```
Применить новую миграцию
```
migrate -path internal/storage/migrations -database postgres://user:password@localhost:5432/dbname?sslmode=disable up
```
Откатить миграцию
```
migrate -path internal/storage/migrations -database postgres://user:password@localhost:5432/dbname?sslmode=disable down 1
```
Добавить новую версию. Например имзенить тип колонки value в табличке balances
```
migrate create -ext sql -dir internal/storage/migrations -seq change_value_type_in_balances
```

## Запуск сервисов (docker-compose)
```
docker-compose up
```

## Запуск сервисов (local)
Установка taskfile утилиты
```
go install github.com/go-task/task/v3/cmd/task@latest
```
Конфигурирование
```
vim .env
```
Запуск Gophermart
```
task run
```
Запуск Accrual System
```
task run-accrual
```