import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject, debounceTime, takeUntil } from 'rxjs';
import { v4 as uuidv4 } from 'uuid';

import { 
  SelfAppDisclosureConfig, 
  countryCodes, 
  SelfAppBuilder,
  type Country3LetterCode,
  type SelfApp,
  SelfQRcodeComponent
} from '@selfxyz/qrcode-angular';

import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-playground',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule,
    SelfQRcodeComponent
  ],
  templateUrl: './playground.component.html',
  styleUrls: ['./playground.component.css']
})
export class PlaygroundComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject<void>();
  private disclosuresSubject$ = new Subject<SelfAppDisclosureConfig>();

  userId: string | null = null;
  savingOptions = false;
  selfApp: SelfApp | null = null;

  disclosures: SelfAppDisclosureConfig = {
    // DG1 disclosures
    issuing_state: false,
    name: false,
    nationality: true,
    date_of_birth: false,
    passport_number: false,
    gender: false,
    expiry_date: false,
    // Custom checks
    minimumAge: 18,
    excludedCountries: ['IRN', 'IRQ', 'PRK', 'RUS', 'SYR', 'VEN'] as any,
    ofac: true,
  };

  showCountryModal = false;
  selectedCountries: Country3LetterCode[] = ['IRN', 'IRQ', 'PRK', 'RUS', 'SYR', 'VEN'] as any;

  countrySelectionError: string | null = null;
  searchQuery = '';

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.userId = uuidv4();
    this.buildSelfApp();

    this.saveOptionsToServer();

    // Set up debounced saving
    this.disclosuresSubject$
      .pipe(
        debounceTime(500),
        takeUntil(this.destroy$)
      )
      .subscribe(() => {
        if (this.userId) {
          this.saveOptionsToServer();
        }
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onAgeChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    const newAge = parseInt(target.value);
    this.disclosures = { ...this.disclosures, minimumAge: newAge };
    this.buildSelfApp();
    this.disclosuresSubject$.next(this.disclosures);
  }

  onCheckboxChange(field: keyof SelfAppDisclosureConfig): void {
    this.disclosures = {
      ...this.disclosures,
      [field]: !this.disclosures[field]
    };
    this.buildSelfApp();
    this.disclosuresSubject$.next(this.disclosures);
  }

  onCountryToggle(country: string): void {
    const countryCode = country as Country3LetterCode;
    if (this.selectedCountries.includes(countryCode)) {
      this.countrySelectionError = null;
      this.selectedCountries = this.selectedCountries.filter(c => c !== countryCode);
    } else {
      if (this.selectedCountries.length >= 40) {
        this.countrySelectionError = 'Maximum 40 countries can be excluded';
        return;
      }
      this.selectedCountries = [...this.selectedCountries, countryCode];
    }
  }

  isCountrySelected(country: string): boolean {
    return this.selectedCountries.includes(country as Country3LetterCode);
  }

  saveCountrySelection(): void {
    const codes = this.selectedCountries.map(countryName => {
      const entry = Object.entries(countryCodes).find(([_, name]) => name === countryName);
      return entry ? entry[0] : countryName.substring(0, 3).toUpperCase();
    }) as Country3LetterCode[];

    this.disclosures = { ...this.disclosures, excludedCountries: this.selectedCountries as any };
    this.showCountryModal = false;
    this.buildSelfApp();
    this.disclosuresSubject$.next(this.disclosures);
  }

  get filteredCountries() {
    return Object.entries(countryCodes).filter(([_, country]) =>
      (country as string).toLowerCase().includes(this.searchQuery.toLowerCase())
    );
  }

  private buildSelfApp(): void {
    if (!this.userId) return;

    // Create SelfApp object using SelfAppBuilder like in React version
    const app = new SelfAppBuilder({
      appName: "Self Playground",
      scope: "self-playground",
      endpoint: "https://cc10778f114e.ngrok-free.app/api/verify",
      endpointType: "staging_https",
      logoBase64: "https://i.imgur.com/Rz8B3s7.png",
      userId: this.userId,
      disclosures: {
        ...this.disclosures,
        minimumAge: this.disclosures.minimumAge && this.disclosures.minimumAge > 0 
          ? this.disclosures.minimumAge 
          : undefined,
      },
      version: 2,
      userDefinedData: "hello from the playground",
      devMode: false,
    } as Partial<SelfApp>).build();
    
    this.selfApp = app;
    console.log("selfApp built:", this.selfApp);
    console.log("selfApp keys:", Object.keys(this.selfApp));
    console.log("selfApp type:", typeof this.selfApp);
  }

  private async saveOptionsToServer(): Promise<void> {
    if (!this.userId || this.savingOptions) return;

    this.savingOptions = true;
    try {
      const response = await this.apiService.saveOptions({
        userId: this.userId,
        options: {
          minimumAge: this.disclosures.minimumAge && this.disclosures.minimumAge > 0 
            ? this.disclosures.minimumAge 
            : undefined,
          excludedCountries: this.disclosures.excludedCountries,
          ofac: this.disclosures.ofac,
          issuing_state: this.disclosures.issuing_state,
          name: this.disclosures.name,
          nationality: this.disclosures.nationality,
          date_of_birth: this.disclosures.date_of_birth,
          passport_number: this.disclosures.passport_number,
          gender: this.disclosures.gender,
          expiry_date: this.disclosures.expiry_date
        }
      }).toPromise();

      console.log("saved options to server", response);
    } catch (error) {
      console.error('Error saving options:', error);
      // Only show alert if it's a user-facing error
      if (error instanceof Error && error.message) {
        alert(error.message);
      } else {
        alert('Failed to save verification options. Please try again.');
      }
    } finally {
      this.savingOptions = false;
    }
  }

  onSuccess(): void {
    console.log('Verification successful');
  }

  onError(data: { error_code?: string; reason?: string }): void {
    console.error('Error generating QR code', data);
  }


  trackByCountry(index: number, item: [string, unknown]): string {
    return item[0];
  }
}
