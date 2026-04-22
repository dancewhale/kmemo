export type ContextEntityType = 'tree-node' | 'tree-card' | 'tree-extract' | 'article' | 'extract'

export interface ContextMenuContext {
  entityType: ContextEntityType
  entityId: string
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  x: number
  y: number
}

export interface ContextActionItem {
  id: string
  label: string
  icon?: string
  disabled?: boolean
  danger?: boolean
  dividerBefore?: boolean
}

export interface ContextMenuState {
  isOpen: boolean
  x: number
  y: number
  context: ContextMenuContext | null
}

/** Static row from registry (no runtime disabled). */
export interface ContextMenuRegistryRow {
  id: string
  label: string
  icon?: string
  danger?: boolean
  dividerBefore?: boolean
}
