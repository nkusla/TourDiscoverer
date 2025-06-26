const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const morgan = require('morgan');
const { validateJWT, verifyUserExists } = require('./middleware');

const dotenv = require('dotenv');
dotenv.config();

const app = express();
app.use(morgan('dev'));

app.use('/api/auth', createProxyMiddleware({
  target: process.env.AUTH_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '', // Remove /api/auth prefix when forwarding
  },
}));

app.get('/ping', (req, res) => {
  res.status(200).json({
    message: 'pong',
    timestamp: new Date().toISOString(),
    service: 'API Gateway'
  });
});

const PORT = process.env.API_GATEWAY_PORT || 3000;
app.listen(PORT, () => {
  console.log(`API Gateway listening on port ${PORT}`);
});