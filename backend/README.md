# SmartLeague Backend

Бэкенд для SmartLeague.

Что есть сейчас:
- cookie-based авторизация (HttpOnly access/refresh cookies);
- CRUD профиля;
- клубы + участники (1 человек = 1 клуб) + роли в клубе (member/leader/president);
- серии игр от клуба + игры + участники/результаты + лидерборд;
- PostgreSQL (основные данные), Valkey (access token blacklist), MinIO (хранилище файлов/картинок — по мере надобности).

## Публичный доступ (без регистрации/логина)

Следующие ручки доступны без cookies:

- `GET /api/v1/club/all`
- `GET /api/v1/club/{id}`
- `GET /api/v1/club/{id}/members`
- `GET /api/v1/club/{id}/series`
- `GET /api/v1/series/{id}`
- `GET /api/v1/series/all`
- `GET /api/v1/series/{id}/full`
- `GET /api/v1/series/{id}/participants`
- `GET /api/v1/series/{id}/games`
- `GET /api/v1/series/{id}/leaderboard`
- `GET /api/v1/game/{id}`
- `GET /api/v1/game/{id}/full`
- `GET /api/v1/user/{id}` (просмотр аккаунта/профиля)

Остальные ручки на создание/изменение/удаление требуют авторизацию.

Техническая документация: `DEV.md`.

## Быстрый запуск через Docker

Все команды ниже выполняй из папки `backend/`.

### 1. Подготовить окружение

Скопируйте `.env.example` в `.env`:

```bash
cp .env.example .env
```

Скопируйте `config.yaml.example` в `config.yaml`:

```bash
cp config.yaml.example config.yaml
```

Для запуска через `docker compose` в `config.yaml` должны быть хосты сервисов: `postgres`, `valkey`, `minio` (в `config.yaml.example` уже так).

### 2. Сгенерировать RSA-ключи (обязательно)

Сервис использует JWT с RSA-ключами. До сборки контейнера создайте ключи:

```bash
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

### 3. Запустить проект

```bash
docker compose up --build -d
```

### 4. Проверить, что сервис поднялся

API: `http://localhost:8000`  
Swagger: `http://localhost:8000/api/v1/swagger/index.html`

## Как обновлять токены

Токены живут в HttpOnly cookies, руками их дергать не надо.

- В системе **один refresh-токен на пользователя** (на все устройства). Если залогиниться на другом устройстве — refresh-токен на предыдущем будет перезаписан.
- Когда access-токен истёк — дергай `POST /api/v1/auth/refresh`.
- Эндпоинт читает `user_auth_refresh_token` из cookies и в ответ выставляет новые cookies:
  - `user_auth_access_token`
  - `user_auth_refresh_token`

Практика для клиента:
- Если любой защищенный запрос вернул `401` — сначала дерни refresh, потом повтори исходный запрос.

Проверка статуса контейнеров:

```bash
docker compose ps
```

Логи приложения:

```bash
docker compose logs -f app
```

### 5. Остановить проект

```bash
docker compose down
```

## Локальный запуск (без Docker)

1) Подними Postgres/Valkey/MinIO как тебе удобно и пропиши их в `config.yaml`.
2) Запусти:

```bash
GOCACHE=/tmp/gocache go run ./cmd
```
