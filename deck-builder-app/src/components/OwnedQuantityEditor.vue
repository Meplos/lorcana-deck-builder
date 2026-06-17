<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { addCard } from '@/api/client'

const props = defineProps<{
  collectionId: string
  cardId: string
  ownedQuantity: number
  compact?: boolean
}>()

const queryClient = useQueryClient()
const quantity = ref(props.ownedQuantity)
const feedback = ref<{ type: 'success' | 'error'; message: string } | null>(null)

watch(
  () => props.ownedQuantity,
  (value) => {
    quantity.value = value
  },
)

const { mutate: save, isPending } = useMutation({
  mutationFn: addCard,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['collections'] })
    feedback.value = { type: 'success', message: 'Quantité mise à jour.' }
  },
  onError: () => {
    feedback.value = { type: 'error', message: 'Échec de la mise à jour.' }
  },
})

function submit() {
  if (quantity.value < 0) return
  feedback.value = null
  save({
    collectionId: props.collectionId,
    cardId: props.cardId,
    quantity: quantity.value,
  })
}
</script>

<template>
  <div :class="compact ? 'flex flex-col gap-2' : 'flex flex-col gap-3'">
    <p v-if="ownedQuantity > 0" class="text-xs text-slate-500 dark:text-slate-400">
      Possédé : <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ ownedQuantity }}</span>
    </p>

    <div class="flex items-end gap-2">
      <label class="flex-1">
        <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
          Quantité
        </span>
        <input
          v-model.number="quantity"
          type="number"
          min="0"
          class="w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-900 dark:text-white"
        />
      </label>
      <button
        type="button"
        :disabled="isPending || quantity < 0"
        class="rounded-xl bg-emerald-500 px-3 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
        @click="submit"
      >
        {{ isPending ? '…' : ownedQuantity > 0 ? 'Modifier' : 'Ajouter' }}
      </button>
    </div>

    <p
      v-if="feedback"
      class="text-xs"
      :class="feedback.type === 'success' ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-500'"
    >
      {{ feedback.message }}
    </p>
  </div>
</template>
