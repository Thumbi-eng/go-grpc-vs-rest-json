import dotenv from "dotenv";
dotenv.config();

import express, { Application, Request, Response } from "express";
import authRoutes from "./routes/auth.routes";
import userRoutes from "./routes/user.routes";

const app: Application = express();
const PORT = process.env.PORT || 3000;

/**
 * Global Middlewares
 */
// Parse incoming requests with JSON payloads (essential for reading req.body)
app.use(express.json());

/**
 * Health Check & API Info Root Endpoint
 */
app.get("/", (req: Request, res: Response) => {
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
app.use("/auth", authRoutes);
app.use("/users", userRoutes);

/**
 * 404 Catch-All Handler
 */
app.use((req: Request, res: Response) => {
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

export default app;
