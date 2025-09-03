# Review Service API Documentation

## Overview
The Review Service allows tourists to leave reviews for tours, including ratings, comments, visit dates, and images.

## Endpoints

### Create Review
- **POST** `/api/reviews/`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Body**:
```json
{
  "tour_id": 1,
  "rating": 5,
  "comment": "Amazing tour! Highly recommended.",
  "visit_date": "2025-08-15",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```
- **Response**: `201 Created`
```json
{
  "id": 1,
  "tour_id": 1,
  "tourist_username": "john_doe",
  "rating": 5,
  "comment": "Amazing tour! Highly recommended.",
  "visit_date": "2025-08-15T00:00:00Z",
  "review_date": "2025-09-01T10:30:00Z",
  "images": [
    {
      "id": 1,
      "image_url": "https://example.com/image1.jpg"
    },
    {
      "id": 2,
      "image_url": "https://example.com/image2.jpg"
    }
  ],
  "created_at": "2025-09-01T10:30:00Z",
  "updated_at": "2025-09-01T10:30:00Z"
}
```

### Get Review by ID
- **GET** `/api/reviews/{id}`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `200 OK` (same structure as create response)

### Get Tour Reviews
- **GET** `/api/reviews/tour/{tour_id}?page=1&page_size=10`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `200 OK`
```json
{
  "reviews": [
    {
      "id": 1,
      "tour_id": 1,
      "tourist_username": "john_doe",
      "rating": 5,
      "comment": "Amazing tour!",
      "visit_date": "2025-08-15T00:00:00Z",
      "review_date": "2025-09-01T10:30:00Z",
      "images": [...],
      "created_at": "2025-09-01T10:30:00Z",
      "updated_at": "2025-09-01T10:30:00Z"
    }
  ],
  "total_count": 1,
  "average_rating": 5.0
}
```

### Get My Reviews
- **GET** `/api/reviews/my?page=1&page_size=10`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `200 OK` (same structure as Get Tour Reviews)

### Update Review
- **PUT** `/api/reviews/{id}`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Body**:
```json
{
  "rating": 4,
  "comment": "Updated comment",
  "images": [
    "https://example.com/new_image.jpg"
  ]
}
```
- **Response**: `200 OK` (same structure as create response)

### Delete Review
- **DELETE** `/api/reviews/{id}`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `204 No Content`

### Get Tour Average Rating
- **GET** `/api/reviews/tour/{tour_id}/rating`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `200 OK`
```json
{
  "tour_id": 1,
  "average_rating": 4.5
}
```

## Validation Rules
- Rating must be between 1 and 5
- Comment is required
- Visit date must be in format "YYYY-MM-DD"
- Only the review author can update or delete their review
- Users can only create one review per tour

## Database Schema

### Review Table
- `id` (Primary Key)
- `tour_id` (Foreign Key to Tour)
- `tourist_username` (Foreign Key to User)
- `rating` (1-5)
- `comment` (Text)
- `visit_date` (Date)
- `review_date` (Timestamp, auto-generated)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `deleted_at` (Soft delete timestamp)

### ReviewImage Table
- `id` (Primary Key)
- `review_id` (Foreign Key to Review)
- `image_url` (Text)

## Error Responses
- `400 Bad Request`: Invalid request data or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't own the review
- `404 Not Found`: Review or tour not found
- `409 Conflict`: Review already exists for this tour by this user
- `500 Internal Server Error`: Server error
