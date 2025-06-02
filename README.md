This is a minimal Go backend API service based on: https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service

### project setup
```bash
git clone https://github.com/LeonibelDev/askme/
```

```bash
cd askme
```

```bash
go run cmd/main/main.go 
```

### check documentation
```bash
localhost:[port]/swagger/index.html
```

### project structure
```
askme/
├── cmd/
│   └── main/
│       └── main.go              # Application entry point
├── db/
│   └── conn.go                  # Database connection setup
├── docs/
│   ├── docs.go                  # Swagger docs registration
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── controllers/             # HTTP handler logic
│   │   ├── blogcontroller.go
│   │   ├── newsletter.go
│   │   └── usercontroller.go
│   └── routes/                  # Route definitions
│       ├── admin/
│       │   └── admin.go
│       ├── auth/
│       │   └── auth.go
│       ├── blog/
│       │   ├── blog.go
│       │   └── repos.go
│       └── newsletter/
│           └── newsletter.go
├── pkg/
│   └── utils/                   # Shared utility packages
│       ├── hash/
│       │   └── hash.go
│       ├── models/
│       │   └── models.go
│       └── token/
│           └── jwt.go
├── go.mod
├── go.sum
└── README.md
```
