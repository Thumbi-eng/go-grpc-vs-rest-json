import { Request } from "express";

/**
 * Core User Entity
 * Represents how a user is stored in our database
 */
export interface User {
  id: string;
  username: string;
  email: string;
  password: string; // Stored as a bcrypt hash, never plaintext
  createdAt: Date;
}

/**
 * Safe User Response
 * Excludes sensitive fields like password before sending data back to client
 */
export type UserResponse = Omit<User, "password">;

/**
 * Data Transfer Objects (DTOs) for incoming HTTP request bodies
 */

// Payload for POST /auth/register
export interface RegisterDTO {
  username: string;
  email: string;
  password: string;
}

// Payload for POST /auth/login
export interface LoginDTO {
  email: string;
  password: string;
}

// Payload for PUT /users/:id (all fields optional for partial updates)
export interface UpdateUserDTO {
  username?: string;
  email?: string;
  password?: string;
}

/**
 * JWT Token Payload
 * Data encoded inside the JSON Web Token
 */
export interface JwtPayload {
  userId: string;
  email: string;
}

/**
 * Extended Express Request
 * Populated by auth.middleware with the decoded JWT payload
 */
export interface AuthenticatedRequest extends Request {
  user?: JwtPayload;
}
