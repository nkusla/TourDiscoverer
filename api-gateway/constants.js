const USER_ROLES = Object.freeze({
  ADMIN: 'admin',
  GUIDE: 'guide',
  TOURIST: 'tourist',
});

const JWT_SECRET = process.env.JWT_SECRET || 'your_secret_key';
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL;
const TOUR_SERVICE_URL = process.env.TOUR_SERVICE_URL;

module.exports = {
  USER_ROLES,
  JWT_SECRET,
  AUTH_SERVICE_URL,
  TOUR_SERVICE_URL,
};