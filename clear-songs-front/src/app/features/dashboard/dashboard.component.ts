import { Component, computed, effect, inject, Injector, runInInjectionContext, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NgbModal, NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { filter, finalize, switchMap } from 'rxjs/operators';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { ArtistSummary } from '../../core/models/artist.model';
import { LoadingService } from '../../core/services/loading.service';
import { NotificationService } from '../../core/services/notification.service';
import { TrackService } from '../../core/services/track.service';
import { modalResult$, openConfirmDialog } from '../../core/utils/modal-helper';
import { D3BarChartComponent } from '../../shared/components/d3-bar-chart/d3-bar-chart.component';
import { SkeletonChartComponent, SkeletonStatComponent, SkeletonTableComponent } from '../../shared/components/skeleton/skeleton-components';
import { ArtistTracksModalComponent } from '../tracks/artist-tracks-modal.component';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    D3BarChartComponent,
    SkeletonStatComponent,
    SkeletonTableComponent,
    SkeletonChartComponent,
    NgbModule,
    TranslateModule
  ]
})
export class DashboardComponent {
  private injector = inject(Injector);
  private trackService = inject(TrackService);
  private notificationService = inject(NotificationService);
  public loadingService = inject(LoadingService);
  private modalService = inject(NgbModal);
  private translate = inject(TranslateService);

  searchFilter = signal<string>('');
  selectedGenre = signal<string>('');
  minRange = signal<number>(0);
  maxRange = signal<number>(100);

  rangeMinDraft = signal<number>(0);
  rangeMaxDraft = signal<number>(100);
  
  currentPage = signal<number>(1);
  itemsPerPage = signal<number>(10);
  sortColumn = signal<string>('name');
  sortDirection = signal<'asc' | 'desc'>('asc');

  /** True while fetching debug JSON for download. */
  debugGenresDownloading = signal(false);

  /** Single resource; URL factory reads signals so requests refetch when filters change. */
  private readonly trackSummaryResource = runInInjectionContext(this.injector, () =>
    this.trackService.createTrackSummaryResource({
      min: () => (this.minRange() > 0 ? this.minRange() : undefined),
      max: () => (this.maxRange() < 100 ? this.maxRange() : undefined),
      genre: () => (this.selectedGenre() || undefined),
    })
  );

  constructor() {
    effect(() => {
      this.selectedGenre();
      this.minRange();
      this.maxRange();
      this.currentPage.set(1);
    });
    
    effect(() => {
      if (this.searchFilter()) {
        this.currentPage.set(1);
      }
    });
    
    effect(() => {
      const resource = this.trackSummaryResource;
      if (resource.isLoading()) {
        this.loadingService.show();
      } else {
        this.loadingService.hide();
      }
    });

    effect(() => {
      const resource = this.trackSummaryResource;
      if (resource.error()) {
        this.notificationService.error(this.translate.instant('DASHBOARD.LOAD_ERROR'));
      }
    });
  }

  private getResource() {
    return this.trackSummaryResource;
  }

  isLoading = computed(() => this.getResource()?.isLoading() ?? true);
  
  artists = computed<ArtistSummary[]>(() => this.getResource()?.value()?.data ?? []);

  totalTracks = computed(() => this.artists().reduce((sum, artist) => sum + artist.count, 0));
  totalArtists = computed(() => this.artists().length);

  availableGenres = computed(() => {
    const genreSet = new Set<string>();
    this.artists().forEach(artist => {
      if (artist.genre) {
        genreSet.add(artist.genre);
      }
    });
    return Array.from(genreSet).sort();
  });

  maxTrackCount = computed(() => {
    const max = Math.max(...this.artists().map(a => a.count), 0);
    return max > 0 ? max : 100;
  });

  chartData = computed(() => {
    const data = this.artists();
    const sortedArtists = [...data].sort((a, b) => b.count - a.count).slice(0, 5);
    return sortedArtists.map(artist => ({
      label: artist.name,
      value: artist.count
    }));
  });

  public chartColors: string[] = [
    'rgba(29, 185, 84, 0.8)',
    'rgba(29, 200, 100, 0.8)',
    'rgba(0, 212, 255, 0.8)',
    'rgba(16, 185, 129, 0.8)',
    'rgba(245, 158, 11, 0.8)'
  ];

  filteredArtists = computed(() => {
    let filtered = this.artists();
    const filterValue = this.searchFilter().trim().toLowerCase();
    
    if (filterValue) {
      filtered = filtered.filter(artist => 
        artist.name.toLowerCase().includes(filterValue)
      );
    }
    
    const col = this.sortColumn();
    const dir = this.sortDirection();
    return [...filtered].sort((a, b) => {
      let comparison = 0;
      if (col === 'name') {
        comparison = a.name.localeCompare(b.name);
      } else if (col === 'count') {
        comparison = a.count - b.count;
      }
      return dir === 'asc' ? comparison : -comparison;
    });
  });

  paginatedArtists = computed(() => {
    const page = this.currentPage();
    const items = this.itemsPerPage();
    const start = (page - 1) * items;
    return this.filteredArtists().slice(start, start + items);
  });

  totalPages = computed(() => {
    return Math.ceil(this.filteredArtists().length / this.itemsPerPage());
  });

  loadTrackSummary(): void {
    this.trackService.invalidateLibraryCache().subscribe({
      next: () => this.getResource().reload(),
      error: () => {
        this.notificationService.error(this.translate.instant('DASHBOARD.LOAD_ERROR'));
      },
    });
  }

  /** Downloads all library artists with raw Spotify genre arrays as JSON (debug / strategy work). */
  downloadArtistGenresJson(): void {
    this.debugGenresDownloading.set(true);
    this.trackService
      .getArtistGenresDebug()
      .pipe(finalize(() => this.debugGenresDownloading.set(false)))
      .subscribe({
        next: res => {
          if (!res.success || res.data === undefined) {
            this.notificationService.error(this.translate.instant('DASHBOARD.DEBUG_GENRES_ERROR'));
            return;
          }
          const json = JSON.stringify(res.data, null, 2);
          const blob = new Blob([json], { type: 'application/json;charset=utf-8' });
          const url = URL.createObjectURL(blob);
          const a = document.createElement('a');
          a.href = url;
          a.download = `artist-genres-debug-${new Date().toISOString().replace(/[:.]/g, '-')}.json`;
          a.rel = 'noopener';
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          URL.revokeObjectURL(url);
        },
        error: () => {
          this.notificationService.error(this.translate.instant('DASHBOARD.DEBUG_GENRES_ERROR'));
        },
      });
  }

  applyFilter(event?: Event): void {
    const target = event?.target as HTMLInputElement | null;
    if (target) {
      this.searchFilter.set(target.value);
    }
  }

  clearGenre(event?: Event): void {
    event?.preventDefault();
    event?.stopPropagation();
    this.selectedGenre.set('');
  }

  coerceRangeDraft(value: string | number): number {
    const n = typeof value === 'number' ? value : parseInt(String(value), 10);
    return Number.isFinite(n) ? n : 0;
  }

  applyRangeFilter(): void {
    const step = 1;
    const cap = this.maxTrackCount();
    let minV = this.rangeMinDraft();
    let maxV = this.rangeMaxDraft();

    minV = Math.max(0, Math.min(minV, cap));
    maxV = Math.max(0, Math.min(maxV, cap));

    if (minV > maxV) {
      const t = minV;
      minV = maxV;
      maxV = t;
    }

    minV = Math.max(0, Math.min(minV, maxV - step));
    maxV = Math.min(cap, Math.max(maxV, minV + step));

    this.rangeMinDraft.set(minV);
    this.rangeMaxDraft.set(maxV);
    this.minRange.set(minV);
    this.maxRange.set(maxV);
  }

  resetFilters(): void {
    this.searchFilter.set('');
    this.selectedGenre.set('');
    const cap = this.maxTrackCount();
    this.minRange.set(0);
    this.maxRange.set(cap);
    this.rangeMinDraft.set(0);
    this.rangeMaxDraft.set(cap);
  }

  sortTable(column: string): void {
    if (this.sortColumn() === column) {
      this.sortDirection.set(this.sortDirection() === 'asc' ? 'desc' : 'asc');
    } else {
      this.sortColumn.set(column);
      this.sortDirection.set('asc');
    }
  }

  changePage(page: number): void {
    this.currentPage.set(page);
  }

  openArtistTracks(artist: ArtistSummary): void {
    const modalRef = this.modalService.open(ArtistTracksModalComponent, {
      size: 'lg',
      centered: true,
      scrollable: true,
    });
    modalRef.componentInstance.artist = artist;

    modalResult$<boolean>(modalRef, false)
      .pipe(filter((tracksChanged) => tracksChanged))
      .subscribe(() => this.loadTrackSummary());
  }

  deleteArtistTracks(artist: ArtistSummary): void {
    openConfirmDialog(this.modalService, {
      title: this.translate.instant('DASHBOARD.DELETE_ARTIST_TITLE'),
      message: this.translate.instant('DASHBOARD.DELETE_ARTIST_MSG', { count: artist.count, name: artist.name }),
      confirmText: this.translate.instant('COMMON.DELETE'),
      cancelText: this.translate.instant('COMMON.CANCEL'),
      size: 'md',
      centered: true,
    })
      .pipe(
        filter((confirmed) => confirmed),
        switchMap(() => {
          this.loadingService.show();
          return this.trackService.deleteTracksByArtist(artist.id).pipe(
            finalize(() => this.loadingService.hide())
          );
        })
      )
      .subscribe({
        next: () => {
          this.notificationService.success(this.translate.instant('DASHBOARD.DELETE_ARTIST_SUCCESS', { name: artist.name }));
          this.loadTrackSummary();
        },
        error: () => {
          this.notificationService.error(this.translate.instant('DASHBOARD.DELETE_ARTIST_ERROR'));
        },
      });
  }
}
