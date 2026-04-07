import { TestBed } from '@angular/core/testing';
import { HttpTestingController } from '@angular/common/http/testing';
import { AuthService } from './auth.service';
import { Router } from '@angular/router';
import { ApiResponse } from '../models/api-response.model';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideZonelessChangeDetection } from '@angular/core';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  let router: jasmine.SpyObj<Router>;

  beforeEach(async () => {
    const routerSpy = jasmine.createSpyObj('Router', ['navigate']);

    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        AuthService,
        { provide: Router, useValue: routerSpy },
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });

    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    router = TestBed.inject(Router) as jasmine.SpyObj<Router>;
  });

  afterEach(() => {
    httpMock.verify();
    localStorage.removeItem('isAuthenticated');
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should expose signals for auth state', () => {
    expect(service.isAuthenticated).toBeDefined();
    expect(service.currentUser).toBeDefined();
    expect(service.sessionStatus).toBeDefined();
  });

  it('should handle callback and refresh session', () => {
    const code = 'auth-code';
    const mockResponse: ApiResponse = { success: true };

    service.handleCallback(code).subscribe(response => {
      expect(response.success).toBeTrue();
    });

    const req = httpMock.expectOne(req => req.url.includes('/auth/callback'));
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should navigate to login on logout', () => {
    localStorage.setItem('isAuthenticated', 'true');

    service.logout().subscribe(response => {
      expect(response).toBeDefined();
    });

    const req = httpMock.expectOne(req => req.url.includes('/auth/logout'));
    expect(req.request.method).toBe('GET');
    req.flush({ success: true });

    expect(router.navigate).toHaveBeenCalledWith(['/login']);
  });
});
