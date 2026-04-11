# Guap.ru Framework (Go)

Portfolio Project

Фреймворк для автоматизации тестирования API веб-приложения guap.ru.
Демонстрирует навыки: API-автоматизация, CI/CD, контейнеризация, многоуровневое тестирование.

## Зачем этот проект

Автоматизирует проверку ключевых сценариев образовательного портала ГУАП:

| Сценарий | Тип теста | Зачем |
|----------|-----------|-------|
| Проверка статуса | Smoke | Быстрая проверка доступности API |
| Получение списка студентов | Smoke + Regression | Валидация списочных endpoint'ов |
| Расписание группы | Critical + Regression | Проверка фильтрации по параметрам |
| Оценки студента | Regression | Валидация связанных данных |
| Авторизация | Critical | Проверка security |
| Ошибки 4xx/5xx | Negative | Контрактное тестирование |

## Быстрый старт

```bash
# Клонировать
git clone https://github.com/ssrjkk/guap-test-framework-go.git
cd guap-test-framework-go

# Запустить тесты
go test ./tests/smoke/...        # Smoke тесты
go test ./tests/regression/...   # Regression тесты
go test ./tests/critical/...     # Critical + Negative тесты

# Параллельно
go test -parallel 4 ./tests/...

# Docker
docker build -f docker/Dockerfile -t qa-tests .
docker run --rm qa-tests
```

## Архитектура

```
.
├── core/                   # Базовый слой
│   ├── base/               # HTTP-клиент: retry, логирование, валидация
│   │   ├── client.go       # BaseClient с retry и timeout
│   │   └── validator.go    # Schema validation (required, email, min/max)
│   ├── errors/             # Кастомные ошибки
│   │   └── errors.go       # APIError, ValidationError, RetryableError
│   └── utils/              # Логирование
│       └── logger.go       # Request/Response логирование
├── services/api/           # API сервисы
│   └── services.go         # HealthService, AuthService, StudentService,
│                           # ScheduleService, SubjectService, GradesService
├── fixtures/               # Фикстуры с DI
│   └── api.go              # APIClient, AuthFixture, ScheduleFixture
├── config/                 # Конфигурация
│   └── config.go           # dev/stage окружения из .env
├── tests/                  # Тесты
│   ├── smoke/              # Smoke тесты (< 1 мин)
│   ├── regression/         # Regression тесты (< 5 мин)
│   ├── critical/           # Critical + Negative тесты (< 3 мин)
│   └── tests.go            # Утилиты: метрики, retry, waiters
└── docker/                 # Контейнеризация
    ├── Dockerfile          # Production image
    └── Dockerfile.test     # Test image для CI
```

## Пример теста

```go
// tests/regression/students_test.go
func TestRegressionStudentHasRequiredFields(t *testing.T) {
    client := fixtures.NewAPIClient(fixtures.GetEnv())
    client.Init()

    students, err := client.StudentService().GetAll(ctx, token)
    if err != nil {
        t.Skipf("Students not available: %v", err)
    }

    for _, student := range students {
        if student.ID == "" {
            t.Error("Student has no ID")
        }
        if student.Name == "" {
            t.Errorf("Student %s has no name", student.ID)
        }
    }
}
```

```go
// services/api/services.go
type StudentService struct {
    client *base.Client
}

func (s *StudentService) GetAll(ctx context.Context, token string) ([]Student, error) {
    resp, err := s.client.Get(ctx, "/api/students", nil,
        base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
    if err != nil {
        return nil, err
    }

    var students []Student
    if err := s.client.DecodeJSON(resp, &students); err != nil {
        return nil, err
    }
    return students, nil
}
```

## Тестируемые эндпоинты

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| GET | `/api/health` | Проверка статуса |
| GET | `/api/students` | Список студентов |
| GET | `/api/students/:id` | Студент по ID |
| GET | `/api/schedule` | Расписание |
| GET | `/api/schedule?group=Z3420` | Расписание группы |
| GET | `/api/subjects` | Список предметов |
| GET | `/api/grades` | Все оценки |
| GET | `/api/grades?student_id=1` | Оценки студента |
| GET | `/api/teachers` | Преподаватели |
| POST | `/api/auth/login` | Авторизация |
| POST | `/api/auth/refresh` | Refresh token |

## Стек

| Технология | Зачем |
|-----------|-------|
| Go 1.21+ | Язык с встроенным testing framework |
| go testing | Оркестрация, фикстуры, параллельность |
| Docker | Multi-stage build, воспроизводимое окружение |
| GitHub Actions | CI/CD: lint -> tests -> report |
| Retry + Backoff | Обработка нестабильных endpoint'ов |
| Schema validation | Контрактное тестирование |

## Уровни тестов

```bash
# Smoke — < 1 мин
go test -v -timeout 2m -run "^TestSmoke" ./tests/smoke/...

# Regression — < 5 мин
go test -v -timeout 5m -run "^TestRegression" ./tests/regression/...

# Critical + Negative — < 3 мин
go test -v -timeout 5m -run "^(TestCritical|TestNegative)" ./tests/critical/...

# Nightly — все тесты + race detection
go test -v -race -count=1 -timeout 10m ./tests/...
```

## CI/CD

GitHub Actions автоматически запускает:

```
lint ──→ test (smoke) ──┐
                       ├─→ test (regression) ──┤
                       └─→ test (critical) ─────┴──→ nightly (cron)
```

- **lint**: go vet, gofmt, staticcheck
- **test**: параллельное выполнение по уровням
- **nightly**: полный набор + race detection (cron: 2:00)

## Запуск в Docker

```bash
# Production image (~15MB)
docker build -f docker/Dockerfile -t qa-tests .
docker run --rm qa-tests

# Test image для CI
docker build -f docker/Dockerfile.test -t qa-tests-runner .
docker run --rm qa-tests-runner go test -v ./tests/...
```

## Конфигурация

| Переменная | Описание | По умолчанию |
|-----------|---------|--------------|
| `API_BASE_URL` | Base URL API | guap.ru |
| `API_TIMEOUT` | Таймаут (сек) | 30 |
| `API_MAX_RETRIES` | Retry попыток | 3 |
| `TEST_ENV` | Окружение (dev/stage) | dev |

## Контакты

Telegram: @ssrjkk
Email: ray013lefe@gmail.com
