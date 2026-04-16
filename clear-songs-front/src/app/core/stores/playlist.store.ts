import { Injectable, signal, computed } from '@angular/core';
import { UserPlaylist } from '../models/artist.model';

@Injectable({
  providedIn: 'root',
})
export class PlaylistStore {
  private _playlists = signal<UserPlaylist[]>([]);
  private _loading = signal<boolean>(false);
  private _error = signal<string | null>(null);
  private _selectedPlaylist = signal<UserPlaylist | null>(null);

  public readonly playlists = this._playlists.asReadonly();
  public readonly loading = this._loading.asReadonly();
  public readonly error = this._error.asReadonly();
  public readonly selectedPlaylist = this._selectedPlaylist.asReadonly();

  public readonly totalPlaylists = computed(() => this._playlists().length);

  // Sets playlists.
  setPlaylists(playlists: UserPlaylist[]): void {
    this._playlists.set(playlists);
    this._error.set(null);
  }

  // Sets loading.
  setLoading(loading: boolean): void {
    this._loading.set(loading);
  }

  // Sets error.
  setError(error: string | null): void {
    this._error.set(error);
    if (error) {
      this._loading.set(false);
    }
  }

  // Selects playlist.
  selectPlaylist(playlist: UserPlaylist | null): void {
    this._selectedPlaylist.set(playlist);
  }

  // Removes playlist.
  removePlaylist(playlistId: string): void {
    this._playlists.update(playlists => playlists.filter(p => p.id !== playlistId));
  }

  // Resets.
  reset(): void {
    this._playlists.set([]);
    this._loading.set(false);
    this._error.set(null);
    this._selectedPlaylist.set(null);
  }
}
