import { createStore } from 'solid-js/store'
import { IFile } from './App'

export interface IUser {
  name:string
  id: string
  code: string
}
const [user, setUser] = createStore({
  name: "",
  id: "",
  code: ""
})
export interface IOdev {
  name: string;
  lesson: string;
  files: IFile[];
}
export interface IProgress {
  state: boolean;
  value: number;
}


const [odevler, setOdevler] = createStore<IOdev[]>([])
const [progress, setProgress] = createStore({
  state:false,
  value:0
})



export const useOdevler = (): [IOdev[], (odevler: IOdev[]) => void] => [odevler, setOdevler];
export const useProgress = (): [IProgress, (progress: IProgress) => void] => [progress, setProgress];
export const useUser = (): [IUser, (user: IUser) => void] => [user, setUser];
