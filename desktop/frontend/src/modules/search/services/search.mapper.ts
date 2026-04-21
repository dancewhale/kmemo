import type { ReaderArticle } from '@/modules/reader/types'
import type { ExtractItem } from '@/modules/extract/types'
import type { KnowledgeNode } from '@/modules/knowledge-tree/types'
import type { CardItem } from '@/modules/card/types'
import type { ReviewItem } from '@/modules/review/types'
import type { SearchDocument } from '../types'

export function mapArticlesToSearchDocuments(articles: ReaderArticle[]): SearchDocument[] {
  return articles.map((article) => ({
    id: article.id,
    entityType: 'article',
    title: article.title,
    subtitle: article.summary,
    content: `${article.summary}\n${article.content}`,
    articleId: article.id,
    sourceId: article.sourceUrl ?? null,
    sourceTitle: article.sourceType,
    updatedAt: article.updatedAt,
    keywords: [article.status, article.sourceType, ...(article.tags ?? [])],
  }))
}

export function mapExtractsToSearchDocuments(extracts: ExtractItem[]): SearchDocument[] {
  return extracts.map((extract) => ({
    id: extract.id,
    entityType: 'extract',
    title: extract.title,
    subtitle: extract.sourceArticleTitle,
    content: `${extract.quote}\n${extract.note}`,
    sourceId: extract.sourceArticleId,
    sourceTitle: extract.sourceArticleTitle,
    nodeId: extract.treeNodeId ?? null,
    articleId: extract.sourceArticleId,
    extractId: extract.id,
    updatedAt: extract.updatedAt,
    keywords: ['extract', 'quote', 'note'],
  }))
}

export function mapCardsToSearchDocuments(cards: CardItem[]): SearchDocument[] {
  return cards.map((card) => ({
    id: card.id,
    entityType: 'card' as const,
    title: card.title,
    subtitle: card.prompt.slice(0, 80),
    content: `${card.title}\n${card.prompt}\n${card.answer}`,
    nodeId: card.nodeId,
    cardId: card.id,
    updatedAt: card.updatedAt,
    keywords: ['card', 'prompt', 'answer'],
  }))
}

export function mapNodesToSearchDocuments(nodes: KnowledgeNode[]): SearchDocument[] {
  return nodes.map((node) => ({
    id: node.id,
    entityType: 'node',
    title: node.title,
    subtitle: `Knowledge ${node.type}`,
    content: `${node.title}\n${node.type}`,
    sourceId: node.parentId,
    sourceTitle: node.parentId ? `Parent: ${node.parentId}` : 'Root node',
    nodeId: node.id,
    articleId: node.sourceArticleId ?? null,
    extractId: node.sourceExtractId ?? null,
    updatedAt: node.updatedAt ?? null,
    keywords: [node.type, node.parentId ?? 'root'],
  }))
}

export function mapReviewsToSearchDocuments(reviews: ReviewItem[]): SearchDocument[] {
  return reviews.map((review) => ({
    id: review.id,
    entityType: 'review',
    title: review.title,
    subtitle: review.prompt,
    content: `${review.prompt}\n${review.answer}`,
    sourceId: review.sourceArticleId ?? null,
    sourceTitle: review.sourceArticleTitle ?? null,
    nodeId: review.nodeId ?? null,
    articleId: review.sourceArticleId ?? null,
    extractId: review.extractId ?? null,
    cardId: review.cardId ?? null,
    reviewId: review.id,
    updatedAt: review.lastReviewedAt ?? review.dueAt,
    keywords: [review.type, review.status],
  }))
}
