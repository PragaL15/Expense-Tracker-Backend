
```markdown
# 📊 Expense Tracker Backend  

A simple and efficient **backend API** for managing expenses, incomes, and balances.  
Built with **Go (Golang)** and **Fiber** for fast, lightweight, and scalable performance.  

---

## ⚡ Tech Stack  

- Go (Golang)  
- Fiber Framework  
- PostgreSQL  
- JWT Authentication  
- GORM / SQLC / pgx (any ORM or query builder of choice)  

---

## 📂 Project Structure  

```

Expense-Tracker-Backend/
│── cmd/                 # Entry point for app
│── config/              # Configuration files (env, DB)
│── controllers/         # Route handlers
│── middlewares/         # JWT auth, logging, etc.
│── models/              # Database models
│── routes/              # API routes
│── services/            # Business logic
│── utils/               # Helpers (validators, error handlers)
│── main.go              # Main server entry
│── go.mod               # Dependencies
│── go.sum

````

---

## 🚀 Getting Started  

### Prerequisites  
- Go **1.21+**  
- PostgreSQL (or another database)  
- Git  

### Clone the Repository  
```bash
git clone https://github.com/PragaL15/Expense-Tracker-Backend.git
cd Expense-Tracker-Backend
````

### Install Dependencies

```bash
go mod tidy
```

### Setup Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=expensetracker
JWT_SECRET=your_jwt_secret
```

### Run the Server

```bash
go run main.go
```

Server runs at:
👉 `http://localhost:3000`

---

## 📡 API Endpoints

### Authentication

| Method | Endpoint         | Description         |
| ------ | ---------------- | ------------------- |
| POST   | `/register` | Register a new user |
| POST   | `/login`    | Login & get token   |

### Transactions

| Method | Endpoint            | Description                  |
| ------ | ------------------- | ---------------------------- |
| GET    | `/transactions`     | Get all transactions (auth)  |
| POST   | `/transactions`     | Add a new transaction (auth) |
| GET    | `/transactions/:id` | Get a single transaction     |
| PUT    | `/transactions/:id` | Update transaction           |
| DELETE | `/transactions/:id` | Delete transaction           |

### Balance

| Method | Endpoint   | Description         |
| ------ | ---------- | ------------------- |
| GET    | `/balance` | Get current balance |

---

## 🛠️ Example Requests

**Register User**

```bash
curl -X POST http://localhost:8080/auth/register \
   -H "Content-Type: application/json" \
   -d '{"username":"john", "password":"123456"}'
```

**Add Expense**

```bash
curl -X POST http://localhost:8080/transactions \
   -H "Authorization: Bearer <token>" \
   -H "Content-Type: application/json" \
   -d '{"amount":500, "type":"expense", "category":"Food", "description":"Lunch"}'
```

---

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for details.

```

---

Do you also want me to add a **Docker setup section** (with `Dockerfile` + `docker-compose.yml` usage) inside this README so that the backend is production-ready?
```

