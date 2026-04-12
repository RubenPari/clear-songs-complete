import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, ReactiveFormsModule, ValidationErrors, Validators } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { filter, finalize, switchMap } from 'rxjs/operators';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

import { LoadingService } from '../../../core/services/loading.service';
import { NotificationService } from '../../../core/services/notification.service';
import { TrackService } from '../../../core/services/track.service';
import { openConfirmDialog } from '../../../core/utils/modal-helper';

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

  rangeForm: FormGroup<{ min: FormControl<number | null>; max: FormControl<number | null> }>;

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

    openConfirmDialog(this.modalService, {
      title: this.translate.instant('TRACKS.DELETE_TITLE'),
      message,
      confirmText: this.translate.instant('COMMON.DELETE'),
      cancelText: this.translate.instant('COMMON.CANCEL'),
      size: 'md',
      centered: true,
    })
      .pipe(
        filter((confirmed) => confirmed),
        switchMap(() => {
          this.loadingService.show();
          const minValue = min ?? undefined;
          const maxValue = max ?? undefined;
          return this.trackService.deleteTracksByRange(minValue, maxValue).pipe(
            finalize(() => this.loadingService.hide())
          );
        })
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

  applyPreset(preset: PresetRange): void {
    this.rangeForm.patchValue({
      min: preset.min,
      max: preset.max,
    });
  }
}
