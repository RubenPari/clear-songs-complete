import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection, signal, WritableSignal } from '@angular/core';
import { DashboardComponent } from './dashboard.component';
import { TrackService } from '../../core/services/track.service';
import { NotificationService } from '../../core/services/notification.service';
import { LoadingService } from '../../core/services/loading.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { TranslateModule } from '@ngx-translate/core';
import { ApiResponse } from '../../core/models/api-response.model';
import { ArtistSummary } from '../../core/models/artist.model';
import { of } from 'rxjs';


describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let trackService: jasmine.SpyObj<TrackService>;
  let resourceValue: WritableSignal<ApiResponse<ArtistSummary[]>>;
  let mockResource: {
    value: WritableSignal<ApiResponse<ArtistSummary[]>>;
    isLoading: ReturnType<typeof signal<boolean>>;
    error: ReturnType<typeof signal<unknown>>;
    status: ReturnType<typeof signal<number>>;
    reload: jasmine.Spy;
  };

  beforeEach(async () => {
    const trackServiceSpy = jasmine.createSpyObj('TrackService', [
      'createTrackSummaryResource',
      'deleteTracksByArtist',
      'invalidateLibraryCache',
    ]);
    trackServiceSpy.invalidateLibraryCache.and.returnValue(of({ success: true }));
    const notificationServiceSpy = jasmine.createSpyObj('NotificationService', ['success', 'error']);
    const loadingServiceSpy = jasmine.createSpyObj('LoadingService', ['show', 'hide']);
    const modalServiceSpy = jasmine.createSpyObj('NgbModal', ['open']);

    // Mock the resource API
    resourceValue = signal<ApiResponse<ArtistSummary[]>>({ success: true, data: [] });
    mockResource = {
      value: resourceValue,
      isLoading: signal(false),
      error: signal(null),
      status: signal(3), // Resolved
      reload: jasmine.createSpy('reload'),
    };
    trackServiceSpy.createTrackSummaryResource.and.returnValue(mockResource as never);

    await TestBed.configureTestingModule({
      imports: [DashboardComponent, TranslateModule.forRoot()],
      providers: [
        provideZonelessChangeDetection(),
        { provide: TrackService, useValue: trackServiceSpy },
        { provide: NotificationService, useValue: notificationServiceSpy },
        { provide: LoadingService, useValue: loadingServiceSpy as LoadingService },
        { provide: NgbModal, useValue: modalServiceSpy },
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
    trackService = TestBed.inject(TrackService) as jasmine.SpyObj<TrackService>;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should calculate statistics correctly from signals', () => {
    const mockData = [
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 20 }
    ];
    
    resourceValue.set({ success: true, data: mockData });
    
    expect(component.totalArtists()).toBe(2);
    expect(component.totalTracks()).toBe(30);
  });

  it('should update search filter and filter artists', () => {
    const mockData = [
      { id: '1', name: 'Alice', count: 10 },
      { id: '2', name: 'Bob', count: 20 }
    ];
    
    resourceValue.set({ success: true, data: mockData });
    
    component.searchFilter.set('ali');
    
    expect(component.filteredArtists().length).toBe(1);
    expect(component.filteredArtists()[0].name).toBe('Alice');
  });

  it('should invalidate cache then reload data when loadTrackSummary is called', () => {
    component.loadTrackSummary();
    expect(trackService.invalidateLibraryCache).toHaveBeenCalled();
    expect(mockResource.reload).toHaveBeenCalled();
  });
});
