const jwt = require('jsonwebtoken');
const { trace, context, SpanStatusCode } = require('@opentelemetry/api');
const { JWT_SECRET } = require('./constants');

function validateJWT(req, res, next) {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return res.status(401).json({ message: 'Missing Authorization header' });
  }

  const token = authHeader.split(' ')[1];
  if (!token) {
    return res.status(401).json({ message: 'Token missing' });
  }

  try {
    const decoded = jwt.verify(token, JWT_SECRET);

		if (!decoded.username || !decoded.role) {
			return res.status(403).json({
				message: 'Invalid token: missing required fields (username and role)'
			});
		}

    req.headers['x-user-role'] = decoded.role;
    req.headers['x-username'] = decoded.username;

    req.user = decoded;
    next();
  } catch (err) {
    return res.status(403).json({ message: 'Invalid or expired token' });
  }
}

function blockInternalRoutes(req, res, next) {
  if (req.path.includes('/internal')) {
    return res.status(403).json({
      error: 'Forbidden',
      message: 'Internal routes are not accessible externally'
    });
  }
  next();
}

function tracingMiddleware(req, res, next) {
  const tracer = trace.getTracer('api-gateway');

  // Create new trace/span
  const span = tracer.startSpan(`${req.method} ${req.path}`);

  // Add span to request context
  req.span = span;
  req.traceContext = trace.setSpan(context.active(), span);

  // Add event
  span.addEvent('HTTP request received');
  span.setAttributes({
    'http.method': req.method,
    'http.url': req.originalUrl,
    'http.user_agent': req.get('User-Agent') || '',
    'user.name': req.user?.username || 'anonymous'
  });

  // End span when response finishes
  res.on('finish', () => {
    span.setAttributes({
      'http.status_code': res.statusCode
    });

    if (res.statusCode >= 400) {
      span.setStatus({
        code: SpanStatusCode.ERROR,
        message: `HTTP ${res.statusCode}`
      });
    } else {
      span.setStatus({
        code: SpanStatusCode.OK
      });
    }

    span.end();
  });

  next();
}

module.exports = {
  validateJWT,
  blockInternalRoutes,
  tracingMiddleware
};