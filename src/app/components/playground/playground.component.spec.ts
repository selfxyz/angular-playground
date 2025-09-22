import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FormsModule } from '@angular/forms';

import { PlaygroundComponent } from './playground.component';
import { ApiService } from '../../services/api.service';

describe('PlaygroundComponent', () => {
  let component: PlaygroundComponent;
  let fixture: ComponentFixture<PlaygroundComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        PlaygroundComponent,
        HttpClientTestingModule,
        FormsModule
      ],
      providers: [ApiService]
    }).compileComponents();

    fixture = TestBed.createComponent(PlaygroundComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with default disclosures', () => {
    expect(component.disclosures.nationality).toBe(true);
    expect(component.disclosures.minimumAge).toBe(18);
    expect(component.disclosures.ofac).toBe(true);
  });

  it('should generate userId on init', () => {
    expect(component.userId).toBeTruthy();
    expect(typeof component.userId).toBe('string');
  });

  it('should build selfApp when userId is available', () => {
    expect(component.selfApp).toBeTruthy();
    expect(component.selfApp?.appName).toBe('Self Playground');
    expect(component.selfApp?.userId).toBe(component.userId);
  });

  it('should toggle checkbox values', () => {
    const initialValue = component.disclosures.name;
    component.onCheckboxChange('name');
    expect(component.disclosures.name).toBe(!initialValue);
  });

  it('should update age from slider', () => {
    const mockEvent = {
      target: { value: '25' }
    } as any;
    
    component.onAgeChange(mockEvent);
    expect(component.disclosures.minimumAge).toBe(25);
  });

  it('should filter countries based on search query', () => {
    component.searchQuery = 'united';
    const filtered = component.filteredCountries;
    expect(filtered.length).toBeGreaterThan(0);
    expect(filtered.some(([_, name]) => name.toLowerCase().includes('united'))).toBe(true);
  });
});

