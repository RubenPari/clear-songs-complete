/**
 * Track service for managing saved tracks via the backend API.
 * Provides reactive resource fetching, deletion operations, and cache invalidation.
 */
import { Injectable, inject } from '@angular/core';
import { HttpClient, httpResource } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ArtistSummary, Track } from '../models/artist.model';
import { ApiResponse } from '../models/api-response.model';
import { buildRangeParams } from '../utils/http-params.helper';

/**
 * Service for managing saved tracks.
 * Provides track summary fetching, deletion by artist/range, and cache management.
 */
@Injectable({
  providedIn: 'root',
})
export class TrackService {
  private apiUrl = `${environment.apiUrl}/track`;
  private http = inject(HttpClient);

  /**
   * Creates a reactive HTTP resource for fetching the track summary.
   * The resource automatically refetches when min, max, or genre signals change.
   */
  createTrackSummaryResource(deps: {
    min: () => number | undefined;
    max: () => number | undefined;
    genre: () => string | undefined;
  }) {
    return httpResource<ApiResponse<ArtistSummary[]>>(() => {
      const params = buildRangeParams(deps.min(), deps.max(), deps.genre());
      return `${this.apiUrl}/summary?${params.toString()}`;
    });
  }

  /** Clears the cached user tracks and track summary from Redis. */
  invalidateLibraryCache(): Observable<ApiResponse<unknown>> {
    return this.http.post<ApiResponse<unknown>>(`${this.apiUrl}/library-cache/invalidate`, {});
  }

  /**
   * Deletes all saved tracks by the specified artist.
   * @param artistId - The Spotify ID of the artist
   */
  deleteTracksByArtist(artistId: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.apiUrl}/by-artist/${artistId}`);
  }

  /**
   * Deletes all saved tracks whose primary artist has a track count within the specified range.
   * @param min - Minimum track count (inclusive)
   * @param max - Maximum track count (inclusive)
   */
  deleteTracksByRange(min?: number, max?: number): Observable<ApiResponse> {
    const params = buildRangeParams(min, max);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/by-range`, { params });
  }

  /**
   * Retrieves all saved tracks by the specified artist.
   * @param artistId - The Spotify ID of the artist
   */
  getTracksByArtist(artistId: string): Observable<ApiResponse<Track[]>> {
    return this.http.get<ApiResponse<Track[]>>(`${this.apiUrl}/by-artist/${artistId}`);
  }

  /**
   * Deletes a single track from the user's library.
   * @param trackId - The Spotify ID of the track
   */
  deleteTrack(trackId: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.apiUrl}/${trackId}`);
  }
}
