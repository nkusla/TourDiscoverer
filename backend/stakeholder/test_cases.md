# Stakeholder Service Test Cases

## Test Case 1: Create Stakeholder Profile

```bash
curl -X POST http://localhost:3000/api/stakeholder \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "profile_picture": "https://example.com/profile.jpg",
    "biography": "A passionate tour guide with 10 years of experience.",
    "motto": "Life is a journey, not a destination."
  }'
```

Expected Response: 201 Created with stakeholder data

## Test Case 2: Get Profile (requires authentication)

First, get JWT token from auth service:
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

Then use the token to get profile:
```bash
curl -X GET http://localhost:3000/api/profile \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Expected Response: 200 OK with profile data

## Test Case 3: Update Profile (requires authentication)

```bash
curl -X PUT http://localhost:3000/api/profile \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Johnny",
    "biography": "Updated biography with more details."
  }'
```

Expected Response: 200 OK with updated profile data

## Test Case 4: Error Cases

### Profile not found (no profile created for user)
```bash
curl -X GET http://localhost:3000/api/profile \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Expected Response: 404 Not Found

### Missing authentication
```bash
curl -X GET http://localhost:3000/api/profile
```

Expected Response: 401 Unauthorized
