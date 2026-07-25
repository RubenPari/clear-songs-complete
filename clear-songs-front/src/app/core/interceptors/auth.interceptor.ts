/**
 * HTTP interceptor for authentication and error handling.
 * Injects credentials, handles 401 redirects, and displays error notifications.
 */
import { HttpErrorResponse, HttpInterceptorFn, HttpResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, of, throwError } from 'rxjs';

import { NotificationService } from '../services/notification.service';

/**
 * Intercepts all HTTP requests to add credentials and handle errors.
 * - 401 on /auth/is-auth: returns success=false instead of error
 * - 401 on other endpoints: shows notification and redirects to login
 * - 500: shows server error notification
 * - 0 (network error): shows connection error notification
 */
export const authInterceptor: HttpInterceptorFn = (request, next) => {
  const router = inject(Router);
  const notificationService = inject(NotificationService);

  const modifiedRequest = request.clone({ withCredentials: true });

  return next(modifiedRequest).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        const isAuthCheck = request.url.includes('/auth/is-auth');

        if (isAuthCheck) {
          return of(new HttpResponse({ status: 200, body: { success: false } }));
        }

        notificationService.error('Session expired. Please login again.');
        router.navigate(['/login']);
      } else if (error.status === 500) {
        notificationService.error('Server error occurred. Please try again.');
      } else if (error.status === 0) {
        notificationService.error('Unable to connect to server. Please check your connection.');
      }

      return throwError(() => error);
    })
  );
};
