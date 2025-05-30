# REST API todo application

## Introduction

This is a microservice application with simple CRUD operations with tasks and lists

## Technology Stack
- **Language** Go (1.23)
- **gRPC** interaction between microservices
- **Framework** `gin` - API route
- **Authentication** JWT lib to identity users (JSON web token)
- **Database** PostgreSQL using `pgx` lib
- **Migrations** using `github.com/golang-migrate/migrate/v4`

## Usage

1. Cloning repository

```sh
git clone https://github.com/dwuiem/todo-app.git
```

2. Docker Compose

```sh
docker-compose up -d
```

## Project Structure
```
.
├── api
│   ├── gen
│   │   └── sso
│   └── proto
├── sso
│   ├── cmd
│   │   ├── app
│   │   └── migrate
│   ├── config
│   ├── gen
│   │   └── sso
│   ├── internal
│   │   ├── adapter
│   │   ├── app
│   │   ├── config
│   │   ├── domain
│   │   ├── lib
│   │   ├── service
│   │   └── storage
│   └── migrations
└── todo
    ├── cmd
    │   ├── app
    │   └── migrate
    ├── config
    ├── gen
    │   └── sso
    ├── internal
    │   ├── adapter
    │   ├── app
    │   └── domain
    └── migrations
```