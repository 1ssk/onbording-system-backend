# RoleLeader Service API

Микросервис для управления встречами и обратной связью между лидерами и пользователями.

## API Endpoints

### 1. Создание обратной связи для встречи

**Endpoint:**  
`POST /api/feedback`

**Request Body (JSON):**

```json
{
  "meeting_id": "meeting_123",
  "message": "Отличная встреча! Всё четко по делу."
}
```

**Response (Success):**

```json
{
  "status": "success"
}
```

---

### 2. Получение информации о встрече

**Endpoint:**  
`GET /api/leaders/meetings/{meeting_id}`

**Пример запроса:**

```
GET /api/leaders/meetings/meeting_123
```

**Response:**

```json
{
  "meeting": {
    "meeting_id": "meeting_123",
    "user_id": "user_456",
    "leader_id": "leader_789",
    "title": "Планирование Q4",
    "start_time": "2024-03-20T14:00:00Z",
    "status": "completed",
    "feedback": "Обсудили ключевые метрики"
  }
}
```

---

### 3. Получение всех встреч лидера

**Endpoint:**  
`GET /api/leaders/{leader_id}/meetings`

**Пример запроса:**

```
GET /api/leaders/leader_789/meetings
```

**Response:**

```json
{
  "meetings": [
    {
      "meeting_id": "meeting_123",
      "user_id": "user_456",
      "leader_id": "leader_789",
      "title": "Планирование Q4",
      "start_time": "2024-03-20T14:00:00Z",
      "status": "completed",
      "feedback": "Обсудили ключевые метрики"
    },
    {
      "meeting_id": "meeting_124",
      "user_id": "user_457",
      "leader_id": "leader_789",
      "title": "Разбор итогов спринта",
      "start_time": "2024-03-21T10:00:00Z",
      "status": "scheduled",
      "feedback": ""
    }
  ]
}
```

---

## Примеры cURL

### Создание обратной связи

```bash
curl -X POST "http://localhost:8080/api/feedback" \
-H "Content-Type: application/json" \
-d '{
  "meeting_id": "meeting_123",
  "message": "Отличная встреча!"
}'
```

### Получение встречи

```bash
curl "http://localhost:8080/api/leaders/meetings/meeting_123"
```

### Получение встреч лидера

```bash
curl "http://localhost:8080/api/leaders/leader_789/meetings"
```

---

## Структуры данных

### Объект Meeting

| Поле       | Тип    | Описание                          |
| ---------- | ------ | --------------------------------- |
| meeting_id | string | Уникальный ID встречи             |
| user_id    | string | ID пользователя                   |
| leader_id  | string | ID лидера                         |
| title      | string | Название встречи                  |
| start_time | string | Время начала (ISO 8601)           |
| status     | string | Статус (planned/active/completed) |
| feedback   | string | Комментарий обратной связи        |
