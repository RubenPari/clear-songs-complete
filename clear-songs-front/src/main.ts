/**
 * Angular application entry point.
 * Bootstraps the root App component with the application configuration.
 */
import { bootstrapApplication } from '@angular/platform-browser';

import { App } from './app/app';
import { appConfig } from './app/app.config';

bootstrapApplication(App, appConfig).catch((err: unknown) => {
  console.error('Error bootstrapping application:', err);
});
