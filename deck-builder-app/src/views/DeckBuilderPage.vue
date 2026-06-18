<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { buildDeck, fetchCollections, saveDeck } from '@/api/client'
import CardPreviewModal from '@/components/CardPreviewModal.vue'
import {
  DECK_LEVELS,
  INK_COLORS,
  type DeckCard,
  type DeckLevel,
  type GeneratedDeck,
  type InkColor,
} from '@/types/api'
import { INK_COLOR_ACCENT, INK_COLOR_INFO } from '@/data/inkColors'
import { normalizeRarity } from '@/utils/rarity'

const { data: collectionsData } = useQuery({
  queryKey: ['collections'],
  queryFn: () => fetchCollections({ limit: 50 }),
  retry: false,
})

const collectionOptions = computed(
  () => collectionsData.value?.docs.map((c) => c.name) ?? [],
)

const collectionName = ref('')
const colors = ref<InkColor[]>([])
const level = ref<DeckLevel>('beginner')
const error = ref('')
const deck = ref<GeneratedDeck | null>(null)
const previewCard = ref<DeckCard | null>(null)
const saveFeedback = ref<{ type: 'success' | 'error'; message: string } | null>(null)
const deckSaved = ref(false)

watch(
  collectionOptions,
  (options) => {
    if (!collectionName.value && options[0]) {
      collectionName.value = options[0]
    }
  },
  { immediate: true },
)

function toggleColor(c: InkColor) {
  const idx = colors.value.indexOf(c)
  if (idx >= 0) {
    colors.value.splice(idx, 1)
  } else if (colors.value.length < 2) {
    colors.value.push(c)
  }
}

function isColorSelected(c: InkColor) {
  return colors.value.includes(c)
}

function isColorDisabled(c: InkColor) {
  return colors.value.length >= 2 && !isColorSelected(c)
}

const { mutate: build, isPending: isBuilding } = useMutation({
  mutationFn: buildDeck,
  onSuccess: (data) => {
    error.value = ''
    deck.value = data
    saveFeedback.value = null
    deckSaved.value = false
  },
  onError: () => {
    error.value = 'La génération du deck a échoué.'
  },
})

const { mutate: save, isPending: isSaving } = useMutation({
  mutationFn: saveDeck,
  onSuccess: () => {
    deckSaved.value = true
    saveFeedback.value = { type: 'success', message: 'Deck sauvegardé.' }
  },
  onError: () => {
    saveFeedback.value = { type: 'error', message: 'La sauvegarde a échoué.' }
  },
})

const canSubmit = computed(
  () => Boolean(collectionName.value) && colors.value.length === 2,
)

const levelLabel = computed(
  () => DECK_LEVELS.find((l) => l.value === level.value)?.label ?? level.value,
)

function submit() {
  error.value = ''
  if (!collectionName.value) {
    error.value = 'Sélectionne une collection.'
    return
  }
  if (colors.value.length !== 2) {
    error.value = 'Choisis exactement 2 couleurs.'
    return
  }
  build({
    collection: collectionName.value,
    colors: colors.value,
    level: level.value,
  })
}

function restart() {
  deck.value = null
  previewCard.value = null
  error.value = ''
  saveFeedback.value = null
  deckSaved.value = false
  colors.value = []
  level.value = 'beginner'
}

function saveCurrentDeck() {
  if (!deck.value || deck.value.deck.length === 0) return
  saveFeedback.value = null
  save({
    name: deck.value.name,
    size: deck.value.size,
    strategy: deck.value.strategy,
    deck: deck.value.deck.map((card) => ({
      id: card.id,
      quantity: card.quantity,
    })),
  })
}

function openPreview(card: DeckCard) {
  previewCard.value = card
}

function closePreview() {
  previewCard.value = null
}
</script>

<template>
  <div class="flex min-h-screen w-full flex-col gap-6 px-3 py-6 sm:gap-8 sm:px-6 sm:py-8 lg:px-8">
    <section
      v-if="!deck && !isBuilding"
      class="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm sm:p-6 dark:border-slate-700 dark:bg-slate-800"
    >
      <h2 class="text-sm font-semibold text-slate-900 dark:text-white">
        Identités des encres
      </h2>
      <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
        Rappel des caractéristiques principales de chaque couleur pour t'aider à choisir ton duo.
      </p>
      <div class="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <article
          v-for="c in INK_COLORS"
          :key="`info-${c}`"
          class="rounded-xl border p-3 transition"
          :class="
            isColorSelected(c)
              ? `${INK_COLOR_ACCENT[c].border} ${INK_COLOR_ACCENT[c].bg} shadow-sm`
              : 'border-slate-200 bg-slate-50/80 dark:border-slate-700 dark:bg-slate-900/40'
          "
        >
          <div class="flex items-center gap-2">
            <span
              class="size-2.5 shrink-0 rounded-full"
              :class="INK_COLOR_ACCENT[c].dot"
              aria-hidden="true"
            />
            <h3
              class="text-sm font-semibold"
              :class="isColorSelected(c) ? INK_COLOR_ACCENT[c].text : 'text-slate-900 dark:text-white'"
            >
              {{ INK_COLOR_INFO[c].label }}
            </h3>
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide"
              :class="
                isColorSelected(c)
                  ? `${INK_COLOR_ACCENT[c].bg} ${INK_COLOR_ACCENT[c].text}`
                  : 'bg-slate-200 text-slate-600 dark:bg-slate-700 dark:text-slate-300'
              "
            >
              {{ INK_COLOR_INFO[c].archetype }}
            </span>
          </div>
          <ul class="mt-2 space-y-1 pl-4">
            <li
              v-for="trait in INK_COLOR_INFO[c].traits"
              :key="trait"
              class="list-disc text-xs leading-relaxed text-slate-600 dark:text-slate-400"
            >
              {{ trait }}
            </li>
          </ul>
        </article>
      </div>
    </section>

    <header>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Créer un <span class="text-emerald-500">deck</span>
      </h1>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
        Choisis une collection, deux couleurs et un niveau pour générer un deck.
      </p>
    </header>

    <form
      v-if="!deck && !isBuilding"
      class="flex flex-col gap-6 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800"
      @submit.prevent="submit"
    >
      <div>
        <label class="mb-2 block text-sm font-semibold text-slate-900 dark:text-white">
          Collection
        </label>
        <p
          v-if="collectionOptions.length === 0"
          class="text-sm text-slate-500 dark:text-slate-400"
        >
          Aucune collection. Crée-en une dans
          <RouterLink to="/collection" class="font-semibold text-emerald-600 hover:underline dark:text-emerald-400">
            Ma collection
          </RouterLink>.
        </p>
        <select
          v-else
          v-model="collectionName"
          class="w-full rounded-xl border border-slate-300 bg-white px-3 py-2.5 text-sm text-slate-900 shadow-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-900 dark:text-white"
        >
          <option v-for="name in collectionOptions" :key="name" :value="name">
            {{ name }}
          </option>
        </select>
      </div>

      <div>
        <label class="mb-2 block text-sm font-semibold text-slate-900 dark:text-white">
          Couleurs
          <span class="font-normal text-slate-500 dark:text-slate-400">
            ({{ colors.length }}/2)
          </span>
        </label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="c in INK_COLORS"
            :key="c"
            type="button"
            :disabled="isColorDisabled(c)"
            class="rounded-full border px-4 py-1.5 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-40"
            :class="
              isColorSelected(c)
                ? `${INK_COLOR_ACCENT[c].border} ${INK_COLOR_ACCENT[c].bg} ${INK_COLOR_ACCENT[c].text} ring-2 ring-offset-1 ring-current dark:ring-offset-slate-800`
                : 'border-slate-300 text-slate-700 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-700'
            "
            @click="toggleColor(c)"
          >
            {{ INK_COLOR_INFO[c].label }}
          </button>
        </div>
      </div>

      <div>
        <label class="mb-2 block text-sm font-semibold text-slate-900 dark:text-white">
          Niveau du deck
        </label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="lvl in DECK_LEVELS"
            :key="lvl.value"
            type="button"
            class="rounded-xl border px-4 py-2 text-sm font-medium transition"
            :class="
              level === lvl.value
                ? 'border-emerald-500 bg-emerald-500 text-white'
                : 'border-slate-300 text-slate-700 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-700'
            "
            @click="level = lvl.value"
          >
            {{ lvl.label }}
          </button>
        </div>
      </div>

      <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

      <button
        type="submit"
        :disabled="!canSubmit"
        class="self-start rounded-xl bg-emerald-500 px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
      >
        Générer le deck
      </button>
    </form>

    <section
      v-else-if="isBuilding"
      class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-emerald-300 bg-emerald-50/50 px-6 py-20 text-center dark:border-emerald-800 dark:bg-emerald-950/20"
      aria-live="polite"
      aria-busy="true"
    >
      <div class="relative mb-8 h-28 w-20">
        <div
          class="absolute inset-0 animate-pulse rounded-lg border border-emerald-200 bg-white shadow-md dark:border-emerald-800 dark:bg-slate-800"
          style="animation-delay: 0ms"
        />
        <div
          class="absolute inset-0 translate-x-2 translate-y-2 animate-pulse rounded-lg border border-emerald-300 bg-white shadow-md dark:border-emerald-700 dark:bg-slate-800"
          style="animation-delay: 200ms"
        />
        <div
          class="absolute inset-0 translate-x-4 translate-y-4 animate-pulse rounded-lg border border-emerald-400 bg-white shadow-lg dark:border-emerald-600 dark:bg-slate-800"
          style="animation-delay: 400ms"
        />
        <div
          class="absolute -right-3 -top-3 size-10 animate-spin rounded-full border-4 border-emerald-200 border-t-emerald-500 bg-white dark:border-emerald-900 dark:border-t-emerald-400 dark:bg-slate-900"
        />
      </div>

      <h2 class="text-lg font-semibold text-slate-900 dark:text-white">
        Génération du deck en cours…
      </h2>
      <p class="mt-2 max-w-sm text-sm text-slate-500 dark:text-slate-400">
        L'IA sélectionne les meilleures cartes de ta collection. Cela peut prendre quelques instants.
      </p>

      <div class="mt-6 flex gap-1.5">
        <span
          v-for="i in 3"
          :key="i"
          class="size-2 animate-bounce rounded-full bg-emerald-500"
          :style="{ animationDelay: `${(i - 1) * 150}ms` }"
        />
      </div>
    </section>

    <section v-else-if="deck" class="flex flex-col gap-6">
      <div class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h2 class="text-2xl font-bold text-slate-900 dark:text-white">{{ deck.name }}</h2>
            <div class="mt-2 flex flex-wrap gap-2 capitalize">
              <span
                v-for="c in colors"
                :key="c"
                class="rounded-full px-3 py-1 text-xs font-semibold"
                :class="`${INK_COLOR_ACCENT[c].bg} ${INK_COLOR_ACCENT[c].text}`"
              >
                {{ INK_COLOR_INFO[c].label }}
              </span>
            </div>
            <p class="mt-3 text-sm text-slate-600 dark:text-slate-300">
              {{ deck.size }} cartes · {{ levelLabel }}
            </p>
            <p v-if="deck.strategy" class="mt-2 text-sm text-slate-500 dark:text-slate-400">
              {{ deck.strategy }}
            </p>
          </div>
          <div class="flex flex-col gap-2 self-start sm:items-end">
            <div class="flex flex-wrap gap-2">
              <button
                type="button"
                :disabled="isSaving || deckSaved || deck.deck.length === 0"
                class="rounded-xl bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
                @click="saveCurrentDeck"
              >
                {{ isSaving ? 'Sauvegarde…' : deckSaved ? 'Deck sauvegardé' : 'Sauvegarder le deck' }}
              </button>
              <button
                type="button"
                class="rounded-xl border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-700"
                @click="restart"
              >
                Nouveau deck
              </button>
            </div>
            <p
              v-if="saveFeedback"
              class="text-xs"
              :class="saveFeedback.type === 'success' ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-500'"
            >
              {{ saveFeedback.message }}
            </p>
          </div>
        </div>
      </div>

      <div>
        <h3 class="mb-3 text-sm font-semibold text-slate-900 dark:text-white">
          Cartes ({{ deck.deck.length }})
        </h3>

        <div
          v-if="deck.deck.length === 0"
          class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
        >
          Aucune carte dans ce deck.
        </div>

        <div
          v-else
          class="card-grid"
        >
          <article
            v-for="card in deck.deck"
            :key="`${card.id}-${card.quantity}`"
            tabindex="0"
            role="button"
            :aria-label="`Preview ${card.name}`"
            class="card-tile group cursor-pointer"
            :data-rarity="normalizeRarity(card.rarity)"
            @click="openPreview(card)"
            @keydown.enter="openPreview(card)"
            @keydown.space.prevent="openPreview(card)"
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
                class="size-full object-cover transition duration-300 group-hover:scale-105"
              />
            </div>
            <div class="p-3">
              <p class="truncate text-sm font-semibold text-slate-100" :title="card.name">
                {{ card.name }}
              </p>
              <p
                v-if="card.title && card.title !== card.name"
                class="truncate text-xs italic text-slate-400"
                :title="card.title"
              >
                {{ card.title }}
              </p>
              <p
                v-if="card.number || card.set"
                class="truncate text-xs text-slate-400"
              >
                {{ [card.number, card.set].filter(Boolean).join(' · ') }}
              </p>
              <p
                v-if="card.type || card.rarity"
                class="truncate text-xs capitalize text-slate-500"
              >
                {{ [card.type, card.rarity].filter(Boolean).join(' · ') }}
              </p>
            </div>
          </article>
        </div>
      </div>
    </section>
  </div>

  <CardPreviewModal :card="previewCard" @close="closePreview" />
</template>
