import type {
  AddCardBody,
  BuildDeckBody,
  CreateCollectionBody,
  GeneratedDeck,
  ListCollectionsParams,
  PaginatedCards,
  PaginateParams,
  PaginatedCollections,
} from '@/types/api'
import { normalizeGeneratedDeck, type RawGeneratedDeck } from '@/api/deck'

export const API_BASE = 'http://localhost:9090/api/v1'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  })
  if (!res.ok) {
    throw new Error(`Request failed: ${res.status} ${res.statusText}`)
  }
  if (res.status === 204) {
    return undefined as T
  }
  const text = await res.text()
  return (text ? JSON.parse(text) : undefined) as T
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

/** POST /collections/card */
export function addCard(body: AddCardBody): Promise<void> {
  return request<void>('/collections/card', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

/** POST /deck/build */
export async function buildDeck(body: BuildDeckBody): Promise<GeneratedDeck> {
  const raw = await request<RawGeneratedDeck>('/deck/build', {
    method: 'POST',
    body: JSON.stringify(body),
  })
  return normalizeGeneratedDeck(raw)
}

/** GET /collections/export — télécharge un CSV de toutes les collections */
export async function exportCollections(): Promise<void> {
  const res = await fetch(`${API_BASE}/collections/export`)
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
