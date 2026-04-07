/**
 * Main Layout Component
 * 
 * @component
 * @selector app-main-layout
 * @standalone true
 * @author Clear Songs Development Team
 */
import { Component, Renderer2, signal, inject, PLATFORM_ID, effect, untracked, TemplateRef } from '@angular/core';
import { DOCUMENT, isPlatformBrowser } from '@angular/common';
import { CommonModule } from '@angular/common';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { NgbOffcanvas, NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { AuthService } from '../../core/services/auth.service';
import { LoadingService } from '../../core/services/loading.service';

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
  public authService = inject(AuthService);
  public loadingService = inject(LoadingService);

  isHandset = signal(false);
  isDarkTheme = signal(false);
  currentLang = signal('en');

  constructor() {
    // Initialize translate
    this.translate.addLangs(['en', 'it']);

    if (isPlatformBrowser(this.platformId)) {
      // Setup handset detection
      this.isHandset.set(window.innerWidth < 768);
      window.addEventListener('resize', () => {
        this.isHandset.set(window.innerWidth < 768);
      });

      // Load theme preference
      const savedTheme = localStorage.getItem(this.THEME_KEY);
      if (savedTheme) {
        this.isDarkTheme.set(savedTheme === 'dark');
      } else {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        this.isDarkTheme.set(prefersDark);
      }

      // Load language preference
      const savedLang = localStorage.getItem(this.LANG_KEY);
      if (savedLang && ['en', 'it'].includes(savedLang)) {
        this.currentLang.set(savedLang);
        this.translate.use(savedLang);
      } else {
        this.translate.use('en');
      }
    }

    // Effect for applying theme
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

  toggleTheme(): void {
    this.isDarkTheme.update(value => !value);
  }

  switchLanguage(): void {
    const newLang = this.currentLang() === 'en' ? 'it' : 'en';
    this.currentLang.set(newLang);
    this.translate.use(newLang);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(this.LANG_KEY, newLang);
    }
  }

  logout(): void {
    this.authService.logout().subscribe();
  }

  openSidebar(content: TemplateRef<unknown>): void {
    this.offcanvasService.open(content, { position: 'start' });
  }
}
