# Stakeholder Service

This service manages user profiles for the TourDiscoverer application.

## Features

- Create stakeholder profiles with personal information
- View user profiles
- Update profile information
- Link profiles to users created in the auth service

## Endpoints

### POST /stakeholder
Create a new stakeholder profile.

**Request Body:**
```json
{
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "profile_picture": "https://example.com/profile.jpg",
    "biography": "A passionate tour guide with 10 years of experience.",
    "motto": "Life is a journey, not a destination."
}
```

### GET /profile
Get the current user's profile (requires JWT authentication).

**Headers:**
- `Authorization: Bearer <jwt_token>`

**Response:**
```json
{
    "id": 1,
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "profile_picture": "https://example.com/profile.jpg",
    "biography": "A passionate tour guide with 10 years of experience.",
    "motto": "Life is a journey, not a destination.",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
}
```

### PUT /profile
Update the current user's profile (requires JWT authentication).

**Headers:**
- `Authorization: Bearer <jwt_token>`

**Request Body:**
```json
{
    "first_name": "John",
    "last_name": "Smith",
    "profile_picture": "https://example.com/new-profile.jpg",
    "biography": "Updated biography",
    "motto": "New motto"
}
```

## Database Schema

The stakeholder service uses PostgreSQL with the following table structure:

```sql
CREATE TABLE stakeholders (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    profile_picture VARCHAR(255),
    biography TEXT,
    motto VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Environment Variables

- `STAKEHOLDER_DB_HOST`: Database host
- `STAKEHOLDER_DB_PORT`: Database port (default: 5432)
- `STAKEHOLDER_DB_NAME`: Database name
- `STAKEHOLDER_DB_USER`: Database username
- `STAKEHOLDER_DB_PASSWORD`: Database password
- `PORT`: Service port (default: 8081)

## Usage Flow

1. User registers via auth service
2. Client creates stakeholder profile using POST /stakeholder
3. User can view profile using GET /profile (with JWT)
4. User can update profile using PUT /profile (with JWT)
