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
4. Run the server (`go run main.go`) or (`make server`) if your device has the "make" installed

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
- **Sample Response (Success)**: 
    ```json 
    { 
      "message": "Send otp successfully"
    }

### Verify Register
- **Description**:Verify your registered account with the link sent to your email
- **Sample Response (Success)**: 
    ```json 
    {"message":"Email verified successfully"}

### Login
- **Description**: Logs in a user
- **URL**: POST `/login`
- **Request Body**: 
    ```json 
    { 
      "username": "yourName",
      "password": "yourPassword"
    }
- **Sample Response (Success)**: 
    ```json 
    { 
      "message": "Send OTP to verify login"
    }

### Verify Login
- **Description**: Verifies a user's login using an OTP (One-Time Password)
- **URL**: GET `/verifyLogin`
- **Query Parameters**: 
    ```json 
    { 
      "otp_login": "YourOTP",
      "email": "YourEmail@gmail.com"
    }
- **Sample Response (Success)**: 
    ```json 
    { 
      "accessToken": "YourAccessToken",
      "refreshToken": "YourrefreshToken"
    }

### User Information
- **Description**: Retrieves the information of the currently authenticated user
- **URL**: GET `/user`
- **Headers**: 
    ```json 
    { 
      "accessToken": "YourAccessToken"
    }
- **Sample Response (Success)**: 
    ```json 
    { 
      "createdAt": "2024-03-15T00:00:00Z",
      "email": "yourEmail@gmail.com",
      "username": "yourUsername"
    }

### Refresh Token
- **Description**: Refreshes the user's access token
- **URL**: POST `/refresh`
- **Headers**: 
    ```json 
    { 
      "refreshToken": "YourRefreshToken"
    }
- **Sample Response (Success)**: 
    ```json 
    { 
      "accessToken": "NewAccessToken"
    }