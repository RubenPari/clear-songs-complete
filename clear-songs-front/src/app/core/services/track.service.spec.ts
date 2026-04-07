import { TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { HttpTestingController } from '@angular/common/http/testing';
import { TrackService } from './track.service';
import { TrackStore } from '../stores/track.store';
import { environment } from '../../../environments/environment';
import { ArtistSummary } from '../models/artist.model';
import { ApiResponse } from '../models/api-response.model';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TranslateModule } from '@ngx-translate/core';

describe('TrackService', () => {
  let service: TrackService;
  let httpMock: HttpTestingController;
  let trackStore: jasmine.SpyObj<TrackStore>;

  beforeEach(async () => {
    const trackStoreSpy = jasmine.createSpyObj('TrackStore', ['setArtists', 'removeArtist']);

    TestBed.configureTestingModule({
      imports: [TranslateModule.forRoot()],
      providers: [
        provideZonelessChangeDetection(),
        TrackService,
        { provide: TrackStore, useValue: trackStoreSpy },
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });

    service = TestBed.inject(TrackService);
    httpMock = TestBed.inject(HttpTestingController);
    trackStore = TestBed.inject(TrackStore) as jasmine.SpyObj<TrackStore>;
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', async () => {
    expect(service).toBeTruthy();
  });

  it('should fetch track summary using httpResource', async () => {
    const mockData: ArtistSummary[] = [
      { id: '1', name: 'Artist 1', count: 10 },
      { id: '2', name: 'Artist 2', count: 5 }
    ];

    const resource = TestBed.runInInjectionContext(() => service.getTrackSummaryResource());
    
    // Trigger the resource to make the request
    resource.value();
    TestBed.flushEffects();
    
    // In Angular 19/20, httpResource returns a ResourceRef
    // We expect a call to the summary endpoint
    const req = httpMock.expectOne(req => req.url.includes('/track/summary'));
    expect(req.request.method).toBe('GET');
    req.flush({ success: true, data: mockData } satisfies ApiResponse<ArtistSummary[]>);

    // Wait for the resource to update
    TestBed.flushEffects();
    await Promise.resolve();
    const value = resource.value();
    if (Array.isArray(value)) {
      expect(value).toEqual(mockData);
    } else {
      expect(value?.data).toEqual(mockData);
    }
  });

  it('should delete tracks by artist', async () => {
    const artistId = '123';
    const mockResponse: ApiResponse = { success: true };

    service.deleteTracksByArtist(artistId).subscribe(response => {
      expect(response.success).toBeTrue();
      expect(trackStore.removeArtist).toHaveBeenCalledWith(artistId);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/track/by-artist/${artistId}`);
    expect(req.request.method).toBe('DELETE');
    req.flush(mockResponse);
  });
});
