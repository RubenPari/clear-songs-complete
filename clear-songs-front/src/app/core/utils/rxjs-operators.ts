/**
 * RxJS Custom Operators
 * 
 * Reusable RxJS operators for common patterns in the application.
 * These operators help reduce code duplication and standardize
 * error handling and loading state management.
 * 
 * @file rxjs-operators.ts
 * @author Clear Songs Development Team
 */
import { Observable, finalize, tap } from 'rxjs';
import { LoadingService } from '../services/loading.service';
import { NotificationService } from '../services/notification.service';

/**
 * Operator that automatically manages loading state
 * 
 * Shows loading indicator when the observable starts and hides it when it completes
 * (either successfully or with an error).
 * 
 * @param loadingService - The loading service instance
 * @returns RxJS operator function
 * 
 * @example
 * this.trackService.getTrackSummary()
 *   .pipe(withLoading(this.loadingService))
 *   .subscribe(...)
 */
export function withLoading<T>(
  loadingService: LoadingService
) {
  return (source: Observable<T>) => {
    return source.pipe(
      tap(() => loadingService.show()),
      finalize(() => loadingService.hide())
    );
  };
}

/**
 * Operator that handles errors with notifications
 * 
 * Automatically shows an error notification when the observable errors.
 * 
 * @param notificationService - The notification service instance
 * @param errorMessage - The error message to display
 * @returns RxJS operator function
 * 
 * @example
 * this.trackService.getTrackSummary()
 *   .pipe(withErrorHandling(this.notificationService, 'Failed to load tracks'))
 *   .subscribe(...)
 */
export function withErrorHandling<T>(
  notificationService: NotificationService,
  errorMessage: string
) {
  return (source: Observable<T>) => {
    return source.pipe(
      tap({
        error: () => notificationService.error(errorMessage)
      })
    );
  };
}

/**
 * Combined operator for loading and error handling
 * 
 * Combines both loading state management and error handling in a single operator.
 * 
 * @param loadingService - The loading service instance
 * @param notificationService - The notification service instance
 * @param errorMessage - The error message to display
 * @returns RxJS operator function
 * 
 * @example
 * this.trackService.getTrackSummary()
 *   .pipe(
 *     withLoadingAndError(
 *       this.loadingService,
 *       this.notificationService,
 *       'Failed to load tracks'
 *     )
 *   )
 *   .subscribe(...)
 */
export function withLoadingAndError<T>(
  loadingService: LoadingService,
  notificationService: NotificationService,
  errorMessage: string
) {
  return (source: Observable<T>) => {
    return source.pipe(
      withLoading(loadingService),
      withErrorHandling(notificationService, errorMessage)
    );
  };
}
