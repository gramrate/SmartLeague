# SmartLeague Monorepo

SmartLeague — монорепозиторий с backend API и frontend приложением для клубов, серий и игр.

## Структура

- `backend/` — Go backend (API, миграции, docker-compose, swagger).
- `frontend/` — React + TypeScript + Vite frontend.

## Что в проекте

- cookie-based авторизация;
- клубы с ролями (`member/resident/leader/president`);
- серии (на рейтинг/без рейтинга, club-only, платные, закрытие регистрации);
- игры, участники, результаты, leaderboard;
- бан-лист клуба;
- управление оплатой участников платной серии (для лидеров/президента).

## Документация

- backend: [backend/README.md](backend/README.md)
- frontend: [frontend/README.md](frontend/README.md)
- swagger UI (после запуска backend): `http://localhost:8000/api/v1/swagger/index.html`

## Быстрый старт

1. Поднять backend по `backend/README.md`.
2. Поднять frontend по `frontend/README.md`.
