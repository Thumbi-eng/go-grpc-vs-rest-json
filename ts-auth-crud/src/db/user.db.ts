import { randomUUID } from "crypto";
import { User, UserResponse } from "../types/user.types";

/**
 * In-Memory Database
 * Simulates a real database table using a TypeScript array.
 * Zero database setup required for learners!
 */
export const usersTable: User[] = [];

/**
 * Strips sensitive data (like hashed passwords) before sending responses
 */
export const toUserResponse = (user: User): UserResponse => {
  const { password, ...safeUser } = user;
  return safeUser;
};

/**
 * Database Helper Functions (Simulating Database Repository/ORM methods)
 */
export const userDB = {
  // Find all users (returns safe objects without passwords)
  findAll: (): UserResponse[] => {
    return usersTable.map(toUserResponse);
  },

  // Find a user by their unique ID
  findById: (id: string): User | undefined => {
    return usersTable.find((user) => user.id === id);
  },

  // Find a user by their email address (useful for login and uniqueness checks)
  findByEmail: (email: string): User | undefined => {
    return usersTable.find(
      (user) => user.email.toLowerCase() === email.toLowerCase()
    );
  },

  // Insert a new user record
  create: (data: { username: string; email: string; password: string }): User => {
    const newUser: User = {
      id: randomUUID(),
      username: data.username,
      email: data.email.toLowerCase(),
      password: data.password,
      createdAt: new Date(),
    };
    usersTable.push(newUser);
    return newUser;
  },

  // Update an existing user record
  update: (
    id: string,
    updates: Partial<{ username: string; email: string; password: string }>
  ): User | undefined => {
    const userIndex = usersTable.findIndex((u) => u.id === id);
    if (userIndex === -1) return undefined;

    const existing = usersTable[userIndex];
    usersTable[userIndex] = {
      ...existing,
      username: updates.username ?? existing.username,
      email: updates.email ? updates.email.toLowerCase() : existing.email,
      password: updates.password ?? existing.password,
    };

    return usersTable[userIndex];
  },

  // Delete a user by ID
  deleteById: (id: string): boolean => {
    const userIndex = usersTable.findIndex((u) => u.id === id);
    if (userIndex === -1) return false;

    usersTable.splice(userIndex, 1);
    return true;
  },
};
