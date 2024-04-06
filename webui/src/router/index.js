import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'
import ProfileView from '../views/ProfileView.vue'
import ChangeUsername from '../components/ChangeUsernameComponent.vue';

const routes = [
  {path: '/', name: 'Login', component: LoginView},
  {path: '/stream', name: 'Stream', component: StreamView},
  {path: '/users/:username?', name: 'UserProfile', component: ProfileView},
  {path: '/change-username', name: 'ChangeUsername', component: ChangeUsername}
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router