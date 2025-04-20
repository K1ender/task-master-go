# Task Master Go

A simple and extendable backend for task management, built with Go. REST API with authentication, user roles, and PostgreSQL data storage.

## 🚀 Features
- User registration and JWT-based login
- Full CRUD for tasks
- Role system: user / admin
- PostgreSQL integration
- Input validation
- Swagger documentation

## 📦 Tech Stack
- Go 1.24+
- PostgreSQL
- Chi (router)
- GORM (ORM)
- Cleanenv (configuration)
- Docker / Docker Compose

## 📌 API Endpoints
```http
POST   /register      # user registration
POST   /login         # get JWT token
GET    /tasks         # fetch all tasks
POST   /tasks         # create a new task
PUT    /tasks/{id}    # update a task
DELETE /tasks/{id}    # delete a task
```

## 🧠 TODO
- [ ] Unit tests
- [ ] Integration tests
- [ ] CI/CD pipeline

---

> This project was built for learning purposes — to improve backend development skills using Go 💪

