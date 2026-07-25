/**
 * HTTP parameter utilities for building query strings.
 */
import { HttpParams } from '@angular/common/http';

/**
 * Builds HttpParams for range-based queries with optional min, max, and genre filters.
 * @example buildRangeParams(5, 10, 'rock') // ?min=5&max=10&genre=rock
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

