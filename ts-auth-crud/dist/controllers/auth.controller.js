"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.login = exports.register = void 0;
const bcryptjs_1 = __importDefault(require("bcryptjs"));
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const user_db_1 = require("../db/user.db");
const auth_middleware_1 = require("../middleware/auth.middleware");
/**
 * Auth Controller
 * Handles user account creation (Registration) and credential validation (Login)
 */
// POST /auth/register
const register = async (req, res) => {
    try {
        const { username, email, password } = req.body;
        // 1. Validate required fields
        if (!username || !email || !password) {
            res.status(400).json({
                success: false,
                message: "Username, email, and password are all required.",
            });
            return;
        }
        // 2. Check if user already exists
        const existingUser = user_db_1.userDB.findByEmail(email);
        if (existingUser) {
            res.status(409).json({
                success: false,
                message: "A user with this email already exists.",
            });
            return;
        }
        // 3. Hash the plain password using bcrypt (cost factor of 10)
        const saltRounds = 10;
        const hashedPassword = await bcryptjs_1.default.hash(password, saltRounds);
        // 4. Save new user into the database
        const newUser = user_db_1.userDB.create({
            username,
            email,
            password: hashedPassword,
        });
        // 5. Generate signed JWT token
        const token = jsonwebtoken_1.default.sign({ userId: newUser.id, email: newUser.email }, auth_middleware_1.JWT_SECRET, { expiresIn: "1h" });
        // 6. Return 201 Created with token and safe user payload
        res.status(201).json({
            success: true,
            message: "User registered successfully!",
            token,
            user: (0, user_db_1.toUserResponse)(newUser),
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Internal server error during registration.",
        });
    }
};
exports.register = register;
// POST /auth/login
const login = async (req, res) => {
    try {
        const { email, password } = req.body;
        // 1. Validate input
        if (!email || !password) {
            res.status(400).json({
                success: false,
                message: "Email and password are required.",
            });
            return;
        }
        // 2. Find user by email
        const user = user_db_1.userDB.findByEmail(email);
        if (!user) {
            res.status(401).json({
                success: false,
                message: "Invalid email or password.",
            });
            return;
        }
        // 3. Compare supplied password against hashed password
        const isPasswordValid = await bcryptjs_1.default.compare(password, user.password);
        if (!isPasswordValid) {
            res.status(401).json({
                success: false,
                message: "Invalid email or password.",
            });
            return;
        }
        // 4. Issue new JWT token
        const token = jsonwebtoken_1.default.sign({ userId: user.id, email: user.email }, auth_middleware_1.JWT_SECRET, { expiresIn: "1h" });
        // 5. Send back 200 OK with token and user details
        res.status(200).json({
            success: true,
            message: "Login successful!",
            token,
            user: (0, user_db_1.toUserResponse)(user),
        });
    }
    catch (error) {
        res.status(500).json({
            success: false,
            message: "Internal server error during login.",
        });
    }
};
exports.login = login;
