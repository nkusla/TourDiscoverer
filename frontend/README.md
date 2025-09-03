# Tour Discoverer Frontend

Vue.js frontend application for the Tour Discoverer platform featuring interactive map-based tour creation and management.

## Features

- 🗺️ **Interactive Maps**: Create and edit tours using Leaflet maps
- 📍 **Key Points Management**: Click-to-add, drag-to-move key points
- 🎯 **Tour Route Visualization**: Automatic route drawing between key points
- 📱 **Responsive Design**: Bootstrap-based responsive UI
- 🔐 **Authentication**: Login/Register functionality
- 🎨 **Modern UI**: Clean and intuitive interface

## Technology Stack

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Build tool and dev server
- **Vue Router** - Client-side routing
- **Pinia** - State management
- **Leaflet** - Interactive map library
- **Bootstrap 5** - CSS framework
- **Axios** - HTTP client

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Running backend services (API Gateway, Tour Service, Auth Service)

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure environment variables:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` file with your API Gateway URL.

4. Start development server:
   ```bash
   npm run dev
   ```

5. Open browser to `http://localhost:3000`

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run serve` - Serve production build

## Map Features

### Tour Creation

1. Navigate to "Create Tour" page
2. Fill in tour details (name, description, difficulty, etc.)
3. Click on map to add key points
4. Drag markers to adjust positions
5. Edit key point details in popup or sidebar
6. Save tour

### Key Point Management

- **Add**: Click anywhere on the map
- **Move**: Drag existing markers
- **Edit**: Click marker popup "Edit" button
- **Delete**: Click marker popup "Remove" button
- **Reorder**: Modify order in edit dialog

### Map Controls

- **Zoom**: Mouse wheel or +/- buttons
- **Pan**: Click and drag map
- **Clear All**: Remove all key points
- **Save Route**: Save current tour configuration

## Project Structure

```
src/
├── components/           # Reusable components
│   ├── Map/             # Map-related components
│   │   └── TourMap.vue  # Main map component
│   └── Navbar.vue       # Navigation component
├── views/               # Page components
│   ├── Home.vue         # Landing page
│   ├── Tours.vue        # Tours listing
│   ├── TourEditor.vue   # Tour creation/editing
│   └── Login.vue        # Authentication
├── stores/              # Pinia stores
│   ├── user.js          # User authentication
│   └── tour.js          # Tour management
├── services/            # API services
│   └── api.js           # Axios configuration
├── router/              # Vue Router config
│   └── index.js         # Route definitions
└── main.js              # Application entry point
```

## Environment Variables

- `VITE_API_URL` - Backend API Gateway URL (default: http://localhost:8080)
- `VITE_APP_TITLE` - Application title
- `VITE_APP_VERSION` - Application version

## API Integration

The frontend communicates with the backend through the API Gateway:

- `/auth/*` - Authentication endpoints
- `/tours/*` - Tour management endpoints
- `/users/*` - User management endpoints

## Map Configuration

The map component uses OpenStreetMap tiles by default. You can configure:

- **Tile Provider**: Change `tileLayerUrl` in TourMap.vue
- **Default Location**: Modify `initialCenter` prop (default: Belgrade, Serbia)
- **Zoom Levels**: Adjust `initialZoom` and auto-zoom logic

## Building for Production

1. Build the application:
   ```bash
   npm run build
   ```

2. Build Docker image:
   ```bash
   docker build -t td-frontend .
   ```

3. Run with Docker Compose:
   ```bash
   docker-compose --profile service up frontend
   ```

## Contributing

1. Follow Vue.js style guide
2. Use TypeScript for complex components
3. Maintain responsive design principles
4. Test map functionality across different screen sizes
5. Ensure accessibility standards

## Troubleshooting

### Map Not Loading
- Check network connectivity
- Verify tile server accessibility
- Check browser console for errors

### API Connection Issues
- Verify `VITE_API_URL` in environment variables
- Ensure backend services are running
- Check CORS configuration

### Build Issues
- Clear node_modules and reinstall: `rm -rf node_modules package-lock.json && npm install`
- Check Node.js version compatibility
- Verify all dependencies are installed
