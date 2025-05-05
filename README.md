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
├── cmd
│   └── main
│       └── main.go
├── db
│   └── conn.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── controllers
│   │   ├── blogcontroller.go
│   │   ├── newsletter.go
│   │   └── usercontroller.go
│   └── routes
│       ├── admin
│       │   └── admin.go
│       ├── auth
│       │   └── auth.go
│       ├── blog
│       │   ├── blog.go
│       │   └── repos.go
│       └── newsletter
│           └── newsletter.go
├── pkg
│   └── utils
│       ├── hash
│       │   └── hash.go
│       ├── models
│       │   └── models.go
│       └── token
│           └── jwt.go
└── README.md
```
