import { CommonModule, DOCUMENT, isPlatformBrowser } from '@angular/common';
import { Component, DestroyRef, effect, inject, PLATFORM_ID, Renderer2, signal, untracked } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';
import { EMPTY, fromEvent } from 'rxjs';
import { catchError, finalize, take } from 'rxjs/operators';

import { AuthService } from '../../core/services/auth.service';
import { LoadingService } from '../../core/services/loading.service';
import { NotificationService } from '../../core/services/notification.service';

@Component({
  selector: 'app-main-layout',
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.scss'],
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, TranslateModule, NgIcon],
})
export class MainLayoutComponent {
  private readonly THEME_KEY = 'app-theme-preference';
  private readonly LANG_KEY = 'app-lang-preference';
  private renderer = inject(Renderer2);
  private document = inject(DOCUMENT);
  private platformId = inject(PLATFORM_ID);
  private translate = inject(TranslateService);
  private destroyRef = inject(DestroyRef);
  private notificationService = inject(NotificationService);
  public authService = inject(AuthService);
  public loadingService = inject(LoadingService);

  isHandset = signal(false);
  isDarkTheme = signal(false);
  currentLang = signal('en');
  mobileNavOpen = signal(false);
  userMenuOpen = signal(false);

  readonly navItems = [
    { route: '/dashboard', icon: 'lucideLayoutDashboard', label: 'NAV.DASHBOARD' },
    { route: '/tracks', icon: 'lucideListMusic', label: 'NAV.TRACKS' },
    { route: '/playlists', icon: 'lucideLibrary', label: 'NAV.PLAYLISTS' },
  ];

  constructor() {
    this.translate.addLangs(['en', 'it']);

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

      const savedTheme = localStorage.getItem(this.THEME_KEY);
      if (savedTheme) {
        this.isDarkTheme.set(savedTheme === 'dark');
      } else {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        this.isDarkTheme.set(prefersDark);
      }

      const savedLang = localStorage.getItem(this.LANG_KEY);
      if (savedLang && ['en', 'it'].includes(savedLang)) {
        this.currentLang.set(savedLang);
        this.translate.use(savedLang);
      } else {
        this.translate.use('en');
      }
    }

    effect(() => {
      const isDark = this.isDarkTheme();
      const root = this.document.documentElement;
      if (isDark) {
        this.renderer.addClass(root, 'dark');
      } else {
        this.renderer.removeClass(root, 'dark');
      }

      if (isPlatformBrowser(this.platformId)) {
        untracked(() => {
          localStorage.setItem(this.THEME_KEY, isDark ? 'dark' : 'light');
        });
      }
    });
  }

  // Toggles theme.
  toggleTheme(): void {
    this.isDarkTheme.update((value) => !value);
  }

  // Switches language.
  switchLanguage(): void {
    const newLang = this.currentLang() === 'en' ? 'it' : 'en';
    this.currentLang.set(newLang);
    this.translate.use(newLang);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(this.LANG_KEY, newLang);
    }
  }

  // Toggles the mobile navigation drawer.
  toggleMobileNav(open?: boolean): void {
    this.mobileNavOpen.set(open ?? !this.mobileNavOpen());
  }

  // Toggles the user menu.
  toggleUserMenu(open?: boolean): void {
    this.userMenuOpen.set(open ?? !this.userMenuOpen());
  }

  // Logs out.
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
