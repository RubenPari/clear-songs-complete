/**
 * Playlists component for managing user playlists.
 * Provides options to clear playlist tracks or clear both playlist and library tracks.
 */
import { CommonModule } from '@angular/common';
import { Component, computed, effect, inject, signal } from '@angular/core';
import { Dialog } from '@angular/cdk/dialog';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';

import { ApiError } from '../../core/models/api-response.model';
import { UserPlaylist } from '../../core/models/artist.model';
import { LoadingService } from '../../core/services/loading.service';
import { NotificationService } from '../../core/services/notification.service';
import { PlaylistService } from '../../core/services/playlist.service';
import { confirmAndRunWithNotify } from '../../core/utils/confirm-run.helper';
import { ButtonDirective } from '../../shared/ui/button.directive';
import { CardDirective } from '../../shared/ui/card.directive';

/** Type of playlist deletion action. */
type PlaylistAction = 'playlist' | 'playlistAndLibrary';

/**
 * Component for viewing and managing user playlists.
 * Supports clearing tracks from playlists with optional library cleanup.
 */
@Component({
  selector: 'app-playlists',
  templateUrl: './playlists.component.html',
  styleUrls: ['./playlists.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule,
    NgIcon,
    ButtonDirective,
    CardDirective,
  ],
})
export class PlaylistsComponent {
  private playlistService = inject(PlaylistService);
  private notificationService = inject(NotificationService);
  public loadingService = inject(LoadingService);
  private dialog = inject(Dialog);
  private translate = inject(TranslateService);

  lastOperation = signal<{ playlistId: string; action: PlaylistAction; timestamp: number } | undefined>(undefined);
  
  private playlistsResource = this.playlistService.getUserPlaylistsResource();
  userPlaylists = computed<UserPlaylist[]>(() => this.playlistsResource.value()?.data ?? []);
  loadingPlaylists = computed(() => this.playlistsResource.isLoading());
  
  selectedPlaylistId = signal<string | null>(null);
  failedPlaylistImageIds = signal<Set<string>>(new Set());

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

  /** Selects a playlist for operations. */
  selectPlaylist(playlist: UserPlaylist): void {
    this.selectedPlaylistId.set(playlist.id);
  }

  /** Clears the selected playlist. */
  resetForm(): void {
    this.selectedPlaylistId.set(null);
  }

  /** Checks if the playlist has a valid image that hasn't failed to load. */
  hasPlaylistImage(playlist: UserPlaylist): boolean {
    return Boolean(playlist.image_url) && !this.failedPlaylistImageIds().has(playlist.id);
  }

  /** Marks a playlist image as failed to load. */
  onPlaylistImageError(playlistId: string): void {
    this.failedPlaylistImageIds.update((ids) => new Set(ids).add(playlistId));
  }

  /**
   * Executes a playlist action (clear playlist or clear playlist + library).
   * Extracts nested error messages from the response structure.
   */
  handleAction(action: PlaylistAction): void {
    const playlistId = this.selectedPlaylistId();
    if (!playlistId) {
      return;
    }

    const copy = this.actionCopy()[action];

    confirmAndRunWithNotify(
      {
        title: copy.title,
        message: `${copy.message}\n\n${this.translate.instant('PLAYLISTS.PLAYLIST_ID')}: ${playlistId}`,
        confirmText: copy.confirmText,
        cancelText: this.translate.instant('PLAYLISTS.ACTION_CANCEL'),
        danger: action === 'playlistAndLibrary',
      },
      () =>
        action === 'playlist'
          ? this.playlistService.deleteAllPlaylistTracks(playlistId)
          : this.playlistService.deleteAllPlaylistAndUserTracks(playlistId),
      {
        dialog: this.dialog,
        loadingService: this.loadingService,
        notificationService: this.notificationService,
        successMessage: copy.success,
        errorMessage: copy.error,
        onSuccess: () => {
          this.lastOperation.set({ playlistId, action, timestamp: Date.now() });
          this.selectedPlaylistId.set(null);
        },
        onError: (error) => {
          const rawError: ApiError | string | undefined = (error as { error?: { error?: ApiError | string } })?.error?.error;
          const serverMessage = typeof rawError === 'string' ? rawError : rawError?.message;
          if (serverMessage) {
            this.notificationService.error(serverMessage);
          }
        },
      },
    ).subscribe();
  }
}
