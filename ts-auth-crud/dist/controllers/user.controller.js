"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deleteUser = exports.updateUser = exports.getUserById = exports.getCurrentUser = exports.getAllUsers = void 0;
const bcryptjs_1 = __importDefault(require("bcryptjs"));
const user_db_1 = require("../db/user.db");
/**
 * User CRUD Controller
 * Handles Read, Update, and Delete operations for user entities
 */
// GET /users - Retrieve list of all registered users (passwords excluded)
const getAllUsers = async (req, res) => {
    try {
        const users = user_db_1.userDB.findAll();
        res.status(200).json({
            success: true,
            count: users.length,
            users,
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Failed to fetch users.",
        });
    }
};
exports.getAllUsers = getAllUsers;
// GET /users/me - Retrieve profile of the currently logged-in user
const getCurrentUser = async (req, res) => {
    try {
        const userId = req.user?.userId;
        if (!userId) {
            res.status(401).json({ success: false, message: "Unauthorized." });
            return;
        }
        const user = user_db_1.userDB.findById(userId);
        if (!user) {
            res.status(404).json({ success: false, message: "User not found." });
            return;
        }
        res.status(200).json({
            success: true,
            user: (0, user_db_1.toUserResponse)(user),
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Failed to fetch profile.",
        });
    }
};
exports.getCurrentUser = getCurrentUser;
// GET /users/:id - Retrieve a single user by their ID
const getUserById = async (req, res) => {
    try {
        const { id } = req.params;
        const user = user_db_1.userDB.findById(id);
        if (!user) {
            res.status(404).json({
                success: false,
                message: `User with ID ${id} not found.`,
            });
            return;
        }
        res.status(200).json({
            success: true,
            user: (0, user_db_1.toUserResponse)(user),
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Failed to fetch user.",
        });
    }
};
exports.getUserById = getUserById;
// PUT /users/:id - Update user details (username, email, or password)
const updateUser = async (req, res) => {
    try {
        const { id } = req.params;
        const { username, email, password } = req.body;
        // 1. Check if the target user exists
        const existingUser = user_db_1.userDB.findById(id);
        if (!existingUser) {
            res.status(404).json({
                success: false,
                message: `User with ID ${id} not found.`,
            });
            return;
        }
        // 2. If changing email, ensure it's not already claimed by another user
        if (email && email.toLowerCase() !== existingUser.email.toLowerCase()) {
            const emailTaken = user_db_1.userDB.findByEmail(email);
            if (emailTaken && emailTaken.id !== id) {
                res.status(409).json({
                    success: false,
                    message: "Email is already taken by another account.",
                });
                return;
            }
        }
        // 3. If updating password, hash the new password
        let hashedPassword = undefined;
        if (password) {
            hashedPassword = await bcryptjs_1.default.hash(password, 10);
        }
        // 4. Perform update in DB
        const updatedUser = user_db_1.userDB.update(id, {
            username,
            email,
            ...(hashedPassword && { password: hashedPassword }),
        });
        res.status(200).json({
            success: true,
            message: "User updated successfully.",
            user: (0, user_db_1.toUserResponse)(updatedUser),
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Failed to update user.",
        });
    }
};
exports.updateUser = updateUser;
// DELETE /users/:id - Remove a user by ID
const deleteUser = async (req, res) => {
    try {
        const { id } = req.params;
        const deleted = user_db_1.userDB.deleteById(id);
        if (!deleted) {
            res.status(404).json({
                success: false,
                message: `User with ID ${id} not found.`,
            });
            return;
        }
        res.status(200).json({
            success: true,
            message: `User with ID ${id} was deleted successfully.`,
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Failed to delete user.",
        });
    }
};
exports.deleteUser = deleteUser;
