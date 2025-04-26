# Meeting Service API Documentation

Протокол API для управления встречами, пользователями и лидерами. Реализует CRUD для встреч, регистрацию пользователей и дополнительные функции. Сервис также содержит в себе кеширование данных.

## Содержание

- [Обзор методов](#обзор-методов)
- [Детализация методов](#детализация-методов)
- [Структуры сообщений](#структуры-сообщений)

---

## Обзор методов

| Метод         | HTTP-метод | Путь                               | Описание                        |
| ------------- | ---------- | ---------------------------------- | ------------------------------- |
| AddMeeting    | POST       | /api/meetings                      | Создать новую встречу           |
| GetMeeting    | GET        | /api/meetings/{meeting_id}         | Получить встречу по ID          |
| UpdateMeeting | PATCH      | /api/meetings/{meeting_id}         | Обновить встречу                |
| DeleteMeeting | DELETE     | /api/meetings/{meeting_id}         | Удалить встречу                 |
| GetMeetings   | GET        | /api/meetings                      | Получить все встречи            |
| AddUser       | POST       | /api/meetings/users                | Зарегистрировать пользователя   |
| GetUser       | GET        | /api/meetings/users/{email}        | Получить данные пользователя    |
| AddLeader     | POST       | /api/meetings/leaders              | Добавить лидера                 |
| FinishCourse  | POST       | /api/meetings/users/{email}/finish | Завершить курс для пользователя |

---

## Детализация методов

### ➕ AddMeeting

**POST /api/meetings**

Пример запроса:

```json
{
  "user_id": "user-123",
  "leader_id": "leader-456",
  "title": "Project Planning",
  "start_time": "2024-03-20"
}
```

Ответ (201 Created):

```json
{
  "id": "meeting-789"
}
```

### 📝 GetMeeting

**GET /api/meetings/meeting-789**

Ответ (200 OK):

```json
{
  "meeting": {
    "id": "meeting-789",
    "user_id": "user-123",
    "leader_id": "leader-456",
    "title": "Project Planning",
    "start_time": "2024-03-20T15:00:00Z"
  }
}
```

### ✏️ UpdateMeeting

**PATCH /api/meetings/meeting-789**

Пример запроса:

```json
{
  "title": "Updated Project Scope",
  "leader_id": "leader-999"
}
```

Ответ (200 OK):

```json
{
  "meeting": {
    "id": "meeting-789",
    "user_id": "user-123",
    "leader_id": "leader-999",
    "title": "Updated Project Scope",
    "start_time": "2024-03-20T15:00:00Z"
  }
}
```

### 🗑️ DeleteMeeting

**DELETE /api/meetings/meeting-789**

Ответ (200 OK):

```json
{
  "meeting_id": "meeting-789"
}
```

### 📜 GetMeetings

**GET /api/meetings**

Ответ (200 OK):

```json
{
  "meetings": [
    {
      "id": "meeting-789",
      "user_id": "user-123",
      "leader_id": "leader-456",
      "title": "Project Planning",
      "start_time": "2024-03-20T15:00:00Z"
    }
  ]
}
```

### 👤 AddUser

**POST /api/meetings/users**

Пример запроса:

```json
{
  "email": "user@example.com",
  "fullname": "John Doe",
  "password": "secure123",
  "current_stage": "onboarding"
}
```

Ответ (201 Created):

```json
{
  "id": "user-123"
}
```

### 🔍 GetUser

**GET /api/meetings/users/user@example.com**

Ответ (200 OK):

```json
{
  "id": "user-123",
  "email": "user@example.com",
  "fullname": "John Doe",
  "current_stage": "onboarding"
}
```

### 👑 AddLeader

**POST /api/meetings/leaders**

Пример запроса:

```json
{
  "email": "leader@company.com",
  "fullname": "Alice Smith",
  "password": "admin123"
}
```

Ответ (201 Created):

```json
{
  "id": "leader-456"
}
```

### 🎓 FinishCourse

**POST /api/meetings/users/user@example.com/finish**

Ответ (200 OK):

```json
{
  "id": "user-123"
}
```

---

## Структуры сообщений

### Meeting

```json
{
  "id": "string",
  "user_id": "string",
  "leader_id": "string",
  "title": "string",
  "start_time": "string (ISO 8601)"
}
```

### AddUserRequest

```json
{
  "email": "string",
  "fullname": "string",
  "password": "string",
  "current_stage": "string"
}
```

### Error Response (Пример)

```json
{
  "error": {
    "code": "number",
    "message": "string"
  }
}
```

---

> ⚠️ Все временные метки должны быть в формате **ISO 8601** (пример: `2024-03-20T15:00:00Z`).  
> 🔑 Для методов, требующих аутентификации, добавьте заголовок `Authorization: Bearer <token>`.
