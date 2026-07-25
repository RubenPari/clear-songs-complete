import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { Dialog, DialogRef, DIALOG_DATA } from '@angular/cdk/dialog';
import { finalize, forkJoin } from 'rxjs';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';

import { Track, ArtistSummary } from '../../core/models/artist.model';
import { TrackService } from '../../core/services/track.service';
import { NotificationService } from '../../core/services/notification.service';
import { openConfirmDialog } from '../../core/utils/modal-helper';
import { ApiResponse } from '../../core/models/api-response.model';
import { ButtonDirective } from '../../shared/ui/button.directive';

interface AlbumGroup {
  album: string;
  imageUrl: string;
  tracks: Track[];
}

export interface ArtistTracksDialogData {
  artist: ArtistSummary;
}

@Component({
  selector: 'app-artist-tracks-modal',
  standalone: true,
  imports: [TranslateModule, NgIcon, ButtonDirective],
  templateUrl: './artist-tracks-modal.component.html',
})
export class ArtistTracksModalComponent implements OnInit {
  private readonly data = inject<ArtistTracksDialogData>(DIALOG_DATA);
  private readonly dialogRef = inject<DialogRef<boolean>>(DialogRef);
  private readonly trackService = inject(TrackService);
  private readonly notificationService = inject(NotificationService);
  private readonly dialog = inject(Dialog);
  public readonly translate = inject(TranslateService);

  readonly artist = this.data.artist;

  tracks = signal<Track[]>([]);
  isLoading = signal<boolean>(true);
  tracksChanged = signal<boolean>(false);
  deletingTrackId = signal<string | null>(null);
  deletingAlbum = signal<string | null>(null);
  private collapsedAlbums = signal<Set<string>>(new Set());

  albumGroups = computed<AlbumGroup[]>(() => {
    const grouped = new Map<string, Track[]>();
    for (const track of this.tracks()) {
      const albumKey = track.album || 'Unknown Album';
      if (!grouped.has(albumKey)) {
        grouped.set(albumKey, []);
      }
      grouped.get(albumKey)!.push(track);
    }

    return Array.from(grouped.entries()).map(([album, tracks]) => ({
      album,
      imageUrl: tracks[0]?.image_url || '',
      tracks,
    }));
  });

  // Runs on component initialization.
  ngOnInit(): void {
    this.loadTracks();
  }

  // Closes the dialog, reporting whether tracks changed.
  close(): void {
    this.dialogRef.close(this.tracksChanged());
  }

  // Loads tracks.
  loadTracks(): void {
    this.isLoading.set(true);
    this.trackService
      .getTracksByArtist(this.artist.id)
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (response: ApiResponse<Track[]>) => {
          this.tracks.set(Array.isArray(response.data) ? response.data : []);
        },
        error: () => this.notificationService.error(this.translate.instant('ARTIST_MODAL.LOAD_ERROR')),
      });
  }

  // Toggles album.
  toggleAlbum(albumName: string): void {
    this.collapsedAlbums.update((set) => {
      const next = new Set(set);
      if (next.has(albumName)) {
        next.delete(albumName);
      } else {
        next.add(albumName);
      }
      return next;
    });
  }

  // Checks whether album collapsed.
  isAlbumCollapsed(albumName: string): boolean {
    return this.collapsedAlbums().has(albumName);
  }

  // Deletes album tracks.
  deleteAlbumTracks(group: AlbumGroup, event: Event): void {
    event.stopPropagation();

    openConfirmDialog(this.dialog, {
      title: this.translate.instant('ARTIST_MODAL.DELETE_ALBUM_TITLE'),
      message: this.translate.instant('ARTIST_MODAL.DELETE_ALBUM_MSG', {
        count: group.tracks.length,
        album: group.album,
      }),
      confirmText: this.translate.instant('ARTIST_MODAL.DELETE_ALBUM_CONFIRM'),
      cancelText: this.translate.instant('COMMON.CANCEL'),
      danger: true,
    }).subscribe((confirmed) => {
      if (!confirmed) {
        return;
      }
      this.deletingAlbum.set(group.album);
      const deleteObs = group.tracks.map((t) => this.trackService.deleteTrack(t.id));

      forkJoin(deleteObs)
        .pipe(finalize(() => this.deletingAlbum.set(null)))
        .subscribe({
          next: () => {
            const deletedIds = new Set(group.tracks.map((t) => t.id));
            this.tracks.update((t) => t.filter((item) => !deletedIds.has(item.id)));
            this.tracksChanged.set(true);
            this.notificationService.success(
              this.translate.instant('ARTIST_MODAL.DELETE_ALBUM_SUCCESS', {
                count: group.tracks.length,
                album: group.album,
              }),
            );

            if (this.tracks().length === 0) {
              this.dialogRef.close(true);
            }
          },
          error: () => this.notificationService.error(this.translate.instant('ARTIST_MODAL.DELETE_ALBUM_ERROR')),
        });
    });
  }

  // Deletes track.
  deleteTrack(track: Track): void {
    openConfirmDialog(this.dialog, {
      title: this.translate.instant('ARTIST_MODAL.DELETE_TRACK_TITLE'),
      message: this.translate.instant('ARTIST_MODAL.DELETE_TRACK_MSG', { name: track.name }),
      confirmText: this.translate.instant('COMMON.DELETE'),
      cancelText: this.translate.instant('COMMON.CANCEL'),
      danger: true,
    }).subscribe((confirmed) => {
      if (!confirmed) {
        return;
      }
      this.deletingTrackId.set(track.id);
      this.trackService
        .deleteTrack(track.id)
        .pipe(finalize(() => this.deletingTrackId.set(null)))
        .subscribe({
          next: () => {
            this.notificationService.success(this.translate.instant('ARTIST_MODAL.DELETE_TRACK_SUCCESS'));
            this.tracks.update((t) => t.filter((item) => item.id !== track.id));
            this.tracksChanged.set(true);

            if (this.tracks().length === 0) {
              this.dialogRef.close(true);
            }
          },
          error: () => this.notificationService.error(this.translate.instant('ARTIST_MODAL.DELETE_TRACK_ERROR')),
        });
    });
  }
}
