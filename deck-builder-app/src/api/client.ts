import type {
  AddCardBody,
  BuildDeckBody,
  CreateCollectionBody,
  GeneratedDeck,
  ListCollectionsParams,
  PaginatedCards,
  PaginateParams,
  PaginatedCollections,
  RegisterBody,
  SaveDeckBody,
  SavedDecksList,
} from '@/types/api'

export const API_BASE = import.meta.env.VITE_API_BASE ?? '/api/v1'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  })
  const text = await res.text()
  if (!res.ok) {
    let message = `Request failed: ${res.status} ${res.statusText}`
    if (text) {
      try {
        const body = JSON.parse(text) as { error?: string }
        if (body.error) message = body.error
      } catch {
        // ignore malformed error body
      }
    }
    throw new Error(message)
  }
  if (res.status === 204 || res.status === 202) {
    return undefined as T
  }
  return (text ? JSON.parse(text) : undefined) as T
}

/** POST /auth/register */
export function registerUser(body: RegisterBody): Promise<void> {
  return request<void>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/** GET /cards */
export function fetchCards(params: PaginateParams): Promise<PaginatedCards> {
  const query = new URLSearchParams({
    page: String(params.page),
    limit: String(params.limit),
  })
  if (params.search) query.set('search', params.search)
  if (params.color) query.set('color', params.color)
  return request<PaginatedCards>(`/cards?${query}`)
}

/** GET /collections */
export function fetchCollections(
  params: ListCollectionsParams = {},
): Promise<PaginatedCollections> {
  const query = new URLSearchParams()
  if (params.page) query.set('page', String(params.page))
  if (params.limit) query.set('limit', String(params.limit))
  if (params.name) query.set('name', params.name)
  const qs = query.toString()
  return request<PaginatedCollections>(`/collections${qs ? `?${qs}` : ''}`)
}

/** POST /collections */
export function createCollection(body: CreateCollectionBody): Promise<void> {
  return request<void>('/collections', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/** GET /collections/add-card */
export function addCard(body: AddCardBody): Promise<void> {
  const query = new URLSearchParams({
    collectionId: body.collectionId,
    cardId: body.cardId,
    quantity: String(body.quantity),
  })
  return request<void>(`/collections/add-card?${query}`)
}

/** POST /deck/build */
export function buildDeck(body: BuildDeckBody): Promise<GeneratedDeck> {
  return request<GeneratedDeck>('/deck/build', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/** POST /deck */
export function saveDeck(body: SaveDeckBody): Promise<void> {
  return request<void>('/deck', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/** GET /deck */
export function fetchDecks(): Promise<SavedDecksList> {
  return request<SavedDecksList>('/deck')
}

/** GET /collections/export — télécharge un CSV de toutes les collections */
export async function exportCollections(): Promise<void> {
  const res = await fetch(`${API_BASE}/collections/export`, {
    credentials: 'include',
  })
  if (!res.ok) {
    throw new Error(`Request failed: ${res.status} ${res.statusText}`)
  }
  const blob = await res.blob()
  const disposition = res.headers.get('Content-Disposition')
  const filename = disposition?.match(/filename="(.+)"/)?.[1] ?? 'export.csv'
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}
