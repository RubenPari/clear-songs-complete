import { Component, OnInit, DestroyRef, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from '../../../core/services/auth.service';
import { NotificationService } from '../../../core/services/notification.service';

@Component({
  selector: 'app-callback',
  template: `
    <div class="callback-container">
      <div class="spinner-border text-primary" role="status" style="width: 3rem; height: 3rem;">
        <span class="visually-hidden">{{ 'COMMON.LOADING' | translate }}</span>
      </div>
      <p>{{ 'CALLBACK.AUTHENTICATING' | translate }}</p>
    </div>
  `,
  styles: [
    `
      .callback-container {
        height: 100vh;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        p {
          margin-top: 20px;
          font-size: 18px;
          color: #666;
        }
      }
    `,
  ],
  standalone: true,
  imports: [CommonModule, TranslateModule]
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
