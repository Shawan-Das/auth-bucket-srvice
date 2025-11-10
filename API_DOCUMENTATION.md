# Authentication API Documentation

## Overview
This system provides user registration, login, JWT token management, and refresh token functionality for ADMIN users.

## Configuration
Token expiry times can be configured in `config/local-config.json`:
- `access_token_duration_minutes`: Access token expiry time (default: 15 minutes)
- `refresh_token_duration_hours`: Refresh token expiry time (default: 168 hours / 7 days)
- `access_secret`: Secret key for JWT signing

## API Endpoints

### 1. Register User
**Endpoint:** `POST /api/auth/register`

**Description:** Creates a new user with ADMIN role

**Request Body:**
```json
{
  "user_name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "password": "SecurePass@123"
}
```

**Password Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

**Success Response (201):**
```json
{
  "isSuccess": true,
  "serviceMessage": "User registered successfully",
  "payload": {
    "user_id": 1,
    "user_name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890",
    "role": "ADMIN",
    "created_at": "2025-11-10 10:30:00"
  },
  "ts": "2025-11-10-10:30:00.000"
}
```

---

### 2. Login
**Endpoint:** `POST /api/auth/login`

**Description:** Authenticates user and returns access and refresh tokens

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass@123"
}
```

**Success Response (200):**
```json
{
  "isSuccess": true,
  "serviceMessage": "Login successful",
  "payload": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 1,
    "user_name": "John Doe",
    "email": "john@example.com",
    "role": "ADMIN",
    "expires_in": 900
  },
  "ts": "2025-11-10-10:30:00.000"
}
```

**Notes:**
- Only users with role "ADMIN" can login
- `expires_in` is in seconds

---

### 3. Refresh Token
**Endpoint:** `POST /api/auth/refresh`

**Description:** Generates new access and refresh tokens using a valid refresh token

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "isSuccess": true,
  "serviceMessage": "Token refreshed successfully",
  "payload": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  },
  "ts": "2025-11-10-10:30:00.000"
}
```

---

### 4. Logout
**Endpoint:** `POST /api/auth/logout`

**Description:** Invalidates the user's refresh token

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "isSuccess": true,
  "serviceMessage": "Logout successful",
  "ts": "2025-11-10-10:30:00.000"
}
```

---

## Protected Routes

All routes except the following require authentication:
- `/api/auth/register`
- `/api/auth/login`
- `/api/auth/refresh`
- `/`
- `/swagger/*`

### Authorization Header Format
For protected routes, include the JWT token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Context Variables
After successful authentication, the following user information is available in the Gin context:
- `userID` (int32)
- `email` (string)
- `username` (string)
- `role` (string)

---

## Error Responses

### 400 Bad Request
```json
{
  "isSuccess": false,
  "serviceMessage": "Invalid request body",
  "ts": "2025-11-10-10:30:00.000"
}
```

### 401 Unauthorized
```json
{
  "isSuccess": false,
  "serviceMessage": "Invalid email or password",
  "ts": "2025-11-10-10:30:00.000"
}
```

### 403 Forbidden
```json
{
  "isSuccess": false,
  "serviceMessage": "Access denied. Admin role required",
  "ts": "2025-11-10-10:30:00.000"
}
```

### 409 Conflict
```json
{
  "isSuccess": false,
  "serviceMessage": "User with this email already exists",
  "ts": "2025-11-10-10:30:00.000"
}
```

### 500 Internal Server Error
```json
{
  "isSuccess": false,
  "serviceMessage": "Failed to create user",
  "ts": "2025-11-10-10:30:00.000"
}
```

---

## Database Schema

The system uses the following table structure:

```sql
CREATE TABLE IF NOT EXISTS common.users (
    user_id serial4 PRIMARY KEY NOT NULL,
    user_name text NOT NULL,
    email text NOT NULL UNIQUE,
    phone text NOT NULL,
    pass text NOT NULL,
    pss_valid bool DEFAULT true NOT NULL,
    otp text NOT NULL DEFAULT '',
    otp_valid bool DEFAULT false NOT NULL,
    otp_exp timestamp NULL,
    role text NOT NULL,
    refresh_token text DEFAULT '',
    refresh_token_exp timestamp NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL
);
```

---

## Security Features

1. **Password Hashing:** All passwords are hashed using Argon2id algorithm
2. **JWT Tokens:** Separate access and refresh tokens with different expiry times
3. **Refresh Token Storage:** Refresh tokens are stored in database with expiry
4. **Token Validation:** Both token type and expiry are validated
5. **Role-Based Access:** Only ADMIN role users can access the system
6. **Email Validation:** Email format is validated before registration

---

## Usage Example

### Complete Authentication Flow

1. **Register a new user:**
```bash
curl -X POST http://localhost:7070/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890",
    "password": "SecurePass@123"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:7070/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass@123"
  }'
```

3. **Access protected resource:**
```bash
curl -X GET http://localhost:7070/api/protected-route \
  -H "Authorization: Bearer <access_token>"
```

4. **Refresh token when access token expires:**
```bash
curl -X POST http://localhost:7070/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

5. **Logout:**
```bash
curl -X POST http://localhost:7070/api/auth/logout \
  -H "Authorization: Bearer <access_token>"
```
