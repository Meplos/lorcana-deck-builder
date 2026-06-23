<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { useAuthSession } from '@/auth/session'
import { useSearch } from '@/composables/useSearch'

const route = useRoute()
const { search } = useSearch()
const { isAuthenticated, isInitialized } = useAuthSession()

const showSearch = computed(
  () => route.name === 'home' || route.name === 'collection',
)
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-slate-200 bg-white/80 backdrop-blur dark:border-slate-800 dark:bg-slate-900/80">
    <nav class="flex w-full flex-wrap items-center gap-4 px-3 py-3 sm:px-6 lg:px-8">
      <RouterLink to="/" class="shrink-0 text-lg font-bold tracking-tight text-slate-900 dark:text-white">
        Lorcana <span class="text-emerald-500">Deck Builder</span>
      </RouterLink>

      <div class="flex shrink-0 items-center gap-1 text-sm font-medium">
        <RouterLink to="/" class="nav-link">Cartes</RouterLink>
        <template v-if="isAuthenticated">
          <RouterLink to="/collection" class="nav-link">Ma collection</RouterLink>
          <RouterLink to="/deck" class="nav-link">Mes decks</RouterLink>
          <RouterLink to="/deck/build" class="nav-link">Créer un deck</RouterLink>
          <RouterLink to="/account" class="nav-link">Mon compte</RouterLink>
        </template>
        <template v-else-if="isInitialized">
          <RouterLink to="/login" class="nav-link">Connexion</RouterLink>
          <RouterLink to="/register" class="nav-link">Créer un compte</RouterLink>
        </template>
      </div>

      <div v-if="showSearch" class="relative min-w-0 flex-1 sm:max-w-xs lg:ml-auto lg:max-w-sm">
        <svg
          class="pointer-events-none absolute left-3 top-1/2 size-5 -translate-y-1/2 text-slate-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
          />
        </svg>
        <input
          v-model="search"
          type="search"
          placeholder="Rechercher une carte..."
          class="w-full rounded-xl border border-slate-300 bg-white py-2 pl-10 pr-4 text-sm text-slate-900 shadow-sm transition focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-800 dark:text-white dark:placeholder:text-slate-400"
        />
      </div>
    </nav>
  </header>

  <RouterView />
</template>

<style scoped>
@reference "tailwindcss";

.nav-link {
  @apply rounded-lg px-3 py-2 text-slate-600 transition hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-white;
}

.nav-link.router-link-exact-active {
  @apply bg-emerald-500/10 text-emerald-600 dark:text-emerald-400;
}
</style>
