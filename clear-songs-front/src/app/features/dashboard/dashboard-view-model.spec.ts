import { ArtistSummary } from '../../core/models/artist.model';
import { filterAndSortArtists, normalizeRange, paginate, visiblePageNumbers } from './dashboard-view-model';

describe('dashboard view model', () => {
  const artists: ArtistSummary[] = [
    { id: '1', name: 'Beta', count: 4 },
    { id: '2', name: 'Alpha', count: 12 },
    { id: '3', name: 'Alpine', count: 8 },
  ];

  it('filters and sorts without mutating the source', () => {
    expect(filterAndSortArtists(artists, 'alp', 'count', 'desc').map(({ name }) => name)).toEqual(['Alpha', 'Alpine']);
    expect(artists.map(({ name }) => name)).toEqual(['Beta', 'Alpha', 'Alpine']);
  });

  it('paginates and builds compact page numbers', () => {
    expect(paginate(artists, 2, 2)).toEqual([artists[2]]);
    expect(visiblePageNumbers(7, 5)).toEqual([5, 6, 7]);
    expect(visiblePageNumbers(7, 1)).toEqual([1, 2, 3, '...', 7]);
  });

  it('clamps and orders range inputs', () => {
    expect(normalizeRange(20, -1, 10)).toEqual([0, 10]);
  });
});
