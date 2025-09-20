const TracingManager = require('./tracing');
const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, blockInternalRoutes, tracingMiddleware } = require('./middleware');
const { AUTH_SERVICE_URL, STAKEHOLDER_SERVICE_URL, TOUR_SERVICE_URL, BLOG_SERVICE_URL, REVIEW_SERVICE_URL } = require('./constants');

TracingManager.initTracer();

const dotenv = require('dotenv');
dotenv.config();

const api = express();

// CORS configuration
api.use(cors({
  origin: ['http://localhost:3000', 'http://127.0.0.1:3000'],
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

// Protected position routes (for tourist position simulator)
api.post('/api/stakeholder/position', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder/position': '/position',
  }
}));

api.get('/api/stakeholder/position', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder/position': '/position',
  }
}));

api.delete('/api/stakeholder/position', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder/position': '/position',
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

// TourExecution routes - protected for tourists
api.use('/api/tour', validateJWT, createProxyMiddleware({
  target: TOUR_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/tour': '',
  },
}));

// Public blog routes (for reading blogs and comments)
api.get('/api/blogs', createProxyMiddleware({
  target: BLOG_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/blogs': '',
  },
}));

// Protected blog routes for authenticated users to get personalized blogs
api.get('/api/blogs/personalized', validateJWT, createProxyMiddleware({
  target: BLOG_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/blogs/personalized': '',
  },
}));

api.get('/api/blogs/comments', createProxyMiddleware({
  target: BLOG_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/blogs': '',
  },
}));

// Protected blog routes (for creating, liking, commenting)
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