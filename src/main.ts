import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';
import { provideSelfLottie } from '@selfxyz/qrcode-angular';

bootstrapApplication(App, {
  ...appConfig,
  providers: [
    ...appConfig.providers,
    provideSelfLottie(),
  ]
})
  .catch((err) => console.error(err));
