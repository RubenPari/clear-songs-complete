/** Development environment configuration. */
import { environmentAuto } from './environment.auto';

/** Environment settings for local development (non-production). */
export const environment = {
  ...environmentAuto,
  production: false,
};
