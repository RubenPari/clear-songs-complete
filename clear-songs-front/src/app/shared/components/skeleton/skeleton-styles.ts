/**
 * Skeleton Components Module
 *
 * Provides skeleton loading placeholders for various UI elements.
 * Skeleton screens improve perceived performance and reduce layout shift.
 */

// Skeleton Card Component
export const SkeletonCardStyles = `
  .skeleton-card {
    background: var(--card-bg, #ffffff);
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }

  .skeleton-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;
  }

  .skeleton-circle {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }

  .skeleton-line {
    height: 16px;
    border-radius: 8px;
    background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }

  .skeleton-line.short {
    width: 60%;
  }

  .skeleton-line.medium {
    width: 80%;
  }

  .skeleton-line.long {
    width: 100%;
  }

  .skeleton-body {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }

  /* Dark theme support */
  .dark-theme .skeleton-circle,
  .dark-theme .skeleton-line {
    background: linear-gradient(90deg, #334155 25%, #475569 50%, #334155 75%);
    background-size: 200% 100%;
  }
`;

// Skeleton Stat Card
export const SkeletonStatCard = `
  <div class="skeleton-card stat-skeleton">
    <div class="skeleton-header">
      <div class="skeleton-circle"></div>
      <div style="flex: 1;">
        <div class="skeleton-line short" style="height: 24px; margin-bottom: 8px;"></div>
        <div class="skeleton-line" style="width: 40%;"></div>
      </div>
    </div>
  </div>
`;

// Skeleton Table
export const SkeletonTable = `
  <div class="skeleton-card">
    <div class="skeleton-header">
      <div class="skeleton-circle" style="width: 24px; height: 24px;"></div>
      <div class="skeleton-line short" style="flex: 1; max-width: 200px;"></div>
    </div>
    <div class="skeleton-body">
      <div class="skeleton-row" style="display: flex; gap: 16px; padding: 16px 0; border-bottom: 1px solid rgba(0,0,0,0.05);">
        <div class="skeleton-line" style="flex: 2;"></div>
        <div class="skeleton-line" style="flex: 1;"></div>
        <div class="skeleton-line" style="width: 60px;"></div>
      </div>
      <div class="skeleton-row" style="display: flex; gap: 16px; padding: 16px 0; border-bottom: 1px solid rgba(0,0,0,0.05);">
        <div class="skeleton-line" style="flex: 2;"></div>
        <div class="skeleton-line" style="flex: 1;"></div>
        <div class="skeleton-line" style="width: 60px;"></div>
      </div>
      <div class="skeleton-row" style="display: flex; gap: 16px; padding: 16px 0; border-bottom: 1px solid rgba(0,0,0,0.05);">
        <div class="skeleton-line" style="flex: 2;"></div>
        <div class="skeleton-line" style="flex: 1;"></div>
        <div class="skeleton-line" style="width: 60px;"></div>
      </div>
      <div class="skeleton-row" style="display: flex; gap: 16px; padding: 16px 0;">
        <div class="skeleton-line" style="flex: 2;"></div>
        <div class="skeleton-line" style="flex: 1;"></div>
        <div class="skeleton-line" style="width: 60px;"></div>
      </div>
    </div>
  </div>
`;

// Skeleton Chart
export const SkeletonChart = `
  <div class="skeleton-card">
    <div class="skeleton-header">
      <div class="skeleton-circle" style="width: 24px; height: 24px;"></div>
      <div class="skeleton-line short" style="flex: 1; max-width: 200px;"></div>
    </div>
    <div style="padding: 40px 0; display: flex; align-items: end; justify-content: space-around; gap: 16px; height: 200px;">
      <div class="skeleton-bar" style="width: 40px; height: 60%; border-radius: 4px 4px 0 0; background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite;"></div>
      <div class="skeleton-bar" style="width: 40px; height: 80%; border-radius: 4px 4px 0 0; background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite;"></div>
      <div class="skeleton-bar" style="width: 40px; height: 100%; border-radius: 4px 4px 0 0; background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite;"></div>
      <div class="skeleton-bar" style="width: 40px; height: 70%; border-radius: 4px 4px 0 0; background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite;"></div>
      <div class="skeleton-bar" style="width: 40px; height: 50%; border-radius: 4px 4px 0 0; background: linear-gradient(90deg, #e2e8f0 25%, #f1f5f9 50%, #e2e8f0 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite;"></div>
    </div>
  </div>
`;
