// Types miroir des DTOs exposés par l'API Go.

export type InkColor =
  | 'amber'
  | 'amethyst'
  | 'emerald'
  | 'ruby'
  | 'sapphire'
  | 'steel'

/** CardDTO — renvoyé par GET /cards */
export interface Card {
  id: string
  name: string
  title: string
  colors: InkColor[]
  number: string
  set: string
  rarity: string
  franchise: string
  filepath: string
}

/** OwnedCardDTO — carte possédée dans une collection */
export interface OwnedCard extends Card {
  quantity: number
}

/** PaginateResponse — enveloppe paginée de GET /cards */
export interface PaginatedCards {
  page: number
  docs: Card[]
  total: number
  size: number
}

/** Query params de GET /cards */
export interface PaginateParams {
  page: number
  limit: number
  search?: string
  color?: InkColor | ''
}

/** Couleurs d'encre disponibles (cf. internal/domain/ink.go) */
export const INK_COLORS: InkColor[] = [
  'amber',
  'amethyst',
  'emerald',
  'ruby',
  'sapphire',
  'steel',
]

/** OwnedCollectionDTO — renvoyé dans docs de GET /collections */
export interface OwnedCollection {
  id: string
  name: string
  cards: Record<string, OwnedCard>
}

/** PaginateResponse — enveloppe paginée de GET /collections */
export interface PaginatedCollections {
  page: number
  total: number
  docs: OwnedCollection[]
  size: number
}

/** Query params de GET /collections */
export interface ListCollectionsParams {
  page?: number
  limit?: number
  name?: string
}

/** CreateBodyRequest — body de POST /collections */
export interface CreateCollectionBody {
  name: string
}

/** Body de POST /collections/add-card */
export interface AddCardBody {
  collectionId: string
  cardId: string
  quantity: number
}

/** Niveau de deck souhaité — body de POST /deck/build */
export type DeckLevel = 'beginner' | 'intermediate' | 'advanced'

export const DECK_LEVELS: { value: DeckLevel; label: string }[] = [
  { value: 'beginner', label: 'Débutant' },
  { value: 'intermediate', label: 'Intermédiaire' },
  { value: 'advanced', label: 'Avancé' },
]

/** Body de POST /deck/build */
export interface BuildDeckBody {
  collection: string
  colors: InkColor[]
  level: string
}

/** Carte dans la réponse POST /deck/build */
export interface DeckCard {
  id: string
  name: string
  title: string
  number: string
  set: string
  type: string
  rarity: string
  filepath: string
  quantity: number
}

/** Réponse POST /deck/build */
export interface GeneratedDeck {
  name: string
  size: number
  strategy: string
  deck: DeckCard[]
}

/** Carte dans le body de POST /deck */
export interface SaveDeckCard {
  id: string
  quantity: number
}

/** Body de POST /deck */
export interface SaveDeckBody {
  name: string
  size: number
  strategy: string
  deck: SaveDeckCard[]
}

/** Deck sauvegardé — renvoyé dans docs de GET /deck */
export interface SavedDeck extends GeneratedDeck {
  id: string
}

/** Réponse GET /deck/export */
export interface DeckExportResponse {
  content: string
}

/** Réponse GET /deck */
export interface SavedDecksList {
  total: number
  docs: SavedDeck[]
}

/** Body de POST /auth/register */
export interface RegisterBody {
  name: string
  email: string
  password: string
  confirm_password: string
}

/** Body de POST /auth/login */
export interface LoginBody {
  email: string
  password: string
}
