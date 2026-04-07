import { Injectable, inject } from '@angular/core';
import { HttpClient, httpResource } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ArtistSummary, Track } from '../models/artist.model';
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

  getTrackSummaryResource(min?: number, max?: number, genre?: string) {
    return httpResource<ApiResponse<ArtistSummary[]>>(() => {
      const params = buildRangeParams(min, max, genre);
      return `${this.apiUrl}/summary?${params.toString()}`;
    });
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
