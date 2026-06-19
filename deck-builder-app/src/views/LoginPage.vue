<script lang="ts" setup>
import { ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { RouterLink, useRouter } from 'vue-router'
import { loginUser } from '@/api/client'

const router = useRouter()

const email = ref('')
const password = ref('')
const error = ref('')

function validate(): string | null {
  const trimmedEmail = email.value.trim()

  if (!trimmedEmail) return "L'email est obligatoire."
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(trimmedEmail)) return 'Email invalide.'
  if (!password.value) return 'Le mot de passe est obligatoire.'
  return null
}

const { mutate: login, isPending } = useMutation({
  mutationFn: loginUser,
  onSuccess: () => {
    router.push('/')
  },
  onError: (err) => {
    error.value =
      err instanceof Error
        ? err.message
        : 'Impossible de se connecter. Vérifie tes informations ou réessaie plus tard.'
  },
})

function submit() {
  const validationError = validate()
  if (validationError) {
    error.value = validationError
    return
  }
  error.value = ''
  login({
    email: email.value.trim(),
    password: password.value,
  })
}

const inputClass =
  'w-full rounded-xl border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 shadow-sm transition focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:placeholder:text-slate-400'
</script>

<template>
  <div class="flex min-h-[calc(100vh-4rem)] w-full items-center justify-center px-3 py-8 sm:px-6 lg:px-8">
    <form
      class="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800 sm:p-8"
      @submit.prevent="submit"
    >
      <header class="mb-6 text-center">
        <h1 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
          Se <span class="text-emerald-500">connecter</span>
        </h1>
        <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
          Accède à tes collections et decks Lorcana.
        </p>
      </header>

      <div class="flex flex-col gap-4">
        <label class="block">
          <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
            Email
          </span>
          <input
            v-model="email"
            type="email"
            autocomplete="email"
            placeholder="toi@exemple.com"
            :class="inputClass"
          />
        </label>

        <label class="block">
          <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
            Mot de passe
          </span>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            placeholder="••••••••••••"
            :class="inputClass"
          />
        </label>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-500">{{ error }}</p>

      <button
        type="submit"
        :disabled="isPending || !email.trim() || !password"
        class="mt-6 w-full rounded-xl bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ isPending ? 'Connexion…' : 'Se connecter' }}
      </button>

      <p class="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">
        Pas encore de compte ?
        <RouterLink to="/register" class="font-semibold text-emerald-600 hover:underline dark:text-emerald-400">
          Créer un compte
        </RouterLink>
      </p>
    </form>
  </div>
</template>
