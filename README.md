# GO Framework GUAP.RU

> Go | API Testing | CI/CD | guap.ru

## Статус

![CI](https://github.com/ssrjkk/go-framework-guap/actions/workflows/ci.yml/badge.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)

## Обзор

Фреймворк для автоматизированного тестирования API портала [guap.ru](https://guap.ru) — системы Санкт-Петербургского государственного университета аэрокосмического приборостроения.

Тестируемые модули:
- `/api/auth` — авторизация, refresh token
- `/api/profile` — профиль студента
- `/api/schedule` — расписание занятий
- `/api/teachers` — преподаватели
- `/api/students/*/grades` — оценки студентов

## Архитектура

```
.
├── core/               # Базовый слой
│   ├── base/          # HTTP клиент с retry, валидатор
│   ├── errors/        # APIError, ValidationError, RetryableError
│   └── utils/         # Request/Response логирование
├── services/api/      # AuthService, ScheduleService, GradesService, ProfileService
├── fixtures/          # APIClient, AuthFixture, ScheduleFixture (DI)
├── config/            # dev/stage окружения
├── tests/             # smoke, regression, critical
└── docker/            # Multi-stage build
```

## Тестируемые сценарии

| Модуль | Методы | Тесты |
|--------|--------|-------|
| Auth | POST /api/auth/login, refresh, logout | Вход, выход, refresh token |
| Profile | GET /api/profile, PATCH | Получение и обновление профиля |
| Schedule | GET /api/schedule/{group} | Расписание группы, по дате |
| Teachers | GET /api/teachers, /{id}/schedule | Список, расписание преподавателя |
| Grades | GET /api/students/{id}/grades, /gpa | Оценки, средний балл |

## Уровни тестов

| Уровень | Описание | Время |
|---------|----------|-------|
| **Smoke** | Базовая доступность эндпоинтов | < 1 мин |
| **Critical** | Авторизация, получение данных | < 3 мин |
| **Regression** | Полное функциональное покрытие | < 5 мин |
| **Nightly** | Полный набор + race detection | Ежедневно |

## Фичи

- **Retry logic**: Exponential backoff на 5xx/429
- **Логирование**: Request/Response с headers и body
- **Schema validation**: Required, email, min/max length
- **Fixtures**: DI без `new` в тестах
- **Fail-fast**: false (все параллельные jobs завершаются)
- **Docker**: Multi-stage ~15MB

## Запуск

```bash
# Все тесты
go test ./tests/...

# Smoke
go test -v -run "^TestSmoke" ./tests/smoke/...

# Regression
go test -v -run "^TestRegression" ./tests/regression/...

# Critical + Negative
go test -v -run "^(TestCritical|TestNegative)" ./tests/critical/...

# Параллельно
go test -parallel 4 ./tests/...
```

## CI/CD Pipeline

```
lint → test (smoke) ─┐
                    ├→ test (regression) ─┤
                    └→ test (critical) ────┴→ nightly (cron)
```

## Docker

```bash
docker build -f docker/Dockerfile -t go-framework-guap .
docker run go-framework-guap
```

## Конфигурация

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `API_BASE_URL` | Base URL | guap.ru |
| `API_TIMEOUT` | Таймаут (сек) | 30 |
| `API_MAX_RETRIES` | Retry попыток | 3 |

## Контакты

- Telegram: @ssrjkk
- Email: ray013lefe@gmail.com
