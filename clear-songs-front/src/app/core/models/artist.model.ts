/**
 * Artist and Track Data Models
 * 
 * This file contains TypeScript interfaces that define the structure of
 * artist, track, and playlist data used throughout the application.
 * 
 * These interfaces ensure type safety and provide IntelliSense support
 * when working with data from the backend API or Spotify.
 * 
 * @file artist.model.ts
 * @author Clear Songs Development Team
 */

/**
 * Artist Summary Interface
 * 
 * Represents a summary of an artist with the count of tracks in the user's library.
 * This is used in the dashboard to display artist statistics and allow bulk
 * deletion of tracks by artist.
 * 
 * @interface ArtistSummary
 */
export interface ArtistSummary {
  /** Unique identifier for the artist (database ID) */
  id: string;
  
  /** Display name of the artist */
  name: string;
  
  /** Number of tracks from this artist in the user's library */
  count: number;
  
  /** URL of the artist's profile image (optional) */
  image_url?: string;
  
  /** Array of genre names associated with the artist (optional) */
  genres?: string[];

  /** Resolved canonical genre (e.g., "Hip Hop", "Rock", "Electronic") */
  genre?: string;
}

/**
 * Playlist Interface
 * 
 * Represents a simplified playlist object with essential information
 * for display and selection purposes.
 * 
 * @interface PlaylistRecord
 */
export interface PlaylistRecord {
  /** Unique identifier for the playlist (Spotify ID) */
  id: string;
  
  /** Display name of the playlist */
  name: string;
  
  /** URL of the playlist's cover image (optional) */
  image_url?: string;
}

/**
 * Track Interface
 * 
 * Represents a complete track object with all metadata from Spotify and the database.
 * This interface matches the structure returned by the backend API and stored
 * in the database.
 * 
 * @interface Track
 */
export interface Track {
  /** Unique identifier for the track (Spotify ID or database ID) */
  id: string;
  
  /** Name of the track */
  name: string;
  
  /** Array of artist names associated with this track */
  artists: string[];
  
  /** Name of the album containing this track */
  album: string;
  
  /** Duration of the track in milliseconds */
  duration: number;
  
  /** Popularity score from Spotify (0-100) (optional) */
  popularity?: number;
  
  /** Whether the track contains explicit content (optional) */
  explicit?: boolean;
  
  /** Spotify track ID (optional) */
  spotify_id?: string;
  
  /** Spotify URL (optional) */
  spotify_url?: string;
  
  /** Image URL (optional) */
  image_url?: string;
  
  /** ISO 8601 timestamp when the track was added to the database (optional) */
  created_at?: string;
  
  /** ISO 8601 timestamp when the track was last updated (optional) */
  updated_at?: string;
}

/**
 * User Playlist Interface
 * 
 * Represents a simplified playlist object for user playlist selection.
 * This interface matches the structure returned by the backend API.
 * 
 * @interface UserPlaylist
 */
export interface UserPlaylist {
  /** Unique identifier for the playlist (Spotify ID) */
  id: string;
  
  /** Display name of the playlist */
  name: string;
  
  /** URL of the playlist's cover image (optional) */
  image_url?: string;
}

/**
 * Playlist Interface (Database Model)
 * 
 * Represents a complete playlist object stored in the database.
 * This interface is used when working with playlist data from the database.
 * 
 * @interface Playlist
 */
export interface Playlist {
  /** Unique identifier for the playlist (database ID) */
  id: string;
  
  /** Name of the playlist */
  name: string;
  
  /** Optional description of the playlist */
  description?: string;
  
  /** Spotify playlist ID (for API calls and links) */
  spotify_id: string;
  
  /** Number of tracks in the playlist */
  tracks_count: number;
  
  /** ISO 8601 timestamp when the playlist was added to the database */
  created_at: string;
  
  /** ISO 8601 timestamp when the playlist was last updated */
  updated_at: string;
}