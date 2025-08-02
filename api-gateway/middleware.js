const jwt = require('jsonwebtoken');
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

module.exports = {
  validateJWT,
  blockInternalRoutes
};