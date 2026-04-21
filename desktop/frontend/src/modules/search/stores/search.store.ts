import { defineStore } from 'pinia'
import { router } from '@/app/router'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useCardStore } from '@/modules/card/stores/card.store'
import { buildSearchIndex } from '../services/search.index'
import { matchSearchDocuments } from '../services/search.matcher'
import type { SearchDocument, SearchFilterType, SearchResultItem } from '../types'

interface SearchState {
  query: string
  filter: SearchFilterType
  results: SearchResultItem[]
  selectedResultId: string | null
  loading: boolean
  initialized: boolean
  documents: SearchDocument[]
  focusToken: number
}

async function goContext(context: 'reading' | 'knowledge' | 'review' | 'search') {
  const workspace = useWorkspaceStore()
  workspace.setContext(context)
  await router.push({ name: ROUTE_NAMES[context] })
}

export const useSearchStore = defineStore('search', {
  state: (): SearchState => ({
    query: '',
    filter: 'all',
    results: [],
    selectedResultId: null,
    loading: false,
    initialized: false,
    documents: [],
    focusToken: 0,
  }),
  getters: {
    hasQuery(): boolean {
      return this.query.trim().length > 0
    },
    hasResults(): boolean {
      return this.results.length > 0
    },
    selectedResult(): SearchResultItem | null {
      if (!this.selectedResultId) {
        return null
      }
      return this.results.find((item) => item.id === this.selectedResultId) ?? null
    },
    resultCount(): number {
      return this.results.length
    },
  },
  actions: {
    refreshIndex() {
      const reader = useReaderStore()
      const tree = useTreeStore()
      const extract = useExtractStore()
      const review = useReviewStore()
      const card = useCardStore()
      this.documents = buildSearchIndex({
        articles: reader.articles,
        extracts: extract.items,
        nodes: tree.rawNodes,
        cards: card.items,
        reviews: review.items,
      })
    },
    async initialize() {
      if (this.initialized) {
        this.runSearch()
        return
      }
      this.loading = true
      const reader = useReaderStore()
      const tree = useTreeStore()
      const extract = useExtractStore()
      const review = useReviewStore()
      const card = useCardStore()
      await Promise.all([
        reader.initialize(),
        tree.initialize(),
        extract.initialize(),
        review.initialize(),
        card.initialize(),
      ])
      this.refreshIndex()
      this.initialized = true
      this.loading = false
      this.runSearch()
    },
    setQuery(query: string) {
      this.query = query
      this.runSearch()
    },
    setFilter(filter: SearchFilterType) {
      this.filter = filter
      this.runSearch()
    },
    runSearch() {
      this.refreshIndex()
      this.results = matchSearchDocuments(this.documents, this.query, this.filter)
      this.selectedResultId = this.results[0]?.id ?? null
    },
    setSelectedResult(id: string | null) {
      this.selectedResultId = id
    },
    requestFocus() {
      this.focusToken += 1
    },
    async openResult(result: SearchResultItem) {
      const reader = useReaderStore()
      const tree = useTreeStore()
      const extract = useExtractStore()
      const review = useReviewStore()
      const editor = useEditorStore()
      const card = useCardStore()

      if (result.entityType === 'article') {
        const articleId = result.articleId ?? result.id
        await reader.initialize()
        reader.openArticleById(articleId)
        const article = reader.getArticleById(articleId)
        if (article) {
          editor.openDocument({
            id: article.id,
            title: article.title,
            content: article.content,
            contentType: 'article',
            updatedAt: article.updatedAt,
            sourceUrl: article.sourceUrl,
            tags: article.tags,
          })
        }
        await goContext('reading')
      } else if (result.entityType === 'card') {
        const cardId = result.cardId ?? result.id
        await Promise.all([tree.initialize(), card.initialize()])
        const c = card.getCardById(cardId)
        if (c) {
          tree.setSelectedNode(c.nodeId)
          card.setSelectedCard(c.id)
        }
        await goContext('knowledge')
      } else if (result.entityType === 'extract') {
        const extractId = result.extractId ?? result.id
        await Promise.all([extract.initialize(), tree.initialize()])
        extract.openExtract(extractId)
        const nodeId = result.nodeId ?? extract.getExtractById(extractId)?.treeNodeId ?? null
        if (nodeId) {
          tree.setSelectedNode(nodeId)
        }
        await goContext('knowledge')
      } else if (result.entityType === 'node') {
        const nodeId = result.nodeId ?? result.id
        await Promise.all([tree.initialize(), extract.initialize()])
        tree.setSelectedNode(nodeId)
        const node = tree.findNodeById(nodeId)
        if (node?.sourceExtractId) {
          extract.openExtract(node.sourceExtractId)
        }
        await goContext('knowledge')
      } else if (result.entityType === 'review') {
        const reviewId = result.reviewId ?? result.id
        await review.initialize()
        review.setSelectedItem(reviewId)
        await goContext('review')
      }
    },
    clearSearch() {
      this.query = ''
      this.filter = 'all'
      this.runSearch()
    },
  },
})
