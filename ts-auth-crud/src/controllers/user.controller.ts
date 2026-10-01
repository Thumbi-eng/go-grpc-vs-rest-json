import { Response } from "express";
import bcrypt from "bcryptjs";
import { AuthenticatedRequest, UpdateUserDTO } from "../types/user.types";
import { userDB, toUserResponse } from "../db/user.db";

/**
 * User CRUD Controller
 * Handles Read, Update, and Delete operations for user entities
 */

// GET /users - Retrieve list of all registered users (passwords excluded)
export const getAllUsers = async (
  req: AuthenticatedRequest,
  res: Response
): Promise<void> => {
  try {
    const users = userDB.findAll();
    res.status(200).json({
      success: true,
      count: users.length,
      users,
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Failed to fetch users.",
    });
  }
};

// GET /users/me - Retrieve profile of the currently logged-in user
export const getCurrentUser = async (
  req: AuthenticatedRequest,
  res: Response
): Promise<void> => {
  try {
    const userId = req.user?.userId;
    if (!userId) {
      res.status(401).json({ success: false, message: "Unauthorized." });
      return;
    }

    const user = userDB.findById(userId);
    if (!user) {
      res.status(404).json({ success: false, message: "User not found." });
      return;
    }

    res.status(200).json({
      success: true,
      user: toUserResponse(user),
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Failed to fetch profile.",
    });
  }
};

// GET /users/:id - Retrieve a single user by their ID
export const getUserById = async (
  req: AuthenticatedRequest,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;
    const user = userDB.findById(id);

    if (!user) {
      res.status(404).json({
        success: false,
        message: `User with ID ${id} not found.`,
      });
      return;
    }

    res.status(200).json({
      success: true,
      user: toUserResponse(user),
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Failed to fetch user.",
    });
  }
};

// PUT /users/:id - Update user details (username, email, or password)
export const updateUser = async (
  req: AuthenticatedRequest,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;
    const { username, email, password }: UpdateUserDTO = req.body;

    // 1. Check if the target user exists
    const existingUser = userDB.findById(id);
    if (!existingUser) {
      res.status(404).json({
        success: false,
        message: `User with ID ${id} not found.`,
      });
      return;
    }

    // 2. If changing email, ensure it's not already claimed by another user
    if (email && email.toLowerCase() !== existingUser.email.toLowerCase()) {
      const emailTaken = userDB.findByEmail(email);
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
      hashedPassword = await bcrypt.hash(password, 10);
    }

    // 4. Perform update in DB
    const updatedUser = userDB.update(id, {
      username,
      email,
      ...(hashedPassword && { password: hashedPassword }),
    });

    res.status(200).json({
      success: true,
      message: "User updated successfully.",
      user: toUserResponse(updatedUser!),
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Failed to update user.",
    });
  }
};

// DELETE /users/:id - Remove a user by ID
export const deleteUser = async (
  req: AuthenticatedRequest,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    const deleted = userDB.deleteById(id);
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
  } catch (error) {
    res.status(500).json({
      success: false,
      message: "Failed to delete user.",
    });
  }
};
