import { Request, Response } from "express";
import bcrypt from "bcryptjs";
import jwt from "jsonwebtoken";
import { userDB, toUserResponse } from "../db/user.db";
import { RegisterDTO, LoginDTO } from "../types/user.types";
import { JWT_SECRET } from "../middleware/auth.middleware";

/**
 * Auth Controller
 * Handles user account creation (Registration) and credential validation (Login)
 */

// POST /auth/register
export const register = async (req: Request, res: Response): Promise<void> => {
  try {
    const { username, email, password }: RegisterDTO = req.body;

    // 1. Validate required fields
    if (!username || !email || !password) {
      res.status(400).json({
        success: false,
        message: "Username, email, and password are all required.",
      });
      return;
    }

    // 2. Check if user already exists
    const existingUser = userDB.findByEmail(email);
    if (existingUser) {
      res.status(409).json({
        success: false,
        message: "A user with this email already exists.",
      });
      return;
    }

    // 3. Hash the plain password using bcrypt (cost factor of 10)
    const saltRounds = 10;
    const hashedPassword = await bcrypt.hash(password, saltRounds);

    // 4. Save new user into the database
    const newUser = userDB.create({
      username,
      email,
      password: hashedPassword,
    });

    // 5. Generate signed JWT token
    const token = jwt.sign(
      { userId: newUser.id, email: newUser.email },
      JWT_SECRET,
      { expiresIn: "1h" }
    );

    // 6. Return 201 Created with token and safe user payload
    res.status(201).json({
      success: true,
      message: "User registered successfully!",
      token,
      user: toUserResponse(newUser),
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Internal server error during registration.",
    });
  }
};

// POST /auth/login
export const login = async (req: Request, res: Response): Promise<void> => {
  try {
    const { email, password }: LoginDTO = req.body;

    // 1. Validate input
    if (!email || !password) {
      res.status(400).json({
        success: false,
        message: "Email and password are required.",
      });
      return;
    }

    // 2. Find user by email
    const user = userDB.findByEmail(email);
    if (!user) {
      res.status(401).json({
        success: false,
        message: "Invalid email or password.",
      });
      return;
    }

    // 3. Compare supplied password against hashed password
    const isPasswordValid = await bcrypt.compare(password, user.password);
    if (!isPasswordValid) {
      res.status(401).json({
        success: false,
        message: "Invalid email or password.",
      });
      return;
    }

    // 4. Issue new JWT token
    const token = jwt.sign(
      { userId: user.id, email: user.email },
      JWT_SECRET,
      { expiresIn: "1h" }
    );

    // 5. Send back 200 OK with token and user details
    res.status(200).json({
      success: true,
      message: "Login successful!",
      token,
      user: toUserResponse(user),
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Internal server error during login.",
    });
  }
};
