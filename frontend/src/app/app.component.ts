import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, timer } from 'rxjs';
import { switchMap, map } from 'rxjs/operators'
import { environment } from '../environments/environment';

import { FormControl } from '@angular/forms';

export interface Question {
  Body: string
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
    this.questions$ = <any>timer(0, 2000).pipe(
      switchMap(() => http.get(environment.apiURL).pipe(map(res => (<any>res).Content))),
    )
    console.log(environment.apiURL)
    console.log(this.questions$)
  }

  submitQuestion() {
    this.http.put(environment.apiURL, {
      body: this.question.value
    }).subscribe()
  }

  get questions(): Observable<Question> {
    return this.questions$
  }
}
