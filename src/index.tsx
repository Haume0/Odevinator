/* @refresh reload */
import { render } from 'solid-js/web'

import './index.css'
import {App, AppLayout, Duzenle, DuzenleOdev} from './App'
import {Login, LoginLayout} from './Login'
import { Route, Router } from '@solidjs/router'
import { createEffect, onMount } from 'solid-js'
import { useUser } from './Store'
const [user,setUser] = useUser()
onMount(() => {
  const user = sessionStorage.getItem('user')
  if(user){
    setUser(JSON.parse(user))
  }
})
createEffect(() => {
  sessionStorage.setItem('user', JSON.stringify(user))
})
const root = document.getElementById('root')
render(() => (
  <>
   <Router>
    <Route path="/" component={AppLayout}>
      <Route path="" component={App} />
      <Route path="duzenle/:index" component={DuzenleOdev} />
      <Route path="duzenle" component={Duzenle} />
    </Route>
    <Route path="/login" component={LoginLayout}>
      <Route path="" component={Login} />
    </Route>
   </Router>
  </>
), root!)