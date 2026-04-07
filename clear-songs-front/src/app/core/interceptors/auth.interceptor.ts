/**
 * Authentication HTTP Interceptor
 * 
 * This interceptor is automatically applied to all HTTP requests made through
 * Angular's HttpClient. It performs two main functions:
 * 
 * 1. Request Modification:
 *    - Adds credentials (cookies) to all HTTP requests
 *    - This is essential for session-based authentication with the backend
 * 
 * 2. Error Handling:
 *    - Intercepts HTTP errors and provides user-friendly notifications
 *    - Handles authentication errors (401) by redirecting to login
 *    - Handles server errors (500) with appropriate error messages
 *    - Handles network errors (status 0) when server is unreachable
 * 
 * The interceptor uses Angular's functional interceptor approach (Angular 15+),
 * which is more tree-shakeable and performant than class-based interceptors.
 * 
 * Usage:
 * This interceptor is automatically applied to all HTTP requests via the
 * app.config.ts configuration. No manual setup is required in components.
 * 
 * Error Status Codes Handled:
 * - 401: Unauthorized - Session expired, redirects to login
 * - 500: Internal Server Error - Shows generic error message
 * - 0: Network Error - Server unreachable or CORS issue
 * 
 * @file auth.interceptor.ts
 * @author Clear Songs Development Team
 */
import { inject } from '@angular/core';
import { HttpInterceptorFn, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { catchError, throwError, of } from 'rxjs';
import { Router } from '@angular/router';
import { NotificationService } from '../services/notification.service';

/**
 * Authentication Interceptor Function
 * 
 * Functional interceptor that modifies HTTP requests and handles errors.
 * Uses Angular's inject() function for dependency injection, which is the
 * modern approach for functional interceptors.
 * 
 * @param request - The outgoing HTTP request
 * @param next - Function to pass the request to the next interceptor or backend
 * @returns Observable of the HTTP response or error
 * 
 * @example
 * // Automatically applied to all HTTP requests:
 * this.http.get('/api/tracks').subscribe(...)
 * // Request will include credentials and errors will be handled
 */
export const authInterceptor: HttpInterceptorFn = (request, next) => {
  // Inject dependencies using Angular's inject() function
  // This is the modern way to get services in functional interceptors
  const router = inject(Router);
  const notificationService = inject(NotificationService);

  /**
   * Request Modification
   * 
   * Clones the request and adds withCredentials: true to include cookies
   * in cross-origin requests. This is necessary for session-based authentication
   * where the backend sets authentication cookies.
   * 
   * The request is cloned rather than modified directly to maintain immutability,
   * which is a best practice in Angular.
   */
  const modifiedRequest = request.clone({
    withCredentials: true,
  });

  /**
   * Error Handling
   * 
   * Intercepts HTTP errors and provides appropriate handling:
   * - Shows user-friendly error messages
   * - Redirects to login on authentication failures
   * - Re-throws errors so components can still handle them if needed
   */
  return next(modifiedRequest).pipe(
    catchError((error: HttpErrorResponse) => {
      /**
       * 401 Unauthorized - Authentication Required
       * 
       * The user's session has expired or they are not authenticated.
       * Redirect to login page and show a notification.
       * 
       * Exception: Don't show error or redirect for /auth/is-auth endpoint,
       * as this is used to check authentication status and 401 is expected
       * when the user is not authenticated.
       */
      if (error.status === 401) {
        const isAuthCheck = request.url.includes('/auth/is-auth');

        if (isAuthCheck) {
          // Return a resolved response with success: false so httpResource
          // handles it gracefully instead of surfacing an unhandled error
          return of(new HttpResponse({ status: 200, body: { success: false } }));
        }

        notificationService.error('Session expired. Please login again.');
        router.navigate(['/login']);
      } 
      /**
       * 500 Internal Server Error
       * 
       * The server encountered an unexpected error. This could be due to:
       * - Database connection issues
       * - Application bugs
       * - External service failures (Spotify API)
       */
      else if (error.status === 500) {
        notificationService.error('Server error occurred. Please try again.');
      } 
      /**
       * Status 0 - Network Error
       * 
       * The request failed before reaching the server. Common causes:
       * - Server is not running
       * - Network connectivity issues
       * - CORS configuration problems
       * - Request timeout
       */
      else if (error.status === 0) {
        notificationService.error(
          'Unable to connect to server. Please check your connection.'
        );
      }

      // Re-throw the error so components can still handle it if needed
      // Using throwError with a factory function is the RxJS best practice
      return throwError(() => error);
    })
  );
};
