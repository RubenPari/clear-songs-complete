import { CommonModule, DOCUMENT, isPlatformBrowser } from '@angular/common';
import { Component, DestroyRef, effect, inject, PLATFORM_ID, Renderer2, signal, TemplateRef, untracked } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { NgbModule, NgbOffcanvas } from '@ng-bootstrap/ng-bootstrap';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
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
  imports: [
    CommonModule,
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    NgbModule,
    TranslateModule
  ]
})
export class MainLayoutComponent {
  private readonly THEME_KEY = 'app-theme-preference';
  private readonly LANG_KEY = 'app-lang-preference';
  private renderer = inject(Renderer2);
  private document = inject(DOCUMENT);
  private platformId = inject(PLATFORM_ID);
  private offcanvasService = inject(NgbOffcanvas);
  private translate = inject(TranslateService);
  private destroyRef = inject(DestroyRef);
  private notificationService = inject(NotificationService);
  public authService = inject(AuthService);
  public loadingService = inject(LoadingService);

  isHandset = signal(false);
  isDarkTheme = signal(false);
  currentLang = signal('en');

  constructor() {
    this.translate.addLangs(['en', 'it']);

    if (isPlatformBrowser(this.platformId)) {
      this.isHandset.set(window.innerWidth < 768);
      fromEvent(window, 'resize')
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe(() => this.isHandset.set(window.innerWidth < 768));

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
      if (isDark) {
        this.renderer.addClass(this.document.body, 'dark-theme');
      } else {
        this.renderer.removeClass(this.document.body, 'dark-theme');
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
    this.isDarkTheme.update(value => !value);
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

  // Logs out.
  logout(): void {
    this.loadingService.show();
    this.authService.logout()
      .pipe(
        take(1),
        catchError(() => {
          this.notificationService.error('Unable to log out. Please try again.');
          return EMPTY;
        }),
        finalize(() => this.loadingService.hide())
      )
      .subscribe();
  }

  // Opens sidebar.
  openSidebar(content: TemplateRef<unknown>): void {
    this.offcanvasService.open(content, { position: 'start' });
  }
}
