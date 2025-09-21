const TracingManager = require('./tracing');
const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, blockInternalRoutes, tracingMiddleware } = require('./middleware');

// Load environment variables before requiring constants so they use .env values
const dotenv = require('dotenv');
dotenv.config();

const {
  AUTH_SERVICE_URL,
  TOUR_SERVICE_URL,
  BLOG_SERVICE_URL,
  FOLLOWER_SERVICE_URL,
  STAKEHOLDER_SERVICE_URL,
  REVIEW_SERVICE_URL,
  PURCHASE_SERVICE_URL
} = require('./constants');

TracingManager.initTracer();

// Debug logging for service URLs
console.log('🔧 Service URLs Configuration:');
console.log('AUTH_SERVICE_URL:', AUTH_SERVICE_URL);
console.log('TOUR_SERVICE_URL:', TOUR_SERVICE_URL);
console.log('BLOG_SERVICE_URL:', BLOG_SERVICE_URL);
console.log('FOLLOWER_SERVICE_URL:', FOLLOWER_SERVICE_URL);
console.log('STAKEHOLDER_SERVICE_URL:', STAKEHOLDER_SERVICE_URL);
console.log('REVIEW_SERVICE_URL:', REVIEW_SERVICE_URL);
console.log('PURCHASE_SERVICE_URL:', PURCHASE_SERVICE_URL);

const api = express();

// CORS configuration
api.use(cors({
  origin: [
    'http://localhost:3000', 
    'http://127.0.0.1:3000',
    'http://localhost:5173',  // Vite dev server
    'http://127.0.0.1:5173'   // Vite dev server
  ],
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization']
}));

api.use(morgan('dev'));

// Add tracing middleware
api.use(tracingMiddleware);

// Add filtering middleware
api.use(blockInternalRoutes);

api.post('/api/auth/login', createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '',
  }
}));

api.post('/api/auth/register', createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '',
  }
}));

api.use('/api/auth', validateJWT, createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '',
  }
}));

// Protected stakeholder profile routes (must come first for specificity)
api.get('/api/stakeholder/profile', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder/profile': '/profile',
  }
}));

api.put('/api/stakeholder/profile', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder/profile': '/profile',
  }
}));

// Public stakeholder routes (for creating profiles)
api.use('/api/stakeholder', createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder': '',
  }
}));

api.use('/api/tours', validateJWT, createProxyMiddleware({
  target: TOUR_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/tours': '',
  },
}));

api.use('/api/blogs', validateJWT, createProxyMiddleware({
  target: BLOG_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/blogs': '',
  },
}));

api.use('/api/reviews', validateJWT, createProxyMiddleware({
  target: REVIEW_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/reviews': '',
  },
}));


api.use('/api/purchases', validateJWT, createProxyMiddleware({
  target: PURCHASE_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/purchases': '',
  },
}));

api.get('/ping', (req, res) => {
  res.status(200).json({
    message: 'pong',
    service: 'API Gateway'
  });
});

const PORT = process.env.API_GATEWAY_PORT || 3000;
api.listen(PORT, () => {
  console.log(`API Gateway listening on port ${PORT}`);
});