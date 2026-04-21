import type { CardItem } from '@/modules/card/types'
import type { ExtractItem } from '@/modules/extract/types'
import type { ReviewGrade, ReviewItem, ReviewRecord, ReviewStats } from '../types'

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
