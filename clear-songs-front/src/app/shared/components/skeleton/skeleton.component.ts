/**
 * Skeleton Loading Component
 *
 * Premium dark-themed skeleton placeholders for loading states.
 * Improves perceived performance and reduces layout shift.
 */

import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-skeleton',
  standalone: true,
  imports: [CommonModule],
  template: `
    @switch (type) {
      @case ('stat') {
        <div class="skeleton-stat">
          <div class="skeleton-icon"></div>
          <div class="skeleton-content">
            <div class="skeleton-value"></div>
            <div class="skeleton-label"></div>
          </div>
        </div>
      }
      @case ('table') {
        <div class="skeleton-table">
          <div class="skeleton-header-row">
            <div class="skeleton-cell" style="flex: 2;"></div>
            <div class="skeleton-cell" style="flex: 1;"></div>
            <div class="skeleton-cell" style="width: 80px;"></div>
          </div>
          @for (i of [1,2,3,4,5]; track i) {
            <div class="skeleton-data-row">
              <div class="skeleton-cell" style="flex: 2;">
                <div class="skeleton-avatar-text">
                  <div class="skeleton-avatar"></div>
                  <div class="skeleton-text"></div>
                </div>
              </div>
              <div class="skeleton-cell" style="flex: 1;">
                <div class="skeleton-badge"></div>
              </div>
              <div class="skeleton-cell" style="width: 80px;">
                <div class="skeleton-btn"></div>
              </div>
            </div>
          }
        </div>
      }
      @case ('chart') {
        <div class="skeleton-chart">
          <div class="skeleton-bars">
            @for (height of ['60%', '80%', '100%', '70%', '50%']; track $index) {
              <div class="skeleton-bar" [style.height]="height"></div>
            }
          </div>
        </div>
      }
      @default {
        <div class="skeleton-card">
          <div class="skeleton-line short"></div>
          <div class="skeleton-line"></div>
          <div class="skeleton-line medium"></div>
        </div>
      }
    }
  `,
  styles: [`
    :host {
      display: block;
    }

    /* Shimmer Animation */
    @keyframes shimmer {
      0% {
        background-position: 200% 0;
      }
      100% {
        background-position: -200% 0;
      }
    }

    @keyframes pulse {
      0%, 100% {
        opacity: 1;
      }
      50% {
        opacity: 0.6;
      }
    }

    /* Base Skeleton Mixin */
    .skeleton-base {
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.03) 0%,
        rgba(255, 255, 255, 0.08) 50%,
        rgba(255, 255, 255, 0.03) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    /* Skeleton Card */
    .skeleton-card {
      background: var(--bg-glass, rgba(18, 18, 20, 0.7));
      backdrop-filter: blur(20px);
      border: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
      border-radius: var(--radius-xl, 24px);
      padding: 24px;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    /* Skeleton Stat */
    .skeleton-stat {
      background: var(--bg-glass, rgba(18, 18, 20, 0.7));
      backdrop-filter: blur(20px);
      border: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
      border-radius: var(--radius-xl, 24px);
      padding: 24px;
      display: flex;
      align-items: center;
      gap: 20px;
    }

    .skeleton-icon {
      width: 72px;
      height: 72px;
      border-radius: var(--radius-lg, 16px);
      background: linear-gradient(
        90deg,
        rgba(29, 185, 84, 0.1) 0%,
        rgba(29, 185, 84, 0.2) 50%,
        rgba(29, 185, 84, 0.1) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-content {
      flex: 1;
    }

    .skeleton-value {
      width: 60%;
      height: 36px;
      border-radius: var(--radius-md, 12px);
      margin-bottom: 8px;
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.05) 0%,
        rgba(255, 255, 255, 0.1) 50%,
        rgba(255, 255, 255, 0.05) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-label {
      width: 40%;
      height: 14px;
      border-radius: var(--radius-sm, 8px);
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.03) 0%,
        rgba(255, 255, 255, 0.08) 50%,
        rgba(255, 255, 255, 0.03) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
      animation-delay: 0.2s;
    }

    /* Skeleton Table */
    .skeleton-table {
      background: var(--bg-glass, rgba(18, 18, 20, 0.7));
      backdrop-filter: blur(20px);
      border: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
      border-radius: var(--radius-xl, 24px);
      overflow: hidden;
    }

    .skeleton-header-row {
      display: flex;
      gap: 16px;
      padding: 16px 24px;
      border-bottom: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
      background: var(--surface-subtle, rgba(255, 255, 255, 0.03));
    }

    .skeleton-data-row {
      display: flex;
      gap: 16px;
      padding: 12px 24px;
      border-bottom: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
    }

    .skeleton-cell {
      display: flex;
      align-items: center;
    }

    .skeleton-avatar-text {
      display: flex;
      align-items: center;
      gap: 14px;
      width: 100%;
    }

    .skeleton-avatar {
      width: 44px;
      height: 44px;
      border-radius: 50%;
      flex-shrink: 0;
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.05) 0%,
        rgba(255, 255, 255, 0.1) 50%,
        rgba(255, 255, 255, 0.05) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-text {
      flex: 1;
      height: 16px;
      border-radius: var(--radius-sm, 8px);
      max-width: 150px;
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.03) 0%,
        rgba(255, 255, 255, 0.08) 50%,
        rgba(255, 255, 255, 0.03) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-badge {
      width: 50px;
      height: 26px;
      border-radius: var(--radius-full, 9999px);
      background: linear-gradient(
        90deg,
        rgba(29, 185, 84, 0.1) 0%,
        rgba(29, 185, 84, 0.2) 50%,
        rgba(29, 185, 84, 0.1) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-btn {
      width: 40px;
      height: 40px;
      border-radius: var(--radius-md, 12px);
      background: linear-gradient(
        90deg,
        rgba(255, 71, 87, 0.1) 0%,
        rgba(255, 71, 87, 0.2) 50%,
        rgba(255, 71, 87, 0.1) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    /* Skeleton Chart */
    .skeleton-chart {
      background: var(--bg-glass, rgba(18, 18, 20, 0.7));
      backdrop-filter: blur(20px);
      border: 1px solid var(--border-subtle, rgba(255, 255, 255, 0.06));
      border-radius: var(--radius-xl, 24px);
      padding: 24px;
      height: 300px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .skeleton-bars {
      display: flex;
      align-items: end;
      justify-content: space-around;
      gap: 24px;
      height: 200px;
      width: 100%;
      max-width: 400px;
    }

    .skeleton-bar {
      flex: 1;
      max-width: 60px;
      border-radius: var(--radius-md, 12px) var(--radius-md, 12px) 0 0;
      background: linear-gradient(
        180deg,
        rgba(29, 185, 84, 0.3) 0%,
        rgba(0, 212, 255, 0.2) 50%,
        rgba(29, 185, 84, 0.1) 100%
      );
      background-size: 100% 200%;
      animation: pulse 2s ease-in-out infinite;
      animation-delay: calc(var(--index, 0) * 0.1s);
    }

    .skeleton-bars .skeleton-bar:nth-child(1) { --index: 1; animation-delay: 0s; }
    .skeleton-bars .skeleton-bar:nth-child(2) { --index: 2; animation-delay: 0.15s; }
    .skeleton-bars .skeleton-bar:nth-child(3) { --index: 3; animation-delay: 0.3s; }
    .skeleton-bars .skeleton-bar:nth-child(4) { --index: 4; animation-delay: 0.45s; }
    .skeleton-bars .skeleton-bar:nth-child(5) { --index: 5; animation-delay: 0.6s; }

    /* Common Skeleton Lines */
    .skeleton-line {
      height: 16px;
      border-radius: var(--radius-md, 12px);
      background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.03) 0%,
        rgba(255, 255, 255, 0.08) 50%,
        rgba(255, 255, 255, 0.03) 100%
      );
      background-size: 200% 100%;
      animation: shimmer 2s ease-in-out infinite;
    }

    .skeleton-line.short {
      width: 40%;
    }

    .skeleton-line.medium {
      width: 70%;
    }
  `]
})
export class SkeletonComponent {
  @Input() type: 'card' | 'stat' | 'table' | 'chart' = 'card';
}