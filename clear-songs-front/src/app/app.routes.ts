import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { MainLayoutComponent } from './layouts/main-layout/main-layout.component';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login/login.component').then(m => m.LoginComponent),
  },
  {
    path: 'callback',
    loadComponent: () => import('./features/auth/callback/callback.component').then(m => m.CallbackComponent),
  },
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () => import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent),
      },
      {
        path: 'tracks',
        loadComponent: () => import('./features/tracks/track-management/track-management.component').then(m => m.TrackManagementComponent),
      },
      {
        path: 'playlists',
        loadComponent: () => import('./features/playlists/playlists.component').then(m => m.PlaylistsComponent),
      },
      {
        path: '',
        redirectTo: '/dashboard',
        pathMatch: 'full',
      },
    ],
  },
  {
    path: '**',
    redirectTo: '/dashboard',
  },
];
