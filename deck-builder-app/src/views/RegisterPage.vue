<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { RouterLink, useRouter } from 'vue-router'
import { registerUser } from '@/api/client'

const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')

const passwordChecks = computed(() => ({
  length: password.value.length >= 12,
  upper: /[A-Z]/.test(password.value),
  lower: /[a-z]/.test(password.value),
  digit: /\d/.test(password.value),
  special: /[^A-Za-z0-9]/.test(password.value),
  match:
    confirmPassword.value.length > 0 && password.value === confirmPassword.value,
}))

const allPasswordChecksOk = computed(() =>
  Object.values(passwordChecks.value).every(Boolean),
)

function validate(): string | null {
  const trimmedName = name.value.trim()
  const trimmedEmail = email.value.trim()

  if (!trimmedName) return 'Le nom est obligatoire.'
  if (!trimmedEmail) return "L'email est obligatoire."
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(trimmedEmail)) return 'Email invalide.'
  if (password.value !== confirmPassword.value) {
    return 'Les mots de passe ne correspondent pas.'
  }
  if (password.value.length < 12) {
    return 'Le mot de passe doit contenir au moins 12 caractères.'
  }
  if (!passwordChecks.value.upper) {
    return 'Le mot de passe doit contenir une majuscule.'
  }
  if (!passwordChecks.value.lower) {
    return 'Le mot de passe doit contenir une minuscule.'
  }
  if (!passwordChecks.value.digit) {
    return 'Le mot de passe doit contenir un chiffre.'
  }
  if (!passwordChecks.value.special) {
    return 'Le mot de passe doit contenir un caractère spécial.'
  }
  return null
}

const { mutate: register, isPending } = useMutation({
  mutationFn: registerUser,
  onSuccess: () => {
    router.push('/')
  },
  onError: (err) => {
    error.value =
      err instanceof Error
        ? err.message
        : 'Impossible de créer le compte. Vérifie tes informations ou réessaie plus tard.'
  },
})

function submit() {
  const validationError = validate()
  if (validationError) {
    error.value = validationError
    return
  }
  error.value = ''
  register({
    name: name.value.trim(),
    email: email.value.trim(),
    password: password.value,
    confirm_password: confirmPassword.value,
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
          Créer un <span class="text-emerald-500">compte</span>
        </h1>
        <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
          Rejoins le Lorcana Deck Builder pour sauvegarder tes collections et decks.
        </p>
      </header>

      <div class="flex flex-col gap-4">
        <label class="block">
          <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
            Nom
          </span>
          <input
            v-model="name"
            type="text"
            autocomplete="name"
            placeholder="Ton pseudo"
            :class="inputClass"
          />
        </label>

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
            autocomplete="new-password"
            placeholder="••••••••••••"
            :class="inputClass"
          />
        </label>

        <label class="block">
          <span class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400">
            Confirmer le mot de passe
          </span>
          <input
            v-model="confirmPassword"
            type="password"
            autocomplete="new-password"
            placeholder="••••••••••••"
            :class="inputClass"
          />
        </label>
      </div>

      <ul class="mt-4 space-y-1 text-xs">
        <li
          v-for="(check, key) in {
            length: 'Au moins 12 caractères',
            upper: 'Une majuscule',
            lower: 'Une minuscule',
            digit: 'Un chiffre',
            special: 'Un caractère spécial',
            match: 'Les mots de passe correspondent',
          }"
          :key="key"
          class="flex items-center gap-2"
          :class="passwordChecks[key as keyof typeof passwordChecks] ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-400'"
        >
          <svg
            class="size-3.5 shrink-0"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2.5"
            stroke="currentColor"
          >
            <path
              v-if="passwordChecks[key as keyof typeof passwordChecks]"
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m4.5 12.75 6 6 9-13.5"
            />
            <circle v-else cx="12" cy="12" r="3" />
          </svg>
          {{ check }}
        </li>
      </ul>

      <p v-if="error" class="mt-4 text-sm text-red-500">{{ error }}</p>

      <button
        type="submit"
        :disabled="isPending || !allPasswordChecksOk || !name.trim() || !email.trim()"
        class="mt-6 w-full rounded-xl bg-emerald-500 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ isPending ? 'Création du compte…' : 'Créer mon compte' }}
      </button>

      <p class="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">
        Déjà un compte ?
        <RouterLink to="/" class="font-semibold text-emerald-600 hover:underline dark:text-emerald-400">
          Retour à l'accueil
        </RouterLink>
      </p>
    </form>
  </div>
</template>
