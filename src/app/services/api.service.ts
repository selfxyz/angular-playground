import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { SelfAppDisclosureConfig } from '@selfxyz/qrcode-angular';

export interface SaveOptionsRequest {
  userId: string;
  options: SelfAppDisclosureConfig;
}

export interface SaveOptionsResponse {
  message: string;
  userId?: string;
  savedAt?: string;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private http: HttpClient) {}

  saveOptions(request: SaveOptionsRequest): Observable<SaveOptionsResponse> {
    return this.http.post<SaveOptionsResponse>('https://8e252df08be5.ngrok-free.app/api/saveOptions', request);
  }
}
