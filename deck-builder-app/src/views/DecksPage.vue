<script lang="ts" setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { fetchDecks } from '@/api/client'
import CardPreviewModal from '@/components/CardPreviewModal.vue'
import type { DeckCard, GeneratedDeck } from '@/types/api'
import { normalizeRarity } from '@/utils/rarity'

const {
  data: decksData,
  isPending,
  isError,
} = useQuery({
  queryKey: ['decks'],
  queryFn: fetchDecks,
})

const decks = computed(() => decksData.value?.docs ?? [])
const selectedDeck = ref<GeneratedDeck | null>(null)
const previewCard = ref<DeckCard | null>(null)

function selectDeck(deck: GeneratedDeck) {
  selectedDeck.value =
    selectedDeck.value?.name === deck.name && selectedDeck.value?.size === deck.size
      ? null
      : deck
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
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
          Mes <span class="text-emerald-500">decks</span>
        </h1>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          Consulte tes decks sauvegardés.
        </p>
      </div>

      <RouterLink
        to="/deck/build"
        class="inline-flex self-start rounded-xl bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-emerald-600"
      >
        Créer un deck
      </RouterLink>
    </header>

    <div
      v-if="isPending"
      class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
    >
      Chargement des decks…
    </div>

    <div
      v-else-if="isError"
      class="rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-600 dark:border-red-900 dark:bg-red-950/30 dark:text-red-400"
    >
      Impossible de charger les decks.
    </div>

    <div
      v-else-if="decks.length === 0"
      class="rounded-2xl border border-dashed border-slate-300 p-12 text-center dark:border-slate-700"
    >
      <p class="text-slate-500 dark:text-slate-400">Aucun deck sauvegardé.</p>
      <RouterLink
        to="/deck/build"
        class="mt-4 inline-block text-sm font-semibold text-emerald-600 hover:underline dark:text-emerald-400"
      >
        Générer ton premier deck
      </RouterLink>
    </div>

    <template v-else>
      <ul class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <li v-for="(deck, index) in decks" :key="`${deck.name}-${index}`">
          <button
            type="button"
            class="w-full rounded-2xl border p-4 text-left transition"
            :class="
              selectedDeck?.name === deck.name && selectedDeck?.size === deck.size
                ? 'border-emerald-500 bg-emerald-500/5 dark:bg-emerald-500/10'
                : 'border-slate-200 bg-white hover:border-slate-300 dark:border-slate-700 dark:bg-slate-800 dark:hover:border-slate-600'
            "
            @click="selectDeck(deck)"
          >
            <h2 class="font-semibold text-slate-900 dark:text-white">{{ deck.name }}</h2>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
              {{ deck.size }} cartes · {{ deck.deck.length }} types
            </p>
            <p
              v-if="deck.strategy"
              class="mt-2 line-clamp-2 text-xs text-slate-500 dark:text-slate-400"
            >
              {{ deck.strategy }}
            </p>
          </button>
        </li>
      </ul>

      <section v-if="selectedDeck" class="flex flex-col gap-6">
        <div class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800">
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white">{{ selectedDeck.name }}</h2>
          <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">
            {{ selectedDeck.size }} cartes
          </p>
          <p v-if="selectedDeck.strategy" class="mt-2 text-sm text-slate-500 dark:text-slate-400">
            {{ selectedDeck.strategy }}
          </p>
        </div>

        <div>
          <h3 class="mb-3 text-sm font-semibold text-slate-900 dark:text-white">
            Cartes ({{ selectedDeck.deck.length }})
          </h3>

          <div
            v-if="selectedDeck.deck.length === 0"
            class="rounded-2xl border border-dashed border-slate-300 p-12 text-center text-slate-500 dark:border-slate-700 dark:text-slate-400"
          >
            Aucune carte dans ce deck.
          </div>

          <div v-else class="card-grid">
            <article
              v-for="card in selectedDeck.deck"
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
    </template>
  </div>

  <CardPreviewModal :card="previewCard" @close="closePreview" />
</template>
