const RARITY_ALIASES: Record<string, string> = {
  common: 'common',
  commune: 'common',
  uncommon: 'uncommon',
  'peu commune': 'uncommon',
  rare: 'rare',
  'super rare': 'super-rare',
  'super-rare': 'super-rare',
  legendary: 'legendary',
  légendaire: 'legendary',
  legendaire: 'legendary',
  epic: 'epic',
  épique: 'epic',
  epique: 'epic',
  enchanted: 'enchanted',
  enchantée: 'enchanted',
  enchantee: 'enchanted',
  iconic: 'iconic',
  iconique: 'iconic',
  promo: 'promo',
}

/** Normalise une rareté API vers un slug CSS (common, super-rare, …). */
export function normalizeRarity(rarity: string | undefined): string {
  if (!rarity?.trim()) return 'unknown'
  const key = rarity.trim().toLowerCase()
  return RARITY_ALIASES[key] ?? key.replace(/\s+/g, '-')
}
