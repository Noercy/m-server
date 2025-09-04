import { Router, Route } from 'preact-router';

import './app.css'
import SeriesView from './components/SeriesView';
import Reader from './components/Reader/Reader';
import Home from './components/Home/Home';


export const App = () => (
  
 <div id="main-content">
    <Router>
      <Route path="/" component={Home} />
      <Route path="/series/:id" component={SeriesView} />
      <Route path='/series/:id/reader/:vId' component={Reader}></Route>
    </Router>
  </div>
)