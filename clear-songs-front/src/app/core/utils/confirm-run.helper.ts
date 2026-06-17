import { Dialog } from '@angular/cdk/dialog';
import { Observable, of } from 'rxjs';
import { catchError, filter, finalize, switchMap } from 'rxjs/operators';

import { LoadingService } from '../services/loading.service';
import { NotificationService } from '../services/notification.service';
import { ConfirmDialogOptions, openConfirmDialog } from './modal-helper';

export interface ConfirmAndRunDeps {
  dialog: Dialog;
  loadingService: LoadingService;
  notificationService: NotificationService;
}

/**
 * Opens a confirm dialog and, if confirmed, runs an async operation while
 * managing the global loading state.
 *
 * @returns Observable that emits the operation result only when confirmed,
 *          otherwise completes without emitting.
 */
export function confirmAndRun<T>(
  options: ConfirmDialogOptions,
  run: () => Observable<T>,
  deps: ConfirmAndRunDeps,
): Observable<T> {
  return openConfirmDialog(deps.dialog, options).pipe(
    filter((confirmed) => confirmed),
    switchMap(() => {
      deps.loadingService.show();
      return run().pipe(finalize(() => deps.loadingService.hide()));
    }),
  );
}

/**
 * Variant of confirmAndRun that also shows a success or error notification.
 */
export function confirmAndRunWithNotify<T>(
  options: ConfirmDialogOptions,
  run: () => Observable<T>,
  deps: ConfirmAndRunDeps & {
    successMessage: string;
    errorMessage: string;
    onSuccess?: (value: T) => void;
    onError?: (error: unknown) => void;
  },
): Observable<T | undefined> {
  return confirmAndRun(options, run, deps).pipe(
    switchMap((value) => {
      deps.notificationService.success(deps.successMessage);
      deps.onSuccess?.(value);
      return of(value);
    }),
    catchError((error) => {
      deps.notificationService.error(deps.errorMessage);
      deps.onError?.(error);
      return of(undefined);
    }),
  );
}
