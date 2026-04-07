/**
 * Track Management Component
 * 
 * This component provides tools for bulk deletion of tracks based on count ranges.
 * Users can delete tracks from artists that have a specific number of tracks,
 * making it easy to clean up singles, EPs, albums, or large collections.
 * 
 * Features:
 * - Range-based track deletion (min/max track count)
 * - Quick preset buttons for common scenarios (singles, EPs, albums, etc.)
 * - Custom range input with validation
 * - Range validation (max must be >= min)
 * - Confirmation dialogs before deletion
 * - Loading states and error handling
 * 
 * Use Cases:
 * - Remove all singles (artists with 1 track)
 * - Remove EPs (artists with 2-5 tracks)
 * - Remove albums (artists with 11-20 tracks)
 * - Remove large collections (artists with 20+ tracks)
 * - Custom ranges for specific needs
 * 
 * The component uses reactive forms with a custom validator to ensure
 * range validity. All operations create backups before deletion.
 * 
 * @component
 * @selector app-track-management
 * @standalone true
 * @author Clear Songs Development Team
 */
import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators, AbstractControl, ValidationErrors, FormControl } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { finalize } from 'rxjs/operators';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { TrackService } from '../../../core/services/track.service';
import { NotificationService } from '../../../core/services/notification.service';
import { LoadingService } from '../../../core/services/loading.service';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';


/**
 * Preset Range Interface
 * 
 * Defines a preset range configuration with a human-readable label
 * and min/max values. Used for quick selection buttons.
 */
interface PresetRange {
  readonly label: string;
  readonly min: number | null;
  readonly max: number | null;
}

@Component({
  selector: 'app-track-management',
  templateUrl: './track-management.component.html',
  styleUrls: ['./track-management.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule
  ]
})
export class TrackManagementComponent {
  private fb = inject(FormBuilder);
  private trackService = inject(TrackService);
  private notificationService = inject(NotificationService);
  public loadingService = inject(LoadingService);
  private modalService = inject(NgbModal);
  private translate = inject(TranslateService);

  /**
   * Reactive form for range input
   */
  rangeForm: FormGroup<{ min: FormControl<number | null>; max: FormControl<number | null> }>;
  
  /**
   * Preset range configurations for quick selection
   */
  readonly presetRanges: readonly PresetRange[] = [
    { label: 'TRACKS.PRESET_SINGLES', min: 1, max: 1 },
    { label: 'TRACKS.PRESET_EPS', min: 2, max: 5 },
    { label: 'TRACKS.PRESET_SMALL', min: 6, max: 10 },
    { label: 'TRACKS.PRESET_ALBUMS', min: 11, max: 20 },
    { label: 'TRACKS.PRESET_LARGE', min: 20, max: null },
  ] as const;

  constructor() {
    this.rangeForm = this.fb.group(
      {
        min: this.fb.control<number | null>(null, [Validators.min(0)]),
        max: this.fb.control<number | null>(null, [Validators.min(0)]),
      },
      { validators: TrackManagementComponent.rangeValidator }
    );
  }

  static rangeValidator(form: AbstractControl): ValidationErrors | null {
    const min = form.get('min')?.value;
    const max = form.get('max')?.value;

    if (min !== null && max !== null && min > max) {
      return { invalidRange: true };
    }
    return null;
  }

  deleteByRange(): void {
    if (this.rangeForm.invalid) {
      return;
    }

    const { min, max } = this.rangeForm.value;
    let message = '';

    if (min !== null && max !== null) {
      message = this.translate.instant('TRACKS.DELETE_MSG_RANGE', { min, max });
    } else if (min !== null) {
      message = this.translate.instant('TRACKS.DELETE_MSG_MIN', { min });
    } else if (max !== null) {
      message = this.translate.instant('TRACKS.DELETE_MSG_MAX', { max });
    } else {
      this.notificationService.warning(this.translate.instant('TRACKS.NO_RANGE'));
      return;
    }

    const modalRef = this.modalService.open(ConfirmDialogComponent, {
      size: 'md',
      centered: true
    });
    modalRef.componentInstance.title = this.translate.instant('TRACKS.DELETE_TITLE');
    modalRef.componentInstance.message = message;
    modalRef.componentInstance.confirmText = this.translate.instant('COMMON.DELETE');
    modalRef.componentInstance.cancelText = this.translate.instant('COMMON.CANCEL');

    modalRef.result.then(
      (result) => {
        if (result) {
          this.loadingService.show();
          const minValue = min ?? undefined;
          const maxValue = max ?? undefined;
          this.trackService.deleteTracksByRange(minValue, maxValue)
            .pipe(
              finalize(() => this.loadingService.hide())
            )
            .subscribe({
              next: () => {
                this.notificationService.success(this.translate.instant('TRACKS.DELETE_SUCCESS'));
                this.rangeForm.reset();
              },
              error: () => {
                this.notificationService.error(this.translate.instant('TRACKS.DELETE_ERROR'));
              },
            });
        }
      },
      () => {
        // Modal dismissed
      }
    );
  }

  applyPreset(preset: PresetRange): void {
    this.rangeForm.patchValue({
      min: preset.min,
      max: preset.max,
    });
  }
}
