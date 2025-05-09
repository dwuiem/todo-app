# REST API todo application

## Introduction

This is a microservice application with simple CRUD operations with tasks and lists with authentication

## Technology Stack
- **Language** Go (1.23)
- **gRPC** interaction between microservices
- **Framework** `gin` - API route
- **Authentication** JWT lib to identity users (JSON web token)
- **Database** PostgreSQL using `pgx` lib
- **Migrations** using `github.com/golang-migrate/migrate/v4`
- 

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
todo-app
├── api                     # API Protobuf contracts 
│    ├── proto
│    └── Makefile
├── services             
│    ├── sso                # SSO service (GRPC server)
│    └── todo               # TODO service (HTTP server)
├── .gitignore
├── README.md
└── docker-compose.yaml     # Docker Compose
```