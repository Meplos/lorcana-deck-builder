import { readonly, ref } from 'vue'

const isAuthenticated = ref(false)
const isInitialized = ref(false)

export function markAuthenticated() {
  isAuthenticated.value = true
}

export function markUnauthenticated() {
  isAuthenticated.value = false
}

export function markSessionInitialized() {
  isInitialized.value = true
}

export function useAuthSession() {
  return {
    isAuthenticated: readonly(isAuthenticated),
    isInitialized: readonly(isInitialized),
    markAuthenticated,
    markUnauthenticated,
  }
}
