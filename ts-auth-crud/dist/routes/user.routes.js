"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const user_controller_1 = require("../controllers/user.controller");
const auth_middleware_1 = require("../middleware/auth.middleware");
const router = (0, express_1.Router)();
/**
 * Protected User CRUD Routes
 * All endpoints in this router require a valid 'Authorization: Bearer <token>' header
 */
// Apply JWT authentication middleware to every route in this router
router.use(auth_middleware_1.authenticateJWT);
// GET /users - Read all users
router.get("/", user_controller_1.getAllUsers);
// GET /users/me - Read current logged-in user profile (Declared BEFORE /:id so "me" is not treated as an ID)
router.get("/me", user_controller_1.getCurrentUser);
// GET /users/:id - Read a single user by ID
router.get("/:id", user_controller_1.getUserById);
// PUT /users/:id - Update user details
router.put("/:id", user_controller_1.updateUser);
// DELETE /users/:id - Delete a user
router.delete("/:id", user_controller_1.deleteUser);
exports.default = router;
