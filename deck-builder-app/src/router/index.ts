import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/views/HomePage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage,
    },
    {
      path: '/collection',
      name: 'collection',
      component: () => import('@/views/CollectionPage.vue'),
    },
    {
      path: '/deck',
      name: 'deck',
      component: () => import('@/views/DeckBuilderPage.vue'),
    },
  ],
})

export default router
