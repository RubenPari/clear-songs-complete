/**
 * Playlist service for managing user playlists via the backend API.
 * Provides reactive resource fetching and deletion operations.
 */
import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams, httpResource } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../models/api-response.model';
import { UserPlaylist } from '../models/artist.model';

/**
 * Service for managing user playlists.
 * Provides reactive playlist fetching and deletion operations.
 */
@Injectable({
  providedIn: 'root',
})
export class PlaylistService {
  private apiUrl = `${environment.apiUrl}/playlist`;
  private http = inject(HttpClient);

  /**
   * Creates a reactive HTTP resource for fetching the user's playlists.
   * The resource automatically refetches when dependencies change.
   */
  getUserPlaylistsResource() {
    return httpResource<ApiResponse<UserPlaylist[]>>(() => `${this.apiUrl}/list`);
  }

  /**
   * Deletes all tracks from the specified playlist.
   * @param playlistId - The ID of the playlist to clear
   */
  deleteAllPlaylistTracks(playlistId: string): Observable<ApiResponse> {
    const params = new HttpParams().set('id', playlistId);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/delete-tracks`, { params });
  }

  /**
   * Deletes all tracks from the specified playlist and removes them from the user's library.
   * @param playlistId - The ID of the playlist to clear
   */
  deleteAllPlaylistAndUserTracks(playlistId: string): Observable<ApiResponse> {
    const params = new HttpParams().set('id', playlistId);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/delete-tracks-and-library`, { params });
  }
}
