import type { CardWithSRSDTO } from '@/api/types/dto'
import type { CardItem } from '@/modules/card/types'
import type { ExtractItem } from '@/modules/extract/types'
import type { ReviewGrade, ReviewItem, ReviewRecord, ReviewStats, ReviewStatus } from '../types'

function isSameDay(dateIso: string, now: Date): boolean {
  const d = new Date(dateIso)
  return (
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate()
  )
}

export function normalizeReviewItem(input: ReviewItem): ReviewItem {
  return {
    ...input,
    sourceArticleId: input.sourceArticleId ?? null,
    sourceArticleTitle: input.sourceArticleTitle ?? null,
    extractId: input.extractId ?? null,
    nodeId: input.nodeId ?? null,
    cardId: input.cardId ?? null,
    lastReviewedAt: input.lastReviewedAt ?? null,
  }
}

export function buildReviewItemFromExtract(extract: ExtractItem, reviewId: string): ReviewItem {
  const now = new Date().toISOString()
  return normalizeReviewItem({
    id: reviewId,
    type: 'extract',
    title: extract.title,
    prompt: extract.quote,
    answer: extract.note?.trim() || '—',
    sourceArticleId: extract.sourceArticleId,
    sourceArticleTitle: extract.sourceArticleTitle,
    extractId: extract.id,
    nodeId: extract.treeNodeId ?? null,
    cardId: null,
    dueAt: now,
    status: 'new',
    reviewCount: 0,
  })
}

export function buildReviewItemFromCard(card: CardItem, reviewId: string): ReviewItem {
  const now = new Date().toISOString()
  return normalizeReviewItem({
    id: reviewId,
    type: 'card',
    title: card.title,
    prompt: card.prompt,
    answer: card.answer,
    sourceArticleId: null,
    sourceArticleTitle: null,
    extractId: null,
    nodeId: card.nodeId,
    cardId: card.id,
    dueAt: now,
    status: 'new',
    reviewCount: 0,
  })
}

export function buildReviewRecord(itemId: string, grade: ReviewGrade, reviewedAt: string): ReviewRecord {
  return {
    itemId,
    grade,
    reviewedAt,
  }
}

const FSRS_TO_UI_STATUS: Record<string, ReviewStatus> = {
  new: 'new',
  learning: 'learning',
  review: 'review',
  relearning: 'learning',
}

export function mapFsrsStateToReviewStatus(fsrsState: string): ReviewStatus {
  return FSRS_TO_UI_STATUS[fsrsState] ?? 'new'
}

export function mapReviewGradeToRating(grade: ReviewGrade): 1 | 2 | 3 | 4 {
  switch (grade) {
    case 'again':
      return 1
    case 'hard':
      return 2
    case 'good':
      return 3
    case 'easy':
      return 4
    default:
      return 3
  }
}

function toIso(value: string | Date | null | undefined): string {
  if (value == null) {
    return new Date().toISOString()
  }
  if (typeof value === 'string') {
    return value
  }
  return value.toISOString()
}

/** Map backend due-card row to the UI review queue item. */
export function mapCardWithSrsDtoToReviewItem(row: CardWithSRSDTO): ReviewItem {
  const { card, srs } = row
  const dueAt = srs.dueAt ? toIso(srs.dueAt) : toIso(card.updatedAt)
  const plain = card.htmlContent
    ? card.htmlContent.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim()
    : ''
  return normalizeReviewItem({
    id: card.id,
    type: 'card',
    title: card.title,
    prompt: plain.slice(0, 2000) || card.title,
    answer: '',
    sourceArticleId: card.sourceDocumentId,
    sourceArticleTitle: null,
    extractId: null,
    nodeId: card.id,
    cardId: card.id,
    dueAt,
    status: mapFsrsStateToReviewStatus(srs.fsrsState),
    lastReviewedAt: srs.lastReviewAt ? toIso(srs.lastReviewAt) : null,
    reviewCount: srs.reps ?? 0,
  })
}

export function buildReviewStats(input: {
  queueItems: ReviewItem[]
  records: ReviewRecord[]
}): ReviewStats {
  const now = new Date()
  const completedToday = input.records.filter((r) => isSameDay(r.reviewedAt, now)).length
  const againCount = input.records.filter((r) => r.grade === 'again').length
  return {
    completedToday,
    remaining: input.queueItems.length,
    againCount,
  }
}
