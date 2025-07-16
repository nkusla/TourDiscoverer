const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, validateJWTWithRole, verifyUserExists } = require('./middleware');
const { AUTH_SERVICE_URL, TOUR_SERVICE_URL, USER_ROLES } = require('./constants');

const dotenv = require('dotenv');
dotenv.config();

const api = express();
api.use(morgan('dev'));

// AUTH SERVICE PROXIES

api.get('/api/auth/users', validateJWTWithRole(USER_ROLES.ADMIN), createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '',
  }
}));

api.post('/api/auth/block', validateJWTWithRole(USER_ROLES.ADMIN), createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '',
  }
}));

api.use('/api/auth', createProxyMiddleware({
  target: AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '', // Remove /api/auth prefix when forwarding
  },
}));

// TOUR SERVICE PROXIES

api.use('/api/tours', validateJWT, createProxyMiddleware({
  target: TOUR_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/tours': '', 
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