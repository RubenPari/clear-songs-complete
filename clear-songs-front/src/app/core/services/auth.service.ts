import { Injectable, inject, signal, effect, Injector } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { environment } from '../../../environments/environment';
import { Observable, filter, map, take, tap } from 'rxjs';
import { HttpClient, httpResource } from '@angular/common/http';
import { Router } from '@angular/router';
import { ApiResponse, User } from '../models/api-response.model';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = environment.apiUrl;
  private http = inject(HttpClient);
  private router = inject(Router);
  private injector = inject(Injector);

  private _isAuthenticated = signal<boolean>(false);
  private _currentUser = signal<User | null>(null);
  private _sessionStatus = signal<'idle' | 'loading' | 'resolved' | 'error' | 'reloading' | 'local'>('idle');

  public readonly isAuthenticated = this._isAuthenticated.asReadonly();
  public readonly currentUser = this._currentUser.asReadonly();
  public readonly sessionStatus = this._sessionStatus.asReadonly();

  private sessionResource = httpResource<ApiResponse<{ user?: User }>>(() => `${this.apiUrl}/auth/is-auth`);

  constructor() {
    effect(() => {
      const session = this.sessionResource.value();
      const status = this.sessionResource.status();

      this._sessionStatus.set(status);

      if (status === 'resolved') {
        const isAuth = !!session?.success;
        this._isAuthenticated.set(isAuth);

        if (isAuth) {
          localStorage.setItem('isAuthenticated', 'true');
          this._currentUser.set(session?.data?.user ?? null);
        } else {
          localStorage.removeItem('isAuthenticated');
          this._currentUser.set(null);
        }
      } else if (status === 'error') {
        this._sessionStatus.set('error');
        localStorage.removeItem('isAuthenticated');
        this._currentUser.set(null);
      }
    });
  }

  login(): void {
    window.location.href = `${this.apiUrl}/auth/login`;
  }

  handleCallback(code: string): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.apiUrl}/auth/callback?code=${code}`).pipe(
      tap((response) => {
        if (response.success) {
          localStorage.setItem('isAuthenticated', 'true');
          this._isAuthenticated.set(true);
          this.sessionResource.reload();
        }
      }),
    );
  }

  logout(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.apiUrl}/auth/logout`).pipe(
      tap(() => {
        localStorage.removeItem('isAuthenticated');
        this._isAuthenticated.set(false);
        this._currentUser.set(null);
        this.sessionResource.reload();
        this.router.navigate(['/login']);
      }),
    );
  }

  checkAuthStatus(): Observable<boolean> {
    return toObservable(this.sessionStatus, { injector: this.injector }).pipe(
      filter((status) => status === 'resolved' || status === 'error'),
      map(() => this._isAuthenticated()),
      take(1)
    );
  }
}
