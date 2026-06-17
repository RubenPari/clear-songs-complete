import { CommonModule, DOCUMENT, isPlatformBrowser } from '@angular/common';
import { Component, DestroyRef, effect, inject, PLATFORM_ID, Renderer2, signal } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { TranslateModule } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';
import { EMPTY, fromEvent } from 'rxjs';
import { catchError, finalize, take } from 'rxjs/operators';

import { AuthService } from '../../core/services/auth.service';
import { LoadingService } from '../../core/services/loading.service';
import { NotificationService } from '../../core/services/notification.service';
import { PreferencesService } from '../../core/services/preferences.service';

@Component({
  selector: 'app-main-layout',
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss'],
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, TranslateModule, NgIcon],
})
export class MainLayoutComponent {
  private preferencesService = inject(PreferencesService);
  private renderer = inject(Renderer2);
  private document = inject(DOCUMENT);
  private platformId = inject(PLATFORM_ID);
  private destroyRef = inject(DestroyRef);
  private notificationService = inject(NotificationService);
  public authService = inject(AuthService);
  public loadingService = inject(LoadingService);

  isHandset = signal(false);
  isDarkTheme = this.preferencesService.isDarkTheme;
  currentLang = this.preferencesService.currentLang;
  mobileNavOpen = signal(false);
  userMenuOpen = signal(false);

  readonly navItems = [
    { route: '/dashboard', icon: 'lucideLayoutDashboard', label: 'NAV.DASHBOARD' },
    { route: '/tracks', icon: 'lucideListMusic', label: 'NAV.TRACKS' },
    { route: '/playlists', icon: 'lucideLibrary', label: 'NAV.PLAYLISTS' },
  ];

  constructor() {
    if (isPlatformBrowser(this.platformId)) {
      this.isHandset.set(window.innerWidth < 768);
      fromEvent(window, 'resize')
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe(() => {
          this.isHandset.set(window.innerWidth < 768);
          if (window.innerWidth >= 768) {
            this.mobileNavOpen.set(false);
          }
        });
    }

    effect(() => {
      const isDark = this.isDarkTheme();
      const root = this.document.documentElement;
      if (isDark) {
        this.renderer.addClass(root, 'dark');
      } else {
        this.renderer.removeClass(root, 'dark');
      }
    });
  }

  toggleTheme(): void {
    this.preferencesService.toggleTheme();
  }

  switchLanguage(): void {
    this.preferencesService.switchLanguage();
  }

  toggleMobileNav(open?: boolean): void {
    this.mobileNavOpen.set(open ?? !this.mobileNavOpen());
  }

  toggleUserMenu(open?: boolean): void {
    this.userMenuOpen.set(open ?? !this.userMenuOpen());
  }

  logout(): void {
    this.userMenuOpen.set(false);
    this.mobileNavOpen.set(false);
    this.loadingService.show();
    this.authService
      .logout()
      .pipe(
        take(1),
        catchError(() => {
          this.notificationService.error('Unable to log out. Please try again.');
          return EMPTY;
        }),
        finalize(() => this.loadingService.hide()),
      )
      .subscribe();
  }
}
