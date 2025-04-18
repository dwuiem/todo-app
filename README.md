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
│    └── main.go
├── services             
│    ├── sso                # SSO service (GRPC server)
│    └── todo               # TODO service (HTTP server)
├── .gitignore
├── README.md
└── docker-compose.yaml     # Docker Compose
```

## API Usage

### User Authentication
- **Sign Up**
  - URL: `/auth/sign-up`
  - Method: `POST`
  - Request Body:
    ```json
    {
      "username": "your name",
      "password": "your password"
    }
    ```
  - Responce Body:
    ```json
    {
      "id": 1
    }
    ```
- **Sign in**
  - URL: `/auth/sign-in`
  - Method: `POST`
  - Request Body:
    ```json
    {
      "username": "your name",
      "password": "your password"
    }
    ```
  - Responce Body:
    ```json
    {
      "token": "your token"
    }
    ```
### Authorized Requests
To make api requests you need to include JWT token in the `Authorization` header. For example


**Create List Request**
- Authorization Header: "your bearer token"
- URL: `/api/lists`
- Method: `POST`
- Request Body:
  ```json
  {
    "title": "your title"
  }
  ```
- Responce Body:
  ```json
  {
    "id": 1
  }
  ```