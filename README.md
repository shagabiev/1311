# URL Checker Service

Простой HTTP-сервис для проверки доступности ссылок и генерации PDF-отчётов.  
**Без внешних зависимостей** — только стандартная библиотека Go.

---

## Функции

| Endpoint | Метод | Описание |
|---------|-------|--------|
| `POST /check` | JSON | Проверяет список ссылок, возвращает статус и `links_num` |
| `POST /report` | JSON | Генерирует PDF-отчёт по `links_num` |

---

## Запуск

```powershell
go run .
```

---

## Проверка Endpoints в Postman

POST /check
{
  "links": ["google.com", "bad.site"]
}
Ответ:
{
  "links": {
    "google.com": "available",
    "bad.site": "not available"
  },
  "links_num": 1
}

POST /report
{
  "links_num": [1]
}
Ответ: PDF-файл (report.pdf)
