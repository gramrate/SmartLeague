# SmartLeague Backend

Go backend для SmartLeague: авторизация, клубы, серии, игры, участники, лидерборд и модерация.

## Текущее покрытие

- cookie-based auth (`HttpOnly` access/refresh cookies);
- пользователи и профиль;
- клубы: роли (`member/resident/leader/president`), бан-лист, управление участниками;
- серии: рейтинг/нерейтинг, club-only, платные серии, закрытие регистрации;
- игры в серии, участники игр, результаты, leaderboard;
- платежный статус участников **платной серии** (для лидеров/президента).

## Swagger и комментарии

- Swagger UI: `http://localhost:8000/api/v1/swagger/index.html`
- Исходники swagger: `backend/docs/swagger.{yaml,json}`
- Комментарии для новых ручек оплат добавлены:
  - `GET /api/v1/series/{id}/payments`
  - `POST /api/v1/series/{id}/payment/{profile_id}`

Если после изменений API нужно пересобрать swagger, запусти генерацию в проекте (по вашему принятому процессу).

## Доступ без авторизации

Публичные `GET`-ручки:

- `/api/v1/club/all`
- `/api/v1/club/{id}`
- `/api/v1/club/{id}/members`
- `/api/v1/club/{id}/series`
- `/api/v1/series/all`
- `/api/v1/series/{id}`
- `/api/v1/series/{id}/full`
- `/api/v1/series/{id}/participants`
- `/api/v1/series/{id}/games`
- `/api/v1/series/{id}/leaderboard`
- `/api/v1/game/{id}`
- `/api/v1/game/{id}/full`
- `/api/v1/user/{id}`

Остальные mutating endpoints требуют auth и проверку прав.

## Валидация полей

Все основные DTO-строки ограничены `min/max` через теги `validate`.
Дополнительно закрыт обход в `CreateGameDraft`: теперь валидация `CreateGameRequest` вызывается и для draft-ручки.

## Быстрый запуск (Docker)

Все команды выполняются из `backend/`.

1. Подготовка:
```bash
cp .env.example .env
cp config.yaml.example config.yaml
```

2. RSA-ключи:
```bash
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

3. Старт:
```bash
docker compose up --build -d
```

4. Проверка:
```bash
docker compose ps
docker compose logs -f app
```

5. Остановка:
```bash
docker compose down
```

## Локальный запуск

Подними Postgres/Valkey/MinIO и запусти:

```bash
GOCACHE=/tmp/gocache go run ./cmd
```

## Полезно

- Технические детали: `backend/DEV.md`
- Миграции: `backend/internal/adapters/repository/sql/migrate/migrations`
