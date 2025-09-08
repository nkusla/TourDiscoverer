/**
 * Decode JWT token to extract user information
 * @param {string} token - JWT token
 * @returns {object|null} - Decoded user data or null if invalid
 */
export function decodeJWT(token) {
  if (!token) return null
  
  try {
    // JWT has 3 parts separated by dots: header.payload.signature
    const parts = token.split('.')
    if (parts.length !== 3) return null
    
    // Decode the payload (second part)
    const payload = parts[1]
    
    // Add padding if needed for base64 decoding
    const paddedPayload = payload + '='.repeat((4 - payload.length % 4) % 4)
    
    // Decode base64
    const decodedPayload = atob(paddedPayload)
    
    // Parse JSON
    const userData = JSON.parse(decodedPayload)
    
    return userData
  } catch (error) {
    console.error('Error decoding JWT:', error)
    return null
  }
}

/**
 * Check if JWT token is expired
 * @param {string} token - JWT token
 * @returns {boolean} - True if expired, false otherwise
 */
export function isTokenExpired(token) {
  const decoded = decodeJWT(token)
  if (!decoded || !decoded.exp) return true
  
  // exp is in seconds, Date.now() is in milliseconds
  return decoded.exp * 1000 < Date.now()
}

/**
 * Extract user info from JWT token
 * @param {string} token - JWT token
 * @returns {object|null} - User information or null
 */
export function getUserFromToken(token) {
  const decoded = decodeJWT(token)
  if (!decoded) return null
  
  return {
    username: decoded.username,
    role: decoded.role
  }
}
