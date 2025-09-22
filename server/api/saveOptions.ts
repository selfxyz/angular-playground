import { Request, Response } from "express";
import { Redis } from "@upstash/redis";
import dotenv from 'dotenv';
dotenv.config();

// Initialize Redis client
const redis = new Redis({
  url: process.env['KV_REST_API_URL']!,
  token: process.env['KV_REST_API_TOKEN']!,
});

export async function saveOptionsHandler(
  req: Request,
  res: Response
) {
  if (req.method !== "POST") {
    return res.status(405).json({ message: "Method not allowed" });
  }

  try {
    const { userId, options } = req.body;

    if (!userId) {
      return res.status(400).json({ message: "User ID is required" });
    }

    if (!options) {
      return res.status(400).json({ message: "Options are required" });
    }

    console.log("Saving options for user:", userId, options);
    // Store the options in Redis with a 30-minute expiration
    await redis.set(userId, JSON.stringify(options), { ex: 1800 }); // 1800 seconds = 30 minutes

    return res.status(200).json({ message: "Options saved successfully" });
  } catch (error) {
    console.error("Error saving options:", error);
    return res.status(500).json({
      message: "Internal server error",
      error: error instanceof Error ? error.message : "Unknown error",
    });
  }
}
