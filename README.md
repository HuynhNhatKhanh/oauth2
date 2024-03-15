# Project Name
OAuth2

## Description
This project is a backend service written in Golang using the GoFiber framework and MongoDB as the database. 


## Features
- User registration with email verification
- User login with two-step authentication (username/password + OTP from email)
- Refresh token mechanism
- Get user information using access token

## Getting Started
1. Install Golang
2. Clone the repository
3. Install dependencies (`go mod tidy`)
4. Run the server (`go run main.go`)

## Usage Endpoints

### Register User
- **Description**: Registers a new user
- **URL** POST `/register`
- **Request Body**: 
    ```json 
    { 
      "username": "yourName",
      "email": "yourEmail@gmail.com",
      "password": "yourPassword"
    }
- **Response**: 
    ```json 
    { 
      "username": "yourName",
      "email": "yourEmail@gmail.com",
      "password": "yourPassword"
    }

### Verify Register