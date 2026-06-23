<script lang="ts" setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useQuery, keepPreviousData } from '@tanstack/vue-query'
import { RouterLink } from 'vue-router'
import { fetchCards, fetchCollections } from '@/api/client'
import { useAuthSession } from '@/auth/session'
import OwnedQuantityEditor from '@/components/OwnedQuantityEditor.vue'
import { useSearch } from '@/composables/useSearch'
import { INK_COLORS, type Card, type InkColor } from '@/types/api'
import { normalizeRarity } from '@/utils/rarity'

const PAGE_SIZE = 50
const SEARCH_DEBOUNCE_MS = 350

const { search } = useSearch()
const { isAuthenticated } = useAuthSession()

const page = ref(1)
const debouncedSearch = ref('')
const color = ref<InkColor | ''>('')

let searchTimer: ReturnType<typeof setTimeout> | undefined
watch(search, (value) => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    debouncedSearch.value = value.trim()
  }, SEARCH_DEBOUNCE_MS)
})

// Tout changement de filtre ramène à la première page.
watch([debouncedSearch, color], () => {
  page.value = 1
})

const { data: query, isFetching, isError } = useQuery({
  queryKey: ['cards', page, debouncedSearch, color],
  queryFn: () =>
    fetchCards({
      page: page.value,
      limit: PAGE_SIZE,
      search: debouncedSearch.value,
      color: color.value,
    }),
  placeholderData: keepPreviousData,
})

const cards = computed(() => query.value?.docs ?? [])

const totalPages = computed(() => {
  const total = query.value?.total ?? 0
  return Math.max(1, Math.ceil(total / PAGE_SIZE))
})

const pageNumbers = computed(() => {
  const total = totalPages.value
  const current = page.value
  const span = 2
  const start = Math.max(1, current - span)
  const end = Math.min(total, current + span)
  const pages: number[] = []
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const firstPageNumber = computed(() => pageNumbers.value[0] ?? 1)
const lastPageNumber = computed(
  () => pageNumbers.value[pageNumbers.value.length - 1] ?? totalPages.value,
)

function goToPage(p: number) {
  if (p < 1 || p > totalPages.value || p === page.value) return
  page.value = p
}

const selectedCard = ref<Card | null>(null)
const selectedCollection = ref('')

const { data: collectionsData } = useQuery({
  queryKey: ['collections'],
  queryFn: () => fetchCollections({ limit: 50 }),
  enabled: computed(() => isAuthenticated.value),
  retry: false,
})

const collectionOptions = computed(
  () => collectionsData.value?.docs.map((c) => ({ id: c.id, name: c.name })) ?? [],
)

watch(
  collectionOptions,
  (options) => {
    if (!selectedCollection.value && options[0]) {
      selectedCollection.value = options[0].id
    }
  },
  { immediate: true },
)

const activeCollection = computed(() =>
  collectionsData.value?.docs.find((c) => c.id === selectedCollection.value),
)

function ownedQuantityFor(cardId: string): number {
  return activeCollection.value?.cards[cardId]?.quantity ?? 0
}

const ownedQuantity = computed(() => {
  if (!selectedCard.value) return 0
  return ownedQuantityFor(selectedCard.value.id)
})

function openCard(card: Card) {
  selectedCard.value = card
}

function closeCard() {
  selectedCard.value = null
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') closeCard()
}

watch(selectedCard, (card) => {
  if (typeof document === 'undefined') return
  if (card) {
    document.addEventListener('keydown', onKeydown)
    document.body.style.overflow = 'hidden'
  } else {
    document.removeEventListener('keydown', onKeydown)
    document.body.style.overflow = ''
  }
})

onBeforeUnmount(() => {
  clearTimeout(searchTimer)
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <div class="flex min-h-screen w-full flex-col gap-6 px-3 py-6 sm:gap-8 sm:px-6 sm:py-8 lg:px-8">
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Lorcana
        <span class="text-emerald-500">Cards</span>
      </h1>

      <div class="flex w-full flex-col gap-3 sm:w-auto sm:flex-row sm:items-center">
        <select
          v-if="collectionOptions.length > 0"
          v-model="selectedCollection"
          class="w-full rounded-xl border border-slate-300 bg-white py-2.5 pl-3 pr-8 text-sm text-slate-900 shadow-sm transition focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 sm:w-52 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        >
          <option v-for="col in collectionOptions" :key="col.id" :value="col.id">
            {{ col.name }}
          </option>
        </select>

        <select
          v-model="color"
          class="w-full rounded-xl border border-slate-300 bg-white py-2.5 pl-3 pr-8 text-sm capitalize text-slate-900 shadow-sm transition focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 sm:w-44 dark:border-slate-700 dark:bg-slate-800 dark:text-white"
        >
          <option value="">Toutes les couleurs</option>
          <option v-for="c in INK_COLORS" :key="c" :value="c" class="capitalize">
            {{ c }}
          </option>
        </select>
      </div>
    </header>

    <main class="flex-1">
      <div
        v-if="isError"
        class="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-red-600 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-400"
      >
        Failed to load cards. Is the API running on <code>localhost:9090</code>?
      </div>

      <div
        v-else-if="cards.length === 0 && !isFetching"
        class="rounded-xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700"
      >
        Aucune carte trouvée.
      </div>

      <div
        v-else
        class="card-grid"
        :class="{ 'opacity-60': isFetching }"
      >
        <article
          v-for="card in cards"
          :key="card.id"
          tabindex="0"
          role="button"
          :aria-label="`Preview ${card.name}`"
          class="card-tile group cursor-pointer"
          :data-rarity="normalizeRarity(card.rarity)"
          @click="openCard(card)"
          @keydown.enter="openCard(card)"
          @keydown.space.prevent="openCard(card)"
        >
          <div class="card-tile-media relative aspect-[5/7] overflow-hidden">
            <span
              v-if="selectedCollection && ownedQuantityFor(card.id) > 0"
              class="absolute right-2 top-2 z-10 rounded-full bg-emerald-500 px-2 py-0.5 text-xs font-bold text-white shadow"
            >
              ×{{ ownedQuantityFor(card.id) }}
            </span>
            <img
              :src="card.filepath"
              :alt="card.name"
              loading="lazy"
              class="size-full object-cover transition duration-300 group-hover:scale-105"
            />
          </div>
          <div class="p-3">
            <p class="truncate text-sm font-semibold text-slate-100" :title="card.name">
              {{ card.name }}
            </p>
            <p v-if="card.set || card.number" class="truncate text-xs text-slate-400">
              {{ [card.set, card.number].filter(Boolean).join(' · ') }}
            </p>
          </div>
        </article>
      </div>
    </main>

    <footer class="flex items-center justify-center gap-1 pt-2">
      <button
        type="button"
        class="rounded-lg border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 transition enabled:hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-40 dark:border-slate-700 dark:text-slate-200 dark:enabled:hover:bg-slate-800"
        :disabled="page <= 1"
        @click="goToPage(page - 1)"
      >
        Prev
      </button>

      <button
        v-if="firstPageNumber > 1"
        type="button"
        class="rounded-lg border border-slate-300 px-3.5 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="goToPage(1)"
      >
        1
      </button>
      <span v-if="firstPageNumber > 2" class="px-1 text-slate-400">…</span>

      <button
        v-for="p in pageNumbers"
        :key="p"
        type="button"
        class="rounded-lg border px-3.5 py-2 text-sm font-medium transition"
        :class="
          p === page
            ? 'border-emerald-500 bg-emerald-500 text-white'
            : 'border-slate-300 text-slate-700 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800'
        "
        @click="goToPage(p)"
      >
        {{ p }}
      </button>

      <span v-if="lastPageNumber < totalPages - 1" class="px-1 text-slate-400">…</span>
      <button
        v-if="lastPageNumber < totalPages"
        type="button"
        class="rounded-lg border border-slate-300 px-3.5 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="goToPage(totalPages)"
      >
        {{ totalPages }}
      </button>

      <button
        type="button"
        class="rounded-lg border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 transition enabled:hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-40 dark:border-slate-700 dark:text-slate-200 dark:enabled:hover:bg-slate-800"
        :disabled="page >= totalPages"
        @click="goToPage(page + 1)"
      >
        Next
      </button>
    </footer>
  </div>

  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="selectedCard"
        class="fixed inset-0 z-50 overflow-y-auto bg-black/70 p-4 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        @click.self="closeCard"
      >
        <div class="flex min-h-full items-center justify-center py-6 sm:py-8">
          <div class="relative flex w-full max-w-md flex-col items-center gap-4">
          <button
            type="button"
            aria-label="Close preview"
            class="absolute -right-2 -top-2 z-10 flex size-9 items-center justify-center rounded-full bg-white text-slate-700 shadow-lg transition hover:bg-slate-100 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700"
            @click="closeCard"
          >
            <svg
              class="size-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="2"
              stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
            </svg>
          </button>

          <img
            :src="selectedCard.filepath"
            :alt="selectedCard.name"
            class="max-h-[55vh] w-auto max-w-[min(100%,16rem)] rounded-2xl shadow-2xl sm:max-w-[18rem]"
          />

          <div class="text-center">
            <p class="text-lg font-semibold text-white">{{ selectedCard.name }}</p>
            <p
              v-if="selectedCard.title && selectedCard.title !== selectedCard.name"
              class="text-sm italic text-slate-300"
            >
              {{ selectedCard.title }}
            </p>
            <p
              v-if="selectedCard.franchise"
              class="mt-1 text-sm text-slate-300"
            >
              {{ selectedCard.franchise }}
            </p>
            <p
              v-if="selectedCard.set || selectedCard.number || selectedCard.rarity"
              class="mt-1 text-sm capitalize text-slate-400"
            >
              {{ [selectedCard.set, selectedCard.number, selectedCard.rarity].filter(Boolean).join(' · ') }}
            </p>
          </div>

          <div class="w-full rounded-2xl bg-white p-4 shadow-xl dark:bg-slate-800">
            <p
              v-if="!isAuthenticated"
              class="text-center text-sm text-slate-500 dark:text-slate-400"
            >
              <RouterLink
                :to="{ name: 'login', query: { redirect: '/' } }"
                class="font-semibold text-emerald-600 hover:underline dark:text-emerald-400"
              >
                Connecte-toi
              </RouterLink>
              pour ajouter des cartes à ta collection.
            </p>

            <p
              v-else-if="collectionOptions.length === 0"
              class="text-center text-sm text-slate-500 dark:text-slate-400"
            >
              Aucune collection. Crée-en une dans
              <RouterLink to="/collection" class="font-semibold text-emerald-600 hover:underline dark:text-emerald-400">
                Ma collection
              </RouterLink>
              pour ajouter cette carte.
            </p>

            <template v-else>
              <label>
                <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
                  Collection
                </span>
                <select
                  v-model="selectedCollection"
                  class="w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-900 dark:text-white"
                >
                  <option v-for="col in collectionOptions" :key="col.id" :value="col.id">
                    {{ col.name }}
                  </option>
                </select>
              </label>

              <OwnedQuantityEditor
                v-if="selectedCard"
                class="mt-4"
                :collection-id="selectedCollection"
                :card-id="selectedCard.id"
                :owned-quantity="ownedQuantity"
              />
            </template>
          </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
