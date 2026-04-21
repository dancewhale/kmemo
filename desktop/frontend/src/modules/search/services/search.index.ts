import { mockArticles } from '@/mock/articles'
import { mockExtractItems } from '@/mock/extracts'
import { mockKnowledgeNodes } from '@/mock/tree'
import { mockReviewItems } from '@/mock/review'
import type { ReaderArticle } from '@/modules/reader/types'
import type { ExtractItem } from '@/modules/extract/types'
import type { KnowledgeNode } from '@/modules/knowledge-tree/types'
import type { ReviewItem } from '@/modules/review/types'
import type { CardItem } from '@/modules/card/types'
import {
  mapArticlesToSearchDocuments,
  mapCardsToSearchDocuments,
  mapExtractsToSearchDocuments,
  mapNodesToSearchDocuments,
  mapReviewsToSearchDocuments,
} from './search.mapper'
import type { SearchDocument } from '../types'

interface SearchIndexSource {
  articles: ReaderArticle[]
  extracts: ExtractItem[]
  nodes: KnowledgeNode[]
  cards: CardItem[]
  reviews: ReviewItem[]
}

export function buildSearchIndex(source?: Partial<SearchIndexSource>): SearchDocument[] {
  const articles = source?.articles ?? mockArticles
  const extracts = source?.extracts ?? mockExtractItems
  const nodes = source?.nodes ?? mockKnowledgeNodes
  const cards = source?.cards ?? []
  const reviews = source?.reviews ?? mockReviewItems

  return [
    ...mapArticlesToSearchDocuments(articles),
    ...mapExtractsToSearchDocuments(extracts),
    ...mapNodesToSearchDocuments(nodes),
    ...mapCardsToSearchDocuments(cards),
    ...mapReviewsToSearchDocuments(reviews),
  ]
}
