# Tour Discoverer

A microservices-based tour management platform with interactive map functionality for creating and managing tours with key points.

## Project Structure

```
TourDiscoverer/
├── backend/                 # Backend microservices
│   ├── api-gateway/        # Node.js API Gateway (Express)
│   ├── auth/               # Authentication service (Go)
│   ├── blog/               # Blog service (Go)
│   ├── follower/           # Follower service (Go)
│   ├── purchase/           # Purchase service (Go)
│   ├── review/             # Review service (Go)
│   ├── seed/               # Database seeding service (Go)
│   ├── stakeholder/        # Stakeholder service (Go)
│   └── tour/               # Tour management service (Go)
├── frontend/               # Vue.js frontend application
├── docker-compose.yaml     # Docker orchestration
├── .env                    # Environment variables
└── README.md              # This file
```

## Features

### 🗺️ Interactive Map-Based Tour Management
- **Tour Creation**: Create tours by clicking on interactive maps
- **Key Points**: Add, edit, move, and delete key points with drag-and-drop
- **Route Visualization**: Automatic route drawing between key points
- **Distance Calculation**: Real-time distance calculation for tour routes

### 🏗️ Microservices Architecture
- **API Gateway**: Central entry point with routing and authentication
- **Service Isolation**: Independent services for different domains
- **Scalable Design**: Each service can be scaled independently
- **Database per Service**: Dedicated databases for data isolation

### 🎯 Frontend Features
- **Vue.js 3**: Modern reactive frontend framework
- **Leaflet Maps**: Interactive map library with OpenStreetMap
- **Responsive Design**: Bootstrap-based responsive UI
- **Real-time Updates**: Live map interactions and updates

## Technology Stack

### Backend
- **Languages**: Go, Node.js
- **Databases**: PostgreSQL, MongoDB, Neo4j
- **Containerization**: Docker & Docker Compose
- **API Gateway**: Express.js with proxy middleware

### Frontend
- **Framework**: Vue.js 3 with Composition API
- **Build Tool**: Vite
- **State Management**: Pinia
- **Maps**: Leaflet with Vue-Leaflet
- **UI Framework**: Bootstrap 5
- **HTTP Client**: Axios

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.19+ (for local backend development)

### 1. Start Databases
```bash
docker-compose --profile db up -d
```

### 2. Start Backend Services
```bash
docker-compose --profile service up -d
```

### 3. Access the Application
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Individual Services**: Check `.env` file for specific ports

### 4. Development Mode

#### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

#### Backend Development
Each service can be run individually:
```bash
cd backend/[service-name]
go mod tidy
go run .
```

## Services Overview

### API Gateway (Port 8080)
- **Technology**: Node.js + Express
- **Purpose**: Route requests to appropriate services
- **Features**: JWT authentication, request logging, CORS handling

### Auth Service (Port 3001)
- **Technology**: Go + PostgreSQL
- **Purpose**: User authentication and authorization
- **Features**: User registration, login, JWT token management

### Tour Service (Port 3006)
- **Technology**: Go + PostgreSQL
- **Purpose**: Tour and key point management
- **Features**: CRUD operations for tours, key points, distance calculation

### Other Services
- **Blog Service**: Content management for tour blogs
- **Review Service**: Tour reviews and ratings
- **Follower Service**: User following relationships
- **Stakeholder Service**: Stakeholder management
- **Seed Service**: Database initialization and test data

## Database Configuration

The application uses multiple databases:
- **PostgreSQL**: Auth, Tour, Review, Stakeholder services
- **MongoDB**: Blog service
- **Neo4j**: Follower service (graph relationships)

Database configurations are managed through environment variables in `.env`.

## Environment Variables

Key environment variables (see `.env` for complete list):

```bash
# Service URLs and Ports
API_GATEWAY_PORT=8080
FRONTEND_PORT=3000

# Database Settings
AUTH_DB_HOST=auth-db
TOUR_DB_HOST=tour-db
# ... (see .env file for complete configuration)

# Security
JWT_SECRET=your_jwt_secret_key
```

## Development Workflow

### Adding New Features

1. **Backend Changes**:
   - Modify appropriate service in `backend/[service-name]/`
   - Update API endpoints
   - Run service locally or rebuild Docker image

2. **Frontend Changes**:
   - Modify Vue components in `frontend/src/`
   - Update API calls in stores/services
   - Test with `npm run dev`

3. **Database Changes**:
   - Update models in respective service
   - Run migrations if needed
   - Update seed data if applicable

### Docker Commands

```bash
# Start all services
docker-compose --profile service up -d

# Start only databases
docker-compose --profile db up -d

# Rebuild and start specific service
docker-compose up --build [service-name]

# View logs
docker-compose logs [service-name]

# Stop all services
docker-compose down
```

## Map Functionality

### Tour Creation Workflow
1. Navigate to "Create Tour" page
2. Fill in basic tour information
3. Click on map to add key points
4. Drag markers to adjust positions
5. Edit key point details in modal/sidebar
6. Save tour with calculated route

### Key Point Management
- **Add**: Click anywhere on the map
- **Move**: Drag existing markers
- **Edit**: Click marker and use edit button
- **Delete**: Remove from marker popup or sidebar
- **Reorder**: Modify order in edit dialog

### Map Features
- **Interactive Navigation**: Zoom, pan, marker interaction
- **Route Visualization**: Polyline connecting key points
- **Distance Calculation**: Real-time route distance
- **Responsive Design**: Works on mobile and desktop

## API Documentation

### Authentication Endpoints
```
POST /auth/register - User registration
POST /auth/login    - User login
GET  /auth/profile  - Get user profile
PUT  /auth/profile  - Update user profile
```

### Tour Endpoints
```
GET    /tours           - List tours
POST   /tours           - Create tour
GET    /tours/:id       - Get tour details
PUT    /tours/:id       - Update tour
DELETE /tours/:id       - Delete tour
POST   /tours/:id/keypoints     - Add key point
PUT    /tours/:id/keypoints/:kp - Update key point
DELETE /tours/:id/keypoints/:kp - Delete key point
```

## Contributing

1. Fork the repository
2. Create feature branch
3. Make changes in appropriate service/frontend
4. Test locally with Docker Compose
5. Submit pull request

## License

This project is licensed under the terms specified in the LICENSE file.

## Support

For issues and questions:
1. Check existing GitHub issues
2. Create new issue with detailed description
3. Include logs and environment details for bugs
