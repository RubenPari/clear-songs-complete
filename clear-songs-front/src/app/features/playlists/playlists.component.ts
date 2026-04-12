import { CommonModule } from '@angular/common';
import { Component, computed, effect, inject, signal } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { filter, finalize, switchMap } from 'rxjs/operators';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { ApiError } from '../../core/models/api-response.model';
import { UserPlaylist } from '../../core/models/artist.model';
import { LoadingService } from '../../core/services/loading.service';
import { NotificationService } from '../../core/services/notification.service';
import { PlaylistService } from '../../core/services/playlist.service';
import { openConfirmDialog } from '../../core/utils/modal-helper';

type PlaylistAction = 'playlist' | 'playlistAndLibrary';

@Component({
  selector: 'app-playlists',
  templateUrl: './playlists.component.html',
  styleUrls: ['./playlists.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule
  ],
})
export class PlaylistsComponent {
  private playlistService = inject(PlaylistService);
  private notificationService = inject(NotificationService);
  public loadingService = inject(LoadingService);
  private modalService = inject(NgbModal);
  private translate = inject(TranslateService);

  lastOperation = signal<{ playlistId: string; action: PlaylistAction; timestamp: number } | undefined>(undefined);
  
  private playlistsResource = this.playlistService.getUserPlaylistsResource();
  userPlaylists = computed<UserPlaylist[]>(() => this.playlistsResource.value()?.data ?? []);
  loadingPlaylists = computed(() => this.playlistsResource.isLoading());
  
  selectedPlaylistId = signal<string | null>(null);

  private actionCopy = computed(() => ({
    playlist: {
      title: this.translate.instant('PLAYLISTS.ACTION_CLEAR_TITLE'),
      message: this.translate.instant('PLAYLISTS.ACTION_CLEAR_MSG'),
      confirmText: this.translate.instant('PLAYLISTS.ACTION_CLEAR_CONFIRM'),
      success: this.translate.instant('PLAYLISTS.ACTION_CLEAR_SUCCESS'),
      error: this.translate.instant('PLAYLISTS.ACTION_CLEAR_ERROR'),
    },
    playlistAndLibrary: {
      title: this.translate.instant('PLAYLISTS.ACTION_CLEAR_LIB_TITLE'),
      message: this.translate.instant('PLAYLISTS.ACTION_CLEAR_LIB_MSG'),
      confirmText: this.translate.instant('PLAYLISTS.ACTION_CLEAR_LIB_CONFIRM'),
      success: this.translate.instant('PLAYLISTS.ACTION_CLEAR_LIB_SUCCESS'),
      error: this.translate.instant('PLAYLISTS.ACTION_CLEAR_LIB_ERROR'),
    },
  }));

  constructor() {
    effect(() => {
      if (this.playlistsResource.error()) {
        this.notificationService.error(this.translate.instant('PLAYLISTS.LOAD_ERROR'));
      }
    });
  }

  selectPlaylist(playlist: UserPlaylist): void {
    this.selectedPlaylistId.set(playlist.id);
  }

  resetForm(): void {
    this.selectedPlaylistId.set(null);
  }

  handleAction(action: PlaylistAction): void {
    const playlistId = this.selectedPlaylistId();
    if (!playlistId) {
      return;
    }

    const copy = this.actionCopy()[action];

    openConfirmDialog(this.modalService, {
      title: copy.title,
      message: `${copy.message}\n\n${this.translate.instant('PLAYLISTS.PLAYLIST_ID')}: ${playlistId}`,
      confirmText: copy.confirmText,
      cancelText: this.translate.instant('PLAYLISTS.ACTION_CANCEL'),
      size: 'md',
      centered: true,
    })
      .pipe(
        filter((confirmed) => confirmed),
        switchMap(() => {
          this.loadingService.show();
          const request$ =
            action === 'playlist'
              ? this.playlistService.deleteAllPlaylistTracks(playlistId)
              : this.playlistService.deleteAllPlaylistAndUserTracks(playlistId);

          return request$.pipe(finalize(() => this.loadingService.hide()));
        })
      )
      .subscribe({
        next: () => {
          this.notificationService.success(copy.success);
          this.lastOperation.set({ playlistId, action, timestamp: Date.now() });
          this.selectedPlaylistId.set(null);
        },
        error: (error) => {
          const rawError: ApiError | string | undefined = error?.error?.error;
          const serverMessage = typeof rawError === 'string' ? rawError : rawError?.message;
          this.notificationService.error(serverMessage || copy.error);
        },
      });
  }
}
