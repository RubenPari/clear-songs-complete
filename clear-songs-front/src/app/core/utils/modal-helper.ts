/**
 * Modal Helper Utilities
 * 
 * Utility functions for common modal operations, reducing code duplication
 * across components that use confirmation dialogs.
 * 
 * @file modal-helper.ts
 * @author Clear Songs Development Team
 */
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Observable, from } from 'rxjs';
import { ConfirmDialogComponent } from '../../shared/components/confirm-dialog/confirm-dialog.component';

/**
 * Configuration options for confirmation dialogs
 */
export interface ConfirmDialogOptions {
  /** Dialog title */
  title: string;
  /** Dialog message */
  message: string;
  /** Text for confirm button (default: 'Confirm') */
  confirmText?: string;
  /** Text for cancel button (default: 'Cancel') */
  cancelText?: string;
  /** Modal size (default: 'md') */
  size?: 'sm' | 'md' | 'lg' | 'xl';
  /** Whether modal is centered (default: true) */
  centered?: boolean;
}

/**
 * Opens a confirmation dialog and returns an Observable
 * 
 * The Observable emits:
 * - true if user confirms
 * - false if user cancels or dismisses
 * 
 * @param modalService - The NgbModal service instance
 * @param options - Configuration options for the dialog
 * @returns Observable that emits boolean indicating confirmation
 * 
 * @example
 * openConfirmDialog(this.modalService, {
 *   title: 'Delete Tracks',
 *   message: 'Are you sure?',
 *   confirmText: 'Delete',
 *   cancelText: 'Cancel'
 * }).subscribe(confirmed => {
 *   if (confirmed) {
 *     // User confirmed
 *   }
 * });
 */
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
  
  return from(modalRef.result);
}
