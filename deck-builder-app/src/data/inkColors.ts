import type { InkColor } from '@/types/api'

export interface InkColorInfo {
  id: InkColor
  label: string
  archetype: string
  traits: string[]
}

export const INK_COLOR_INFO: Record<InkColor, InkColorInfo> = {
  amber: {
    id: 'amber',
    label: 'Ambre',
    archetype: 'Cœur',
    traits: ['Soutien et soins', 'Chants et gain de lore', 'Héros protecteurs'],
  },
  amethyst: {
    id: 'amethyst',
    label: 'Améthyste',
    archetype: 'Méfait',
    traits: ['Contrôle et renvoi', 'Échange et sorcellerie', 'Personnages vilains'],
  },
  emerald: {
    id: 'emerald',
    label: 'Émeraude',
    archetype: 'Aventure',
    traits: ['Mobilité et flexibilité', 'Shift et tempo', 'Personnages efficaces'],
  },
  ruby: {
    id: 'ruby',
    label: 'Rubis',
    archetype: 'Action',
    traits: ['Aggression et rush', 'Force et défi direct', 'Pression constante'],
  },
  sapphire: {
    id: 'sapphire',
    label: 'Saphir',
    archetype: 'Stratégie',
    traits: ['Pioche et objets', 'Encrage et contrôle', 'Manipulation du deck'],
  },
  steel: {
    id: 'steel',
    label: 'Acier',
    archetype: 'Force',
    traits: ['Personnages robustes', 'Bodyguard et résistance', 'Présence sur le plateau'],
  },
}

export const INK_COLOR_ACCENT: Record<
  InkColor,
  { border: string; bg: string; text: string; dot: string }
> = {
  amber: {
    border: 'border-amber-400 dark:border-amber-500',
    bg: 'bg-amber-50 dark:bg-amber-950/30',
    text: 'text-amber-800 dark:text-amber-300',
    dot: 'bg-amber-500',
  },
  amethyst: {
    border: 'border-violet-400 dark:border-violet-500',
    bg: 'bg-violet-50 dark:bg-violet-950/30',
    text: 'text-violet-800 dark:text-violet-300',
    dot: 'bg-violet-500',
  },
  emerald: {
    border: 'border-emerald-400 dark:border-emerald-500',
    bg: 'bg-emerald-50 dark:bg-emerald-950/30',
    text: 'text-emerald-800 dark:text-emerald-300',
    dot: 'bg-emerald-500',
  },
  ruby: {
    border: 'border-red-400 dark:border-red-500',
    bg: 'bg-red-50 dark:bg-red-950/30',
    text: 'text-red-800 dark:text-red-300',
    dot: 'bg-red-500',
  },
  sapphire: {
    border: 'border-blue-400 dark:border-blue-500',
    bg: 'bg-blue-50 dark:bg-blue-950/30',
    text: 'text-blue-800 dark:text-blue-300',
    dot: 'bg-blue-500',
  },
  steel: {
    border: 'border-slate-400 dark:border-slate-500',
    bg: 'bg-slate-100 dark:bg-slate-800/60',
    text: 'text-slate-700 dark:text-slate-300',
    dot: 'bg-slate-500',
  },
}
