import { TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { PlaylistStore } from './playlist.store';

describe('PlaylistStore', () => {
  let store: PlaylistStore;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideZonelessChangeDetection(), PlaylistStore]
    });
    store = TestBed.inject(PlaylistStore);
  });

  it('should be created', () => {
    expect(store).toBeTruthy();
  });

  it('should have initial empty state', () => {
    expect(store.playlists()).toEqual([]);
    expect(store.loading()).toBeFalse();
    expect(store.error()).toBeNull();
    expect(store.selectedPlaylist()).toBeNull();
  });

  it('should set playlists', () => {
    const playlists = [
      { id: '1', name: 'Playlist 1' },
      { id: '2', name: 'Playlist 2' }
    ];
    store.setPlaylists(playlists);
    expect(store.playlists()).toEqual(playlists);
    expect(store.totalPlaylists()).toBe(2);
  });

  it('should set loading state', () => {
    store.setLoading(true);
    expect(store.loading()).toBeTrue();
  });

  it('should set error state', () => {
    store.setError('Test error');
    expect(store.error()).toBe('Test error');
    expect(store.loading()).toBeFalse();
  });

  it('should select a playlist', () => {
    const playlist = { id: '1', name: 'Playlist 1' };
    store.selectPlaylist(playlist);
    expect(store.selectedPlaylist()).toEqual(playlist);
  });

  it('should remove a playlist', () => {
    const playlists = [
      { id: '1', name: 'Playlist 1' },
      { id: '2', name: 'Playlist 2' }
    ];
    store.setPlaylists(playlists);
    store.removePlaylist('1');
    expect(store.playlists()).toHaveSize(1);
    expect(store.playlists()[0].id).toBe('2');
  });

  it('should reset state', () => {
    store.setPlaylists([{ id: '1', name: 'Playlist 1' }]);
    store.setLoading(true);
    store.setError('Some error');
    store.selectPlaylist({ id: '1', name: 'Playlist 1' });

    store.reset();

    expect(store.playlists()).toEqual([]);
    expect(store.loading()).toBeFalse();
    expect(store.error()).toBeNull();
    expect(store.selectedPlaylist()).toBeNull();
  });
});
