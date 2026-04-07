/**
 * API Response and User Data Models
 * 
 * This file contains TypeScript interfaces for API responses and user-related
 * data structures. These interfaces define the contract between the frontend
 * and backend API, ensuring type safety and consistent data handling.
 * 
 * @file api-response.model.ts
 * @author Clear Songs Development Team
 */

/**
 * Generic API Response Interface
 * 
 * Standard response structure for all API endpoints. This provides a consistent
 * format for success and error responses, making error handling predictable
 * across the application.
 * 
 * The interface is generic, allowing different response types for the data field
 * while maintaining a consistent structure.
 * 
 * @interface ApiResponse
 * @template T - The type of data returned in the response (defaults to unknown)
 * 
 * @example
 * // Success response with data
 * const response: ApiResponse<ArtistSummary[]> = {
 *   success: true,
 *   message: "Tracks retrieved successfully",
 *   data: [...]
 * };
 * 
 * // Error response
 * const error: ApiResponse = {
 *   success: false,
 *   error: "Invalid playlist ID"
 * };
 */
export interface ApiResponse<T = unknown> {
  /** Indicates whether the operation was successful */
  success: boolean;
  
  /** Optional success message from the server */
  message?: string;
  
  /** Optional response data (type depends on the endpoint) */
  data?: T;
  
  /** Optional structured error if the operation failed */
  error?: ApiError | string;
}

/**
 * API Error Interface
 *
 * Structured error payload returned by backend endpoints.
 */
export interface ApiError {
  code: string;
  message: string;
}

/**
 * User Interface
 * 
 * Represents a user account in the system. This includes both database
 * information and Spotify profile data.
 * 
 * @interface User
 */
export interface User {
  /** Unique identifier for the user (database ID) */
  id: string;
  
  /** Spotify user ID (from Spotify API) */
  spotify_id: string;
  
  /** Display name from Spotify profile */
  display_name: string;
  
  /** Optional email address from Spotify account */
  email?: string;
  
  /** Optional URL to the user's Spotify profile image */
  profile_image?: string;
  
  /** ISO 8601 timestamp when the user account was created */
  created_at?: string;
  
  /** ISO 8601 timestamp when the user account was last updated */
  updated_at?: string;
}

/**
 * Authentication Response Interface
 * 
 * Specialized API response for authentication endpoints. Includes user
 * information when authentication is successful.
 * 
 * @interface AuthResponse
 */
export type AuthResponse = ApiResponse<{ user?: User }> & { user?: User };
