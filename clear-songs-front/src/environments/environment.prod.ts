/** Production environment configuration. */
import { environmentAuto } from './environment.auto';

/** Environment settings for production builds. */
export const environment = {
  ...environmentAuto,
  production: true,
};
