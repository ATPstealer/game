import { atom } from 'recoil'

const nickNameState = atom({
  key: 'nickName',
  default: ''
})

export {
  nickNameState
}