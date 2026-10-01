"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.userDB = exports.toUserResponse = exports.usersTable = void 0;
const crypto_1 = require("crypto");
/**
 * In-Memory Database
 * Simulates a real database table using a TypeScript array.
 * Zero database setup required for learners!
 */
exports.usersTable = [];
/**
 * Strips sensitive data (like hashed passwords) before sending responses
 */
const toUserResponse = (user) => {
    const { password, ...safeUser } = user;
    return safeUser;
};
exports.toUserResponse = toUserResponse;
/**
 * Database Helper Functions (Simulating Database Repository/ORM methods)
 */
exports.userDB = {
    // Find all users (returns safe objects without passwords)
    findAll: () => {
        return exports.usersTable.map(exports.toUserResponse);
    },
    // Find a user by their unique ID
    findById: (id) => {
        return exports.usersTable.find((user) => user.id === id);
    },
    // Find a user by their email address (useful for login and uniqueness checks)
    findByEmail: (email) => {
        return exports.usersTable.find((user) => user.email.toLowerCase() === email.toLowerCase());
    },
    // Insert a new user record
    create: (data) => {
        const newUser = {
            id: (0, crypto_1.randomUUID)(),
            username: data.username,
            email: data.email.toLowerCase(),
            password: data.password,
            createdAt: new Date(),
        };
        exports.usersTable.push(newUser);
        return newUser;
    },
    // Update an existing user record
    update: (id, updates) => {
        const userIndex = exports.usersTable.findIndex((u) => u.id === id);
        if (userIndex === -1)
            return undefined;
        const existing = exports.usersTable[userIndex];
        exports.usersTable[userIndex] = {
            ...existing,
            username: updates.username ?? existing.username,
            email: updates.email ? updates.email.toLowerCase() : existing.email,
            password: updates.password ?? existing.password,
        };
        return exports.usersTable[userIndex];
    },
    // Delete a user by ID
    deleteById: (id) => {
        const userIndex = exports.usersTable.findIndex((u) => u.id === id);
        if (userIndex === -1)
            return false;
        exports.usersTable.splice(userIndex, 1);
        return true;
    },
};
