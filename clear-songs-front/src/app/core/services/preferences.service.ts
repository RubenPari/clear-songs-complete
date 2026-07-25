/**
 * Preferences service for managing user preferences (theme, language).
 * Persists preferences to localStorage and applies them to the document.
 */
import { Injectable, PLATFORM_ID, Signal, inject, signal } from '@angular/core';
import { DOCUMENT, isPlatformBrowser } from '@angular/common';
import { TranslateService } from '@ngx-translate/core';

/** Supported theme variants. */
export type Theme = 'dark' | 'light';

const SUPPORTED_LANGS = ['en', 'it'] as const;
type SupportedLang = typeof SUPPORTED_LANGS[number];

/**
 * Service for managing user preferences.
 * Handles theme toggling (dark/light) and language switching (en/it).
 */
@Injectable({
  providedIn: 'root',
})
export class PreferencesService {
  private readonly THEME_KEY = 'app-theme-preference';
  private readonly LANG_KEY = 'app-lang-preference';

  private document = inject(DOCUMENT);
  private platformId = inject(PLATFORM_ID);
  private translate = inject(TranslateService);

  private _currentLang = signal<SupportedLang>('en');
  private _isDarkTheme = signal<boolean>(false);

  readonly currentLang: Signal<SupportedLang> = this._currentLang.asReadonly();
  readonly isDarkTheme: Signal<boolean> = this._isDarkTheme.asReadonly();

  constructor() {
    this.translate.addLangs([...SUPPORTED_LANGS]);
    this.translate.setDefaultLang('en');
  }

  /**
   * Initializes the user's language preference.
   * Resolution order: saved preference > browser language > default (en).
   */
  initLocale(): void {
    if (isPlatformBrowser(this.platformId)) {
      const browserLang = (this.translate.getBrowserLang() || 'en') as string;
      const savedLang = localStorage.getItem(this.LANG_KEY) as SupportedLang | null;
      const langToUse = this.isSupported(savedLang) ? savedLang
        : (this.isSupported(browserLang) ? browserLang : 'en');
      this.useLanguage(langToUse);
    } else {
      this.translate.use('en');
    }
  }

  /**
   * Initializes the user's theme preference.
   * Resolution order: saved preference > system preference (prefers-color-scheme).
   */
  initTheme(): void {
    if (isPlatformBrowser(this.platformId)) {
      const savedTheme = localStorage.getItem(this.THEME_KEY) as Theme | null;
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      this.setDarkTheme(savedTheme ? savedTheme === 'dark' : prefersDark);
    }
  }

  /** Toggles between dark and light themes. */
  toggleTheme(): void {
    this.setDarkTheme(!this._isDarkTheme());
  }

  /** Switches between English and Italian languages. */
  switchLanguage(): void {
    const newLang = this._currentLang() === 'en' ? 'it' : 'en';
    this.useLanguage(newLang);
  }

  private useLanguage(lang: SupportedLang): void {
    this._currentLang.set(lang);
    this.translate.use(lang);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(this.LANG_KEY, lang);
    }
  }

  private setDarkTheme(isDark: boolean): void {
    this._isDarkTheme.set(isDark);
    this.document.documentElement.classList.toggle('dark', isDark);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(this.THEME_KEY, isDark ? 'dark' : 'light');
    }
  }

  private isSupported(lang: string | null): lang is SupportedLang {
    return !!lang && SUPPORTED_LANGS.includes(lang as SupportedLang);
  }
}
