# Bookshop

## Описание

Bookshop — это монолитный сервис на Go для управления каталогом книг с поддержкой ролей пользователей, корзины, заказов, интеграцией с Keycloak, Redis, Kafka и PostgreSQL.  
Документация API (Swagger) доступна по адресу:  
**`http://localhost:8081/swagger/index.html`** (после генерации командой `make swag`).

---

## Быстрый старт

1. Клонируйте репозиторий и перейдите в папку проекта.
2. Установите зависимости:
   ```sh
   make deps
   ```
3. Соберите и запустите сервисы:
   ```sh
   make up
   ```
   - Приложение: http://localhost:8081
   - Keycloak: http://localhost:8080 (realm: `bookshop`)

4. Для перезапуска сервиса после изменений:
   ```sh
   docker-compose up -d --build bookshop
   ```

---

## Пользователи

**Пользователи:**
- user1@bookshop.local / user1pass (роль: user)
- user2@bookshop.local / user2pass (роль: user)

**Администраторы:**
- admin1@bookshop.local / admin1pass (роль: admin)
- admin2@bookshop.local / admin2pass (роль: admin)

---

## Переменные окружения и конфиг

- `CONFIG_PATH` — путь к yaml-конфигу приложения (по умолчанию `/app/configs/config.yaml`)

**Пример config.yaml:**
```yaml
postgres:
  host: postgres
  port: 5432
  user: bookshop
  password: bookshop
  dbname: bookshop
  sslmode: disable
redis:
  addr: redis:6379
  db: 0
kafka:
  brokers:
    - kafka:9092
  order_topic: order_placed
keycloak:
  url: http://keycloak:8080
  realm: bookshop
  client_id: bookshop-api
http:
  addr: :8081
log:
  level: info
```

---

## Основные команды Makefile

- `make up` — запуск всех сервисов через docker-compose
- `make build` — сборка приложения
- `make run` — запуск приложения локально
- `make test` — запуск unit-тестов
- `make lint` — запуск линтеров
- `make mocks` — генерация моков через mockery
- `make migrate` — применение миграций
- `make swag` — генерация swagger-документации (docs/swagger.yaml, docs/swagger.json)
- `make deps` — установка зависимостей

---

## Тесты и моки

- Для запуска тестов:
  ```sh
  make test
  ```
- Для генерации моков:
  ```sh
  make mocks
  ```

---

## Swagger/OpenAPI

- Для генерации документации:
  ```sh
  make swag
  ```
- После генерации документация доступна по адресу:  
  http://localhost:8081/swagger/index.html

---

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

---

## Миграции

Для применения миграций:
```sh
make migrate
```

---

## Структура проекта

- `cmd/bookshop` — точка входа приложения
- `cmd/migrate` — миграции
- `internal/` — бизнес-логика, сервисы, репозитории, интерфейсы, моки
- `configs/` — конфиги
- `migrations/` — миграции БД
- `docs/` — swagger-документация
- `Makefile` — основные команды