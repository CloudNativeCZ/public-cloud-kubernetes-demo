import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, timer } from 'rxjs';
import { switchMap, map } from 'rxjs/operators'
import { environment } from '../environments/environment';

export interface Question {

}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  private readonly questions$: Observable<Question>;

  constructor(private http: HttpClient) {
    this.questions$ = timer(0, 2000).pipe(
      switchMap(() => http.get(environment.apiURL)),
    )
  }

  get questions(): Observable<Question> {
    return this.questions$
  }
}
