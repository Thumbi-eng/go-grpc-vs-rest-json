"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const auth_controller_1 = require("../controllers/auth.controller");
const router = (0, express_1.Router)();
/**
 * Public Authentication Routes
 * These endpoints are open to all clients (no token required)
 */
// POST /auth/register - Create a new user account
router.post("/register", auth_controller_1.register);
// POST /auth/login - Authenticate credentials and receive a JWT token
router.post("/login", auth_controller_1.login);
exports.default = router;
