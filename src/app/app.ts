import { Component } from '@angular/core';
import { PlaygroundComponent } from './components/playground/playground.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [PlaygroundComponent],
  template: `
    <app-playground></app-playground>
  `,
  styleUrl: './app.css'
})
export class App {}
