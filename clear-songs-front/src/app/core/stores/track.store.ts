import { Injectable, signal, computed } from '@angular/core';
import { ArtistSummary } from '../models/artist.model';

@Injectable({
  providedIn: 'root',
})
export class TrackStore {
  private _artists = signal<ArtistSummary[]>([]);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);

  public readonly artists = this._artists.asReadonly();
  public readonly loading = this._loading.asReadonly();
  public readonly error = this._error.asReadonly();

  public readonly totalTracks = computed(() =>
    this._artists().reduce((sum, artist) => sum + artist.count, 0)
  );

  public readonly totalArtists = computed(() => this._artists().length);

  public readonly topArtists = computed(() =>
    [...this._artists()]
      .sort((a, b) => b.count - a.count)
      .slice(0, 5)
  );

  setArtists(artists: ArtistSummary[]): void {
    this._artists.set(artists);
    this._error.set(null);
  }

  setLoading(loading: boolean): void {
    this._loading.set(loading);
  }

  setError(error: string | null): void {
    this._error.set(error);
    if (error) {
      this._loading.set(false);
    }
  }

  removeArtist(artistId: string): void {
    this._artists.update(artists => artists.filter(a => a.id !== artistId));
  }

  updateArtist(artistId: string, updates: Partial<ArtistSummary>): void {
    this._artists.update(artists =>
      artists.map(artist => (artist.id === artistId ? { ...artist, ...updates } : artist))
    );
  }

  reset(): void {
    this._artists.set([]);
    this._loading.set(false);
    this._error.set(null);
  }
}
