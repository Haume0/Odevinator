import { path } from '@tauri-apps/api';
import { createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'

export interface IUser {
  id: string
}
const [user, setUser] = createStore({
  id: "2314716027"
})
export const useUser = (): [IUser, (user: IUser) => void] => [user, setUser];
