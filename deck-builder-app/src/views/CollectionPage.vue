<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { createCollection, exportCollections, fetchCollections } from '@/api/client'
import OwnedQuantityEditor from '@/components/OwnedQuantityEditor.vue'
import { useSearch } from '@/composables/useSearch'
import type { OwnedCollection } from '@/types/api'
import { normalizeRarity } from '@/utils/rarity'

const queryClient = useQueryClient()
const { search } = useSearch()

const {
  data: collectionsData,
  isPending,
  isError,
} = useQuery({
  queryKey: ['collections'],
  queryFn: () => fetchCollections({ limit: 50 }),
})

const collections = computed(() => collectionsData.value?.docs ?? [])
const hasCollections = computed(() => collections.value.length > 0)

const selectedCollection = ref<OwnedCollection | null>(null)
const selectedCards = computed(() =>
  selectedCollection.value ? Object.values(selectedCollection.value.cards) : [],
)
const filteredSelectedCards = computed(() => {
  const term = search.value.trim().toLowerCase()
  if (!term) return selectedCards.value
  return selectedCards.value.filter(
    (card) =>
      card.name?.toLowerCase().includes(term) ||
      card.title?.toLowerCase().includes(term),
  )
})
const selectedCardsCount = computed(() =>
  selectedCollection.value ? Object.keys(selectedCollection.value.cards).length : 0,
)

const isFormOpen = ref(false)
const newName = ref('')
const error = ref('')
const exportError = ref('')

const { mutate: create, isPending: isCreating } = useMutation({
  mutationFn: createCollection,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['collections'] })
    closeForm()
  },
  onError: () => {
    error.value = 'Impossible de créer la collection.'
  },
})

const { mutate: exportCsv, isPending: isExporting } = useMutation({
  mutationFn: exportCollections,
  onSuccess: () => {
    exportError.value = ''
  },
  onError: () => {
    exportError.value = "Échec de l'export."
  },
})

function openForm() {
  isFormOpen.value = true
  newName.value = ''
  error.value = ''
}

function closeForm() {
  isFormOpen.value = false
  error.value = ''
}

function submit() {
  const name = newName.value.trim()
  if (!name) {
    error.value = 'Le nom est obligatoire.'
    return
  }
  error.value = ''
  create({ name })
}

function selectCollection(col: OwnedCollection) {
  selectedCollection.value = selectedCollection.value?.name === col.name ? null : col
}

watch(collectionsData, (data) => {
  if (!selectedCollection.value || !data) return
  const updated = data.docs.find((c) => c.id === selectedCollection.value?.id)
  if (updated) selectedCollection.value = updated
})
</script>

<template>
  <div class="flex min-h-screen w-full flex-col gap-6 px-3 py-6 sm:gap-8 sm:px-6 sm:py-8 lg:px-8">
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
          Mes <span class="text-emerald-500">collections</span>
        </h1>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          Gère tes collections de cartes.
        </p>
      </div>

      <div class="flex flex-wrap gap-3">
        <button
          type="button"
          :disabled="isExporting || !hasCollections"
          class="inline-flex items-center gap-2 self-start rounded-xl border border-slate-300 px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
          @click="exportCsv()"
        >
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
          </svg>
          {{ isExporting ? 'Export…' : 'Exporter' }}
        </button>
        <button
          type="button"
          class="inline-flex items-center gap-2 self-start rounded-xl bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-emerald-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40"
          @click="openForm"
        >
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          Créer une collection
        </button>
      </div>
    </header>

    <p v-if="exportError" class="-mt-4 text-sm text-red-500">{{ exportError }}</p>

    <main class="flex-1">
      <div
        v-if="isPending"
        class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
      >
        Chargement…
      </div>

      <div
        v-else-if="isError"
        class="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-red-600 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-400"
      >
        Échec du chargement. L'API tourne-t-elle sur <code>localhost:9090</code> ?
      </div>

      <div
        v-else-if="!hasCollections"
        class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-slate-300 p-12 text-center dark:border-slate-700"
      >
        <p class="text-slate-500 dark:text-slate-400">
          Tu n'as pas encore de collection.
        </p>
        <button
          type="button"
          class="text-sm font-semibold text-emerald-600 hover:underline dark:text-emerald-400"
          @click="openForm"
        >
          Créer ta première collection
        </button>
      </div>

      <div v-else class="flex flex-col gap-6">
        <ul class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <li
            v-for="col in collections"
            :key="col.name"
          >
            <button
              type="button"
              class="w-full rounded-xl border p-4 text-left shadow-sm transition hover:shadow-md"
              :class="
                selectedCollection?.name === col.name
                  ? 'border-emerald-500 bg-emerald-500/5 dark:bg-emerald-500/10'
                  : 'border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800'
              "
              @click="selectCollection(col)"
            >
              <p class="truncate font-semibold text-slate-900 dark:text-white" :title="col.name">
                {{ col.name }}
              </p>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                {{ Object.keys(col.cards).length }} carte(s)
              </p>
            </button>
          </li>
        </ul>

        <section v-if="selectedCollection">
          <h2 class="mb-4 text-lg font-semibold text-slate-900 dark:text-white">
            {{ selectedCollection.name }}
          </h2>

          <div
            v-if="selectedCardsCount === 0"
            class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
          >
            Cette collection est vide. Ajoute des cartes depuis la page Cartes.
          </div>

          <div
            v-else-if="filteredSelectedCards.length === 0"
            class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
          >
            Aucune carte ne correspond à ta recherche.
          </div>

          <div
            v-else
            class="card-grid"
          >
            <article
              v-for="card in filteredSelectedCards"
              :key="card.id"
              class="card-tile group cursor-pointer"
              :data-rarity="normalizeRarity(card.rarity)"
            >
              <div class="card-tile-media relative aspect-[5/7] overflow-hidden">
                <span
                  class="absolute right-2 top-2 z-10 rounded-full bg-emerald-500 px-2 py-0.5 text-xs font-bold text-white shadow"
                >
                  ×{{ card.quantity }}
                </span>
                <img
                  :src="card.filepath"
                  :alt="card.name"
                  loading="lazy"
                  class="size-full object-cover"
                />
              </div>
              <div class="p-3">
                <p class="truncate text-sm font-semibold text-slate-100" :title="card.name">
                  {{ card.name }}
                </p>
                <p v-if="card.set || card.number" class="truncate text-xs text-slate-400">
                  {{ [card.set, card.number].filter(Boolean).join(' · ') }}
                </p>
                <OwnedQuantityEditor
                  v-if="selectedCollection"
                  class="mt-3"
                  compact
                  :collection-id="selectedCollection.id"
                  :card-id="card.id"
                  :owned-quantity="card.quantity"
                />
              </div>
            </article>
          </div>
        </section>
      </div>
    </main>
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
        v-if="isFormOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        @click.self="closeForm"
      >
        <form
          class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl dark:bg-slate-800"
          @submit.prevent="submit"
        >
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">
            Nouvelle collection
          </h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            Donne un nom à ta collection.
          </p>

          <input
            v-model="newName"
            type="text"
            placeholder="Nom de la collection"
            autofocus
            class="mt-4 w-full rounded-xl border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 shadow-sm transition focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:placeholder:text-slate-400"
          />
          <p v-if="error" class="mt-2 text-sm text-red-500">{{ error }}</p>

          <div class="mt-6 flex justify-end gap-3">
            <button
              type="button"
              class="rounded-xl border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-700"
              @click="closeForm"
            >
              Annuler
            </button>
            <button
              type="submit"
              :disabled="isCreating"
              class="rounded-xl bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {{ isCreating ? 'Création…' : 'Créer' }}
            </button>
          </div>
        </form>
      </div>
    </Transition>
  </Teleport>
</template>
