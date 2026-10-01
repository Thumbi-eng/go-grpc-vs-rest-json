import { Router } from "express";
import { register, login } from "../controllers/auth.controller";

const router = Router();

/**
 * Public Authentication Routes
 * These endpoints are open to all clients (no token required)
 */

// POST /auth/register - Create a new user account
router.post("/register", register);

// POST /auth/login - Authenticate credentials and receive a JWT token
router.post("/login", login);

export default router;
