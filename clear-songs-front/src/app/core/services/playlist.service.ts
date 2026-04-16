import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams, httpResource } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../models/api-response.model';
import { UserPlaylist } from '../models/artist.model';
import { PlaylistStore } from '../stores/playlist.store';

@Injectable({
  providedIn: 'root',
})
export class PlaylistService {
  private apiUrl = `${environment.apiUrl}/playlist`;
  private http = inject(HttpClient);
  private playlistStore = inject(PlaylistStore);

  // Fetches user playlists resource.
  getUserPlaylistsResource() {
    return httpResource<ApiResponse<UserPlaylist[]>>(() => `${this.apiUrl}/list`);
  }

  // Deletes all playlist tracks.
  deleteAllPlaylistTracks(playlistId: string): Observable<ApiResponse> {
    const params = new HttpParams().set('id', playlistId);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/delete-tracks`, { params });
  }

  // Deletes all playlist and user tracks.
  deleteAllPlaylistAndUserTracks(playlistId: string): Observable<ApiResponse> {
    const params = new HttpParams().set('id', playlistId);
    return this.http.delete<ApiResponse>(`${this.apiUrl}/delete-tracks-and-library`, { params });
  }
}
