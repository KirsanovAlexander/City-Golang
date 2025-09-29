# City Management API Documentation

## Обзор

City Management API - это REST API для управления симуляцией города. API позволяет создавать и управлять городом, его зданиями, жителями, ресурсами и событиями.

## Базовый URL

```
http://localhost:8080
```

## Swagger документация

Интерактивная документация доступна по адресу: `http://localhost:8080/swagger/`

## Архитектура

### Структура проекта

```
city/
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── handlers/        # HTTP обработчики (контроллеры)
│   ├── services/        # Бизнес-логика (use case)
│   ├── models/          # Модели данных
│   │   ├── domain.go    # Бизнес-модели
│   │   └── dto.go       # DTO для API
│   ├── storage/         # Слой данных
│   └── router/          # Маршрутизация
└── docs/               # Swagger документация
```

### Поток выполнения запросов

1. **HTTP Request** → **Router** → **Handler** → **Service** → **Storage**
2. **Response** ← **Storage** ← **Service** ← **Handler** ← **Router** ← **HTTP Response**

### Компоненты

- **Handlers**: Обрабатывают HTTP запросы, валидируют входные данные, вызывают сервисы
- **Services**: Содержат бизнес-логику, работают с доменными моделями
- **Storage**: Интерфейс для работы с данными (в памяти, БД и т.д.)
- **Models**: 
  - `domain.go` - бизнес-модели
  - `dto.go` - модели для API запросов/ответов

## Эндпоинты

### Город (City)

#### Создать город
```http
POST /city
Content-Type: application/json

{
  "name": "My City",
  "difficulty": "normal"
}
```

**Ответ:**
```json
{
  "day": 0,
  "settings": {
    "name": "My City",
    "difficulty": "normal",
    "foodPerFarmer": 10,
    "moneyPerWorker": 8,
    "energyPerEngineer": 6,
    "baseFoodUsePerCap": 1,
    "happinessDecay": 0.5
  },
  "resources": {
    "food": 100,
    "energy": 100,
    "money": 100
  },
  "buildings": [],
  "citizens": []
}
```

#### Получить информацию о городе
```http
GET /city
```

#### Сбросить город
```http
DELETE /city/reset
```

#### Обновить настройки города
```http
PATCH /city/settings
Content-Type: application/json

{
  "name": "New City Name",
  "difficulty": "hard"
}
```

### Здания (Buildings)

#### Создать здание
```http
POST /city/buildings
Content-Type: application/json

{
  "type": "farm"
}
```

**Типы зданий:**
- `house` - жилой дом (влияет на мораль)
- `farm` - ферма (+5 еды за уровень в день)
- `factory` - фабрика (+7 денег, -3 энергии за уровень в день)
- `powerplant` - электростанция (+10 энергии за уровень в день)

#### Получить список зданий
```http
GET /city/buildings
```

#### Получить эффекты зданий
```http
GET /city/buildings/effects
```

#### Улучшить здание
```http
PATCH /city/buildings/{id}/upgrade
```

#### Починить здание
```http
PATCH /city/buildings/{id}/repair
Content-Type: application/json

{
  "amount": 10
}
```

#### Удалить здание
```http
DELETE /city/buildings/{id}
```

### Жители (Citizens)

#### Создать жителя
```http
POST /city/citizens
Content-Type: application/json

{
  "name": "John Doe",
  "job": "farmer"
}
```

**Профессии:**
- `unemployed` - безработный
- `farmer` - фермер (+10 еды в день)
- `worker` - рабочий (+8 денег в день)
- `engineer` - инженер (+6 энергии в день)

#### Получить список жителей
```http
GET /city/citizens
```

#### Изменить профессию жителя
```http
PATCH /city/citizens/{id}/job
Content-Type: application/json

{
  "job": "worker"
}
```

#### Изменить счастье жителя
```http
PATCH /city/citizens/{id}/happiness
Content-Type: application/json

{
  "delta": 5.0
}
```

#### Удалить жителя
```http
DELETE /city/citizens/{id}
```

#### Получить статистику по профессиям
```http
GET /city/citizens/jobs
```

**Ответ:**
```json
{
  "unemployed": 2,
  "farmer": 5,
  "worker": 3,
  "engineer": 1
}
```

#### Массовое добавление жителей
```http
POST /city/citizens/mass-add
Content-Type: application/json

{
  "count": 5,
  "job": "farmer",
  "prefix": "Citizen"
}
```

### Ресурсы (Resources)

#### Торговля ресурсами
```http
POST /city/trade
Content-Type: application/json

{
  "resource": "food",
  "amount": 10,
  "price": 2
}
```

**Ресурсы:**
- `food` - еда
- `energy` - энергия
- `money` - деньги

#### Получить историю ресурсов
```http
GET /city/resources/history
```

#### Настроить ресурсы
```http
PATCH /city/resources/adjust
Content-Type: application/json

{
  "food": 150,
  "energy": 200,
  "money": 300
}
```

### Симуляция (Simulation)

#### Прогресс на один день
```http
POST /city/tick
```

**Что происходит при тике:**
1. Производство ресурсов жителями
2. Пассивные эффекты зданий
3. Потребление еды населением
4. Естественное снижение счастья
5. Увеличение дня на 1

#### Случайное событие
```http
POST /city/events/random
```

**Типы событий:**
- `storm` - шторм (повреждения зданий, потеря ресурсов)
- `festival` - фестиваль (повышение счастья, трата ресурсов)
- `population_growth` - рост населения (новые жители)
- `blackout` - отключение электричества (потеря энергии)

#### Пользовательское событие
```http
POST /city/events/custom
Content-Type: application/json

{
  "type": "festival",
  "message": "Custom festival",
  "delta": {
    "food": -5,
    "energy": -5,
    "money": -10
  },
  "happinessDelta": 10.0
}
```

#### История событий
```http
GET /city/events/history
```

#### Статистика города
```http
GET /city/stats
```

**Ответ:**
```json
{
  "day": 5,
  "population": 10,
  "avgHappiness": 75.5,
  "foodBalance": 150,
  "energyBalance": 200,
  "moneyBalance": 300
}
```

## Игровая механика

### Сложность игры

- **Easy**: Стандартные параметры
- **Normal**: Стандартные параметры
- **Hard**: 
  - Меньше производства ресурсов
  - Больше потребления еды
  - Быстрее снижается счастье

### Система ресурсов

1. **Еда**: Производится фермерами и фермами, потребляется населением
2. **Энергия**: Производится инженерами и электростанциями, потребляется фабриками
3. **Деньги**: Производится рабочими и фабриками, тратятся на торговлю

### Система счастья

- Начальное счастье: 80
- Естественное снижение: 0.5 в день
- Влияние голода: дополнительное снижение при нехватке еды
- События могут влиять на счастье

### Система зданий

- Уровень: влияет на эффективность
- Здоровье: 0-100, снижается от событий
- Ремонт: восстанавливает здоровье

## Коды ошибок

- `200` - Успех
- `201` - Создано
- `204` - Нет содержимого
- `400` - Неверный запрос
- `404` - Не найдено
- `500` - Внутренняя ошибка сервера

## Примеры использования

### Создание и развитие города

1. **Создать город:**
```bash
curl -X POST http://localhost:8080/city \
  -H "Content-Type: application/json" \
  -d '{"name": "My City", "difficulty": "normal"}'
```

2. **Добавить жителей:**
```bash
curl -X POST http://localhost:8080/city/citizens \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "job": "farmer"}'

curl -X POST http://localhost:8080/city/citizens \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane", "job": "worker"}'
```

3. **Построить здания:**
```bash
curl -X POST http://localhost:8080/city/buildings \
  -H "Content-Type: application/json" \
  -d '{"type": "farm"}'

curl -X POST http://localhost:8080/city/buildings \
  -H "Content-Type: application/json" \
  -d '{"type": "house"}'
```

4. **Прогресс времени:**
```bash
curl -X POST http://localhost:8080/city/tick
```

5. **Проверить статистику:**
```bash
curl http://localhost:8080/city/stats
```

### Управление ресурсами

1. **Торговля:**
```bash
# Купить еду
curl -X POST http://localhost:8080/city/trade \
  -H "Content-Type: application/json" \
  -d '{"resource": "food", "amount": 20, "price": 2}'

# Продать энергию
curl -X POST http://localhost:8080/city/trade \
  -H "Content-Type: application/json" \
  -d '{"resource": "energy", "amount": -10, "price": 3}'
```

2. **Настройка ресурсов:**
```bash
curl -X PATCH http://localhost:8080/city/resources/adjust \
  -H "Content-Type: application/json" \
  -d '{"food": 200, "money": 500}'
```

### События

1. **Случайное событие:**
```bash
curl -X POST http://localhost:8080/city/events/random
```

2. **Пользовательское событие:**
```bash
curl -X POST http://localhost:8080/city/events/custom \
  -H "Content-Type: application/json" \
  -d '{
    "type": "festival",
    "message": "City Festival",
    "delta": {"food": -10, "money": -20},
    "happinessDelta": 15.0
  }'
```

## Запуск и тестирование

### Запуск сервера

```bash
go run cmd/server/main.go
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# Тесты конкретного пакета
go test ./internal/handlers
go test ./internal/services

# Тесты с покрытием
go test -cover ./...
```

### Генерация Swagger документации

```bash
# Установка swag
go install github.com/swaggo/swag/cmd/swag@latest

# Генерация документации
swag init -g cmd/server/main.go -o docs
```

## Технические детали

### Зависимости

- **chi/v5** - HTTP роутер
- **chi/cors** - CORS middleware
- **swaggo/swag** - Swagger генерация
- **swaggo/http-swagger** - Swagger UI
- **stretchr/testify** - Тестирование
- **google/uuid** - Генерация UUID

### Структура данных

Все данные хранятся в памяти (in-memory storage). При перезапуске сервера данные теряются.

### Безопасность

- CORS настроен для всех источников
- Валидация входных данных
- Обработка ошибок

### Производительность

- In-memory storage для быстрого доступа
- Минимальные зависимости
- Эффективная маршрутизация с chi

## Расширения

### Возможные улучшения

1. **Персистентность данных:**
   - Добавить поддержку PostgreSQL/MySQL
   - Реализовать миграции

2. **Аутентификация:**
   - JWT токены
   - Роли пользователей

3. **Кэширование:**
   - Redis для кэширования
   - Кэширование статистики

4. **Мониторинг:**
   - Метрики Prometheus
   - Логирование

5. **Масштабирование:**
   - Микросервисная архитектура
   - Message queues

### Добавление новых функций

1. **Новые типы зданий:**
   - Добавить в `BuildingType`
   - Реализовать логику в `Tick()`

2. **Новые события:**
   - Добавить в `EventType`
   - Реализовать в `RandomEvent()`

3. **Новые профессии:**
   - Добавить в `Job`
   - Реализовать производство в `Tick()`

## Заключение

City Management API предоставляет полный набор функций для управления симуляцией города. API спроектирован с использованием clean architecture принципов, что обеспечивает хорошую тестируемость и расширяемость.

Архитектура разделяет ответственность между слоями:
- **Handlers** - HTTP обработка
- **Services** - Бизнес-логика  
- **Storage** - Работа с данными
- **Models** - Структуры данных

Такая структура позволяет легко добавлять новые функции, тестировать компоненты и поддерживать код.

