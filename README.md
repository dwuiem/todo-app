# REST API todo application

## Introduction

This is a REST API application with simple CRUD operations with tasks and lists using Go language

## Technology Stack
- **Language** Go (1.20+ SDK)
- **Framework** `gin` - API route
- **Authentication** JWT lib to identity users (JSON web token)
- **Database** PostgreSQL using `sqlx` lib

## API Usage

### User Authentication
- **Sign Up**
  - URL: `/auth/sign-up`
  - Method: `POST`
  - Request Body:
    ```json
    {
      "username": "your name"
      "password": "your password"
    }
    ```
  - Responce Body:
    ```json
    {
      "id": <your id>
    }
    ```
- **Sign in**
  - URL: `/auth/sign-in`
  - Method: `POST`
  - Request Body:
    ```json
    {
      "username": "your name"
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
To make api requests you need to include JWT token in the `Authorization` header
- **Create List**
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
      "id": <list id>
    }
    ```
- **...**