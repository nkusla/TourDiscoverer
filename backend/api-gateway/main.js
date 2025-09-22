const TracingManager = require('./tracing');
const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, blockInternalRoutes, tracingMiddleware } = require('./middleware');
const BlogRPCClient = require('./blog_rpc_client');
const PurchaseRPCClient = require('./purchase_rpc_client');
const StakeholderRPCClient = require('./stakeholder_rpc_client');

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

// Kreiram RPC klijente
const blogRPCClient = new BlogRPCClient(process.env.BLOG_SERVICE_HOST || 'blog-service', 3012);
const purchaseRPCClient = new PurchaseRPCClient(process.env.PURCHASE_SERVICE_HOST || 'purchase-service', 3013);
const stakeholderRPCClient = new StakeholderRPCClient(process.env.STAKEHOLDER_SERVICE_HOST || 'stakeholder-service', 3014);

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

// Protected stakeholder profile route (use RPC)
api.get('/api/stakeholder/profile', validateJWT, async (req, res) => {
  try {
    const username = req.user && req.user.username;
    const profile = await stakeholderRPCClient.getProfile(username);
    res.json(profile || {});
  } catch (err) {
    console.error('Stakeholder RPC getProfile error:', err);
    res.status(500).json({ error: 'Failed to fetch profile via RPC' });
  }
});

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


api.use('/api/purchases', validateJWT, createProxyMiddleware({
  target: PURCHASE_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/purchases': '',
  },
}));

// RPC-backed purchase handlers for cart checkout and tokens (match frontend routes)
api.post('/api/purchases/cart/checkout', validateJWT, express.json(), async (req, res) => {
  try {
    const username = req.user && req.user.username;
    const result = await purchaseRPCClient.checkout(username);
    res.json(result);
  } catch (err) {
    console.error('Purchase RPC checkout error:', err);
    res.status(500).json({ error: 'Failed to checkout via RPC' });
  }
});

api.get('/api/purchases/tokens', validateJWT, async (req, res) => {
  try {
    const username = req.user && req.user.username;
    const result = await purchaseRPCClient.getPurchasedTours(username);
    res.json(result.tokens || []);
  } catch (err) {
    console.error('Purchase RPC getPurchasedTours error:', err);
    res.status(500).json({ error: 'Failed to fetch purchased tokens via RPC' });
  }
});

api.use('/api/followers', validateJWT, createProxyMiddleware({
  target: FOLLOWER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/followers': '',
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