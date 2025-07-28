
# Full Project Mock — Go Auth Service

Продакшен-монолит с авторизацией на Go, реализующий безопасную работу с JWT (access/refresh), Redis-сессиями и PostgreSQL. Проект построен на принципах Clean Architecture и покрыт тестами (unit/integration).

---

## Возможности

- Регистрация и логин пользователя
- Генерация access/refresh токенов (RSA, TTL)
- Обновление access-токена по refresh
- Выход с одного или всех устройств
- Redis-реализация session store с TTL и hash-идентификацией
- Получение информации о текущем пользователе (Me)
- Поддержка graceful shutdown
- Тесты всех слоёв (handler, usecase, repository, middleware)
- Контейнеризация с Docker Compose

---

## Архитектура

Проект построен по принципам **чистой архитектуры**:

```
cmd/                - точка входа (main.go)
internal/app        - запуск и shutdown
internal/bootstrap  - DI и сборка всех зависимостей
internal/module     - модули (handler/usecase/repo)
internal/domain     - контракты, модели, интерфейсы
internal/service    - токены, логгер, traceID
internal/middleware - авторизация, логгирование
pkg/                - переиспользуемые пакеты (hasher, validator и др.)
```

Зависимости между слоями инвертированы, бизнес-логика не зависит от инфраструктуры.

---

## Технологии

- Go 1.22+
- PostgreSQL
- Redis
- JWT (RS256)
- chi router
- slog (structured logging)
- Docker + Docker Compose
- go-redis v9
- bcrypt (хэширование паролей)

---

## API endpoints (v1)

| Метод | Путь               | Описание                    |
|-------|--------------------|-----------------------------|
| POST  | `/register`        | Регистрация пользователя    |
| POST  | `/login`           | Вход и получение токенов    |
| POST  | `/refresh`         | Обновление access-токена    |
| POST  | `/auth/logout`     | Выход с текущего устройства |
| POST  | `/auth/logout_all` | Выход со всех устройств     |
| GET   | `/auth/me`         | Получение ID пользователя   |

---

## Тесты

- Unit-тесты (usecase, middleware)
- Integration-тесты (handler → usecase → repository)
- `TestMain` с rollback транзакций и очисткой Redis

```bash
make test
```

---

## Запуск

```bash
# Поднять окружение
make docker-up

# Применить миграции
make migrate

# Остановить окружение
make docker-down
```

---

## Безопасность

- Refresh-токены хранятся в Redis в виде **хэшей**
- Redis TTL для удаления по времени
- Возможность инвалидации токена по ID
- RSA-ключи для access-токенов
- Пароли хэшируются с bcrypt

---

## Структура проекта

<details>
<summary>Нажми, чтобы развернуть</summary>

```
full-project-mock/
├── cmd/                  # main.go
├── config/               # config.yaml + .env
├── docker/               # миграции и окружение
├── internal/
│   ├── app/              # запуск
│   ├── bootstrap/        # DI
│   ├── config/           # конфиги
│   ├── database/         # подключение к PostgreSQL и Redis
│   ├── delivery/rest/    # роутинг
│   ├── domain/           # модели, интерфейсы
│   ├── middleware/       # middleware
│   ├── module/user/      # handler/usecase/repo
│   └── service/          # токены, логгер, trace
├── migrations/           # SQL
├── pkg/                  # утилиты
├── Makefile
└── docker-compose.yml
```

</details>
