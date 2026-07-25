/**
 * Loading service for managing global loading state.
 * Uses Angular signals for reactive state management.
 */
import { Injectable, signal } from '@angular/core';

/**
 * Service for managing the global loading indicator.
 * Provides a readonly signal for components to observe loading state.
 */
@Injectable({
  providedIn: 'root',
})
export class LoadingService {
  private _loading = signal<boolean>(false);

  /** Readonly signal indicating whether a loading operation is in progress. */
  public readonly loading = this._loading.asReadonly();

  /** Shows the loading indicator. */
  show(): void {
    this._loading.set(true);
  }

  /** Hides the loading indicator. */
  hide(): void {
    this._loading.set(false);
  }
}
