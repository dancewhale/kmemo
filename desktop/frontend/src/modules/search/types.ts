export type SearchEntityType = 'article' | 'extract' | 'node' | 'card' | 'review'
export type SearchFilterType = 'all' | SearchEntityType

export interface SearchDocument {
  id: string
  entityType: SearchEntityType
  title: string
  content: string
  subtitle?: string
  sourceId?: string | null
  sourceTitle?: string | null
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  cardId?: string | null
  updatedAt?: string | null
  keywords?: string[]
}

export interface SearchResultItem {
  id: string
  entityType: SearchEntityType
  title: string
  snippet: string
  subtitle?: string
  score: number
  sourceId?: string | null
  sourceTitle?: string | null
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  cardId?: string | null
  updatedAt?: string | null
}
