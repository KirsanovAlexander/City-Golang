## City Simulation API (Golang)

Ин-мемори сервис для симуляции виртуального города: строительство зданий, управление жителями, ресурсы (еда, энергия, деньги), события и дневные тики.

### Быстрый старт

1) Требования: Go 1.22+

2) Установка зависимостей (в проекте уже выполнено):
```bash
go mod tidy
```

3) Запуск сервера:
```bash
go run ./cmd/server
```
По умолчанию сервер слушает порт 8080. Можно переопределить переменной окружения `PORT`.

Windows PowerShell: команды разделяйте `;`, пример:
```powershell
cd C:\programming\projects\city; go run ./cmd/server
```

### Архитектура проекта

```
cmd/
  server/
    main.go                # вход, поднимает HTTP сервер
internal/
  router/
    router.go              # маршрутизация (Gin)
  models/
    types.go               # DTO/модели: City, Building, Citizen, Event, ...
  storage/
    memory.go              # ин-мемори состояние города, мьютексы, история
  handlers/
    city.go                # ручки города
    buildings.go           # ручки зданий
    citizens.go            # ручки жителей
    resources.go           # ручки экономики/ресурсов
    simulation.go          # ручки симуляции/событий/статистики
```

Хранение данных: в памяти процесса с историей ресурсов и событий; при перезапуске всё сбрасывается.

### Сущности (кратко)

- City: `settings`, `resources {food, energy, money}`, `buildings[]`, `citizens[]`, `day`.
- Building: `id`, `type: house|farm|factory|powerplant`, `level`, `health`, `createdAt`.
- Citizen: `id`, `name`, `job: unemployed|farmer|worker|engineer`, `happiness`, `createdAt`.
- Event: `id`, `type`, `message`, `delta{food,energy,money}`, `happinessDelta`, `createdAt`.

### Настройки и базовая экономика

- По умолчанию `difficulty=normal`.
- Производство за тик:
  - Фермеры дают `Food += FoodPerFarmer (по умолчанию 10)`
  - Рабочие дают `Money += MoneyPerWorker (8)`
  - Инженеры дают `Energy += EnergyPerEngineer (6)`
- Эффекты зданий за тик (на уровень):
  - `farm`: `+5 food`
  - `factory`: `+7 money`, `-3 energy`
  - `powerplant`: `+10 energy`
  - `house`: моральное влияние (косвенно)
- Потребление: `Food -= BaseFoodUsePerCap (1) * население`.
- Нехватка еды снижает счастье и может уменьшать население.

### Полный список эндпоинтов

Базовый URL: `http://localhost:8080`

#### Город
- POST `/city` — создать город
- GET `/city` — текущее состояние города
- DELETE `/city/reset` — сбросить город
- PATCH `/city/settings` — изменить настройки города

#### Здания
- POST `/city/buildings` — построить здание
- PATCH `/city/buildings/{id}/upgrade` — улучшить здание
- PATCH `/city/buildings/{id}/repair` — починить здание
- DELETE `/city/buildings/{id}` — снести здание
- GET `/city/buildings` — список зданий
- GET `/city/buildings/effects` — справка по эффектам зданий

#### Жители
- POST `/city/citizens` — добавить жителя
- GET `/city/citizens` — список жителей
- PATCH `/city/citizens/{id}/job` — сменить профессию
- PATCH `/city/citizens/{id}/happiness` — изменить счастье
- DELETE `/city/citizens/{id}` — удалить жителя
- GET `/city/citizens/jobs` — статистика профессий
- POST `/city/citizens/mass-add` — массовое добавление

#### Экономика и ресурсы
- POST `/city/trade` — торговля ресурсами
- GET `/city/resources/history` — история ресурсов по дням
- PATCH `/city/resources/adjust` — админское изменение ресурсов

#### Симуляция и события
- POST `/city/tick` — прожить один день
- POST `/city/events/random` — случайное событие
- POST `/city/events/custom` — кастомное событие
- GET `/city/events/history` — история событий
- GET `/city/stats` — агрегированная статистика

### Форматы запросов/ответов (основные DTO)

- Создать город
```json
POST /city
{
  "name": "NeoCity",
  "difficulty": "normal" // normal|hard
}
```

- Изменить настройки
```json
PATCH /city/settings
{
  "name": "My City",
  "foodPerFarmer": 12,
  "moneyPerWorker": 9,
  "energyPerEngineer": 7,
  "baseFoodUsePerCap": 1,
  "happinessDecay": 0.5
}
```

- Построить здание
```json
POST /city/buildings
{
  "type": "farm" // farm|house|factory|powerplant
}
```

- Отремонтировать здание
```json
PATCH /city/buildings/{id}/repair
{ "amount": 15 }
```

- Добавить жителя
```json
POST /city/citizens
{
  "name": "Alice",
  "job": "farmer" // unemployed|farmer|worker|engineer
}
```

- Изменить профессию
```json
PATCH /city/citizens/{id}/job
{ "job": "engineer" }
```

- Массовое добавление
```json
POST /city/citizens/mass-add
{
  "count": 5,
  "job": "worker",
  "prefix": "Migrants"
}
```

- Торговля
```json
POST /city/trade
{
  "resource": "food",   // food|energy|money
  "amount": 10,          // >0 купить, <0 продать (money начисляется/списывается)
  "price": 2             // цена за единицу
}
```

- Кастомное событие
```json
POST /city/events/custom
{
  "type": "festival",
  "message": "Big party!",
  "delta": {"food": -10, "energy": -5, "money": -20},
  "happinessDelta": 12
}
```

### Примеры curl

```bash
# Создать город
curl -X POST http://localhost:8080/city \
  -H "Content-Type: application/json" \
  -d '{"name":"NeoCity","difficulty":"normal"}'

# Построить ферму
curl -X POST http://localhost:8080/city/buildings \
  -H "Content-Type: application/json" \
  -d '{"type":"farm"}'

# Добавить жителя-фермера
curl -X POST http://localhost:8080/city/citizens \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","job":"farmer"}'

# Тик (прожить день)
curl -X POST http://localhost:8080/city/tick

# Случайное событие
curl -X POST http://localhost:8080/city/events/random

# Статистика
curl http://localhost:8080/city/stats
```

### Статусы ответов и ошибки

- 201 Created — успешное создание сущности
- 200 OK — успешный запрос
- 204 No Content — успешное удаление/сброс
- 400 Bad Request — неверные данные запроса
- 404 Not Found — город или сущность ещё не создана/не найдена

### Тестирование

Минимальные юнит-тесты можно строить поверх `internal/storage/memory.go`, мокая/инициализируя состояние через публичные методы. Для интеграционных тестов — запускать Gin в тестовом режиме и вызывать ручки через `httptest`.

### Нагрузочные/производительность

Хранилище — один инстанс с `sync.RWMutex`. Под высокую нагрузку и множественные города потребуется вынести состояние в БД/кэш и добавить шардинг/мульти-инстансы.

### Возможные расширения

- Персистентность (SQLite/PostgreSQL)
- Авторизация и роли
- Ограничения ресурсов/стоимости строительства
- Планировщик автотиков
- Метрики Prometheus и профилирование

### Лицензия

MIT (по желанию можно заменить).


