# AI Search Visibility Platform

Платформа проверяет, как ИИ-поиск (ChatGPT и другие нейросети) отвечает про бренд и конкурентов, откуда ИИ берёт ответы, и подсказывает, что изменить, чтобы ИИ рекомендовал именно ваш бренд.

> Хакатон, **Кейс 7 · «Как быть первым в ИИ-поиске»** (GEO / AI-SEO). Команда: **AbsoluteCinema**

## Что делает

1. Проверяет, что ИИ отвечает про бренд и конкурентов.
2. Находит источники, откуда ИИ берёт ответы.
3. Предлагает правки: тексты, карточки, отзывы, структуру сайта.
4. Делает повторный замер и показывает **до/после**.

## Команда

| Участник | Роль | Что сделал |
|---|---|---|
| **Абзал** | Go Backend | Сервер на Go (Gin), база PostgreSQL, авторизация (JWT), CRUD компаний, управление анализами, связь с AI Service, Docker Compose |
| **Ертай** | Python AI Engineer | AI Service на FastAPI, поиск через Tavily, запросы к OpenAI, определение упоминаний бренда, конкуренты, Visibility Score, рекомендации |
| **Айбар** | Frontend | React + TypeScript, вход и регистрация, Dashboard, компании, страницы анализа, конкурентов и рекомендаций, графики |

Вместе: контракт API, интеграция сервисов, тесты, README, демо.

## Стек

| Часть | Технологии |
|---|---|
| Backend | Go, Gin, PostgreSQL, pgx, sqlc, Goose |
| AI Service | Python, FastAPI, OpenAI API, Tavily API |
| Frontend | React, TypeScript, Vite, Tailwind CSS, Recharts |
| Инфраструктура | Docker, Docker Compose |

## Архитектура

```
Frontend (React)  ->  Go Backend (Gin)  ->  Python AI Service (FastAPI)
                           |                        |
                       PostgreSQL            OpenAI + Tavily
```

Frontend общается только с Go Backend. К Python напрямую не обращается.

## Запуск

```bash
cp .env.example .env
docker-compose -f deployments/docker-compose.yaml up --build
```

Внутри Docker сервисы видят друг друга по имени: база `postgres`, AI-сервис `http://ai-service:8000`.

| Сервис | Адрес |
|---|---|
| Backend | http://localhost:8080 (если занят, поменять на `8081:8080` в `docker-compose.yaml`) |
| AI Service | http://localhost:8000/docs |

Остановка: `docker-compose -f deployments/docker-compose.yaml down`

## Главные endpoints

| Метод | Endpoint | Назначение |
|---|---|---|
| POST | `/api/auth/register`, `/api/auth/login` | Регистрация и вход |
| GET/POST | `/api/companies` | Список и добавление компаний |
| POST | `/api/analyses` | Запустить анализ |
| GET | `/api/analyses/:id` | Результат анализа |
| GET | `/api/analyses/:id/competitors` | Конкуренты |
| GET | `/api/analyses/:id/responses` | Ответы ИИ |
| GET | `/api/analyses/:id/recommendations` | Рекомендации |
| POST | `/internal/analyze` | Go -> Python: полный анализ |

Статусы анализа: `pending` -> `processing` -> `completed` / `failed`.

## Visibility Score

```
Visibility Score = (M / Q) x 100%
```

M: ответы ИИ, где упомянут бренд. Q: всего проверенных ответов. Пример: 7 из 20 ответов = **35%**. Это внутренняя метрика MVP, а не официальный рейтинг.

## Демо (3 минуты)

1. Добавляем компанию.
2. Запускаем анализ и смотрим ответы ИИ.
3. Показываем источники и рекомендации.
4. Вносим правки и делаем повторный замер.
5. Сравниваем Visibility Score **до и после**.