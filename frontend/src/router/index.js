
import { createRouter, createWebHistory } from 'vue-router';
import Welcome from '../components/Welcome.vue';
import SampleBuild from '../components/SampleBuild.vue';
import Configuration from '../components/Configuration.vue';

const routes = [
  {
    path: '/',
    name: 'Welcome',
    component: Welcome,
  },
  {
    path: '/build',
    name: 'Build',
    component: SampleBuild,
  },
  {
    path: '/config',
    name: 'Configuration',
    component: Configuration,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
