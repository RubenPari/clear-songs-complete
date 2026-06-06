import { Component, inject } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { NgIcon } from '@ng-icons/core';

import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [TranslateModule, NgIcon],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  private authService = inject(AuthService);

  readonly features = [
    { icon: 'lucideTrash2', title: 'LOGIN.FEATURE_ARTIST_TITLE', desc: 'LOGIN.FEATURE_ARTIST_DESC' },
    { icon: 'lucideListFilter', title: 'LOGIN.FEATURE_STATS_TITLE', desc: 'LOGIN.FEATURE_STATS_DESC' },
    { icon: 'lucideLibrary', title: 'LOGIN.FEATURE_PLAYLIST_TITLE', desc: 'LOGIN.FEATURE_PLAYLIST_DESC' },
    { icon: 'lucideShieldCheck', title: 'LOGIN.FEATURE_BACKUP_TITLE', desc: 'LOGIN.FEATURE_BACKUP_DESC' },
  ];

  // Starts with spotify.
  loginWithSpotify(): void {
    this.authService.login();
  }
}
