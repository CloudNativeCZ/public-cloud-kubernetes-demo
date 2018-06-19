import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, timer } from 'rxjs';
import { switchMap, map } from 'rxjs/operators'
import { environment } from '../environments/environment';

import { FormControl } from '@angular/forms';

export interface Question {

}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  private readonly questions$: Observable<Question>;

  question = new FormControl()

  constructor(private http: HttpClient) {
    this.questions$ = timer(0, 2000).pipe(
      switchMap(() => http.get(environment.apiURL)),
    )
  }

  submitQuestion() {
    this.http.post(environment.apiURL, {
      question: this.question
    })
  }

  get questions(): Observable<Question> {
    return this.questions$
  }
}
