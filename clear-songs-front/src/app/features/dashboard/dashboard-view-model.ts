/**
 * Pure functions for transforming and filtering artist data in the dashboard.
 * Contains sorting, pagination, and range normalization logic.
 */
import { ArtistSummary } from '../../core/models/artist.model';

/** Column to sort artists by. */
export type ArtistSortColumn = 'name' | 'count';

/** Sort direction. */
export type SortDirection = 'asc' | 'desc';

/**
 * Filters artists by search term and sorts by the specified column.
 * @param artists - The artists to filter and sort
 * @param search - Search term to filter by name (case-insensitive)
 * @param column - Column to sort by (name or count)
 * @param direction - Sort direction (asc or desc)
 */
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

/**
 * Paginates an array of items.
 * @param items - The items to paginate
 * @param page - The current page number (1-indexed)
 * @param pageSize - The number of items per page
 */
export function paginate<T>(items: T[], page: number, pageSize: number): T[] {
  return items.slice((page - 1) * pageSize, page * pageSize);
}

/**
 * Generates an array of page numbers for pagination UI with ellipsis.
 * Returns a window of pages around the current page, with '...' for gaps.
 * @param total - Total number of pages
 * @param current - Current page number
 */
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

/**
 * Normalizes a min/max range, ensuring min <= max and both are within [0, cap].
 * Swaps values if min > max.
 * @param min - Minimum value
 * @param max - Maximum value
 * @param cap - Maximum allowed value
 */
export function normalizeRange(min: number, max: number, cap: number): [number, number] {
  const normalizedMin = Math.max(0, Math.min(min, cap));
  const normalizedMax = Math.max(0, Math.min(max, cap));
  return normalizedMin <= normalizedMax
    ? [normalizedMin, normalizedMax]
    : [normalizedMax, normalizedMin];
}
