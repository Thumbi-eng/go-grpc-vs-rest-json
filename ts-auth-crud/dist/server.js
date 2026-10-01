"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
const express_1 = __importDefault(require("express"));
const auth_routes_1 = __importDefault(require("./routes/auth.routes"));
const user_routes_1 = __importDefault(require("./routes/user.routes"));
const app = (0, express_1.default)();
const PORT = process.env.PORT || 3000;
/**
 * Global Middlewares
 */
// Parse incoming requests with JSON payloads (essential for reading req.body)
app.use(express_1.default.json());
/**
 * Health Check & API Info Root Endpoint
 */
app.get("/", (req, res) => {
    res.status(200).json({
        success: true,
        message: "TypeScript Auth & CRUD REST API is up and running!",
        endpoints: {
            public: {
                register: "POST /auth/register",
                login: "POST /auth/login",
            },
            protected: {
                getAllUsers: "GET /users",
                getProfile: "GET /users/me",
                getUserById: "GET /users/:id",
                updateUser: "PUT /users/:id",
                deleteUser: "DELETE /users/:id",
            },
        },
    });
});
/**
 * Mount Routers
 */
app.use("/auth", auth_routes_1.default);
app.use("/users", user_routes_1.default);
/**
 * 404 Catch-All Handler
 */
app.use((req, res) => {
    res.status(404).json({
        success: false,
        message: `Cannot ${req.method} ${req.originalUrl} - Route not found.`,
    });
});
/**
 * Start the HTTP Server
 */
app.listen(PORT, () => {
    console.log(`===============================================`);
    console.log(`🚀 Server running on: http://localhost:${PORT}`);
    console.log(`📋 Base Auth URLs:   http://localhost:${PORT}/auth`);
    console.log(`👥 Base User URLs:   http://localhost:${PORT}/users`);
    console.log(`===============================================`);
});
exports.default = app;
