import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

export interface DialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
}

@Component({
  selector: 'app-confirm-dialog',
  template: `
    <div class="modal-header">
      <h5 class="modal-title">{{ title }}</h5>
      <button type="button" class="btn-close" aria-label="Close" (click)="onCancel()"></button>
    </div>
    <div class="modal-body">
      <p>{{ message }}</p>
    </div>
    <div class="modal-footer">
      <button type="button" class="btn btn-secondary" (click)="onCancel()">
        {{ cancelText || 'Cancel' }}
      </button>
      <button type="button" class="btn btn-danger" (click)="onConfirm()">
        {{ confirmText || 'Confirm' }}
      </button>
    </div>
  `,
  styles: [
    `
      .modal-body p {
        margin: 0;
        font-size: 16px;
        white-space: pre-line;
      }
    `,
  ],
  standalone: true,
  imports: [CommonModule]
})
export class ConfirmDialogComponent {
  title = '';
  message = '';
  confirmText?: string;
  cancelText?: string;

  public activeModal = inject(NgbActiveModal);

  // On cancel.
  onCancel(): void {
    this.activeModal.close(false);
  }

  // On confirm.
  onConfirm(): void {
    this.activeModal.close(true);
  }
}
