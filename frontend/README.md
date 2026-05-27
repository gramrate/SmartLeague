# SmartLeague Frontend

Frontend для SmartLeague на `React + TypeScript + Vite` с интеграцией в backend API.

## Стек
- React 19
- TypeScript
- Vite
- TanStack Router
- TanStack Query
- Zustand (store авторизации)
- Cookie-based auth (через backend cookies)

## Запуск

Требования:
- Node.js 20+
- Запущенный backend SmartLeague (по `backend/README.md`)

1. Установка зависимостей:
```bash
cd frontend
npm install
```

2. Настройка API (опционально):
```bash
# по умолчанию используется http://localhost:8000
export VITE_API_BASE_URL=http://localhost:8000
```

3. Dev-режим:
```bash
npm run dev
```

4. Production build:
```bash
npm run build
```

## Что реализовано

- публичные страницы клубов/серий/игроков/игр;
- кабинет пользователя;
- управление клубом для `leader/president`;
- управление серией, включая:
  - регистрация открыта/закрыта;
  - серия на рейтинг / без рейтинга;
  - club-only серия;
  - удаление серии;
  - блок оплат участников в **платной серии** с переключением `Оплатили / Не оплатили`;
- управление играми и результатами;
- бан-лист клуба (поиск, пагинация, бан/разбан, переключение статуса без дергания экрана).

## Архитектура
- `src/lib/api.ts` — единый API-клиент и все вызовы backend
- `src/types/api.ts` — DTO/типы API
- `src/lib/auth-store.ts` — централизованное состояние `me` + auth bootstrap
- `src/routes/*` — страницы
- `src/components/site/*` — общие блоки UI

## Сессия и авторизация
- Backend использует cookie auth (`credentials: include`).
- После логина/регистрации загружается и хранится `me`.
- При `401` фронт делает `POST /api/v1/auth/refresh` и повторяет запрос.
- Публичные страницы доступны без регистрации.

## Иерархия ролей клуба
Используется строго:
- `0` — none
- `1` — member
- `2` — resident
- `3` — leader
- `4` — president

Человеко-читаемые лейблы отображаются в UI (`src/lib/roles.ts`).

## Актуальные лимиты и валидация на UI

- лимиты списков:
  - серии: `10`
  - клубы/игроки/прочие списки: `15`
- в формах добавлены счетчики `текущее/максимум` и остаток символов;
- для полей с лимитами проставлен `maxLength`;
- кнопки отправки/сохранения блокируются при превышении лимитов.

## Права и доступ

- управление клубом/серией/играми доступно только `leader/president` своего клуба;
- club-only серии/игры не видны и не открываются посторонним;
- ошибки permission обрабатываются в UI (`403`, сообщение про запрет доступа).

## Страницы

- `/`
- `/login`, `/register`
- `/clubs`, `/clubs/create`, `/clubs/:id`, `/clubs/:id/manage`, `/clubs/:id/members`, `/clubs/:id/series`, `/clubs/:id/games`
- `/series`, `/series/create`, `/series/:id`, `/series/:id/manage`
- `/game/:id`, `/game/:id/manage`
- `/players`, `/user/:id`, `/user/:id/series`, `/user/:id/games`
- `/account`

## API

Frontend использует `src/lib/api.ts` и работает поверх backend swagger (`backend/docs/swagger.json`).
