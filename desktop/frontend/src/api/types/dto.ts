/** Types aligned with Go `internal/app` JSON (camelCase tags). */

export interface KnowledgeDTO {
  id: string
  name: string
  description: string
  parentId: string | null
  cardCount: number
  dueCount: number
  createdAt: string
  updatedAt: string
  archivedAt: string | null
}

export interface KnowledgeTreeNode extends KnowledgeDTO {
  children: KnowledgeTreeNode[]
}

export interface TagDTO {
  id: string
  name: string
  slug: string
  color: string
  icon: string
  description: string
  cardCount: number
  createdAt: string
}

export interface SRSDTO {
  cardId: string
  fsrsState: string
  dueAt: string | null
  lastReviewAt: string | null
  stability: number | null
  difficulty: number | null
  reps: number
  lapses: number
}

export interface CardDTO {
  id: string
  knowledgeId: string
  knowledgeName: string
  sourceDocumentId: string | null
  parentId: string | null
  title: string
  cardType: string
  htmlPath: string
  htmlContent: string
  status: string
  tags: TagDTO[]
  srs: SRSDTO | null
  createdAt: string
  updatedAt: string
}

export interface CardSummaryDTO {
  id: string
  parentId: string | null
  title: string
  cardType: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface CardDetailDTO extends CardDTO {
  parent: CardSummaryDTO | null
  children: CardSummaryDTO[]
}

export interface ListCardsResult {
  items: CardDTO[]
  total: number
}

export interface CardFilters {
  knowledgeId?: string | null
  cardType?: string
  status?: string
  tagIds?: string[]
  keyword?: string
  parentId?: string | null
  isRoot?: boolean | null
  orderBy?: string
  orderDesc?: boolean
  limit?: number
  offset?: number
}

export interface CardWithSRSDTO {
  card: CardDTO
  srs: SRSDTO
}

export interface SRSStatisticsDTO {
  newCount: number
  learningCount: number
  reviewCount: number
  relearningCount: number
  totalCards: number
  dueToday: number
}

export interface ReviewLogDTO {
  id: string
  cardId: string
  reviewedAt: string
  rating: number
  reviewKind: string
}

export interface ReviewStatisticsDTO {
  TotalReviews: number
  AvgRating: number
  RatingCounts: Record<string, number>
  ReviewsByDay: Record<string, number>
}

export interface FSRSParameterDTO {
  id: string
  name: string
  parametersJson: string
  desiredRetention: number | null
  maximumInterval: number | null
}
