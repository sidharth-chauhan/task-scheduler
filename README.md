# Task Scheduler Backend (Golang)

A simple and reliable Task Scheduler built with Golang and PostgreSQL.  
It allows scheduling one-off or recurring HTTP requests and automatically logs their execution results.

---

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/sidharth-chauhan/task-scheduler.git
cd task-scheduler
```

### 2. Create `.env`

```env
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tasks
DB_PORT=5432
```

### 3. Run with Docker

```bash
docker compose up --build
```

The backend will start at:
**[http://localhost:8080](http://localhost:8080)**

---

## API Endpoints

| Method | Endpoint            | Description               |
| ------ | ------------------- | ------------------------- |
| POST   | /tasks              | Create a new task         |
| GET    | /tasks              | List all tasks            |
| GET    | /tasks/{id}         | Get task details          |
| PUT    | /tasks/{id}         | Update a task             |
| DELETE | /tasks/{id}         | Cancel a task             |
| GET    | /tasks/{id}/results | Get all results of a task |
| GET    | /results            | Get all task results      |

---

## Example Task (POST /tasks)

```json
{
  "name": "Check Google",
  "type": "one-off",
  "utc_datetime": "2025-09-29T06:35:00Z",
  "method": "GET",
  "url": "https://www.google.com"
}
```

---

## Testing

You can test the APIs using Postman or run all tests with Newman:

```bash
newman run tests/task_scheduler.postman_collection.json
```

---

## Project Structure

```
task-scheduler/
├── internal/
│   ├── db/
│   ├── handler/
│   │   ├── runner/
│   │   ├── task.go
│   │   └── results.go
│   └── models/
├── tests/
│   └── task_scheduler.postman_collection.json
├── Dockerfile
├── docker-compose.yml
├── main.go
└── README.md
```

---

## Notes

- Tasks and results are stored persistently in PostgreSQL.
- Scheduled tasks are executed automatically and their responses are saved in `task_results`.
- The system is designed for reliability and simple local testing using Docker.

---

✅ **Start the service with `docker compose up --build` and test APIs via Postman or Newman.**
