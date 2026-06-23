<script lang="ts" setup>
import { onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps<{
  open: boolean
  deckName: string
  content: string
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  close: []
}>()

const copied = ref(false)
let copiedTimeout: ReturnType<typeof setTimeout> | undefined

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

watch(
  () => props.open,
  (isOpen) => {
    if (typeof document === 'undefined') return
    if (isOpen) {
      document.addEventListener('keydown', onKeydown)
      document.body.style.overflow = 'hidden'
      copied.value = false
    } else {
      document.removeEventListener('keydown', onKeydown)
      document.body.style.overflow = ''
    }
  },
)

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
  if (copiedTimeout) clearTimeout(copiedTimeout)
})

async function copyToClipboard() {
  if (!props.content) return
  try {
    await navigator.clipboard.writeText(props.content)
    copied.value = true
    if (copiedTimeout) clearTimeout(copiedTimeout)
    copiedTimeout = setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    copied.value = false
  }
}
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
        v-if="open"
        class="fixed inset-0 z-50 overflow-y-auto bg-black/70 p-4 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        :aria-label="`Export texte — ${deckName}`"
        @click.self="emit('close')"
      >
        <div class="flex min-h-full items-center justify-center py-6 sm:py-8">
          <div
            class="relative flex w-full max-w-lg flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-2xl dark:border-slate-700 dark:bg-slate-800"
          >
            <button
              type="button"
              aria-label="Fermer"
              class="absolute right-4 top-4 flex size-9 items-center justify-center rounded-full text-slate-500 transition hover:bg-slate-100 hover:text-slate-700 dark:text-slate-400 dark:hover:bg-slate-700 dark:hover:text-slate-200"
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

            <div class="pr-10">
              <h2 class="text-lg font-semibold text-slate-900 dark:text-white">
                Export texte
              </h2>
              <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                {{ deckName }}
              </p>
            </div>

            <div
              v-if="loading"
              class="rounded-xl border border-dashed border-slate-300 p-8 text-center text-sm text-slate-500 dark:border-slate-600 dark:text-slate-400"
            >
              Chargement…
            </div>

            <p v-else-if="error" class="text-sm text-red-500">{{ error }}</p>

            <template v-else>
              <pre
                class="max-h-80 overflow-auto rounded-xl border border-slate-200 bg-slate-50 p-4 text-sm leading-relaxed text-slate-800 dark:border-slate-600 dark:bg-slate-900 dark:text-slate-200"
              >{{ content }}</pre>

              <button
                type="button"
                :disabled="!content"
                class="self-start rounded-xl bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
                @click="copyToClipboard"
              >
                {{ copied ? 'Copié !' : 'Copier dans le presse-papier' }}
              </button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
