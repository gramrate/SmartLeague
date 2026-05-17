# SmartLeague Frontend

Реализация по ТЗ: `frontend/SmartLeague_Frontend_TZ.docx`.

## Стек

- React 18 + TypeScript (Vite)
- React Router v6
- Zustand (auth)
- TanStack Query v5 (серверный стейт)
- React Hook Form + Zod (формы)
- Tailwind CSS (базовые стили)

## Быстрый старт

```bash
cd frontend
cp .env.example .env
npm i
npm run dev
```

По умолчанию:
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8000`

Swagger backend: `http://localhost:8000/api/v1/swagger/index.html`

## Настройки

Переменные окружения: `frontend/.env`

- `VITE_API_BASE_URL` — базовый URL backend (используется и для dev-proxy в Vite).
- `VITE_APP_TITLE` — заголовок вкладки.

## Авторизация

Auth — cookie-based, запросы ходят с `credentials: "include"`.
Если защищенная ручка вернула `401`, клиент пытается вызвать `POST /api/v1/auth/refresh` и повторить запрос 1 раз.

## Публичные страницы (без регистрации)

Доступны без логина:

- `/clubs`
- `/clubs/:id`
- `/clubs/:id/series`
- `/series`
- `/series/:id`
- `/game/:id`
- `/user/:id` (просмотр аккаунта)

Требуют логин:

- `/account`
- `/clubs/create`
- `/series/create`
- `/game/:id/manage`
