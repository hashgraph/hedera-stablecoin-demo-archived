import Vue from 'vue'
import Router from 'vue-router'
import Token from '../components/Token'
import Operate from '../components/Operate'
import InputKey from '../components/InputKey'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Token',
      component: Token
    },
    {
      path: '/operate',
      name: 'Operate',
      component: Operate
    },
    {
      path: '/keyInput',
      name: 'KeyInput',
      component: InputKey
    }
  ]
})
