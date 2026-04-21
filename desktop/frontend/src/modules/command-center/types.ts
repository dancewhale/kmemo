import type { CommandId } from '@/shared/constants/commands'

export type CommandGroup =
  | 'navigation'
  | 'focus'
  | 'reading'
  | 'extract'
  | 'review'
  | 'knowledge'
  | 'capture'
  | 'system'

export interface CommandItem {
  id: CommandId
  title: string
  subtitle?: string
  group: CommandGroup
  keywords?: string[]
  shortcut?: string[]
  enabled?: boolean
}

export interface CommandMatchResult extends CommandItem {
  score: number
}

export interface CommandState {
  isOpen: boolean
  query: string
  activeIndex: number
}
