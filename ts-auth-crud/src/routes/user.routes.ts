import { Router } from "express";
import {
  getAllUsers,
  getCurrentUser,
  getUserById,
  updateUser,
  deleteUser,
} from "../controllers/user.controller";
import { authenticateJWT } from "../middleware/auth.middleware";

const router = Router();

/**
 * Protected User CRUD Routes
 * All endpoints in this router require a valid 'Authorization: Bearer <token>' header
 */

// Apply JWT authentication middleware to every route in this router
router.use(authenticateJWT);

// GET /users - Read all users
router.get("/", getAllUsers);

// GET /users/me - Read current logged-in user profile (Declared BEFORE /:id so "me" is not treated as an ID)
router.get("/me", getCurrentUser);

// GET /users/:id - Read a single user by ID
router.get("/:id", getUserById);

// PUT /users/:id - Update user details
router.put("/:id", updateUser);

// DELETE /users/:id - Delete a user
router.delete("/:id", deleteUser);

export default router;
