import type { DeckCard, GeneratedDeck } from '@/types/api'

type RawSchemaCard = {
  ID?: string
  Name?: string
  Title?: string
  Number?: string
  Set?: string
  Type?: string
  Rarity?: string
  id?: string
  name?: string
  title?: string
  number?: string
  set?: string
  type?: string
  rarity?: string
}

type RawDeckCard = {
  Card?: RawSchemaCard
  card?: RawSchemaCard
  FilePath?: string
  filepath?: string
  Quantity?: number
  quantity?: number
}

export type RawGeneratedDeck = {
  name: string
  size: number
  strategy: string
  deck: RawDeckCard[]
}

function normalizeDeckCard(raw: RawDeckCard): DeckCard {
  const card = raw.card ?? raw.Card ?? {}
  const id = card.id ?? card.ID ?? ''
  return {
    id,
    name: card.name ?? card.Name ?? '',
    title: card.title ?? card.Title ?? '',
    number: card.number ?? card.Number ?? '',
    set: card.set ?? card.Set ?? '',
    type: card.type ?? card.Type ?? '',
    rarity: card.rarity ?? card.Rarity ?? '',
    filepath:
      raw.filepath ??
      raw.FilePath ??
      `https://cdn.dreamborn.ink/images/fr/cards/${id}`,
    quantity: raw.quantity ?? raw.Quantity ?? 1,
  }
}

export function normalizeGeneratedDeck(raw: RawGeneratedDeck): GeneratedDeck {
  return {
    name: raw.name,
    size: raw.size,
    strategy: raw.strategy,
    deck: (raw.deck ?? []).map(normalizeDeckCard),
  }
}
