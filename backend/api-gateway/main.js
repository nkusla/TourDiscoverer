const TracingManager = require('./tracing');
const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, blockInternalRoutes, tracingMiddleware } = require('./middleware');
const { AUTH_SERVICE_URL, STAKEHOLDER_SERVICE_URL, TOUR_SERVICE_URL, BLOG_SERVICE_URL, REVIEW_SERVICE_URL } = require('./constants');
const BlogRPCClient = require('./blog_rpc_client');

TracingManager.initTracer();

const dotenv = require('dotenv');
dotenv.config();

const api = express();

// Kreiram Blog RPC klijent
const blogRPCClient = new BlogRPCClient(process.env.BLOG_SERVICE_HOST || 'blog-service', 3012);

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

// RPC endpoint za kreiranje bloga
api.post('/api/blogs', express.json(), validateJWT, async (req, res) => {
  try {
    const username = req.user.username; // iz JWT middleware
    const blogData = {
      title: req.body.title,
      description: req.body.description,
      images: req.body.images || [],
      author: username
    };
    
    const result = await blogRPCClient.createBlog(blogData);
    res.status(201).json(result);
  } catch (error) {
    console.error('Blog creation error:', error);
    res.status(500).json({ error: 'Failed to create blog' });
  }
});

// RPC endpoint za dobavljanje personalizovanih blogova  
api.get('/api/blogs/personalized', validateJWT, async (req, res) => {
  try {
    const username = req.user.username; // iz JWT middleware
    const result = await blogRPCClient.getPersonalizedBlogs(username);
    res.json(result.blogs || []);
  } catch (error) {
    console.error('Get personalized blogs error:', error);
    res.status(500).json({ error: 'Failed to fetch personalized blogs' });
  }
});

// Protected blog routes (for comments and likes) - ostaju HTTP proxy
api.get('/api/blogs/comments', validateJWT, createProxyMiddleware({
  target: BLOG_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/blogs': '',
  },
}));

// Protected blog routes (for liking, commenting) - ostaju HTTP proxy
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