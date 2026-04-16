import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Observable, from, of } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { ConfirmDialogComponent } from '../../shared/components/confirm-dialog/confirm-dialog.component';

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  centered?: boolean;
}

export function modalResult$<T>(modalRef: NgbModalRef, dismissedValue: T): Observable<T> {
  return from(modalRef.result as Promise<T>).pipe(catchError(() => of(dismissedValue)));
}

// Opens confirm dialog.
export function openConfirmDialog(
  modalService: NgbModal,
  options: ConfirmDialogOptions
): Observable<boolean> {
  const modalRef = modalService.open(ConfirmDialogComponent, {
    size: options.size || 'md',
    centered: options.centered !== false,
  });

  modalRef.componentInstance.title = options.title;
  modalRef.componentInstance.message = options.message;
  modalRef.componentInstance.confirmText = options.confirmText || 'Confirm';
  modalRef.componentInstance.cancelText = options.cancelText || 'Cancel';

  return modalResult$<boolean>(modalRef, false);
}
