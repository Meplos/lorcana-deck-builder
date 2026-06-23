import { fetchCollections } from '@/api/client'
import { isApiError } from '@/api/errors'
import {
  markAuthenticated,
  markSessionInitialized,
  markUnauthenticated,
  useAuthSession,
} from '@/auth/session'

let initPromise: Promise<boolean> | null = null

export function resetSessionProbe() {
  initPromise = null
}

export async function ensureSession(): Promise<boolean> {
  const { isAuthenticated, isInitialized } = useAuthSession()

  if (isInitialized.value) {
    return isAuthenticated.value
  }

  if (!initPromise) {
    initPromise = (async () => {
      try {
        await fetchCollections({ limit: 1 })
        markAuthenticated()
      } catch (error) {
        if (isApiError(error) && error.status === 403) {
          markUnauthenticated()
        } else {
          markUnauthenticated()
        }
      } finally {
        markSessionInitialized()
      }

      return useAuthSession().isAuthenticated.value
    })()
  }

  return initPromise
}
