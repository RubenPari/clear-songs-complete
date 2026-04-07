/**
 * Root Application Component
 * 
 * This is the main entry point component for the Clear Songs Angular application.
 * It serves as the root component that wraps the entire application and provides
 * the router outlet for navigation between different views.
 * 
 * The component is configured as standalone, meaning it doesn't require an NgModule
 * and can directly import its dependencies. This follows Angular's modern standalone
 * component architecture introduced in Angular 14+.
 * 
 * @component
 * @selector app-root
 * @standalone true
 */
import { Component, inject, PLATFORM_ID } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  styleUrls: ['./app.css'],
  standalone: true,
  imports: [RouterOutlet]
})
export class App {
  /**
   * Application title
   * Used for identification and potentially in the browser title bar
   */
  title = 'clear-songs-front';
  
  private translate = inject(TranslateService);
  private platformId = inject(PLATFORM_ID);

  constructor() {
    this.translate.addLangs(['en', 'it']);
    this.translate.setDefaultLang('en');
    
    // Check if preference exists, otherwise use browser language or default
    if (isPlatformBrowser(this.platformId)) {
      const browserLang = this.translate.getBrowserLang() || 'en';
      const savedLang = localStorage.getItem('app-lang-preference');
      
      const langToUse = savedLang && ['en', 'it'].includes(savedLang) 
        ? savedLang 
        : (['en', 'it'].includes(browserLang) ? browserLang : 'en');
        
      this.translate.use(langToUse);
    } else {
      this.translate.use('en');
    }
  }
}
