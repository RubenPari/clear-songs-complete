/**
 * Root Application Component
 *
 * Standalone root component that bootstraps the application and provides
 * the router outlet. Initial locale and theme preferences are applied here
 * so the login/callback routes already respect the user's choices.
 */
import { Component, inject, PLATFORM_ID } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';
import { PreferencesService } from './core/services/preferences.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  styleUrls: ['./app.css'],
  standalone: true,
  imports: [RouterOutlet]
})
export class App {
  title = 'clear-songs-front';

  private preferencesService = inject(PreferencesService);
  private platformId = inject(PLATFORM_ID);

  constructor() {
    this.preferencesService.initLocale();

    if (isPlatformBrowser(this.platformId)) {
      this.preferencesService.initTheme();
    }
  }
}
