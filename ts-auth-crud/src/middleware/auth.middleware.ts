import { Response, NextFunction } from "express";
import jwt from "jsonwebtoken";
import { AuthenticatedRequest, JwtPayload } from "../types/user.types";

/**
 * Secret key used to sign and verify JSON Web Tokens.
 * Falls back to a development secret if JWT_SECRET is not set in environment variables.
 */
export const JWT_SECRET = process.env.JWT_SECRET || "supersecret_jwt_key_12345";

/**
 * JWT Authentication Middleware (The Bouncer)
 * Intercepts incoming requests to protected routes:
 * 1. Checks for 'Authorization: Bearer <token>' header
 * 2. Verifies token validity and expiration
 * 3. Injects decoded user payload into 'req.user'
 * 4. Calls next() or rejects with HTTP 401 Unauthorized
 */
export const authenticateJWT = (
  req: AuthenticatedRequest,
  res: Response,
  next: NextFunction
): void => {
  const authHeader = req.headers.authorization;

  // 1. Verify header presence and Bearer scheme
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    res.status(401).json({
      success: false,
      message: "Access denied. No Bearer token provided in Authorization header.",
    });
    return;
  }

  // 2. Extract token string
  const token = authHeader.split(" ")[1];

  try {
    // 3. Verify token signature and expiration
    const decoded = jwt.verify(token, JWT_SECRET) as JwtPayload;

    // 4. Attach decoded payload to request object for downstream controllers
    req.user = decoded;

    // 5. Pass control to the next handler
    next();
  } catch (error) {
    res.status(401).json({
      success: false,
      message: "Invalid or expired token. Please log in again.",
    });
    return;
  }
};
