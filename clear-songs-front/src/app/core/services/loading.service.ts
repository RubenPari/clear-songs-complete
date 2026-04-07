/**
 * Loading Service
 * 
 * A global service that manages the loading state across the entire application.
 * This service uses Angular Signals (introduced in Angular 16+) for reactive state
 * management, providing better performance and simpler API compared to RxJS Observables.
 * 
 * The service provides:
 * - A readonly signal that components can subscribe to for loading state changes
 * - Methods to show and hide the loading indicator
 * - Automatic cleanup when components are destroyed (thanks to signals)
 * 
 * Usage:
 * Components can inject this service and use the `loading` signal to display
 * loading indicators. The signal automatically triggers change detection when
 * the loading state changes.
 * 
 * Example:
 * ```typescript
 * constructor(public loadingService: LoadingService) {}
 * 
 * // In template:
 * @if (loadingService.loading()) {
 *   <div class="spinner-border"></div>
 * }
 * ```
 * 
 * @service
 * @providedIn root
 * @author Clear Songs Development Team
 */
import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LoadingService {
  /**
   * Private signal that holds the loading state
   * 
   * This is a writable signal that can only be modified internally
   * through the show() and hide() methods. The public API exposes
   * a readonly version to prevent external modifications.
   */
  private _loading = signal<boolean>(false);
  
  /**
   * Public readonly signal for loading state
   * 
   * Components can read this signal to reactively display loading indicators.
   * The signal is readonly to prevent external components from directly
   * modifying the loading state, ensuring controlled state management.
   */
  public readonly loading = this._loading.asReadonly();

  /**
   * Shows the loading indicator
   * 
   * Sets the loading state to true, which will trigger change detection
   * in all components that are reading the loading signal.
   * 
   * This method should be called at the start of async operations
   * (HTTP requests, file operations, etc.) to provide user feedback.
   */
  show(): void {
    this._loading.set(true);
  }

  /**
   * Hides the loading indicator
   * 
   * Sets the loading state to false, indicating that async operations
   * have completed. This should be called when operations finish,
   * either successfully or with an error.
   * 
   * Note: Consider using RxJS finalize() operator in HTTP requests
   * to ensure this is called even if an error occurs.
   */
  hide(): void {
    this._loading.set(false);
  }
}
