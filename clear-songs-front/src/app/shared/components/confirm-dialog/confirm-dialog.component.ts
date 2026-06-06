import { Component, inject } from '@angular/core';
import { DIALOG_DATA, DialogRef } from '@angular/cdk/dialog';
import { NgIcon } from '@ng-icons/core';

import { ButtonDirective } from '../../ui/button.directive';

export interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  /** When true the confirm button uses the destructive style. */
  danger?: boolean;
}

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  imports: [ButtonDirective, NgIcon],
  template: `
    <div
      class="w-[min(92vw,28rem)] rounded-xl border bg-card text-card-foreground shadow-2xl"
    >
      <div class="flex items-start gap-3 p-6 pb-4">
        @if (data.danger) {
          <span
            class="flex size-10 shrink-0 items-center justify-center rounded-full bg-destructive/10 text-destructive"
          >
            <ng-icon name="lucideTriangleAlert" size="20" />
          </span>
        }
        <div class="min-w-0 flex-1">
          <h2 class="text-lg font-semibold leading-tight">{{ data.title }}</h2>
          <p class="mt-1.5 whitespace-pre-line text-sm text-muted-foreground">{{ data.message }}</p>
        </div>
      </div>
      <div class="flex justify-end gap-2 px-6 pb-6">
        <button appBtn variant="outline" (click)="cancel()">
          {{ data.cancelText || 'Cancel' }}
        </button>
        <button appBtn [variant]="data.danger ? 'destructive' : 'default'" (click)="confirm()">
          {{ data.confirmText || 'Confirm' }}
        </button>
      </div>
    </div>
  `,
})
export class ConfirmDialogComponent {
  protected readonly data = inject<ConfirmDialogData>(DIALOG_DATA);
  private readonly dialogRef = inject<DialogRef<boolean>>(DialogRef);

  protected cancel(): void {
    this.dialogRef.close(false);
  }

  protected confirm(): void {
    this.dialogRef.close(true);
  }
}
