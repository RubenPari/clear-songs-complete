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
  template: `
    <div class="flex max-h-[85vh] w-[min(94vw,42rem)] flex-col rounded-xl border bg-card text-card-foreground shadow-2xl">
      <!-- Header -->
      <div class="flex items-center justify-between gap-3 border-b px-6 py-4">
        <h2 class="truncate text-lg font-semibold">
          {{ translate.instant('ARTIST_MODAL.TITLE', { name: artist.name }) }}
        </h2>
        <button
          appBtn
          variant="ghost"
          size="icon"
          class="h-8 w-8 shrink-0"
          [attr.aria-label]="'COMMON.CLOSE' | translate"
          (click)="close()"
        >
          <ng-icon name="lucideX" size="18" />
        </button>
      </div>

      <!-- Body -->
      <div class="min-h-0 flex-1 overflow-y-auto p-4">
        @if (isLoading()) {
          <div class="flex flex-col items-center justify-center gap-3 py-16 text-muted-foreground">
            <ng-icon name="lucideLoaderCircle" size="32" class="animate-spin text-primary" />
            <p class="text-sm">{{ 'ARTIST_MODAL.LOADING' | translate }}</p>
          </div>
        } @else if (tracks().length === 0) {
          <div class="flex flex-col items-center justify-center gap-3 py-16 text-muted-foreground">
            <ng-icon name="lucideMusic" size="40" class="opacity-40" />
            <p class="text-sm">{{ 'ARTIST_MODAL.NO_TRACKS' | translate }}</p>
          </div>
        } @else {
          <div class="space-y-2">
            @for (group of albumGroups(); track group.album) {
              <div class="overflow-hidden rounded-lg border">
                <!-- Album header -->
                <div
                  class="flex cursor-pointer items-center justify-between gap-3 bg-muted/40 px-3 py-2.5 transition-colors hover:bg-muted"
                  role="button"
                  tabindex="0"
                  (click)="toggleAlbum(group.album)"
                  (keyup.enter)="toggleAlbum(group.album)"
                >
                  <div class="flex min-w-0 items-center gap-3">
                    <ng-icon
                      [name]="isAlbumCollapsed(group.album) ? 'lucideChevronRight' : 'lucideChevronDown'"
                      size="16"
                      class="shrink-0 text-muted-foreground"
                    />
                    @if (group.imageUrl) {
                      <img [src]="group.imageUrl" class="size-10 shrink-0 rounded object-cover" [alt]="group.album" />
                    } @else {
                      <div class="flex size-10 shrink-0 items-center justify-center rounded bg-muted text-muted-foreground">
                        <ng-icon name="lucideDisc3" size="18" />
                      </div>
                    }
                    <div class="flex min-w-0 flex-col">
                      <span class="truncate text-sm font-semibold">{{ group.album }}</span>
                      <span class="text-xs text-muted-foreground">
                        {{ group.tracks.length }}
                        {{ group.tracks.length === 1 ? ('ARTIST_MODAL.TRACK' | translate) : ('ARTIST_MODAL.TRACKS' | translate) }}
                      </span>
                    </div>
                  </div>
                  <button
                    appBtn
                    variant="ghost"
                    size="sm"
                    class="shrink-0 text-destructive hover:bg-destructive/10 hover:text-destructive"
                    (click)="deleteAlbumTracks(group, $event)"
                    [disabled]="deletingAlbum() === group.album"
                    [title]="('COMMON.DELETE_ALL_FROM' | translate) + group.album"
                  >
                    @if (deletingAlbum() === group.album) {
                      <ng-icon name="lucideLoaderCircle" size="15" class="animate-spin" />
                    } @else {
                      <ng-icon name="lucideTrash2" size="15" />
                    }
                    <span>{{ 'ARTIST_MODAL.DELETE_ALBUM_BTN' | translate }}</span>
                  </button>
                </div>

                <!-- Tracks -->
                @if (!isAlbumCollapsed(group.album)) {
                  <div class="divide-y divide-border/60">
                    @for (track of group.tracks; track track.id) {
                      <div class="flex items-center justify-between gap-3 py-2 pl-[4.25rem] pr-3 transition-colors hover:bg-muted/40">
                        <span class="truncate text-sm">{{ track.name }}</span>
                        <button
                          appBtn
                          variant="ghost"
                          size="icon"
                          class="size-8 shrink-0 text-destructive hover:bg-destructive/10 hover:text-destructive"
                          (click)="deleteTrack(track)"
                          [disabled]="deletingTrackId() === track.id"
                          [title]="'COMMON.DELETE' | translate"
                        >
                          @if (deletingTrackId() === track.id) {
                            <ng-icon name="lucideLoaderCircle" size="15" class="animate-spin" />
                          } @else {
                            <ng-icon name="lucideTrash2" size="15" />
                          }
                        </button>
                      </div>
                    }
                  </div>
                }
              </div>
            }
          </div>
        }
      </div>

      <!-- Footer -->
      <div class="flex justify-end border-t px-6 py-4">
        <button appBtn variant="outline" (click)="close()">{{ 'COMMON.CLOSE' | translate }}</button>
      </div>
    </div>
  `,
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
