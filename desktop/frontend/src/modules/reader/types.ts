export type ReaderSourceType = 'web' | 'pdf' | 'text'

export type ReaderArticleStatus = 'inbox' | 'reading' | 'processed'

export interface ReaderArticle {
  id: string
  title: string
  summary: string
  content: string
  sourceType: ReaderSourceType
  sourceUrl?: string
  status: ReaderArticleStatus
  updatedAt: string
  tags?: string[]
}

export type ReaderStatusFilter = 'all' | ReaderArticleStatus

export interface ReaderFilterState {
  keyword: string
  status: ReaderStatusFilter
}
