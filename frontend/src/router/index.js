import Vue from 'vue'
import Router from 'vue-router'

import Account from '../views/Account'
import Main from '../components/Main'
import Task from '../views/Task'
import ProxyComponent from '../views/Proxy'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      component: Main,
      children: [
        {
          path: '',
          component: Account
        },
        {
          path: 'task',
          component: Task
        },
        {
          path: 'proxy',
          component: ProxyComponent
        }
      ]
    }
  ]
})
