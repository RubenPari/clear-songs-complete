import { TestBed } from '@angular/core/testing';
import { HttpTestingController } from '@angular/common/http/testing';
import { PlaylistService } from './playlist.service';
import { environment } from '../../../environments/environment';
import { UserPlaylist } from '../models/artist.model';
import { ApiResponse } from '../models/api-response.model';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

describe('PlaylistService', () => {
  let service: PlaylistService;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    TestBed.configureTestingModule({
      imports: [TranslateModule.forRoot()],
      providers: [
        PlaylistService,
        provideZonelessChangeDetection(),
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });

    service = TestBed.inject(PlaylistService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', async () => {
    expect(service).toBeTruthy();
  });

  it('should fetch user playlists using httpResource', async () => {
    const mockPlaylists: UserPlaylist[] = [
      { id: '1', name: 'Playlist 1' },
      { id: '2', name: 'Playlist 2' }
    ];

    const resource = TestBed.runInInjectionContext(() => service.getUserPlaylistsResource());
    
    // Trigger the resource
    resource.value();
    TestBed.flushEffects();
    
    // httpResource will trigger a request
    const req = httpMock.expectOne(req => req.url.includes('/playlist/list'));
    expect(req.request.method).toBe('GET');
    req.flush({ success: true, data: mockPlaylists } satisfies ApiResponse<UserPlaylist[]>);

    // Wait for the resource to update its value
    TestBed.flushEffects();
    await Promise.resolve();
    const value = resource.value();
    if (Array.isArray(value)) {
      expect(value).toEqual(mockPlaylists);
    } else {
      expect(value?.data).toEqual(mockPlaylists);
    }
  });

  it('should delete playlist tracks', () => {
    const playlistId = 'playlist123';
    const mockResponse = { success: true };

    service.deleteAllPlaylistTracks(playlistId).subscribe(response => {
      expect(response.success).toBeTrue();
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/playlist/delete-tracks?id=${playlistId}`);
    expect(req.request.method).toBe('DELETE');
    req.flush(mockResponse);
  });
});
