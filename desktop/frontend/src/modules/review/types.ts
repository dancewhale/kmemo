export type ReviewItemType = 'extract' | 'card'
export type ReviewStatus = 'new' | 'learning' | 'review'
export type ReviewGrade = 'again' | 'hard' | 'good' | 'easy'

export interface ReviewItem {
  id: string
  type: ReviewItemType
  title: string
  prompt: string
  answer: string
  sourceArticleId?: string | null
  sourceArticleTitle?: string | null
  extractId?: string | null
  nodeId?: string | null
  cardId?: string | null
  dueAt: string
  status: ReviewStatus
  lastReviewedAt?: string | null
  reviewCount: number
}

export interface ReviewRecord {
  itemId: string
  grade: ReviewGrade
  reviewedAt: string
}

export interface ReviewStats {
  completedToday: number
  remaining: number
  againCount: number
}
