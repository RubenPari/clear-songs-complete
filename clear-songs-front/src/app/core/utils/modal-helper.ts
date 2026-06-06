import { Dialog, DialogConfig } from '@angular/cdk/dialog';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import {
  ConfirmDialogComponent,
  ConfirmDialogData,
} from '../../shared/components/confirm-dialog/confirm-dialog.component';

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  danger?: boolean;
}

/** Shared CDK dialog config: centered overlay with a dark, dismissable backdrop. */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function baseDialogConfig<D>(data: D): DialogConfig<D, any> {
  return {
    data,
    hasBackdrop: true,
    backdropClass: 'app-dialog-backdrop',
    panelClass: 'app-dialog-panel',
  };
}

// Opens confirm dialog.
export function openConfirmDialog(
  dialog: Dialog,
  options: ConfirmDialogOptions,
): Observable<boolean> {
  const ref = dialog.open<boolean, ConfirmDialogData>(
    ConfirmDialogComponent,
    baseDialogConfig<ConfirmDialogData>({
      title: options.title,
      message: options.message,
      confirmText: options.confirmText ?? 'Confirm',
      cancelText: options.cancelText ?? 'Cancel',
      danger: options.danger ?? false,
    }),
  );
  return ref.closed.pipe(map((result) => result === true));
}
