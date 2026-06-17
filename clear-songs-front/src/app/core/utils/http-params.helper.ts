/**
 * HTTP Parameters Helper Utilities
 * 
 * This module provides utility functions for building HTTP query parameters
 * in a consistent and reusable way across the application.
 * 
 * The helper functions abstract away the repetitive logic of creating
 * HttpParams objects, making the code more maintainable and reducing
 * duplication.
 * 
 * @file http-params.helper.ts
 * @author Clear Songs Development Team
 */
import { HttpParams } from '@angular/common/http';

/**
 * Builds HttpParams for range-based queries
 * 
 * Creates an HttpParams object with optional min, max, and genre values for range queries.
 * This is commonly used for filtering tracks by count ranges (e.g., "artists with
 * 5-10 tracks" or "artists with at least 20 tracks") and by genre.
 * 
 * The function only adds parameters that are defined, allowing for flexible queries:
 * - Both min and max: Range query (e.g., min=5, max=10)
 * - Only min: Minimum threshold (e.g., min=20)
 * - Only max: Maximum threshold (e.g., max=5)
 * - Neither: No range filtering
 * - Genre: Filter by artist genre (e.g., genre=rock)
 * 
 * @param min - Optional minimum value for the range. If undefined, not added to params.
 * @param max - Optional maximum value for the range. If undefined, not added to params.
 * @param genre - Optional genre filter. If undefined or empty, not added to params.
 * @returns HttpParams object with min/max/genre query parameters set if provided
 * 
 * @example
 * // Range query: artists with 5-10 tracks
 * const params = buildRangeParams(5, 10);
 * // Result: ?min=5&max=10
 * 
 * // Minimum only: artists with at least 20 tracks
 * const params = buildRangeParams(20, undefined);
 * // Result: ?min=20
 * 
 * // Maximum only: artists with at most 5 tracks
 * const params = buildRangeParams(undefined, 5);
 * // Result: ?max=5
 * 
 * // Genre filter
 * const params = buildRangeParams(undefined, undefined, 'rock');
 * // Result: ?genre=rock
 * 
 * // Combined filters
 * const params = buildRangeParams(1, 10, 'rock');
 * // Result: ?min=1&max=10&genre=rock
 * 
 * // Usage in HTTP request:
 * this.http.get('/api/tracks/summary', { params: buildRangeParams(5, 10) })
 */
export function buildRangeParams(min?: number, max?: number, genre?: string): HttpParams {
  let params = new HttpParams();
  
  if (min !== undefined) {
    params = params.set('min', min.toString());
  }
  
  if (max !== undefined) {
    params = params.set('max', max.toString());
  }
  
  if (genre && genre.trim() !== '') {
    params = params.set('genre', genre.trim());
  }
  
  return params;
}

