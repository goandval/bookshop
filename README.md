# Bookshop

## Описание

Bookshop — это монолитный сервис на Go для управления каталогом книг с поддержкой ролей пользователей, корзины, заказов, интеграцией с Keycloak, Redis, Kafka и PostgreSQL.

## Запуск

1. Клонируйте репозиторий и перейдите в папку проекта.
2. Соберите и запустите сервисы:

```sh
docker-compose up --build
```

3. Приложение будет доступно на `http://localhost:8081`.
4. Keycloak доступен на `http://localhost:8080` (realm: `bookshop`).

## Доступы пользователей

**Пользователи:**
- user1@bookshop.local / user1pass (роль: user)
- user2@bookshop.local / user2pass (роль: user)

**Администраторы:**
- admin1@bookshop.local / admin1pass (роль: admin)
- admin2@bookshop.local / admin2pass (роль: admin)

## Основные переменные окружения

- `CONFIG_PATH` — путь к yaml-конфигу приложения (по умолчанию `/app/configs/config.yaml`)

## Миграции

Миграции БД выполняются автоматически при старте приложения.

## Структура проекта

- `cmd/bookshop` — точка входа приложения
- `cmd/migrate` — точка входа для миграций
- `internal/` — бизнес-логика, сервисы, репозитории, интерфейсы
- `api/` — описание API (контракты, схемы)
- `configs/` — конфиги
- `migrations/` — миграции БД
- `pkg/` — переиспользуемые пакеты

## Тесты и моки

Для генерации моков используйте:

```sh
make mocks
```

Для запуска тестов:

```sh
make test
```

## Линтеры

Для запуска линтеров:

```sh
make lint
```

## Примеры запросов к API

### Получить список книг (публично)
```sh
curl http://localhost:8081/books
```

### Получить книгу по id (публично)
```sh
curl http://localhost:8081/books/1
```

### Получить список категорий (публично)
```sh
curl http://localhost:8081/categories
```

### Добавить книгу в корзину (требуется JWT)
```sh
curl -X POST http://localhost:8081/cart \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"book_id": 1}'
```

### Оформить заказ (требуется JWT)
```sh
curl -X POST http://localhost:8081/orders \
  -H "Authorization: Bearer <JWT>"
```

### CRUD для книг и категорий (только для админов, требуется JWT с ролью admin)
- POST /books
- PUT /books/{id}
- DELETE /books/{id}
- POST /categories
- PUT /categories/{id}
- DELETE /categories/{id}

## Миграции

Для применения миграций:
```sh
make migrate
```

## Тесты

Для запуска unit-тестов:
```sh
make test
```

Для генерации моков:
```sh
make mocks
```

## Установка зависимостей

Перед запуском любых команд выполните:
```sh
make deps
``` 