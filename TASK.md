# Техническое задание: Voice Git Assistant

## 1. Описание проекта

**Voice Git Assistant** — консольное Java-приложение для управления git-репозиториями через голосовые и текстовые команды с использованием искусственного интеллекта. Система обеспечивает автоматизацию git-операций, ведение истории действий, аналитику активности и управление множественными репозиториями через личный кабинет пользователя.

---

## 2. Цели и задачи

### Основные цели:
- Упростить работу с git через естественный язык
- Автоматизировать рутинные git-операции
- Предоставить аналитику по активности в репозиториях
- Обеспечить безопасное хранение credentials и токенов

### Задачи:
1. Реализовать систему авторизации и управления пользователями
2. Создать CRUD-интерфейс для управления репозиториями
3. Интегрировать STT/TTS для голосового взаимодействия
4. Подключить AI для понимания команд на естественном языке
5. Реализовать все базовые git-операции через JGit
6. Создать систему логирования и аналитики
7. Обеспечить хранение данных в SQLite

---

## 3. Требования к реализации

### 3.1 Обязательные требования курса
- ✅ Private репозиторий на GitHub
- ✅ Язык: Java 17+
- ✅ Консольный интерфейс
- ✅ Хранение данных в БД (SQLite)
- ✅ Документация в README
- ✅ Покрытие юнит-тестами (JUnit 5 + Mockito)
- ✅ Код по принципам SOLID и GRASP
- ✅ Использование библиотек для парсинга и БД разрешено

### 3.2 Технологический стек

| Компонент | Технология | Назначение |
|-----------|------------|------------|
| **Язык** | Java 17+ | Основной язык разработки |
| **Сборка** | Maven 3.8+ | Управление зависимостями |
| **БД** | SQLite 3 + JDBC | Хранение данных |
| **Git** | JGit 6.7+ | Git-операции без системного git |
| **STT** | OpenAI Whisper API | Распознавание речи |
| **TTS** | ElevenLabs API / Google TTS | Синтез речи |
| **AI** | Anthropic Claude API | Понимание команд |
| **HTTP** | Java HttpClient | Работа с API |
| **Шифрование** | Java Crypto API | Шифрование токенов |
| **Хеширование** | BCrypt / PBKDF2 | Хеширование паролей |
| **Тестирование** | JUnit 5, Mockito | Юнит-тесты |
| **Логирование** | SLF4J + Logback | Логирование |

---

## 4. Функциональные требования

### 4.1 Модуль авторизации

#### 4.1.1 Регистрация пользователя
**Описание**: Создание нового аккаунта в системе

**Входные данные**:
- Username (уникальный, 3-20 символов, только латиница и цифры)
- Password (минимум 8 символов, обязательны: заглавная буква, цифра, спецсимвол)
- Email (опционально, валидация формата)

**Бизнес-логика**:
1. Проверка уникальности username
2. Валидация пароля по требованиям безопасности
3. Хеширование пароля (BCrypt, cost factor 12)
4. Сохранение в таблицу `users`
5. Создание дефолтных настроек в `user_settings`

**Выходные данные**:
- Успех: "Регистрация успешна! Войдите в систему"
- Ошибка: описание проблемы (username занят / слабый пароль / невалидный email)

#### 4.1.2 Вход в систему
**Описание**: Аутентификация пользователя

**Входные данные**:
- Username
- Password

**Бизнес-логика**:
1. Поиск пользователя в БД по username
2. Проверка пароля через BCrypt.verify()
3. Создание сессии (SessionManager хранит текущего User)
4. Логирование успешного входа

**Выходные данные**:
- Успех: переход в личный кабинет
- Ошибка: "Неверный логин или пароль" (без детализации)

#### 4.1.3 Выход из системы
**Описание**: Завершение сессии

**Бизнес-логика**:
1. Очистка SessionManager
2. Возврат в главное меню

---

### 4.2 Модуль управления репозиториями (CRUD)

#### 4.2.1 Создание репозитория (Create)
**Описание**: Добавление нового репозитория в систему

**Входные данные**:
- Название репозитория (уникальное в рамках пользователя)
- Локальный путь (абсолютный путь к папке .git)
- Remote URL (GitHub/GitLab, опционально)
- Ветка по умолчанию (default: main)
- Автокоммит: включен/выключен
- Интервал автокоммита (минуты, если включен)

**Валидация**:
1. Путь существует и содержит .git папку
2. Remote URL валидный (если указан)
3. Название уникально для пользователя
4. Интервал > 0 если автокоммит включен

**Бизнес-логика**:
1. Проверка доступности пути через Files.exists()
2. Проверка наличия .git через JGit.open()
3. Сохранение в таблицу `repositories`
4. Если автокоммит включен — запуск ScheduledExecutorService

**Выходные данные**:
- Успех: "Репозиторий 'backend-api' добавлен"
- Ошибка: описание проблемы

#### 4.2.2 Просмотр репозиториев (Read)
**Описание**: Отображение списка репозиториев пользователя

**Фильтры**:
- По названию (поиск подстроки)
- По статусу автокоммита (включен/выключен)
- По дате добавления

**Сортировка**:
- По названию (A-Z)
- По дате создания (новые сначала)
- По количеству операций (самые активные)

**Отображаемые данные**:
```
ID  | Название        | Путь                     | Ветка | Автокоммит | Операций
1   | backend-api     | /home/user/backend       | main  | ✅ (30 мин)| 145
2   | frontend-app    | /home/user/frontend      | dev   | ❌         | 89
3   | docs           | /home/user/docs          | main  | ✅ (60 мин)| 12
```

**Действия**:
- Выбрать для работы
- Редактировать
- Удалить
- Показать статистику

#### 4.2.3 Редактирование репозитория (Update)
**Описание**: Изменение параметров репозитория

**Редактируемые поля**:
- Название
- Remote URL
- Ветка по умолчанию
- Настройки автокоммита (вкл/выкл, интервал)

**Ограничения**:
- Локальный путь изменить нельзя (только удаление и пересоздание)
- Нельзя переименовать в существующее название

**Бизнес-логика**:
1. Валидация новых значений
2. Если изменился интервал автокоммита — перезапуск планировщика
3. UPDATE в БД
4. Логирование изменений

#### 4.2.4 Удаление репозитория (Delete)
**Описание**: Удаление записи о репозитории из системы

**Важно**:
- Удаляется только запись в БД
- Локальные файлы остаются нетронутыми
- История операций сохраняется (repo_id становится NULL)

**Подтверждение**:
```
Вы уверены, что хотите удалить репозиторий 'backend-api'?
Это действие удалит:
- Настройки репозитория
- Планировщик автокоммита

НЕ будет удалено:
- Локальные файлы
- История операций (останется в логах)

Введите 'YES' для подтверждения:
```

**Бизнес-логика**:
1. Остановка планировщика автокоммита (если был)
2. Удаление из таблицы `repositories`
3. UPDATE `operation_log` SET repo_id = NULL WHERE repo_id = X

---

### 4.3 Модуль голосового ассистента

#### 4.3.1 Распознавание голоса (STT)
**Описание**: Преобразование аудио в текст

**Входные данные**:
- Аудиофайл (WAV/MP3/OGG)
- Язык (ru/en, берется из user_settings)

**Процесс**:
1. Запись аудио через Java Sound API (до нажатия Enter/Timeout)
2. Отправка в OpenAI Whisper API
3. Получение транскрипции

**API запрос**:
```java
POST https://api.openai.com/v1/audio/transcriptions
Content-Type: multipart/form-data

{
  "file": audio_file,
  "model": "whisper-1",
  "language": "ru"
}
```

**Выходные данные**:
- Текст команды: "Закоммить изменения с сообщением фикс багов"
- Ошибка: "Не удалось распознать речь"

**Fallback**: Если STT недоступен → режим текстового ввода

#### 4.3.2 Понимание команды (AI)
**Описание**: Парсинг естественного языка в структурированную команду

**Входные данные**:
- Текст от пользователя (голос или клавиатура)
- Контекст: текущий репозиторий, ветка, последние операции

**Промпт для Claude API**:
```
Ты ассистент для git-команд. Пользователь сказал: "{user_input}"

Текущий контекст:
- Репозиторий: {repo_name}
- Ветка: {branch}
- Незакоммиченные файлы: {modified_files}

Верни JSON с командой:
{
  "action": "commit|push|branch|merge|log|diff|status",
  "parameters": {
    "message": "commit message",
    "branch_name": "название ветки",
    "count": 5,
    ...
  },
  "confidence": 0.95
}

Если команда непонятна, верни {"action": "clarify", "question": "..."}
```

**Примеры парсинга**:

| Ввод пользователя | Результат |
|-------------------|-----------|
| "Закоммить всё, я починил баг в логине" | `{action: "commit", parameters: {message: "fix: bug in login"}}` |
| "Покажи последние 5 коммитов" | `{action: "log", parameters: {count: 5}}` |
| "Создай ветку feature-payment" | `{action: "branch", parameters: {name: "feature-payment"}}` |
| "Что изменилось с утра?" | `{action: "diff", parameters: {since: "today 00:00"}}` |
| "Запуши всё на мастер" | `{action: "push", parameters: {branch: "master"}}` |

**Бизнес-логика**:
1. Отправка запроса в Claude API
2. Парсинг JSON-ответа
3. Валидация команды
4. Если confidence < 0.7 → запрос уточнения у пользователя

#### 4.3.3 Выполнение git-команды
**Описание**: Выполнение операции через JGit

**Поддерживаемые команды**:

##### COMMIT
```java
Parameters: {message: String, addAll: boolean}
Действия:
1. git add . (если addAll = true)
2. git commit -m "{message}"
3. Сохранение commit hash в БД
```

##### PUSH
```java
Parameters: {branch: String, remote: String = "origin"}
Действия:
1. git push {remote} {branch}
2. Авторизация через токен из api_tokens
3. Обработка конфликтов
```

##### LOG
```java
Parameters: {count: int = 10, author: String, since: Date, until: Date}
Действия:
1. git log --oneline -n {count}
2. Фильтрация по параметрам
3. Форматирование вывода
```

##### BRANCH
```java
Parameters: {name: String, checkout: boolean}
Действия:
1. git branch {name}
2. git checkout {name} (если checkout = true)
```

##### DIFF
```java
Parameters: {since: Date, files: List<String>}
Действия:
1. git diff [files]
2. Опционально отправка в AI для summarization
```

##### STATUS
```java
Parameters: {}
Действия:
1. git status --porcelain
2. Парсинг измененных файлов
3. Группировка по типам (modified, added, deleted)
```

**Обработка ошибок**:
- Merge conflict → "Обнаружен конфликт в файле X, разрешите его вручную"
- Auth failed → "Проверьте GitHub токен в настройках"
- No changes → "Нет изменений для коммита"

#### 4.3.4 Озвучивание результата (TTS)
**Описание**: Преобразование текста в речь

**Входные данные**:
- Текст результата операции
- Настройка включения TTS из user_settings

**API запрос (ElevenLabs)**:
```java
POST https://api.elevenlabs.io/v1/text-to-speech/{voice_id}
{
  "text": "Готово! Коммит c4f2b9a успешно запушен в ветку main",
  "model_id": "eleven_multilingual_v2",
  "voice_settings": {
    "stability": 0.5,
    "similarity_boost": 0.75
  }
}
```

**Воспроизведение**:
- Получение MP3
- Воспроизведение через Java Sound API

**Fallback**: Если TTS отключен или недоступен → вывод в консоль

---

### 4.4 Модуль истории операций

#### 4.4.1 Логирование операций
**Описание**: Автоматическая запись всех действий

**Логируемые данные**:
- User ID
- Repository ID
- Тип операции (COMMIT, PUSH, BRANCH, etc)
- Входная команда (голос/текст)
- Распарсенная команда (JSON)
- Commit hash (если применимо)
- Статус (SUCCESS/FAILED)
- Текст ошибки (если FAILED)
- Способ вызова (VOICE/TEXT/AUTO/SCHEDULE)
- Timestamp

**Бизнес-логика**:
```java
public void log(OperationLog entry) {
    entry.setTimestamp(LocalDateTime.now());
    entry.setUserId(sessionManager.getCurrentUser().getId());
    repository.insert(entry);
}
```

#### 4.4.2 Просмотр истории
**Описание**: Отображение лога операций

**Фильтры**:
- По репозиторию
- По типу операции
- По статусу (только успешные / только ошибки)
- По дате (диапазон)
- По способу вызова

**Пример вывода**:
```
Дата/Время          | Репозиторий    | Операция | Команда                              | Статус
--------------------|----------------|----------|--------------------------------------|--------
20.05.2026 15:34    | backend-api    | COMMIT   | "Закоммить изменения"                | ✅
20.05.2026 15:35    | backend-api    | PUSH     | "Запуши на main"                     | ✅
20.05.2026 14:22    | frontend-app   | BRANCH   | "Создай ветку feature-X"             | ✅
20.05.2026 12:10    | docs           | COMMIT   | auto: scheduled commit               | ❌
```

**Детальный просмотр**:
```
Операция #145
────────────────────────────────
Дата: 20.05.2026 15:34:12
Репозиторий: backend-api (/home/user/backend)
Тип: COMMIT
Способ: VOICE
Голосовая команда: "Закоммить изменения, я починил баг в авторизации"
Распознанная команда: {"action":"commit","parameters":{"message":"fix: bug in auth"}}
Commit hash: c4f2b9a3e1f
Статус: SUCCESS
Время выполнения: 1.2s
```

#### 4.4.3 Экспорт истории
**Описание**: Выгрузка данных в файл

**Форматы**:
- CSV (для Excel)
- JSON (для программной обработки)
- HTML (читаемый отчет)

**CSV пример**:
```csv
timestamp,repo,operation,status,commit_hash,triggered_by
2026-05-20 15:34:12,backend-api,COMMIT,SUCCESS,c4f2b9a3e1f,VOICE
2026-05-20 15:35:01,backend-api,PUSH,SUCCESS,,VOICE
```

#### 4.4.4 Очистка истории
**Описание**: Удаление старых записей

**Параметры**:
- Удалить записи старше N дней
- Оставить только ошибки (для отладки)
- Полная очистка (с подтверждением)

**Подтверждение**:
```
Будет удалено 1247 записей старше 90 дней.
Это действие необратимо. Продолжить? (yes/no):
```

---

### 4.5 Модуль аналитики

#### 4.5.1 Статистика по репозиторию
**Описание**: Аналитика активности в конкретном репозитории

**Метрики**:
```
Статистика: backend-api
───────────────────────────────────────
За последние 7 дней:
• Коммитов: 23
• Успешных push: 21
• Ошибок: 2 (8.7%)
• Создано веток: 3
• Слияний: 1

По дням недели:
ПН ████████░░ 8
ВТ ██████████ 10
СР ██████░░░░ 6
ЧТ ████░░░░░░ 4
ПТ ████████░░ 8
СБ ██░░░░░░░░ 2
ВС ░░░░░░░░░░ 0

Топ-5 файлов по изменениям:
1. src/api/AuthController.java (12 коммитов)
2. src/service/UserService.java (8 коммитов)
3. pom.xml (5 коммитов)
4. README.md (3 коммита)
5. src/config/SecurityConfig.java (2 коммита)

Средний размер коммита: 47 строк
Самый большой коммит: c4f2b9a (342 строки)
```

#### 4.5.2 Общая статистика пользователя
**Описание**: Агрегированная статистика по всем репозиториям

**Метрики**:
```
Общая статистика (все репозитории)
───────────────────────────────────────
Всего репозиториев: 3
Всего операций: 456

За последний месяц:
• Коммитов: 89
• Push: 76
• Создано веток: 12
• Слияний: 5
• Ошибок: 8 (1.8%)

Самый активный репозиторий: backend-api (234 операции)
Самый активный день: 15.05.2026 (18 операций)
Средний интервал между коммитами: 3.4 часа

Использование:
• Голосовых команд: 45% (206 операций)
• Текстовых команд: 35% (160 операций)
• Автокоммитов: 20% (90 операций)
```

#### 4.5.3 Трендовый анализ
**Описание**: Определение паттернов активности

**Анализируемые паттерны**:
- Самые активные часы работы
- Дни недели с максимальной активностью
- Корреляция между типами операций
- Рост/падение активности по неделям

**Пример вывода**:
```
Анализ паттернов активности
───────────────────────────────────────
Ваш пик продуктивности: 14:00-17:00 (62% коммитов)
Самые продуктивные дни: ВТ, СР (73% всех операций)

Наблюдения:
⚠️ Количество ошибок push выросло на 35% за последнюю неделю
   Рекомендация: проверьте GitHub токен

✓ Средний размер коммитов уменьшился на 40%
   Хороший признак: меньше, но чаще

📈 Активность выросла на 25% по сравнению с прошлым месяцем
```

---

### 4.6 Модуль настроек пользователя

#### 4.6.1 Управление профилем

**Изменение пароля**:
```
Текущий пароль: ****
Новый пароль: ****
Повторите новый пароль: ****

Валидация:
- Минимум 8 символов
- Заглавная буква
- Цифра
- Спецсимвол
- Не совпадает со старым
```

**Изменение email**:
```
Новый email: user@example.com
Валидация формата через regex
```

#### 4.6.2 Управление API токенами

**Структура**:
```
API токены
───────────────────────────────────────
1. GitHub
   Статус: ✅ Активен
   Добавлен: 01.05.2026
   [Обновить] [Удалить]

2. OpenAI (Whisper)
   Статус: ❌ Не настроен
   [Добавить токен]

3. ElevenLabs (TTS)
   Статус: ✅ Активен
   Добавлен: 05.05.2026
   [Обновить] [Удалить]
```

**Добавление токена**:
```
Введите GitHub токен (скрыт): ********************
Проверка токена... ✓
Токен сохранен (зашифрован AES-256)
```

**Шифрование**:
- Алгоритм: AES-256-GCM
- Ключ: производный от мастер-пароля пользователя (PBKDF2)
- Salt: уникальный для каждого токена

#### 4.6.3 Настройки голоса

```
Настройки голосового ввода
───────────────────────────────────────
STT (распознавание речи)
  [✓] Включено
  Язык: [Русский ▼]
  Сервис: [OpenAI Whisper ▼]

TTS (озвучивание)
  [✓] Включено
  Голос: [Rachel (ElevenLabs) ▼]
  Скорость: [1.0x]
  Озвучивать:
    [✓] Результаты команд
    [✓] Ошибки
    [ ] Подтверждения
```

#### 4.6.4 Настройки AI

```
Настройки AI-ассистента
───────────────────────────────────────
Провайдер: [Claude (Anthropic) ▼]
           Альтернативы: GPT-4, Local LLM

Режим работы:
  ( ) Быстрый - короткие ответы
  (•) Сбалансированный
  ( ) Детальный - с объяснениями

Автоматические действия:
  [✓] Автоматический push после commit
  [ ] Автоматическое создание веток по паттернам
  [✓] Умные commit messages (AI генерирует на основе diff)
```

---

### 4.7 Модуль автокоммита

#### 4.7.1 Настройка автокоммита
**Описание**: Автоматическое создание коммитов по расписанию

**Параметры**:
- Включить/выключить
- Интервал (минуты): 15, 30, 60, 120, 240
- Время работы: только рабочие часы / круглосуточно
- Автопуш: включить/выключить
- Шаблон сообщения: "auto: {timestamp}" или кастомный

**Конфигурация**:
```
Автокоммит для: backend-api
───────────────────────────────────────
Статус: [✓ Включен]
Интервал: [30 минут ▼]
Активен: [09:00 - 18:00 ПН-ПТ]

Шаблон commit message:
  (•) Умный (AI генерирует на основе изменений)
  ( ) Фиксированный: auto: {timestamp}
  ( ) Кастомный: _______________________

Действия после коммита:
  [✓] Автоматический push
  [ ] Уведомление в TTS

Последний автокоммит: 20.05.2026 15:30
Следующий запланирован: 20.05.2026 16:00
```

#### 4.7.2 Логика работы

**Планировщик**:
```java
ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);

scheduler.scheduleAtFixedRate(() -> {
    List<Repository> repos = repoService.getWithAutoCommitEnabled();
    for (Repository repo : repos) {
        if (isInWorkingHours(repo) && hasChanges(repo)) {
            performAutoCommit(repo);
        }
    }
}, 0, 5, TimeUnit.MINUTES); // Проверка каждые 5 минут
```

**Процесс автокоммита**:
1. Проверка наличия изменений (git status)
2. Если изменений нет → пропуск
3. Генерация commit message:
   - Если "умный режим" → отправка `git diff` в AI
   - Иначе → использование шаблона
4. Выполнение commit
5. Если включен автопуш → выполнение push
6. Логирование в БД с triggered_by = "AUTO"

**Умная генерация сообщения**:
```
Промпт для AI:
"Проанализируй изменения и создай короткое commit message (max 50 символов):

git diff:
{diff_output}

Формат: type: description
Типы: feat, fix, refactor, docs, test, chore
Примеры: 'fix: null pointer in login', 'feat: add user roles'"
```

---

## 5. Нефункциональные требования

### 5.1 Производительность
- Запуск приложения: < 3 секунды
- Отклик на команду (без AI): < 500ms
- Распознавание голоса: < 5 секунд
- AI-парсинг команды: < 3 секунды
- Git-операция (commit): < 2 секунды
- Запрос к БД: < 100ms

### 5.2 Безопасность

#### Хранение паролей:
- Хеширование: BCrypt (cost factor 12)
- Никогда не хранятся в plaintext
- Не логируются

#### Хранение токенов:
- Шифрование: AES-256-GCM
- Ключ шифрования: производный от мастер-пароля через PBKDF2
- Salt: уникальный для каждого токена, 16 байт
- IV: случайный для каждого шифрования, 12 байт

#### SQL Injection:
- Использование PreparedStatement для всех запросов
- Валидация всех пользовательских вводов
- Параметризованные запросы

#### API Keys:
- Не хранятся в коде (только в БД зашифрованными)
- Не логируются
- Можно отозвать и обновить через настройки

### 5.3 Надежность
- Обработка всех исключений
- Graceful degradation (если STT недоступен → текстовый ввод)
- Retry mechanism для API (3 попытки с exponential backoff)
- Транзакции для критических БД операций
- Бэкапы БД (автоматическое копирование раз в день)

### 5.4 Удобство использования (UX)
- Интуитивное меню навигации
- Подсказки для каждой команды
- Подтверждение критических операций (удаление)
- Прогресс-бары для долгих операций
- Цветной вывод в консоль (ANSI colors)
- История команд (стрелка вверх)

### 5.5 Масштабируемость
- Поддержка неограниченного количества репозиториев
- Поддержка неограниченного количества пользователей
- История операций: партиционирование по дате (опционально)
- Архивация старых логов (> 1 года)

---

## 6. Архитектура системы

### 6.1 Архитектурные принципы

**SOLID**:
- **S**ingle Responsibility: каждый класс отвечает за одну задачу
- **O**pen/Closed: расширение через интерфейсы, не модификацию
- **L**iskov Substitution: реализации интерфейсов взаимозаменяемы
- **I**nterface Segregation: узкие специализированные интерфейсы
- **D**ependency Inversion: зависимость от абстракций

**GRASP**:
- **Information Expert**: логика в классах с нужными данными
- **Creator**: фабрики для создания объектов
- **Low Coupling**: минимизация зависимостей через DI
- **High Cohesion**: связанные методы в одном классе
- **Controller**: тонкий слой между UI и бизнес-логикой

### 6.2 Слоистая архитектура

```
┌─────────────────────────────────────────┐
│  Presentation Layer (CLI)               │  ← Консольный интерфейс
├─────────────────────────────────────────┤
│  Application Layer (Services)           │  ← Бизнес-логика
├─────────────────────────────────────────┤
│  Domain Layer (Models)                  │  ← Сущности
├─────────────────────────────────────────┤
│  Infrastructure Layer (DB, APIs)        │  ← Внешние системы
└─────────────────────────────────────────┘
```

### 6.3 Структура пакетов

```
com.voicegit/
│
├── Main.java                              # Точка входа
│
├── cli/                                   # PRESENTATION LAYER
│   ├── ConsoleInterface.java              # Главное меню
│   ├── AuthMenu.java                      # Авторизация
│   ├── DashboardMenu.java                 # Личный кабинет
│   ├── RepositoryMenu.java                # CRUD репозиториев
│   ├── VoiceAssistantMenu.java            # Голосовой ассистент
│   ├── HistoryMenu.java                   # Просмотр истории
│   ├── SettingsMenu.java                  # Настройки
│   ├── StatsMenu.java                     # Статистика
│   └── utils/
│       ├── ConsoleColors.java             # ANSI цвета
│       ├── TableFormatter.java            # Форматирование таблиц
│       └── InputValidator.java            # Валидация ввода
│
├── service/                               # APPLICATION LAYER
│   ├── auth/
│   │   ├── AuthService.java               # Регистрация/вход
│   │   ├── SessionManager.java            # Управление сессией
│   │   └── PasswordHasher.java            # BCrypt хеширование
│   │
│   ├── repository/
│   │   ├── RepositoryService.java         # CRUD для репозиториев
│   │   └── RepositoryValidator.java       # Валидация репозиториев
│   │
│   ├── git/
│   │   ├── GitService.java                # Git операции (JGit)
│   │   ├── GitCommandParser.java          # Парсинг команд
│   │   └── AutoCommitScheduler.java       # Планировщик автокоммита
│   │
│   ├── voice/
│   │   ├── STTService.java                # Speech-to-Text (Whisper)
│   │   └── TTSService.java                # Text-to-Speech (ElevenLabs)
│   │
│   ├── ai/
│   │   ├── AIService.java                 # Работа с Claude API
│   │   ├── CommandParser.java             # Парсинг JSON от AI
│   │   └── PromptBuilder.java             # Построение промптов
│   │
│   ├── analytics/
│   │   ├── StatsService.java              # Вычисление статистики
│   │   └── TrendAnalyzer.java             # Анализ паттернов
│   │
│   └── export/
│       ├── ExportService.java             # Экспорт данных
│       ├── CSVExporter.java               # Экспорт в CSV
│       ├── JSONExporter.java              # Экспорт в JSON
│       └── HTMLExporter.java              # Экспорт в HTML
│
├── domain/                                # DOMAIN LAYER
│   ├── model/
│   │   ├── User.java                      # Пользователь
│   │   ├── Repository.java                # Репозиторий
│   │   ├── OperationLog.java              # Лог операции
│   │   ├── UserSettings.java              # Настройки пользователя
│   │   ├── ApiToken.java                  # API токен
│   │   └── Command.java                   # Распарсенная команда
│   │
│   └── enums/
│       ├── OperationType.java             # COMMIT, PUSH, BRANCH, etc
│       ├── OperationStatus.java           # SUCCESS, FAILED
│       ├── TriggerType.java               # VOICE, TEXT, AUTO, SCHEDULE
│       └── AIProvider.java                # CLAUDE, GPT4, LOCAL
│
├── infrastructure/                        # INFRASTRUCTURE LAYER
│   ├── database/
│   │   ├── Database.java                  # SQLite connection pool
│   │   ├── SchemaInitializer.java         # Создание таблиц
│   │   ├── repositories/
│   │   │   ├── UserRepository.java        # CRUD users
│   │   │   ├── RepositoryRepository.java  # CRUD repositories
│   │   │   ├── OperationLogRepository.java# CRUD operation_log
│   │   │   ├── SettingsRepository.java    # CRUD user_settings
│   │   │   └── ApiTokenRepository.java    # CRUD api_tokens
│   │   └── migration/
│   │       └── DatabaseMigrator.java      # Миграции схемы
│   │
│   ├── api/
│   │   ├── WhisperClient.java             # HTTP клиент для Whisper
│   │   ├── ElevenLabsClient.java          # HTTP клиент для TTS
│   │   ├── ClaudeClient.java              # HTTP клиент для Claude
│   │   └── GitHubClient.java              # HTTP клиент для GitHub API
│   │
│   └── security/
│       ├── TokenEncryption.java           # AES шифрование токенов
│       ├── EncryptionKeyManager.java      # Управление ключами
│       └── SecureStorage.java             # Безопасное хранилище
│
├── config/
│   ├── AppConfig.java                     # Конфигурация приложения
│   ├── DatabaseConfig.java                # Настройки БД
│   └── ApiConfig.java                     # Настройки API
│
└── util/
    ├── Logger.java                        # Обертка над SLF4J
    ├── DateFormatter.java                 # Форматирование дат
    └── FileUtils.java                     # Утилиты для файлов
```

### 6.4 Dependency Injection (DI)

**Ручной DI через конструкторы** (без фреймворка):

```java
// Main.java
public class Main {
    public static void main(String[] args) {
        // Infrastructure
        Database db = new Database(new DatabaseConfig());
        UserRepository userRepo = new UserRepository(db);
        RepositoryRepository repoRepo = new RepositoryRepository(db);
        OperationLogRepository logRepo = new OperationLogRepository(db);
        
        // Security
        TokenEncryption encryption = new TokenEncryption();
        PasswordHasher hasher = new PasswordHasher();
        
        // External APIs
        WhisperClient whisperClient = new WhisperClient(new ApiConfig());
        ClaudeClient claudeClient = new ClaudeClient(new ApiConfig());
        
        // Services
        AuthService authService = new AuthService(userRepo, hasher);
        SessionManager sessionManager = new SessionManager();
        STTService sttService = new STTService(whisperClient);
        AIService aiService = new AIService(claudeClient);
        GitService gitService = new GitService();
        
        // Menus
        AuthMenu authMenu = new AuthMenu(authService, sessionManager);
        DashboardMenu dashboard = new DashboardMenu(
            sessionManager, 
            repoService, 
            statsService
        );
        
        // Start
        ConsoleInterface cli = new ConsoleInterface(authMenu, dashboard);
        cli.start();
    }
}
```

---

## 7. Модель данных (БД)

### 7.1 Схема SQLite

```sql
-- ПОЛЬЗОВАТЕЛИ
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME,
    
    CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 20)
);

CREATE INDEX idx_users_username ON users(username);

-- РЕПОЗИТОРИИ
CREATE TABLE repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    local_path TEXT NOT NULL,
    remote_url TEXT,
    branch TEXT DEFAULT 'main',
    auto_commit_enabled BOOLEAN DEFAULT 0,
    auto_commit_interval INTEGER DEFAULT 0,  -- минуты, 0 = выкл
    auto_push_enabled BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name),
    CHECK (auto_commit_interval >= 0)
);

CREATE INDEX idx_repos_user ON repositories(user_id);
CREATE INDEX idx_repos_auto_commit ON repositories(auto_commit_enabled, auto_commit_interval);

-- ИСТОРИЯ ОПЕРАЦИЙ
CREATE TABLE operation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    repo_id INTEGER,  -- NULL если репозиторий удален
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    operation_type TEXT NOT NULL,  -- COMMIT, PUSH, PULL, BRANCH, MERGE, etc
    voice_input TEXT,              -- Оригинальный голосовой ввод
    parsed_command TEXT,           -- JSON команды от AI
    commit_hash TEXT,              -- SHA коммита (если применимо)
    branch TEXT,                   -- Ветка операции
    files_changed INTEGER,         -- Количество измененных файлов
    lines_added INTEGER,           -- Добавлено строк
    lines_deleted INTEGER,         -- Удалено строк
    status TEXT NOT NULL,          -- SUCCESS, FAILED
    error_message TEXT,            -- Текст ошибки если FAILED
    triggered_by TEXT NOT NULL,    -- VOICE, TEXT, AUTO, SCHEDULE
    execution_time_ms INTEGER,     -- Время выполнения в мс
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (repo_id) REFERENCES repositories(id) ON DELETE SET NULL,
    CHECK (operation_type IN ('COMMIT', 'PUSH', 'PULL', 'BRANCH', 'MERGE', 'CHECKOUT', 'STATUS', 'LOG', 'DIFF')),
    CHECK (status IN ('SUCCESS', 'FAILED')),
    CHECK (triggered_by IN ('VOICE', 'TEXT', 'AUTO', 'SCHEDULE'))
);

CREATE INDEX idx_log_user ON operation_log(user_id);
CREATE INDEX idx_log_repo ON operation_log(repo_id);
CREATE INDEX idx_log_timestamp ON operation_log(timestamp DESC);
CREATE INDEX idx_log_status ON operation_log(status);
CREATE INDEX idx_log_operation ON operation_log(operation_type);

-- НАСТРОЙКИ ПОЛЬЗОВАТЕЛЯ
CREATE TABLE user_settings (
    user_id INTEGER PRIMARY KEY,
    voice_enabled BOOLEAN DEFAULT 1,
    tts_enabled BOOLEAN DEFAULT 1,
    stt_language TEXT DEFAULT 'ru',
    ai_provider TEXT DEFAULT 'claude',
    ai_smart_commits BOOLEAN DEFAULT 1,  -- AI генерирует commit messages
    auto_push_after_commit BOOLEAN DEFAULT 0,
    theme TEXT DEFAULT 'dark',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (ai_provider IN ('claude', 'gpt4', 'local'))
);

-- API ТОКЕНЫ
CREATE TABLE api_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    service_name TEXT NOT NULL,  -- github, openai, elevenlabs, anthropic
    token_encrypted TEXT NOT NULL,
    encryption_iv TEXT NOT NULL,  -- Initialization Vector для AES
    encryption_salt TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_used DATETIME,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, service_name)
);

CREATE INDEX idx_tokens_user ON api_tokens(user_id);

-- ПЛАНИРОВЩИК ЗАДАЧ (опционально, для будущих фич)
CREATE TABLE scheduled_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL,
    task_type TEXT NOT NULL,  -- AUTO_COMMIT, AUTO_BACKUP, etc
    cron_expression TEXT,     -- Для сложных расписаний
    enabled BOOLEAN DEFAULT 1,
    last_run DATETIME,
    next_run DATETIME,
    
    FOREIGN KEY (repo_id) REFERENCES repositories(id) ON DELETE CASCADE
);
```

### 7.2 Примеры запросов

**Регистрация пользователя**:
```sql
INSERT INTO users (username, password_hash, email) 
VALUES (?, ?, ?);

-- Создание дефолтных настроек
INSERT INTO user_settings (user_id) 
VALUES (last_insert_rowid());
```

**Добавление репозитория**:
```sql
INSERT INTO repositories (user_id, name, local_path, remote_url, branch)
VALUES (?, ?, ?, ?, ?);
```

**Логирование операции**:
```sql
INSERT INTO operation_log (
    user_id, repo_id, operation_type, voice_input, parsed_command,
    commit_hash, status, triggered_by, execution_time_ms
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
```

**Получение статистики за неделю**:
```sql
SELECT 
    COUNT(*) as total_operations,
    SUM(CASE WHEN operation_type = 'COMMIT' THEN 1 ELSE 0 END) as commits,
    SUM(CASE WHEN operation_type = 'PUSH' THEN 1 ELSE 0 END) as pushes,
    SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as errors,
    AVG(execution_time_ms) as avg_time
FROM operation_log
WHERE repo_id = ?
  AND timestamp >= datetime('now', '-7 days')
```

**Топ файлов по изменениям** (требует парсинга git diff):
```sql
-- Этот запрос работает если мы храним измененные файлы в JSON
SELECT 
    json_extract(parsed_command, '$.files') as files,
    COUNT(*) as change_count
FROM operation_log
WHERE repo_id = ?
  AND operation_type = 'COMMIT'
  AND timestamp >= datetime('now', '-7 days')
GROUP BY files
ORDER BY change_count DESC
LIMIT 5
```

---

## 8. Интеграции с внешними API

### 8.1 OpenAI Whisper API (STT)

**Документация**: https://platform.openai.com/docs/api-reference/audio

**Эндпоинт**:
```
POST https://api.openai.com/v1/audio/transcriptions
```

**Заголовки**:
```
Authorization: Bearer {OPENAI_API_KEY}
Content-Type: multipart/form-data
```

**Тело запроса**:
```
file: <audio_file>  (MP3, MP4, MPEG, MPGA, M4A, WAV, WEBM)
model: whisper-1
language: ru (опционально, но улучшает качество)
```

**Ответ**:
```json
{
  "text": "Закоммить изменения с сообщением фикс багов в авторизации"
}
```

**Обработка ошибок**:
- 400: Invalid audio format → показать поддерживаемые форматы
- 401: Invalid API key → предложить обновить токен в настройках
- 429: Rate limit → показать retry через N секунд

**Реализация**:
```java
public class STTService {
    private final WhisperClient client;
    
    public String transcribe(File audioFile, String language) {
        try {
            MultipartBody body = new MultipartBody.Builder()
                .addFormDataPart("file", audioFile.getName(),
                    RequestBody.create(audioFile, MediaType.parse("audio/wav")))
                .addFormDataPart("model", "whisper-1")
                .addFormDataPart("language", language)
                .build();
                
            Response response = client.post("/audio/transcriptions", body);
            
            if (response.isSuccessful()) {
                JSONObject json = new JSONObject(response.body().string());
                return json.getString("text");
            } else {
                throw new STTException("Failed: " + response.code());
            }
        } catch (IOException e) {
            throw new STTException("Network error", e);
        }
    }
}
```

---

### 8.2 Anthropic Claude API (AI)

**Документация**: https://docs.anthropic.com/claude/reference

**Эндпоинт**:
```
POST https://api.anthropic.com/v1/messages
```

**Заголовки**:
```
x-api-key: {ANTHROPIC_API_KEY}
anthropic-version: 2023-06-01
content-type: application/json
```

**Тело запроса**:
```json
{
  "model": "claude-3-5-sonnet-20241022",
  "max_tokens": 1024,
  "messages": [
    {
      "role": "user",
      "content": "Пользователь сказал: \"Закоммить изменения\"\n\nКонтекст:\n- Репозиторий: backend-api\n- Ветка: main\n- Незакоммиченные файлы: src/AuthService.java, pom.xml\n\nВерни JSON с командой:\n{\"action\": \"commit|push|branch|merge|log|diff|status\", \"parameters\": {...}, \"confidence\": 0.0-1.0}\n\nЕсли непонятно, верни {\"action\": \"clarify\", \"question\": \"...\"}"
    }
  ]
}
```

**Ответ**:
```json
{
  "id": "msg_01...",
  "type": "message",
  "role": "assistant",
  "content": [
    {
      "type": "text",
      "text": "{\"action\":\"commit\",\"parameters\":{\"message\":\"fix: update auth service\",\"addAll\":true},\"confidence\":0.95}"
    }
  ],
  "model": "claude-3-5-sonnet-20241022",
  "usage": {
    "input_tokens": 145,
    "output_tokens": 28
  }
}
```

**Обработка ответа**:
```java
public Command parseCommand(String userInput, Repository repo, String context) {
    String prompt = buildPrompt(userInput, repo, context);
    
    JSONObject request = new JSONObject()
        .put("model", "claude-3-5-sonnet-20241022")
        .put("max_tokens", 1024)
        .put("messages", new JSONArray()
            .put(new JSONObject()
                .put("role", "user")
                .put("content", prompt)));
    
    Response response = claudeClient.post("/v1/messages", request);
    JSONObject json = new JSONObject(response.body().string());
    
    String text = json.getJSONArray("content")
        .getJSONObject(0)
        .getString("text");
    
    // Парсинг JSON из текста
    JSONObject cmdJson = new JSONObject(text);
    
    if (cmdJson.getString("action").equals("clarify")) {
        throw new ClarificationNeededException(cmdJson.getString("question"));
    }
    
    return Command.fromJSON(cmdJson);
}
```

---

### 8.3 ElevenLabs API (TTS)

**Документация**: https://elevenlabs.io/docs/api-reference

**Эндпоинт**:
```
POST https://api.elevenlabs.io/v1/text-to-speech/{voice_id}
```

**Заголовки**:
```
xi-api-key: {ELEVENLABS_API_KEY}
Content-Type: application/json
```

**Тело запроса**:
```json
{
  "text": "Готово! Коммит c4f2b9a успешно запушен в ветку main",
  "model_id": "eleven_multilingual_v2",
  "voice_settings": {
    "stability": 0.5,
    "similarity_boost": 0.75,
    "style": 0.0,
    "use_speaker_boost": true
  }
}
```

**Ответ**: Binary audio data (MP3)

**Воспроизведение**:
```java
public void speak(String text) {
    byte[] audioData = elevenLabsClient.synthesize(text, voiceId);
    
    // Сохранение временного файла
    File tempFile = File.createTempFile("tts_", ".mp3");
    Files.write(tempFile.toPath(), audioData);
    
    // Воспроизведение через Java Sound API
    AudioInputStream audioStream = AudioSystem.getAudioInputStream(tempFile);
    Clip clip = AudioSystem.getClip();
    clip.open(audioStream);
    clip.start();
    
    // Блокировка до окончания воспроизведения
    Thread.sleep(clip.getMicrosecondLength() / 1000);
    
    // Очистка
    clip.close();
    tempFile.delete();
}
```

---

### 8.4 JGit (Git операции)

**Документация**: https://www.eclipse.org/jgit/

**Maven зависимость**:
```xml
<dependency>
    <groupId>org.eclipse.jgit</groupId>
    <artifactId>org.eclipse.jgit</artifactId>
    <version>6.7.0.202309050840-r</version>
</dependency>
```

**Примеры операций**:

**Commit**:
```java
public String commit(Repository repo, String message) throws GitAPIException {
    try (Git git = Git.open(new File(repo.getLocalPath()))) {
        // git add .
        git.add().addFilepattern(".").call();
        
        // git commit -m "message"
        RevCommit commit = git.commit()
            .setMessage(message)
            .call();
        
        return commit.getName(); // SHA hash
    }
}
```

**Push**:
```java
public void push(Repository repo, String branch) throws GitAPIException {
    String token = apiTokenService.getToken(repo.getUserId(), "github");
    
    try (Git git = Git.open(new File(repo.getLocalPath()))) {
        git.push()
            .setRemote("origin")
            .setRefSpecs(new RefSpec(branch))
            .setCredentialsProvider(
                new UsernamePasswordCredentialsProvider("token", token))
            .call();
    }
}
```

**Log**:
```java
public List<CommitInfo> getLog(Repository repo, int count) {
    try (Git git = Git.open(new File(repo.getLocalPath()))) {
        Iterable<RevCommit> logs = git.log()
            .setMaxCount(count)
            .call();
        
        List<CommitInfo> result = new ArrayList<>();
        for (RevCommit commit : logs) {
            result.add(new CommitInfo(
                commit.getName(),
                commit.getShortMessage(),
                commit.getAuthorIdent().getName(),
                commit.getCommitTime()
            ));
        }
        return result;
    }
}
```

**Status**:
```java
public GitStatus getStatus(Repository repo) {
    try (Git git = Git.open(new File(repo.getLocalPath()))) {
        Status status = git.status().call();
        
        return new GitStatus(
            status.getModified(),
            status.getAdded(),
            status.getRemoved(),
            status.getUntracked()
        );
    }
}
```

**Diff**:
```java
public String getDiff(Repository repo) throws IOException {
    try (Git git = Git.open(new File(repo.getLocalPath()))) {
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        git.diff()
            .setOutputStream(out)
            .call();
        return out.toString("UTF-8");
    }
}
```

---

## 9. Тестирование

### 9.1 Стратегия тестирования

**Покрытие**: минимум 70% code coverage

**Типы тестов**:
- Unit тесты: изолированное тестирование методов
- Integration тесты: взаимодействие с БД
- Mock тесты: внешние API

### 9.2 Инструменты

```xml
<dependencies>
    <!-- JUnit 5 -->
    <dependency>
        <groupId>org.junit.jupiter</groupId>
        <artifactId>junit-jupiter</artifactId>
        <version>5.10.0</version>
        <scope>test</scope>
    </dependency>
    
    <!-- Mockito -->
    <dependency>
        <groupId>org.mockito</groupId>
        <artifactId>mockito-core</artifactId>
        <version>5.5.0</version>
        <scope>test</scope>
    </dependency>
    
    <!-- AssertJ (fluent assertions) -->
    <dependency>
        <groupId>org.assertj</groupId>
        <artifactId>assertj-core</artifactId>
        <version>3.24.2</version>
        <scope>test</scope>
    </dependency>
</dependencies>
```

### 9.3 Примеры тестов

**Unit тест - PasswordHasher**:
```java
@Test
void testPasswordHashing() {
    PasswordHasher hasher = new PasswordHasher();
    String password = "SecurePass123!";
    
    String hash = hasher.hash(password);
    
    assertThat(hash).isNotNull();
    assertThat(hash).isNotEqualTo(password);
    assertThat(hasher.verify(password, hash)).isTrue();
    assertThat(hasher.verify("WrongPass", hash)).isFalse();
}

@Test
void testPasswordHashingIsDeterministic() {
    PasswordHasher hasher = new PasswordHasher();
    String password = "SecurePass123!";
    
    String hash1 = hasher.hash(password);
    String hash2 = hasher.hash(password);
    
    // Хеши должны быть разными из-за random salt
    assertThat(hash1).isNotEqualTo(hash2);
    
    // Но оба должны верифицироваться
    assertThat(hasher.verify(password, hash1)).isTrue();
    assertThat(hasher.verify(password, hash2)).isTrue();
}
```

**Unit тест - CommandParser**:
```java
@Test
void testCommitCommandParsing() {
    CommandParser parser = new CommandParser();
    
    String aiResponse = """
        {"action":"commit","parameters":{"message":"fix: auth bug"},"confidence":0.95}
        """;
    
    Command cmd = parser.parse(aiResponse);
    
    assertThat(cmd.getAction()).isEqualTo(OperationType.COMMIT);
    assertThat(cmd.getParameter("message")).isEqualTo("fix: auth bug");
    assertThat(cmd.getConfidence()).isEqualTo(0.95);
}

@Test
void testClarificationRequest() {
    CommandParser parser = new CommandParser();
    
    String aiResponse = """
        {"action":"clarify","question":"Какую ветку создать?"}
        """;
    
    assertThatThrownBy(() -> parser.parse(aiResponse))
        .isInstanceOf(ClarificationNeededException.class)
        .hasMessageContaining("Какую ветку создать?");
}
```

**Mock тест - AIService**:
```java
@Test
void testAICommandParsing() {
    // Arrange
    ClaudeClient mockClient = mock(ClaudeClient.class);
    AIService aiService = new AIService(mockClient);
    
    String mockResponse = """
        {
          "content": [{
            "type": "text",
            "text": "{\\"action\\":\\"commit\\",\\"parameters\\":{\\"message\\":\\"fix bug\\"}}"
          }]
        }
        """;
    
    when(mockClient.post(any(), any()))
        .thenReturn(mockResponse);
    
    // Act
    Command cmd = aiService.parseCommand("Закоммить изменения", mockRepo);
    
    // Assert
    assertThat(cmd.getAction()).isEqualTo(OperationType.COMMIT);
    verify(mockClient, times(1)).post(any(), any());
}
```

**Integration тест - UserRepository**:
```java
@Test
void testUserCRUD() {
    // Setup in-memory DB
    Database testDb = new Database(":memory:");
    UserRepository userRepo = new UserRepository(testDb);
    
    // Create
    User user = new User("testuser", "hash123", "test@example.com");
    userRepo.save(user);
    
    // Read
    Optional<User> found = userRepo.findByUsername("testuser");
    assertThat(found).isPresent();
    assertThat(found.get().getEmail()).isEqualTo("test@example.com");
    
    // Update
    found.get().setEmail("newemail@example.com");
    userRepo.update(found.get());
    
    Optional<User> updated = userRepo.findByUsername("testuser");
    assertThat(updated.get().getEmail()).isEqualTo("newemail@example.com");
    
    // Delete
    userRepo.delete(updated.get().getId());
    assertThat(userRepo.findByUsername("testuser")).isEmpty();
}
```

**Mock тест - GitService**:
```java
@Test
void testCommitOperation() throws Exception {
    // Arrange
    Git mockGit = mock(Git.class);
    AddCommand mockAdd = mock(AddCommand.class);
    CommitCommand mockCommit = mock(CommitCommand.class);
    RevCommit mockRevCommit = mock(RevCommit.class);
    
    when(mockGit.add()).thenReturn(mockAdd);
    when(mockAdd.addFilepattern(".")).thenReturn(mockAdd);
    when(mockAdd.call()).thenReturn(null);
    
    when(mockGit.commit()).thenReturn(mockCommit);
    when(mockCommit.setMessage(anyString())).thenReturn(mockCommit);
    when(mockCommit.call()).thenReturn(mockRevCommit);
    when(mockRevCommit.getName()).thenReturn("abc123");
    
    GitService gitService = new GitService();
    
    // Act
    String hash = gitService.commit(mockRepo, "test message");
    
    // Assert
    assertThat(hash).isEqualTo("abc123");
    verify(mockAdd).addFilepattern(".");
    verify(mockCommit).setMessage("test message");
}
```

**Parametrized тест - InputValidator**:
```java
@ParameterizedTest
@CsvSource({
    "user, false",          // too short
    "validuser123, true",
    "user-name, false",     // special char
    "ALLUPPER, true",
    "user@name, false",     // @ not allowed
    "a1b2c3d4e5f6g7h8i9j0k1, false"  // too long
})
void testUsernameValidation(String username, boolean expected) {
    InputValidator validator = new InputValidator();
    assertThat(validator.isValidUsername(username)).isEqualTo(expected);
}

@ParameterizedTest
@ValueSource(strings = {
    "weak",
    "NoDigits!",
    "nouppercas3!",
    "NOLOWERCASE3!",
    "NoSpecialChar1"
})
void testWeakPasswords(String password) {
    InputValidator validator = new InputValidator();
    assertThat(validator.isValidPassword(password)).isFalse();
}

@Test
void testStrongPassword() {
    InputValidator validator = new InputValidator();
    assertThat(validator.isValidPassword("Strong1Pass!")).isTrue();
}
```

---

## 10. Конфигурация и настройки

### 10.1 Файл конфигурации

**application.properties**:
```properties
# Database
db.path=./data/voicegit.db
db.backup.enabled=true
db.backup.interval=24h

# API Keys (NOT stored here, only in encrypted DB)
# api.openai.enabled=true
# api.anthropic.enabled=true
# api.elevenlabs.enabled=true

# STT
stt.default.language=ru
stt.timeout.seconds=30

# TTS
tts.default.voice=Rachel
tts.speed=1.0

# Git
git.default.branch=main
git.auto.commit.default.interval=30

# Security
security.password.min.length=8
security.token.encryption.algorithm=AES/GCM/NoPadding
security.session.timeout.minutes=60

# Logging
logging.level=INFO
logging.file.path=./logs/voicegit.log
logging.file.max.size=10MB
```

### 10.2 Environment Variables

Чувствительные данные через env vars:
```bash
# Опционально: дефолтные API ключи для быстрого старта
export VOICEGIT_OPENAI_KEY=sk-...
export VOICEGIT_ANTHROPIC_KEY=sk-ant-...
export VOICEGIT_ELEVENLABS_KEY=...

# Master password для шифрования (если не через GUI)
export VOICEGIT_MASTER_PASSWORD=...
```

---

## 11. Deployment & Distribution

### 11.1 Сборка JAR

**pom.xml**:
```xml
<build>
    <plugins>
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-shade-plugin</artifactId>
            <version>3.5.0</version>
            <executions>
                <execution>
                    <phase>package</phase>
                    <goals>
                        <goal>shade</goal>
                    </goals>
                    <configuration>
                        <transformers>
                            <transformer implementation="org.apache.maven.plugins.shade.resource.ManifestResourceTransformer">
                                <mainClass>com.voicegit.Main</mainClass>
                            </transformer>
                        </transformers>
                        <finalName>voice-git-assistant</finalName>
                    </configuration>
                </execution>
            </executions>
        </plugin>
    </plugins>
</build>
```

**Команда сборки**:
```bash
mvn clean package
# Результат: target/voice-git-assistant.jar
```

### 11.2 Запуск

```bash
# Базовый запуск
java -jar voice-git-assistant.jar

# С кастомной конфигурацией
java -jar voice-git-assistant.jar --config=/path/to/config.properties

# С увеличенной памятью
java -Xmx512m -jar voice-git-assistant.jar
```

### 11.3 Первый запуск

При первом запуске приложение:
1. Создаёт директорию `~/.voicegit/`
2. Инициализирует БД SQLite
3. Создаёт таблицы
4. Предлагает регистрацию первого пользователя

```
╔════════════════════════════════════════════╗
║   Voice Git Assistant v1.0                 ║
║   Первый запуск                            ║
╚════════════════════════════════════════════╝

Создание администратора:
Username: admin
Password: ********
Email (опционально): admin@example.com

✓ Пользователь создан!

Настройка API токенов (можно пропустить):
1. OpenAI (для STT) [Enter = пропустить]: sk-...
2. Anthropic (для AI) [Enter = пропустить]: sk-ant-...
3. ElevenLabs (для TTS) [Enter = пропустить]: ...

✓ Конфигурация завершена!
```

---

## 12. README структура

**README.md** должен содержать:

```markdown
# Voice Git Assistant

## Описание
Консольное Java-приложение для управления git-репозиториями через голосовые и текстовые команды с использованием ИИ.

## Возможности
- 🎤 Голосовой ввод команд (STT через Whisper)
- 🤖 Понимание естественного языка (AI через Claude)
- 📊 Статистика и аналитика по репозиториям
- 🔄 Автокоммиты по расписанию
- 🔒 Безопасное хранение токенов
- 📝 История всех операций

## Требования
- Java 17+
- Git (для локальных репозиториев)
- API ключи: OpenAI, Anthropic, ElevenLabs (опционально)

## Установка

### Из исходников
```bash
git clone https://github.com/username/voice-git-assistant.git
cd voice-git-assistant
mvn clean package
java -jar target/voice-git-assistant.jar
```

### Готовый JAR
```bash
wget https://github.com/username/voice-git-assistant/releases/download/v1.0/voice-git-assistant.jar
java -jar voice-git-assistant.jar
```

## Быстрый старт

1. Регистрация:
```
> Регистрация
Username: yourname
Password: ********
```

2. Добавление репозитория:
```
> Мои репозитории > Добавить
Название: my-project
Путь: /home/user/projects/my-project
Remote URL: https://github.com/user/my-project.git
```

3. Голосовая команда:
```
> Голосовой ассистент > Выбрать репозиторий > my-project
🎤 "Закоммить изменения, я добавил новую фичу"
✅ Готово! Коммит c4f2b9a запушен
```

## Команды

### Поддерживаемые git-операции
- Коммит: "Закоммить изменения"
- Пуш: "Запушить на main"
- Лог: "Покажи последние 10 коммитов"
- Статус: "Что изменилось?"
- Ветки: "Создай ветку feature-X"
- Diff: "Покажи изменения с утра"

## Архитектура

```
Presentation Layer (CLI) → Application Layer (Services) → Domain Layer (Models) → Infrastructure (DB/APIs)
```

## Технологии
- Java 17
- JGit
- SQLite
- OpenAI Whisper API
- Anthropic Claude API
- ElevenLabs API
- JUnit 5 + Mockito

## Конфигурация

`~/.voicegit/application.properties`:
```properties
db.path=./data/voicegit.db
stt.default.language=ru
tts.default.voice=Rachel
```

## Тестирование

```bash
mvn test
mvn jacoco:report  # Coverage report
```

## Лицензия
MIT

## Контакты
[ваш email]


---

## 13. Критерии приёмки

### 13.1 Обязательные функции
- ✅ Регистрация и авторизация
- ✅ CRUD репозиториев
- ✅ Базовые git-операции (commit, push, log, status)
- ✅ Голосовой ввод через STT
- ✅ AI-парсинг команд
- ✅ История операций с фильтрацией
- ✅ Статистика по репозиториям
- ✅ Автокоммиты по расписанию
- ✅ Экспорт данных (CSV/JSON)

### 13.2 Качество кода
- ✅ Минимум 70% test coverage
- ✅ Все методы документированы (JavaDoc)
- ✅ Нет дублирования кода
- ✅ SOLID принципы соблюдены
- ✅ Exception handling на всех уровнях
- ✅ Логирование ошибок

### 13.3 Документация
- ✅ README с полным описанием
- ✅ Инструкция по установке
- ✅ Примеры использования
- ✅ Архитектурная диаграмма
- ✅ Описание API интеграций

### 13.4 Демонстрация
Готовность показать на защите:
1. Регистрацию и вход
2. Добавление репозитория
3. Голосовую команду (коммит)
4. Просмотр истории
5. Статистику
6. Автокоммит в действии
7. Обработку ошибок (неверный токен, нет изменений)

---

## 14. Риски и ограничения

### 14.1 Технические риски

**API Rate Limits**:
- OpenAI: 3 RPM на бесплатном тарифе
- Anthropic: 5 RPM на бесплатном
- Решение: кэширование, локальные альтернативы (Vosk для STT)

**Сетевые проблемы**:
- Все API требуют интернет
- Решение: graceful degradation, fallback на текстовый ввод

**Безопасность токенов**:
- Токены в БД, даже зашифрованные — риск
- Решение: использовать систему credential storage ОС (Keychain/Credential Manager)

### 14.2 Функциональные ограничения

**Git операции**:
- Не поддерживается интерактивный rebase
- Нет визуального merge conflict resolver
- Решение: документировать и предлагать ручное решение

**Языки**:
- STT/TTS поддерживает только русский и английский
- Решение: расширяемая архитектура для добавления языков

**Автокоммит**:
- Может создавать "шумные" коммиты
- Решение: умные commit messages через AI, настройка минимального интервала

---

## 15. Roadmap (будущие версии)

### v1.1
- [ ] Поддержка GitLab/Bitbucket
- [ ] Web UI (Spring Boot + React)
- [ ] Уведомления (email/Telegram)
- [ ] Multi-repository bulk operations

### v1.2
- [ ] AI code review
- [ ] Автоматическое создание PR
- [ ] Integration с CI/CD (GitHub Actions)
- [ ] Команды для Docker operations

### v2.0
- [ ] Локальная LLM (Ollama)
- [ ] Офлайн STT (Vosk)
- [ ] Плагинная система
- [ ] REST API для интеграций

---

## 16. Заключение

Проект **Voice Git Assistant** соответствует всем требованиям курса и превосходит их по сложности благодаря:

1. **Многослойной архитектуре** с чёткой ответственностью
2. **Интеграции 4 внешних API** (Whisper, Claude, ElevenLabs, GitHub)
3. **Сложной бизнес-логике** (AI-парсинг, планировщик, аналитика)
4. **Безопасности** (шифрование, хеширование, валидация)
5. **Масштабируемости** (поддержка множества пользователей и репозиториев)

Проект демонстрирует владение Java, ООП принципами, работой с БД, внешними API, тестированием и документированием.

---

**Дата создания ТЗ**: 20.05.2026  
**Версия**: 1.0  
**Автор**: МишМиш + Клод