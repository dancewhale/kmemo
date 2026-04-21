import { defineStore } from 'pinia'
import { mockArticles } from '@/mock/articles'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import type { ReaderArticle, ReaderArticleStatus, ReaderFilterState, ReaderStatusFilter } from '../types'

interface ReaderState {
  articles: ReaderArticle[]
  selectedArticleId: string | null
  filter: ReaderFilterState
  loading: boolean
  initialized: boolean
}

function matchesKeyword(a: ReaderArticle, keyword: string): boolean {
  const k = keyword.trim().toLowerCase()
  if (!k) {
    return true
  }
  return a.title.toLowerCase().includes(k) || a.summary.toLowerCase().includes(k)
}

function matchesStatus(a: ReaderArticle, status: ReaderStatusFilter): boolean {
  if (status === 'all') {
    return true
  }
  return a.status === status
}

export const useReaderStore = defineStore('reader', {
  state: (): ReaderState => ({
    articles: [],
    selectedArticleId: null,
    filter: { keyword: '', status: 'all' },
    loading: false,
    initialized: false,
  }),

  getters: {
    filteredArticles(): ReaderArticle[] {
      return this.articles.filter(
        (a) => matchesStatus(a, this.filter.status) && matchesKeyword(a, this.filter.keyword),
      )
    },

    selectedArticle(): ReaderArticle | null {
      if (!this.selectedArticleId) {
        return null
      }
      return this.articles.find((a) => a.id === this.selectedArticleId) ?? null
    },

    inboxArticles(): ReaderArticle[] {
      return this.articles.filter((a) => a.status === 'inbox')
    },

    readingArticles(): ReaderArticle[] {
      return this.articles.filter((a) => a.status === 'reading')
    },

    processedArticles(): ReaderArticle[] {
      return this.articles.filter((a) => a.status === 'processed')
    },
    getArticleById(): (id: string) => ReaderArticle | null {
      return (id: string) => this.articles.find((a) => a.id === id) ?? null
    },
  },

  actions: {
    async initialize() {
      if (this.initialized) {
        return
      }
      this.loading = true
      await Promise.resolve()
      this.articles = mockArticles.map((a) => ({ ...a }))
      this.initialized = true
      this.loading = false
    },

    setSelectedArticle(id: string | null) {
      this.selectedArticleId = id
      useWorkspaceStore().selectArticle(id)
    },

    openArticleById(id: string) {
      const found = this.articles.find((a) => a.id === id)
      if (!found) {
        return
      }
      this.setSelectedArticle(found.id)
    },

    addArticle(article: ReaderArticle) {
      this.articles.unshift({ ...article })
    },

    updateArticleStatus(id: string, status: ReaderArticleStatus) {
      this.setArticleStatus(id, status)
    },

    setKeyword(keyword: string) {
      this.filter.keyword = keyword
      this.ensureSelectionVisibleInFilter()
    },

    setStatusFilter(status: ReaderStatusFilter) {
      this.filter.status = status
      this.ensureSelectionVisibleInFilter()
    },

    clearFilters() {
      this.filter.keyword = ''
      this.filter.status = 'all'
      this.ensureSelectionVisibleInFilter()
    },

    ensureSelectionVisibleInFilter() {
      const list = this.filteredArticles
      const sel = this.selectedArticleId
      if (sel && !list.some((a) => a.id === sel)) {
        this.setSelectedArticle(list[0]?.id ?? null)
      }
    },

    async ensureReadyWithSelection(preferredId: string | null) {
      await this.initialize()
      const list = this.filteredArticles
      const exists = preferredId && list.some((a) => a.id === preferredId)
      if (exists && preferredId) {
        this.setSelectedArticle(preferredId)
        return
      }
      const first = list[0] ?? this.articles[0]
      this.setSelectedArticle(first?.id ?? null)
    },

    setArticleStatus(id: string, status: ReaderArticleStatus) {
      const a = this.articles.find((x) => x.id === id)
      if (!a) {
        return
      }
      a.status = status
      a.updatedAt = new Date().toISOString()
      this.ensureSelectionVisibleInFilter()
    },
  },
})
