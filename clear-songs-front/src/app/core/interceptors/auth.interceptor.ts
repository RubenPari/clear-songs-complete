import { HttpErrorResponse, HttpInterceptorFn, HttpResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, of, throwError } from 'rxjs';

import { NotificationService } from '../services/notification.service';

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
