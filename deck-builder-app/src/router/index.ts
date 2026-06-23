import { createRouter, createWebHistory } from 'vue-router'
import { ensureSession } from '@/auth/bootstrap'
import { useAuthSession } from '@/auth/session'
import HomePage from '@/views/HomePage.vue'

const AUTH_REQUIRED = new Set(['collection', 'decks', 'deck-build', 'account'])
const GUEST_ONLY = new Set(['login', 'register'])

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
      meta: { requiresAuth: true },
    },
    {
      path: '/deck',
      name: 'decks',
      component: () => import('@/views/DecksPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/deck/build',
      name: 'deck-build',
      component: () => import('@/views/DeckBuilderPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/account',
      name: 'account',
      component: () => import('@/views/AccountPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginPage.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterPage.vue'),
      meta: { guestOnly: true },
    },
  ],
})

router.beforeEach(async (to) => {
  await ensureSession()
  const { isAuthenticated } = useAuthSession()

  if (to.meta.requiresAuth || AUTH_REQUIRED.has(String(to.name))) {
    if (!isAuthenticated.value) {
      return {
        name: 'login',
        query: { redirect: to.fullPath },
      }
    }
  }

  if (to.meta.guestOnly || GUEST_ONLY.has(String(to.name))) {
    if (isAuthenticated.value) {
      const redirect =
        typeof to.query.redirect === 'string' && to.query.redirect.startsWith('/')
          ? to.query.redirect
          : '/'
      return redirect
    }
  }

  return true
})

export default router
