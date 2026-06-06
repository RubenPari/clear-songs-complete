import { computed, Directive, input } from '@angular/core';
import { cn } from '../utils/cn';

/** Surface container — Spartan/shadcn card. */
@Directive({
  selector: '[appCard]',
  standalone: true,
  host: { '[class]': 'classes()' },
})
export class CardDirective {
  // eslint-disable-next-line @angular-eslint/no-input-rename -- merge consumer's class list
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly classes = computed(() =>
    cn('rounded-xl border bg-card text-card-foreground shadow-sm', this.userClass()),
  );
}
