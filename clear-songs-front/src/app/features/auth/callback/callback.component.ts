/**
 * Callback component for handling Spotify OAuth redirect.
 * Processes the authorization code and completes the authentication flow.
 */
import { Component, OnInit, DestroyRef, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';
import { AuthService } from '../../../core/services/auth.service';
import { NotificationService } from '../../../core/services/notification.service';

/**
 * OAuth callback handler component.
 * Extracts the authorization code from query params and completes login.
 * Falls back to auth status check if no code is present.
 */
@Component({
  selector: 'app-callback',
  template: `
    <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-background text-foreground">
      <ng-icon name="lucideLoaderCircle" size="44" class="animate-spin text-primary" />
      <p class="text-base text-muted-foreground">{{ 'CALLBACK.AUTHENTICATING' | translate }}</p>
    </div>
  `,
  standalone: true,
  imports: [TranslateModule, NgIcon]
})
export class CallbackComponent implements OnInit {
  private readonly destroyRef = inject(DestroyRef);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);
  private readonly notificationService = inject(NotificationService);
  private readonly translate = inject(TranslateService);

  ngOnInit(): void {
    this.route.queryParams
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((params) => {
        const code = params['code'];

        // If authorization code is present, complete the OAuth flow.
        if (code) {
          this.authService.handleCallback(code)
            .pipe(takeUntilDestroyed(this.destroyRef))
            .subscribe({
              next: () => {
                this.notificationService.success(this.translate.instant('CALLBACK.LOGIN_SUCCESS'));
                this.router.navigate(['/dashboard']);
              },
              error: () => {
                this.notificationService.error(this.translate.instant('CALLBACK.LOGIN_FAILED'));
                this.router.navigate(['/login']);
              },
            });
        } else {
          // No code present - verify existing authentication status.
          this.authService.checkAuthStatus()
            .pipe(takeUntilDestroyed(this.destroyRef))
            .subscribe({
              next: (isAuthenticated) => {
                if (isAuthenticated) {
                  this.notificationService.success(this.translate.instant('CALLBACK.LOGIN_SUCCESS'));
                  this.router.navigate(['/dashboard']);
                } else {
                  this.notificationService.error(this.translate.instant('CALLBACK.AUTH_FAILED'));
                  this.router.navigate(['/login']);
                }
              },
              error: () => {
                this.notificationService.error(this.translate.instant('CALLBACK.VERIFY_FAILED'));
                this.router.navigate(['/login']);
              },
            });
        }
      });
  }
}
