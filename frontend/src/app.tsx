import { Router, Route } from 'preact-router';

import './app.css'
import SeriesList from './components/SeriesList';
import SeriesView from './components/SeriesView';
import Reader from './components/Reader/Reader';


export const App = () => (
  
 <div id="main-content">
    <Router>
      <Route path="/" component={SeriesList} />
      <Route path="/series/:id" component={SeriesView} />
      <Route path='/series/:id/reader/:vId' component={Reader}></Route>
    </Router>
  </div>
)