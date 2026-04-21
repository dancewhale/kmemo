import { defineStore } from 'pinia'
import { createEntityId } from '@/shared/utils/id'
import { useToast } from '@/shared/composables/useToast'
import { mockExtractItems } from '@/mock/extracts'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { buildReviewItemFromExtract } from '@/modules/review/services/review.mapper'
import { buildTreeNodeFromExtract } from '../services/extract.mapper'
import type { ExtractItem, ExtractUpdatePayload, PendingExtractPayload } from '../types'

interface ExtractState {
  items: ExtractItem[]
  pending: PendingExtractPayload | null
  isCreateDialogOpen: boolean
  creating: boolean
  selectedExtractId: string | null
  saving: boolean
  initialized: boolean
}

export const useExtractStore = defineStore('extract', {
  state: (): ExtractState => ({
    items: [],
    pending: null,
    isCreateDialogOpen: false,
    creating: false,
    selectedExtractId: null,
    saving: false,
    initialized: false,
  }),

  getters: {
    selectedExtract(): ExtractItem | null {
      if (!this.selectedExtractId) {
        return null
      }
      return this.items.find((item) => item.id === this.selectedExtractId) ?? null
    },
    pendingQuote(): string {
      return this.pending?.quote ?? ''
    },
    hasPendingExtract(): boolean {
      return this.pending != null
    },
    extractCount(): number {
      return this.items.length
    },
    getExtractsByArticleId(): (articleId: string) => ExtractItem[] {
      return (articleId: string) =>
        this.items.filter((item) => item.sourceArticleId === articleId).sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))
    },
    getExtractById(): (id: string) => ExtractItem | null {
      return (id: string) => this.items.find((item) => item.id === id) ?? null
    },
  },

  actions: {
    async initialize() {
      if (this.initialized) {
        return
      }
      this.items = mockExtractItems.map((x) => ({ ...x }))
      this.initialized = true
    },

    openCreateDialog(payload: PendingExtractPayload) {
      this.pending = { ...payload }
      this.isCreateDialogOpen = true
    },

    closeCreateDialog() {
      this.isCreateDialogOpen = false
      this.pending = null
    },

    updatePendingField<K extends keyof PendingExtractPayload>(
      key: K,
      value: PendingExtractPayload[K],
    ) {
      if (!this.pending) {
        return
      }
      this.pending = {
        ...this.pending,
        [key]: value,
      }
    },

    async createExtract(): Promise<ExtractItem | null> {
      if (!this.pending || this.creating) {
        return null
      }
      this.creating = true
      await new Promise((resolve) => window.setTimeout(resolve, 260))
      const now = new Date().toISOString()
      const item: ExtractItem = {
        id: `ext-${Date.now().toString(36)}`,
        sourceArticleId: this.pending.sourceArticleId,
        sourceArticleTitle: this.pending.sourceArticleTitle,
        title: this.pending.title.trim() || this.pending.quote.slice(0, 32),
        quote: this.pending.quote.trim(),
        note: this.pending.note.trim(),
        parentNodeId: this.pending.parentNodeId,
        treeNodeId: null,
        createdAt: now,
        updatedAt: now,
      }
      const node = buildTreeNodeFromExtract(item)
      const itemWithNode: ExtractItem = {
        ...item,
        treeNodeId: node.id,
      }
      this.items.unshift(itemWithNode)

      const treeStore = useTreeStore()
      treeStore.addExtractNode({
        parentNodeId: itemWithNode.parentNodeId,
        title: itemWithNode.title,
        sourceArticleId: itemWithNode.sourceArticleId,
        sourceExtractId: itemWithNode.id,
        nodeId: node.id,
      })
      this.selectedExtractId = itemWithNode.id

      this.creating = false
      this.closeCreateDialog()
      return itemWithNode
    },

    setItems(items: ExtractItem[]) {
      this.items = items.map((x) => ({ ...x }))
    },

    setSelectedExtract(id: string | null) {
      this.selectedExtractId = id
    },

    openExtract(id: string) {
      const found = this.items.find((item) => item.id === id)
      this.selectedExtractId = found ? found.id : null
    },

    async updateExtract(id: string, payload: ExtractUpdatePayload) {
      const idx = this.items.findIndex((item) => item.id === id)
      if (idx < 0) {
        return
      }
      this.saving = true
      await new Promise((resolve) => window.setTimeout(resolve, 180))
      const prev = this.items[idx]
      const updated: ExtractItem = {
        ...prev,
        title: payload.title ?? prev.title,
        note: payload.note ?? prev.note,
        quote: payload.quote ?? prev.quote,
        updatedAt: new Date().toISOString(),
      }
      this.items.splice(idx, 1, updated)
      this.saving = false
      if (updated.treeNodeId && updated.title !== prev.title) {
        useTreeStore().updateNodeTitle(updated.treeNodeId, updated.title)
      }
    },

    getExtractsForArticle(articleId: string): ExtractItem[] {
      return this.items
        .filter((item) => item.sourceArticleId === articleId)
        .sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))
    },

    clearPending() {
      this.pending = null
    },

    addExtract(item: ExtractItem) {
      this.items.unshift({ ...item })
      this.selectedExtractId = item.id
    },

    attachReviewItem(extractId: string, reviewItemId: string) {
      const idx = this.items.findIndex((item) => item.id === extractId)
      if (idx < 0) {
        return
      }
      this.items[idx] = { ...this.items[idx], reviewItemId }
    },

    getReviewItemIdByExtractId(extractId: string): string | null {
      return this.items.find((item) => item.id === extractId)?.reviewItemId ?? null
    },

    async addOrOpenReviewFromSelectedExtract(options?: { openReview?: boolean }): Promise<void> {
      const toast = useToast()
      const ex = this.selectedExtract
      if (!ex) {
        return
      }
      const review = useReviewStore()
      await review.initialize()
      const existing = review.getReviewItemByExtractId(ex.id)
      if (existing) {
        if (!ex.reviewItemId) {
          this.attachReviewItem(ex.id, existing.id)
        }
        toast.info('Opened existing review item')
        if (options?.openReview) {
          await review.openReviewItemInWorkspace(existing.id)
        }
        return
      }
      const revId = createEntityId('rev')
      const item = buildReviewItemFromExtract(ex, revId)
      review.addReviewItem(item, { select: false })
      this.attachReviewItem(ex.id, revId)
      toast.success('Review item created from extract')
      if (options?.openReview) {
        await review.openReviewItemInWorkspace(revId)
      }
    },

    async openLinkedReviewFromSelected(): Promise<boolean> {
      const ex = this.selectedExtract
      if (!ex?.reviewItemId) {
        return false
      }
      const review = useReviewStore()
      await review.initialize()
      await review.openReviewItemInWorkspace(ex.reviewItemId)
      return true
    },

    removeExtractById(id: string) {
      const idx = this.items.findIndex((item) => item.id === id)
      if (idx < 0) {
        return
      }
      const removed = this.items[idx]
      const treeNodeId = removed.treeNodeId
      this.items.splice(idx, 1)
      if (treeNodeId) {
        useTreeStore().removeNodeSubtree(treeNodeId)
      }
      if (this.selectedExtractId === id) {
        this.selectedExtractId = this.items[0]?.id ?? null
      }
    },
  },
})
