/**
 * Authentication Route Guard
 * 
 * This functional guard (Angular 15+) protects routes that require authentication.
 * It checks if the user is authenticated before allowing access to protected routes.
 * 
 * How it works:
 * 1. When a user navigates to a protected route, this guard is executed
 * 2. It calls AuthService.checkAuthStatus() to verify authentication
 * 3. If authenticated: allows navigation (returns true)
 * 4. If not authenticated: redirects to login page (returns UrlTree)
 * 
 * The guard uses Angular's functional guard approach, which is:
 * - More tree-shakeable than class-based guards
 * - Better for performance
 * - Easier to test
 * - Uses inject() for dependency injection
 * 
 * Usage:
 * Applied to routes in app.routes.ts:
 * ```typescript
 * {
 *   path: 'dashboard',
 *   component: DashboardComponent,
 *   canActivate: [authGuard]
 * }
 * ```
 * 
 * @file auth.guard.ts
 * @author Clear Songs Development Team
 */
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map } from 'rxjs';

/**
 * Authentication Guard Function
 * 
 * Functional route guard that protects routes requiring authentication.
 * Uses the Resource-based AuthService to check session status.
 */
export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.checkAuthStatus().pipe(
    map(isAuthenticated => {
      if (isAuthenticated) {
        return true;
      } else {
        return router.createUrlTree(['/login']);
      }
    })
  );
};
