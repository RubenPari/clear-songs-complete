import { Component, Input, Output, EventEmitter, signal, computed, effect, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

export interface RangeValue {
  min: number;
  max: number;
}

@Component({
  selector: 'app-range-slider',
  templateUrl: './range-slider.component.html',
  styleUrls: ['./range-slider.component.scss'],
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule]
})
export class RangeSliderComponent implements AfterViewInit {
  @ViewChild('sliderTrack') sliderTrack!: ElementRef<HTMLDivElement>;
  
  @Input() min = 0;
  @Input() max = 100;
  @Input() step = 1;
  
  private _minValue = signal(0);
  private _maxValue = signal(100);
  
  @Input()
  get minValue(): number {
    return this._minValue();
  }
  set minValue(val: number) {
    this._minValue.set(val);
  }
  
  @Input()
  get maxValue(): number {
    return this._maxValue();
  }
  set maxValue(val: number) {
    this._maxValue.set(val);
  }
  
  @Output() rangeChange = new EventEmitter<RangeValue>();
  
  minInput = signal(0);
  maxInput = signal(100);
  
  constructor() {
    effect(() => {
      this.minInput.set(this._minValue());
    }, { allowSignalWrites: true });
    
    effect(() => {
      this.maxInput.set(this._maxValue());
    }, { allowSignalWrites: true });
  }
  
  ngAfterViewInit(): void {
    this.updateSliderPosition();
  }
  
  leftPercent = computed(() => {
    const range = this.max - this.min;
    if (range === 0) return 0;
    return ((this._minValue() - this.min) / range) * 100;
  });
  
  rightPercent = computed(() => {
    const range = this.max - this.min;
    if (range === 0) return 100;
    return 100 - ((this._maxValue() - this.min) / range) * 100;
  });
  
  onMinInputChange(value: string): void {
    const val = parseInt(value, 10);
    if (isNaN(val)) return;
    
    const newMin = Math.max(this.min, Math.min(val, this._maxValue() - this.step));
    this._minValue.set(newMin);
    this.minInput.set(newMin);
    this.updateSliderPosition();
    this.emitChange();
  }
  
  onMaxInputChange(value: string): void {
    const val = parseInt(value, 10);
    if (isNaN(val)) return;
    
    const newMax = Math.min(this.max, Math.max(val, this._minValue() + this.step));
    this._maxValue.set(newMax);
    this.maxInput.set(newMax);
    this.updateSliderPosition();
    this.emitChange();
  }
  
  onMinSliderChange(value: string): void {
    const val = parseInt(value, 10);
    if (isNaN(val)) return;
    
    const newMin = Math.max(this.min, Math.min(val, this._maxValue() - this.step));
    this._minValue.set(newMin);
    this.minInput.set(newMin);
    this.updateSliderPosition();
    this.emitChange();
  }
  
  onMaxSliderChange(value: string): void {
    const val = parseInt(value, 10);
    if (isNaN(val)) return;
    
    const newMax = Math.min(this.max, Math.max(val, this._minValue() + this.step));
    this._maxValue.set(newMax);
    this.maxInput.set(newMax);
    this.updateSliderPosition();
    this.emitChange();
  }
  
  private updateSliderPosition(): void {
    // Position updates are handled by CSS percentage bindings
  }
  
  private emitChange(): void {
    this.rangeChange.emit({
      min: this._minValue(),
      max: this._maxValue()
    });
  }
}
