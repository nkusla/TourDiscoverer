const JWT_SECRET = process.env.JWT_SECRET || 'your_secret_key';

// Prefer explicit SERVICE_URL env vars when provided, else build from HOST and PORT
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || `http://${process.env.AUTH_SERVICE_HOST || 'auth-service'}:${process.env.AUTH_SERVICE_PORT || '3001'}`;
const TOUR_SERVICE_URL = process.env.TOUR_SERVICE_URL || `http://${process.env.TOUR_SERVICE_HOST || 'tour-service'}:${process.env.TOUR_SERVICE_PORT || '3006'}`;
const BLOG_SERVICE_URL = process.env.BLOG_SERVICE_URL || `http://${process.env.BLOG_SERVICE_HOST || 'blog-service'}:${process.env.BLOG_SERVICE_PORT || '3002'}`;
const FOLLOWER_SERVICE_URL = process.env.FOLLOWER_SERVICE_URL || `http://${process.env.FOLLOWER_SERVICE_HOST || 'follower-service'}:${process.env.FOLLOWER_SERVICE_PORT || '3005'}`;
const STAKEHOLDER_SERVICE_URL = process.env.STAKEHOLDER_SERVICE_URL || `http://${process.env.STAKEHOLDER_SERVICE_HOST || 'stakeholder-service'}:${process.env.STAKEHOLDER_SERVICE_PORT || '3003'}`;
const REVIEW_SERVICE_URL = process.env.REVIEW_SERVICE_URL || `http://${process.env.REVIEW_SERVICE_HOST || 'review-service'}:${process.env.REVIEW_SERVICE_PORT || '3007'}`;
const PURCHASE_SERVICE_URL = process.env.PURCHASE_SERVICE_URL || `http://${process.env.PURCHASE_SERVICE_HOST || 'purchase-service'}:${process.env.PURCHASE_SERVICE_PORT || '8084'}`;
const JAEGER_SERVICE_URL = process.env.JAEGER_SERVICE_URL;

module.exports = {
  JWT_SECRET,
  AUTH_SERVICE_URL,
  TOUR_SERVICE_URL,
  BLOG_SERVICE_URL,
  FOLLOWER_SERVICE_URL,
  STAKEHOLDER_SERVICE_URL,
  REVIEW_SERVICE_URL,
  PURCHASE_SERVICE_URL,
  JAEGER_SERVICE_URL
};