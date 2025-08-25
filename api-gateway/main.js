const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, blockInternalRoutes } = require('./middleware');
const { AUTH_SERVICE_URL, STAKEHOLDER_SERVICE_URL, TOUR_SERVICE_URL, BLOG_SERVICE_URL } = require('./constants');

const dotenv = require('dotenv');
dotenv.config();

const api = express();
api.use(morgan('dev'));

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

// Public stakeholder routes (for creating profiles)
api.use('/api/stakeholder', createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/stakeholder': '',
  }
}));

// Protected stakeholder routes (profile management)
api.use('/api/profile', validateJWT, createProxyMiddleware({
  target: STAKEHOLDER_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/profile': '',
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