/**
 * Skeleton loading placeholder components for reducing perceived latency.
 * Uses Tailwind animations and design-system tokens.
 */
import { Component, Input } from '@angular/core';

/** Skeleton placeholder for stat cards. */
@Component({
  selector: 'app-skeleton-stat',
  standalone: true,
  template: `
    <div class="flex items-center gap-4 rounded-xl border bg-card p-5 shadow-sm">
      <div class="size-12 shrink-0 animate-pulse rounded-lg bg-muted"></div>
      <div class="flex-1 space-y-2">
        <div class="h-7 w-3/5 animate-pulse rounded-md bg-muted"></div>
        <div class="h-4 w-2/5 animate-pulse rounded bg-muted"></div>
      </div>
    </div>
  `,
})
export class SkeletonStatComponent {}

/** Skeleton placeholder for data tables with configurable row count. */
@Component({
  selector: 'app-skeleton-table',
  standalone: true,
  template: `
    <div class="rounded-xl border bg-card p-6 shadow-sm">
      <div class="mb-5 h-6 w-48 animate-pulse rounded bg-muted"></div>
      <div class="flex flex-col">
        @for (row of rows; track $index) {
          <div class="flex items-center gap-4 border-b border-border/60 py-3 last:border-b-0">
            <div class="size-10 shrink-0 animate-pulse rounded-full bg-muted"></div>
            <div class="h-4 flex-[2] animate-pulse rounded bg-muted"></div>
            <div class="h-4 flex-1 animate-pulse rounded bg-muted"></div>
            <div class="h-8 w-16 animate-pulse rounded-md bg-muted"></div>
          </div>
        }
      </div>
    </div>
  `,
})
export class SkeletonTableComponent {
  @Input() rowCount = 5;
  get rows() {
    return new Array(this.rowCount);
  }
}

/** Skeleton placeholder for bar charts with animated bars. */
@Component({
  selector: 'app-skeleton-chart',
  standalone: true,
  template: `
    <div class="rounded-xl border bg-card p-6 shadow-sm">
      <div class="mb-5 h-6 w-48 animate-pulse rounded bg-muted"></div>
      <div class="flex h-64 items-end justify-around gap-6 py-2">
        @for (h of barHeights; track $index) {
          <div
            class="w-full max-w-14 animate-pulse rounded-t-md bg-muted"
            [style.height.%]="h"
          ></div>
        }
      </div>
    </div>
  `,
})
export class SkeletonChartComponent {
  @Input() barCount = 5;
  barHeights = [60, 80, 100, 70, 50];
}
