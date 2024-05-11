import { createRouter, createWebHistory } from 'vue-router'
// import LoginComponent from '../components/LoginComponent.vue'
import StreamView from '../views/StreamView.vue'
import ProfileView from '../views/ProfileView.vue'
import MyProfileView from '../views/MyProfileView.vue'
// import ChangeUsername from '../components/ChangeUsernameComponent.vue';

const routes = [
  //{path: '/', name: 'Login', component: LoginComponent},
  {path: '/stream', name: 'Stream', component: StreamView},
  {path: '/users/:username?', name: 'UserProfile', component: ProfileView},
  {path: '/myprofile', name: 'MyUserProfile', component: MyProfileView},
  //{path: '/change-username', name: 'ChangeUsername', component: ChangeUsername}
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router