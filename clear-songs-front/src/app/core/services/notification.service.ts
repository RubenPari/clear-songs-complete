/**
 * Notification service for displaying toast messages.
 * Wraps ngx-toastr with convenience methods for different notification types.
 */
import { Injectable, inject } from '@angular/core';
import { ToastrService } from 'ngx-toastr';

/**
 * Service for displaying toast notifications.
 * Provides methods for success, error, warning, and info messages.
 */
@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  private toastr = inject(ToastrService);

  /** Displays a success notification. */
  success(message: string, title?: string): void {
    this.toastr.success(message, title || 'Success');
  }

  /** Displays an error notification. */
  error(message: string, title?: string): void {
    this.toastr.error(message, title || 'Error');
  }

  /** Displays a warning notification. */
  warning(message: string, title?: string): void {
    this.toastr.warning(message, title || 'Warning');
  }

  /** Displays an info notification. */
  info(message: string, title?: string): void {
    this.toastr.info(message, title || 'Info');
  }
}
