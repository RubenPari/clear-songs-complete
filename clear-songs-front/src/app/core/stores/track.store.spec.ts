import { TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { TrackStore } from './track.store';

describe('TrackStore', () => {
  let store: TrackStore;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideZonelessChangeDetection(), TrackStore]
    });
    store = TestBed.inject(TrackStore);
  });

  it('should be created', () => {
    expect(store).toBeTruthy();
  });

  it('should have initial empty state', () => {
    expect(store.artists()).toEqual([]);
    expect(store.loading()).toBeFalse();
    expect(store.error()).toBeNull();
  });

  it('should compute totalTracks', () => {
    store.setArtists([
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 5 },
      { id: '3', name: 'Artist 3', count: 15 }
    ]);
    expect(store.totalTracks()).toBe(30);
  });

  it('should compute totalArtists', () => {
    store.setArtists([
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 5 }
    ]);
    expect(store.totalArtists()).toBe(2);
  });

  it('should compute topArtists', () => {
    store.setArtists([
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 50 },
      { id: '3', name: 'Artist 3', count: 30 },
      { id: '4', name: 'Artist 4', count: 20 },
      { id: '5', name: 'Artist 5', count: 40 },
      { id: '6', name: 'Artist 6', count: 60 }
    ]);
    const top = store.topArtists();
    expect(top).toHaveSize(5);
    expect(top[0].count).toBe(60);
    expect(top[1].count).toBe(50);
  });

  it('should remove artist', () => {
    store.setArtists([
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 5 }
    ]);
    store.removeArtist('1');
    expect(store.artists()).toHaveSize(1);
    expect(store.artists()[0].id).toBe('2');
  });

  it('should update artist', () => {
    store.setArtists([
      { id: '1', name: 'Artist 1', count: 10 }
    ]);
    store.updateArtist('1', { count: 20, name: 'Updated Artist' });
    expect(store.artists()[0].count).toBe(20);
    expect(store.artists()[0].name).toBe('Updated Artist');
  });

  it('should reset state', () => {
    store.setArtists([{ id: '1', name: 'Artist 1', count: 10 }]);
    store.setLoading(true);
    store.setError('Some error');
    store.reset();
    expect(store.artists()).toEqual([]);
    expect(store.loading()).toBeFalse();
    expect(store.error()).toBeNull();
  });
});
