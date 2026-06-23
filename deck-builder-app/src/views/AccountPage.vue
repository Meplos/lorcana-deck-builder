<script lang="ts" setup>
import { ref } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { logoutUser } from '@/api/client'
import { resetSessionProbe } from '@/auth/bootstrap'
import { markUnauthenticated } from '@/auth/session'

const router = useRouter()
const queryClient = useQueryClient()
const error = ref('')

const { mutate: logout, isPending } = useMutation({
  mutationFn: logoutUser,
  onSuccess: async () => {
    error.value = ''
    markUnauthenticated()
    resetSessionProbe()
    await queryClient.clear()
    router.push('/login')
  },
  onError: (err) => {
    error.value =
      err instanceof Error ? err.message : 'Impossible de se déconnecter.'
  },
})
</script>

<template>
  <div class="flex min-h-screen w-full flex-col gap-6 px-3 py-6 sm:gap-8 sm:px-6 sm:py-8 lg:px-8">
    <header>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
        Mon <span class="text-emerald-500">compte</span>
      </h1>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
        Paramètres et session.
      </p>
    </header>

    <section class="max-w-lg rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800">
      <h2 class="text-sm font-semibold text-slate-900 dark:text-white">Session</h2>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
        Déconnecte-toi de ton compte sur cet appareil.
      </p>

      <p v-if="error" class="mt-4 text-sm text-red-500">{{ error }}</p>

      <button
        type="button"
        :disabled="isPending"
        class="mt-4 rounded-xl border border-red-300 px-4 py-2.5 text-sm font-semibold text-red-600 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60 dark:border-red-900 dark:text-red-400 dark:hover:bg-red-950/30"
        @click="logout()"
      >
        {{ isPending ? 'Déconnexion…' : 'Se déconnecter' }}
      </button>
    </section>
  </div>
</template>
