import { computed, Directive, input } from '@angular/core';
import { cn } from '../utils/cn';

export type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline' | 'muted';

const VARIANTS: Record<BadgeVariant, string> = {
  default: 'border-transparent bg-primary/10 text-primary',
  secondary: 'border-transparent bg-secondary text-secondary-foreground',
  destructive: 'border-transparent bg-destructive/10 text-destructive',
  outline: 'border-border text-foreground',
  muted: 'border-transparent bg-muted text-muted-foreground',
};

/** Pill badge — Spartan/shadcn style. */
@Directive({
  selector: '[appBadge]',
  standalone: true,
  host: { '[class]': 'classes()' },
})
export class BadgeDirective {
  readonly variant = input<BadgeVariant>('default');
  // eslint-disable-next-line @angular-eslint/no-input-rename -- merge consumer's class list
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly classes = computed(() =>
    cn(
      'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-semibold font-display transition-colors',
      VARIANTS[this.variant()],
      this.userClass(),
    ),
  );
}
