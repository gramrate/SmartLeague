# DEV Guide

Техническая документация для разработки и расширения SmartLeague backend.

## Архитектура проекта

Проект собран по слоям:

- `internal/adapters/controller/api/v1`  
  HTTP-слой: хендлеры, декодинг query/body, валидация, маппинг ошибок в HTTP-коды.

- `internal/domain/service`  
  Бизнес-логика use-case уровня. Здесь не должно быть HTTP-деталей.

- `internal/adapters/repository/sql` / `.../valkey` / `.../minio`  
  Доступ к хранилищам и внешним системам.

- `internal/domain/dto`  
  Контракты запросов/ответов API.

Поток запроса:
1. Handler принимает запрос.
2. Service выполняет бизнес-логику.
3. Repository читает/пишет данные.
4. Service собирает DTO-ответ.
5. Handler возвращает HTTP-ответ.

## Авторизация и refresh

Модель — cookie-based:
- `user_auth_access_token` (короткий)
- `user_auth_refresh_token` (длинный, HttpOnly)

Важно:
- refresh-токен **один на пользователя** (на все устройства), хранится в БД. Новый логин/refresh перезаписывает токен.

Обновление:
- клиент дергает `POST /api/v1/auth/refresh` (без body/query)
- сервер валидирует refresh cookie и выставляет новые cookies (access+refresh)

Поведение клиента:
- если защищенная ручка вернула `401` (нет/просрочен access cookie) — сначала дернуть refresh, затем повторить запрос.

## Service Provider: что это и зачем

`ServiceProvider` (`internal/adapters/app/service_provider`) — центральный контейнер зависимостей приложения.

Он нужен, чтобы:
- создавать зависимости в одном месте;
- лениво инициализировать тяжелые ресурсы (DB, Redis, MinIO, сервисы);
- не размазывать `new(...)` по хендлерам и `server.Setup`.

Как работать правильно:
- новые сервисы и клиенты добавлять через Service Provider;
- в `server.Setup` брать зависимости только через `serviceProvider.*Service()`/`*Config()`/`*Middleware()`;
- не создавать руками репозитории/сервисы в хендлерах.

Это гарантирует единый lifecycle и предсказуемую инициализацию.

## Правила расширения (в ширину)

Когда добавляешь новый модуль (например, `reviews`), делай в таком порядке:

1. DTO  
   Добавь request/response структуры в `internal/domain/dto`.

2. Миграции (если есть новая таблица)  
   Добавь миграцию в `internal/adapters/repository/sql/migrate/migrations`.

3. Repository  
   Реализуй доступ к данным в `internal/adapters/repository/sql/<module>`.

4. Service  
   Добавь бизнес-логику в `internal/domain/service/<module>`.

5. Handler  
   Добавь HTTP-обработчики в `internal/adapters/controller/api/v1/<module>`.

6. Service Provider  
   Зарегистрируй интерфейс и фабрику в `internal/adapters/app/service_provider`.

7. Router setup  
   Подключи handler в `internal/adapters/controller/api/server/server.go`.

8. Swagger  
   Добавь/обнови swagger-комментарии и перегенерируй docs.

Важно:
- Handler не должен ходить в repo напрямую.
- Service не должен знать про Echo/HTTP.
- Repo не должен содержать бизнес-правила.

## Пагинация (единый контракт)

Для list-эндпоинтов используется структура:

```json
{
  "items": [],
  "pagination": {
    "total_items": 0,
    "total_pages": 0,
    "current_page": 1,
    "has_next": false,
    "has_previous": false
  }
}
```

Где:
- размер страницы берется из `limit`;
- смещение берется из `offset`;
- `total_items` считается отдельным `Count` с теми же фильтрами.

## Локальная генерация артефактов

Генерация Swagger:

```bash
swag init -g cmd/main.go -o docs
```

Правило для комментариев над хендлером:
- первая строка должна начинаться с имени функции: `HandlerName Description`

## Миграции

Миграции — обычный SQL, применяются автоматически при старте приложения (см. `internal/adapters/app/service_provider/sql_db.go`).

Где лежат:
- `internal/adapters/repository/sql/migrate/migrations/*.sql`

Правила:
- имя файла: `0003_some_change.sql` (лексикографический порядок = порядок применения);
- миграция должна быть идемпотентной (`IF NOT EXISTS`, `DO $$ ... EXCEPTION WHEN duplicate_object ...`), потому что мигратор простейший.

## Ключи JWT (RSA)

Если ключей нет, сгенерируй:

```bash
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

## Тесты/сборка

Если окружение read-only для go build cache, используй:

```bash
GOCACHE=/tmp/gocache go test ./...
```
