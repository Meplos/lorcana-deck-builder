<script lang="ts" setup>
import { onBeforeUnmount, watch } from 'vue'
import type { DeckCard } from '@/types/api'

const props = defineProps<{
  card: DeckCard | null
}>()

const emit = defineEmits<{
  close: []
}>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

watch(
  () => props.card,
  (card) => {
    if (typeof document === 'undefined') return
    if (card) {
      document.addEventListener('keydown', onKeydown)
      document.body.style.overflow = 'hidden'
    } else {
      document.removeEventListener('keydown', onKeydown)
      document.body.style.overflow = ''
    }
  },
)

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
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
        v-if="card"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        @click.self="emit('close')"
      >
        <div class="relative flex max-h-full w-full max-w-md flex-col items-center gap-4">
          <button
            type="button"
            aria-label="Fermer la preview"
            class="absolute -right-2 -top-2 z-10 flex size-9 items-center justify-center rounded-full bg-white text-slate-700 shadow-lg transition hover:bg-slate-100 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700"
            @click="emit('close')"
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
            :src="card.filepath"
            :alt="card.name"
            class="max-h-[80vh] w-auto rounded-2xl shadow-2xl"
          />

          <div class="text-center">
            <p class="text-lg font-semibold text-white">{{ card.name }}</p>
            <p
              v-if="card.title && card.title !== card.name"
              class="text-sm italic text-slate-300"
            >
              {{ card.title }}
            </p>
            <p
              v-if="card.number || card.set"
              class="mt-1 text-sm text-slate-400"
            >
              {{ [card.number, card.set].filter(Boolean).join(' · ') }}
            </p>
            <p
              v-if="card.type || card.rarity"
              class="mt-1 text-sm capitalize text-slate-400"
            >
              {{ [card.type, card.rarity].filter(Boolean).join(' · ') }}
            </p>
            <p class="mt-2 text-sm font-semibold text-emerald-400">
              ×{{ card.quantity }} dans le deck
            </p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
