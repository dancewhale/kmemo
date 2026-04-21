import { defineStore } from 'pinia'
import { router } from '@/app/router'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { isWailsAvailable } from '@/api/wails'
import { mockReviewItems } from '@/mock/review'
import { useToast } from '@/shared/composables/useToast'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import * as reviewRepository from '../services/review.repository'
import {
  buildReviewRecord,
  buildReviewStats,
  normalizeReviewItem,
} from '../services/review.mapper'
import type { ExtractItem } from '@/modules/extract/types'
import type { ReviewGrade, ReviewItem, ReviewRecord, ReviewStats } from '../types'

interface ReviewState {
  items: ReviewItem[]
  selectedReviewItemId: string | null
  records: ReviewRecord[]
  loading: boolean
  submitting: boolean
  sessionCompletedIds: string[]
  initialized: boolean
}

export const useReviewStore = defineStore('review', {
  state: (): ReviewState => ({
    items: [],
    selectedReviewItemId: null,
    records: [],
    loading: false,
    submitting: false,
    sessionCompletedIds: [],
    initialized: false,
  }),

  getters: {
    queueItems(): ReviewItem[] {
      const done = new Set(this.sessionCompletedIds)
      return this.items.filter((item) => !done.has(item.id))
    },

    selectedItem(): ReviewItem | null {
      if (!this.selectedReviewItemId) {
        return null
      }
      return this.items.find((item) => item.id === this.selectedReviewItemId) ?? null
    },

    remainingCount(): number {
      return this.queueItems.length
    },

    completedTodayCount(): number {
      return buildReviewStats({ queueItems: this.queueItems, records: this.records }).completedToday
    },

    againCount(): number {
      return buildReviewStats({ queueItems: this.queueItems, records: this.records }).againCount
    },

    stats(): ReviewStats {
      return buildReviewStats({ queueItems: this.queueItems, records: this.records })
    },

    isQueueEmpty(): boolean {
      return this.queueItems.length === 0
    },
  },

  actions: {
    async initialize(opts?: { force?: boolean }) {
      if (this.initialized && !opts?.force) {
        return
      }
      this.loading = true
      const toast = useToast()
      try {
        if (isWailsAvailable()) {
          const res = await reviewRepository.fetchDueReviewItems(null, 80)
          if (res.ok) {
            this.items = res.data
            this.sessionCompletedIds = []
          } else {
            toast.warning(res.error.message || '无法加载复习队列')
            this.items = mockReviewItems.map((x) => normalizeReviewItem(x))
          }
        } else {
          await Promise.resolve()
          this.items = mockReviewItems.map((x) => normalizeReviewItem(x))
        }
      } finally {
        this.initialized = true
        this.loading = false
        this.openFirstAvailable()
      }
    },

    setSelectedItem(id: string | null) {
      this.selectedReviewItemId = id
      useWorkspaceStore().selectReview(id)
    },

    async submitGrade(grade: ReviewGrade) {
      const current = this.selectedItem
      if (!current || this.submitting) {
        return
      }
      this.submitting = true
      const toast = useToast()
      const now = new Date().toISOString()
      try {
        if (isWailsAvailable() && current.type === 'card' && current.cardId) {
          const res = await reviewRepository.submitHostReview(current.cardId, grade)
          if (!res.ok) {
            toast.warning(res.error.message || '提交复习失败')
            return
          }
        } else {
          await new Promise((resolve) => window.setTimeout(resolve, 220))
        }

        this.records.unshift(buildReviewRecord(current.id, grade, now))

        const idx = this.items.findIndex((item) => item.id === current.id)
        if (idx >= 0) {
          this.items[idx] = {
            ...this.items[idx],
            reviewCount: this.items[idx].reviewCount + 1,
            lastReviewedAt: now,
          }
        }

        if (!this.sessionCompletedIds.includes(current.id)) {
          this.sessionCompletedIds.push(current.id)
        }

        this.openNext()
      } finally {
        this.submitting = false
      }
    },

    openFirstAvailable() {
      const first = this.queueItems[0]
      this.setSelectedItem(first?.id ?? null)
    },

    openNext() {
      const queue = this.queueItems
      if (queue.length === 0) {
        this.setSelectedItem(null)
        return
      }
      const currentIndex = queue.findIndex((item) => item.id === this.selectedReviewItemId)
      if (currentIndex < 0 || currentIndex + 1 >= queue.length) {
        this.setSelectedItem(queue[0].id)
        return
      }
      this.setSelectedItem(queue[currentIndex + 1].id)
    },

    getItemById(id: string): ReviewItem | null {
      return this.items.find((item) => item.id === id) ?? null
    },

    getReviewItemByCardId(cardId: string, nodeId?: string | null): ReviewItem | null {
      const byCard = this.items.find((item) => item.cardId === cardId)
      if (byCard) {
        return byCard
      }
      if (nodeId) {
        return (
          this.items.find((item) => item.type === 'card' && item.nodeId === nodeId && !item.extractId) ?? null
        )
      }
      return null
    },

    getReviewItemByExtractId(extractId: string): ReviewItem | null {
      return this.items.find((item) => item.extractId === extractId) ?? null
    },

    async openReviewItemInWorkspace(reviewItemId: string): Promise<void> {
      const item = this.getItemById(reviewItemId)
      if (!item) {
        return
      }
      useWorkspaceStore().setContext('review')
      this.setSelectedItem(reviewItemId)
      await router.push({ name: ROUTE_NAMES.review })
    },

    async ensureReadyWithSelection(preferredId: string | null) {
      await this.initialize()
      const queue = this.queueItems
      const preferred = preferredId ? queue.find((x) => x.id === preferredId) : null
      this.setSelectedItem(preferred?.id ?? queue[0]?.id ?? null)
    },

    resetSession() {
      this.records = []
      this.sessionCompletedIds = []
      this.openFirstAvailable()
    },

    addMockReviewFromExtract(extract: ExtractItem) {
      const id = `rev-ctx-${Date.now().toString(36)}`
      const now = new Date().toISOString()
      const item: ReviewItem = normalizeReviewItem({
        id,
        type: 'extract',
        title: `Review: ${extract.title}`,
        prompt: extract.quote.slice(0, 200),
        answer: extract.note?.trim() || extract.quote,
        sourceArticleId: extract.sourceArticleId,
        sourceArticleTitle: extract.sourceArticleTitle,
        extractId: extract.id,
        nodeId: extract.treeNodeId ?? null,
        dueAt: now,
        status: 'new',
        reviewCount: 0,
      })
      this.items.unshift(item)
    },

    addMockReviewFromTreeNode(nodeId: string, title: string) {
      const id = `rev-node-${Date.now().toString(36)}`
      const now = new Date().toISOString()
      const item: ReviewItem = normalizeReviewItem({
        id,
        type: 'card',
        title: `Tree: ${title}`,
        prompt: `Review knowledge node “${title}”.`,
        answer: '',
        sourceArticleId: null,
        sourceArticleTitle: null,
        extractId: null,
        nodeId,
        dueAt: now,
        status: 'new',
        reviewCount: 0,
      })
      this.items.unshift(item)
    },

    addReviewItem(item: ReviewItem, options?: { select?: boolean }) {
      const normalized = normalizeReviewItem(item)
      this.items.unshift(normalized)
      const i = this.sessionCompletedIds.indexOf(normalized.id)
      if (i >= 0) {
        this.sessionCompletedIds.splice(i, 1)
      }
      if (options?.select) {
        this.setSelectedItem(normalized.id)
      }
    },

    removeItemFromSessionQueue(itemId: string) {
      if (!this.sessionCompletedIds.includes(itemId)) {
        this.sessionCompletedIds.push(itemId)
      }
      if (this.selectedReviewItemId === itemId) {
        this.openNext()
      }
    },
  },
})
