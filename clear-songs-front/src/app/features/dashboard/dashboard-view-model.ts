import { ArtistSummary } from '../../core/models/artist.model';

export type ArtistSortColumn = 'name' | 'count';
export type SortDirection = 'asc' | 'desc';

export function filterAndSortArtists(
  artists: ArtistSummary[],
  search: string,
  column: ArtistSortColumn,
  direction: SortDirection,
): ArtistSummary[] {
  const normalizedSearch = search.trim().toLowerCase();
  const filtered = normalizedSearch
    ? artists.filter((artist) => artist.name.toLowerCase().includes(normalizedSearch))
    : artists;

  return [...filtered].sort((left, right) => {
    const comparison = column === 'name'
      ? left.name.localeCompare(right.name)
      : left.count - right.count;
    return direction === 'asc' ? comparison : -comparison;
  });
}

export function paginate<T>(items: T[], page: number, pageSize: number): T[] {
  return items.slice((page - 1) * pageSize, page * pageSize);
}

export function visiblePageNumbers(total: number, current: number): (number | string)[] {
  if (total <= 4) {
    return Array.from({ length: total }, (_, index) => index + 1);
  }
  if (current <= 3) {
    return [1, 2, 3, '...', total];
  }

  const start = Math.min(current, Math.max(1, total - 2));
  const pages: (number | string)[] = [start, start + 1, start + 2].filter((page) => page <= total);
  return pages.at(-1) === total ? pages : [...pages, '...', total];
}

export function normalizeRange(min: number, max: number, cap: number): [number, number] {
  const normalizedMin = Math.max(0, Math.min(min, cap));
  const normalizedMax = Math.max(0, Math.min(max, cap));
  return normalizedMin <= normalizedMax
    ? [normalizedMin, normalizedMax]
    : [normalizedMax, normalizedMin];
}
