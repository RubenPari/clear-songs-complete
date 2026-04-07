/**
 * Playlists Management Component
 * 
 * This component provides tools for managing Spotify playlists, allowing users to
 * perform bulk operations on playlists including clearing tracks and removing tracks
 * from both playlists and the user's library.
 * 
 * Features:
 * - Visual playlist selection from user's library
 * - Two operation modes:
 *   1. Clear playlist only (tracks remain in library)
 *   2. Clear playlist AND library (tracks removed from both, with backup)
 * - Operation history tracking
 * - Confirmation dialogs for destructive operations
 * - Loading states and error handling
 * 
 * @component
 * @selector app-playlists
 * @standalone true
 * @author Clear Songs Development Team
 */
import { Component, inject, signal, computed, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { finalize } from 'rxjs/operators';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { PlaylistService } from '../../core/services/playlist.service';
import { NotificationService } from '../../core/services/notification.service';
import { LoadingService } from '../../core/services/loading.service';
import { ConfirmDialogComponent } from '../../shared/components/confirm-dialog/confirm-dialog.component';
import { UserPlaylist } from '../../core/models/artist.model';
import { ApiError } from '../../core/models/api-response.model';

/**
 * Playlist Action Type
 */
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
    const modalRef = this.modalService.open(ConfirmDialogComponent, {
      size: 'md',
      centered: true
    });
    modalRef.componentInstance.title = copy.title;
    modalRef.componentInstance.message = `${copy.message}\n\n${this.translate.instant('PLAYLISTS.PLAYLIST_ID')}: ${playlistId}`;
    modalRef.componentInstance.confirmText = copy.confirmText;
    modalRef.componentInstance.cancelText = this.translate.instant('PLAYLISTS.ACTION_CANCEL');

    modalRef.result.then(
      (confirmed) => {
        if (!confirmed) {
          return;
        }

        this.loadingService.show();
        const request$ =
          action === 'playlist'
            ? this.playlistService.deleteAllPlaylistTracks(playlistId)
            : this.playlistService.deleteAllPlaylistAndUserTracks(playlistId);

        request$
          .pipe(
            finalize(() => this.loadingService.hide())
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
      },
      () => {
        // Modal dismissed
      }
    );
  }
}
