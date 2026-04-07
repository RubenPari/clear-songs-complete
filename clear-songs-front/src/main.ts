/**
 * Application Bootstrap Entry Point
 * 
 * This is the main entry point file that bootstraps the Angular application.
 * It uses the modern bootstrapApplication API (introduced in Angular 14+) instead
 * of the traditional NgModule-based bootstrap approach.
 * 
 * The bootstrap process:
 * 1. Imports the root App component
 * 2. Imports the application configuration (providers, routes, etc.)
 * 3. Bootstraps the application with error handling
 * 
 * Any errors during bootstrap are caught and logged to the console.
 * 
 * @file main.ts
 * @author Clear Songs Development Team
 */
import { bootstrapApplication } from '@angular/platform-browser';
import { App } from './app/app';
import { appConfig } from './app/app.config';

/**
 * Bootstrap the Angular application
 * 
 * Initializes the Angular application with the root App component and
 * the provided application configuration. If any error occurs during
 * the bootstrap process, it will be caught and logged to the console.
 * 
 * @param App - The root component class
 * @param appConfig - Application configuration containing providers, routes, etc.
 */
bootstrapApplication(App, appConfig).catch((err: unknown) => {
  console.error('❌ Error bootstrapping application:', err);
});
