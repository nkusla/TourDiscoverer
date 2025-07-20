const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, validateJWTWithRole, verifyUserExists } = require('./middleware');
const { AUTH_SERVICE_URL, TOUR_SERVICE_URL, BLOG_SERVICE_URL, USER_ROLES } = require('./constants');

const dotenv = require('dotenv');
dotenv.config();

const api = express();
api.use(morgan('dev'));

const blockInternalRoutes = (req, res, next) => {
  if (req.path.includes('/internal')) {
    return res.status(403).json({
      error: 'Forbidden',
      message: 'Internal routes are not accessible externally'
    });
  }
  next();
};

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