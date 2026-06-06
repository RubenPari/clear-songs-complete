import { computed, Directive, input } from '@angular/core';
import { cn } from '../utils/cn';

/** Form input / select / textarea — Spartan/shadcn style. */
@Directive({
  selector: 'input[appInput], select[appInput], textarea[appInput]',
  standalone: true,
  host: { '[class]': 'classes()' },
})
export class InputDirective {
  // eslint-disable-next-line @angular-eslint/no-input-rename -- merge consumer's class list
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly classes = computed(() =>
    cn(
      'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground ring-offset-background transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
      this.userClass(),
    ),
  );
}
