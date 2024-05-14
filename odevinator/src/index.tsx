/* @refresh reload */
import { render } from 'solid-js/web'

import './index.css'
import {App, AppLayout, Duzenle} from './App'
import {LoginLayout} from './Login'
import { Route, Router } from '@solidjs/router'

const root = document.getElementById('root')
render(() => (
  <>
   <Router>
    <Route path="/" component={AppLayout}>
      <Route path="" component={App} />
      <Route path="duzenle" component={Duzenle} />
    </Route>
    <Route path="/login" component={LoginLayout} />
   </Router>
  </>
), root!)