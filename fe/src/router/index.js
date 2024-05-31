import Vue from 'vue'
import VueRouter from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ChangeUserInfoView from '../views/ChangeUserInfoView.vue'
import AddUserView from '../views/AddUserView.vue'
import ChangeAdminInfoView from '../views/ChangeAdminInfoView.vue'
import CheckUserInfoView from '../views/CheckUserInfoView.vue'
import TestUserLoginView from '../views/TestUserLoginView'
Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'LoginView',
    component: LoginView
  },
  {
    path: '/RegisterView',
    name: 'RegisterView',
    component: RegisterView
  },
  {
    path: '/home',
    name: 'home',
    component: HomeView
  },
  {
    path: '/ChangeUserInfoView',
    name: 'ChangeUserInfoView',
    component: ChangeUserInfoView
  },
  {
    path: '/AddUserView',
    name: 'AddUserView',
    component: AddUserView
  },
  {
    path: '/ChangeAdminInfoView',
    name: 'ChangeAdminInfoView',
    component: ChangeAdminInfoView
  },
  {
    path: '/CheckUserInfoView',
    name: 'CheckUserInfoView',
    component: CheckUserInfoView
  },
  {
    path: '/TestUserLoginView',
    name: 'TestUserLoginView',
    component: TestUserLoginView
  },
]

const router = new VueRouter({
  routes
})

export default router
