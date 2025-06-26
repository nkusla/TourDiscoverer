const jwt = require('jsonwebtoken');
const axios = require('axios');

const JWT_SECRET = process.env.JWT_SECRET || 'your_secret_key';

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

    req.user = decoded;
    next();
  } catch (err) {
    return res.status(403).json({ message: 'Invalid or expired token' });
  }
}

async function verifyUserExists(req, res, next) {
  const username = req.user?.username;
  if (!username) {
    return res.status(400).json({ message: 'Username missing in token' });
  }

  try {
    const response = await axios.get(`${AUTH_SERVICE_URL}/user/${username}/exists`);

    if (response.status === 200) {
      return next();
    } else {
      return res.status(403).json({ message: 'User not found or inactive' });
    }
  } catch (err) {
    return res.status(403).json({ message: 'User verification failed', error: err.message });
  }
}