import { Injectable, inject } from '@angular/core';
import { HttpClient, httpResource } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ArtistGenresDebugEntry, ArtistSummary, Track } from '../models/artist.model';
import { ApiResponse } from '../models/api-response.model';
import { buildRangeParams } from '../utils/http-params.helper';
import { TrackStore } from '../stores/track.store';

@Injectable({
  providedIn: 'root',
})
export class TrackService {
  private apiUrl = `${environment.apiUrl}/track`;
  private http = inject(HttpClient);
  private trackStore = inject(TrackStore);

  /**
   * Track summary as a reactive httpResource. The request URL is recomputed whenever
   * the callbacks read new signal values (must be invoked inside the inner function).
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

  /** Full library: every primary artist with raw Spotify `genres` (for mapping/debug). */
  getArtistGenresDebug(): Observable<ApiResponse<ArtistGenresDebugEntry[]>> {
    return this.http.get<ApiResponse<ArtistGenresDebugEntry[]>>(
      `${this.apiUrl}/debug/artist-genres`
    );
  }

  /** Clears Redis user-tracks + track-summary caches; call before reloading the dashboard summary. */
  invalidateLibraryCache(): Observable<ApiResponse<unknown>> {
    return this.http.post<ApiResponse<unknown>>(`${this.apiUrl}/library-cache/invalidate`, {});
  }

  deleteTracksByArtist(artistId: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.apiUrl}/by-artist/${artistId}`).pipe(
      tap(() => {
        this.trackStore.removeArtist(artistId);
      })
    );
  }

  deleteTracksByRange(min?: number, max?: number): Observable<ApiResponse> {
    const params = buildRangeParams(min, max);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/by-range`, { params });
  }

  getTracksByArtist(artistId: string): Observable<ApiResponse<Track[]>> {
    return this.http.get<ApiResponse<Track[]>>(`${this.apiUrl}/by-artist/${artistId}`);
  }

  deleteTrack(trackId: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.apiUrl}/${trackId}`);
  }
}
