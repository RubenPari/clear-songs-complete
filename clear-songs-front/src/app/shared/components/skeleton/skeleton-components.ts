/**
 * Skeleton Components
 *
 * Reusable skeleton loading placeholders for improved UX during data loading.
 * Provides visual feedback that content is being loaded.
 */
import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-skeleton-stat',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="skeleton-stat">
      <div class="skeleton-icon"></div>
      <div class="skeleton-content">
        <div class="skeleton-value"></div>
        <div class="skeleton-label"></div>
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .skeleton-stat {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 24px;
      background: var(--card-bg, #ffffff);
      border-radius: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    .skeleton-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    .skeleton-content {
      flex: 1;
    }

    .skeleton-value {
      height: 32px;
      width: 60%;
      border-radius: 8px;
      margin-bottom: 8px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    .skeleton-label {
      height: 16px;
      width: 40%;
      border-radius: 4px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    @keyframes shimmer {
      0% { background-position: 200% 0; }
      100% { background-position: -200% 0; }
    }

    /* Dark theme */
    :host-context(.dark-theme) .skeleton-stat {
      background: #1e293b;
    }

    :host-context(.dark-theme) .skeleton-icon,
    :host-context(.dark-theme) .skeleton-value,
    :host-context(.dark-theme) .skeleton-label {
      background: linear-gradient(90deg, #334155 25%, #475569 50%, #334155 75%);
      background-size: 200% 100%;
    }
  `]
})
export class SkeletonStatComponent {}

@Component({
  selector: 'app-skeleton-table',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="skeleton-table-container">
      <div class="skeleton-header">
        <div class="skeleton-title"></div>
      </div>
      <div class="skeleton-rows">
        @for (row of rows; track $index) {
          <div class="skeleton-row">
            <div class="skeleton-cell cell-large"></div>
            <div class="skeleton-cell cell-medium"></div>
            <div class="skeleton-cell cell-small"></div>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .skeleton-table-container {
      background: var(--card-bg, #ffffff);
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    .skeleton-header {
      margin-bottom: 20px;
    }

    .skeleton-title {
      height: 24px;
      width: 200px;
      border-radius: 4px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    .skeleton-rows {
      display: flex;
      flex-direction: column;
      gap: 16px;
    }

    .skeleton-row {
      display: flex;
      gap: 16px;
      padding: 12px 0;
      border-bottom: 1px solid rgba(0, 0, 0, 0.05);
    }

    .skeleton-row:last-child {
      border-bottom: none;
    }

    .skeleton-cell {
      height: 20px;
      border-radius: 4px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    .cell-large { flex: 2; }
    .cell-medium { flex: 1; }
    .cell-small { width: 60px; }

    @keyframes shimmer {
      0% { background-position: 200% 0; }
      100% { background-position: -200% 0; }
    }

    /* Dark theme */
    :host-context(.dark-theme) .skeleton-table-container {
      background: #1e293b;
    }

    :host-context(.dark-theme) .skeleton-title,
    :host-context(.dark-theme) .skeleton-cell {
      background: linear-gradient(90deg, #334155 25%, #475569 50%, #334155 75%);
      background-size: 200% 100%;
    }

    :host-context(.dark-theme) .skeleton-row {
      border-bottom-color: rgba(255, 255, 255, 0.05);
    }
  `]
})
export class SkeletonTableComponent {
  @Input() rowCount = 5;
  get rows() { return new Array(this.rowCount); }
}

@Component({
  selector: 'app-skeleton-chart',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="skeleton-chart-container">
      <div class="skeleton-header">
        <div class="skeleton-title"></div>
      </div>
      <div class="skeleton-chart">
        @for (bar of bars; track $index; let i = $index) {
          <div
            class="skeleton-bar"
            [style.height.%]="barHeights[i]"
          ></div>
        }
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
    }

    .skeleton-chart-container {
      background: var(--card-bg, #ffffff);
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    .skeleton-header {
      margin-bottom: 20px;
    }

    .skeleton-title {
      height: 24px;
      width: 200px;
      border-radius: 4px;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    .skeleton-chart {
      display: flex;
      align-items: flex-end;
      justify-content: space-around;
      gap: 24px;
      height: 280px;
      padding: 20px 0;
    }

    .skeleton-bar {
      flex: 1;
      max-width: 60px;
      border-radius: 4px 4px 0 0;
      background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
    }

    @keyframes shimmer {
      0% { background-position: 200% 0; }
      100% { background-position: -200% 0; }
    }

    /* Dark theme */
    :host-context(.dark-theme) .skeleton-chart-container {
      background: #1e293b;
    }

    :host-context(.dark-theme) .skeleton-title,
    :host-context(.dark-theme) .skeleton-bar {
      background: linear-gradient(90deg, #334155 25%, #475569 50%, #334155 75%);
      background-size: 200% 100%;
    }
  `]
})
export class SkeletonChartComponent {
  @Input() barCount = 5;
  barHeights = [60, 80, 100, 70, 50];
  get bars() { return new Array(this.barCount); }
}
