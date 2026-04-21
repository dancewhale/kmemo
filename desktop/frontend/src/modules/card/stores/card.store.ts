import { defineStore } from 'pinia'
import { isWailsAvailable } from '@/api/wails'
import { createEntityId } from '@/shared/utils/id'
import { useToast } from '@/shared/composables/useToast'
import * as cardRepository from '../services/card.repository'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { buildReviewItemFromCard } from '@/modules/review/services/review.mapper'
import { buildCardItemFromTreeNode } from '../services/card.mapper'
import type { CardItem, CardUpdatePayload } from '../types'

interface CardState {
  items: CardItem[]
  selectedCardId: string | null
  saving: boolean
  initialized: boolean
}

export const useCardStore = defineStore('card', {
  state: (): CardState => ({
    items: [],
    selectedCardId: null,
    saving: false,
    initialized: false,
  }),

  getters: {
    selectedCard(): CardItem | null {
      if (!this.selectedCardId) {
        return null
      }
      return this.items.find((c) => c.id === this.selectedCardId) ?? null
    },

    getCardById(): (id: string) => CardItem | null {
      return (id: string) => this.items.find((c) => c.id === id) ?? null
    },

    getCardByNodeId(): (nodeId: string) => CardItem | null {
      return (nodeId: string) => this.items.find((c) => c.nodeId === nodeId) ?? null
    },
  },

  actions: {
    async initialize() {
      if (this.initialized) {
        return
      }
      const tree = useTreeStore()
      await tree.initialize()
      for (const n of tree.rawNodes) {
        if (n.type === 'card' && !this.items.some((c) => c.nodeId === n.id)) {
          this.items.push(buildCardItemFromTreeNode(n))
        }
      }
      this.initialized = true
    },

    setItems(items: CardItem[]) {
      this.items = items.map((x) => ({ ...x }))
    },

    addCard(card: CardItem) {
      if (this.items.some((c) => c.id === card.id || c.nodeId === card.nodeId)) {
        return
      }
      this.items.push({ ...card })
    },

    setSelectedCard(id: string | null) {
      this.selectedCardId = id
    },

    openCardByNodeId(nodeId: string) {
      const tree = useTreeStore()
      const node = tree.findNodeById(nodeId)
      if (!node || node.type !== 'card') {
        this.setSelectedCard(null)
        return
      }
      let card = this.getCardByNodeId(nodeId)
      if (!card) {
        card = buildCardItemFromTreeNode(node)
        this.addCard(card)
      }
      this.setSelectedCard(card.id)
    },

    async updateCard(id: string, payload: CardUpdatePayload) {
      const idx = this.items.findIndex((c) => c.id === id)
      if (idx < 0) {
        return
      }
      this.saving = true
      const toast = useToast()
      const prev = this.items[idx]
      try {
        if (isWailsAvailable() && !id.startsWith('card-')) {
          const res = await cardRepository.persistCardUpdate(id, payload, prev)
          if (!res.ok) {
            toast.warning(res.error.message ?? '保存失败')
            return
          }
        } else {
          await new Promise((r) => window.setTimeout(r, 120))
        }
        const next: CardItem = {
          ...prev,
          title: payload.title ?? prev.title,
          prompt: payload.prompt ?? prev.prompt,
          answer: payload.answer ?? prev.answer,
          updatedAt: new Date().toISOString(),
        }
        this.items.splice(idx, 1, next)
        if (payload.title != null && payload.title !== prev.title) {
          useTreeStore().updateNodeTitle(next.nodeId, next.title)
        }
      } finally {
        this.saving = false
      }
    },

    attachReviewItem(cardId: string, reviewItemId: string) {
      const idx = this.items.findIndex((c) => c.id === cardId)
      if (idx < 0) {
        return
      }
      this.items[idx] = { ...this.items[idx], reviewItemId }
    },

    getOrCreateFromNode(nodeId: string): CardItem | null {
      const tree = useTreeStore()
      const node = tree.findNodeById(nodeId)
      if (!node || node.type !== 'card') {
        return null
      }
      let card = this.getCardByNodeId(nodeId)
      if (!card) {
        card = buildCardItemFromTreeNode(node)
        this.addCard(card)
      }
      return card
    },

    /** Add review from current card, or open existing. Returns whether navigated to review. */
    async addOrNavigateReviewFromSelected(options?: { openReview?: boolean }): Promise<void> {
      const toast = useToast()
      const card = this.selectedCard
      if (!card) {
        return
      }
      const review = useReviewStore()
      await review.initialize()
      const existing = review.getReviewItemByCardId(card.id, card.nodeId)
      if (existing) {
        if (!card.reviewItemId) {
          this.attachReviewItem(card.id, existing.id)
        }
        toast.info('Opened existing review item')
        if (options?.openReview) {
          await review.openReviewItemInWorkspace(existing.id)
        }
        return
      }
      const revId = createEntityId('rev')
      const item = buildReviewItemFromCard(card, revId)
      review.addReviewItem(item, { select: false })
      this.attachReviewItem(card.id, revId)
      toast.success('Review item created from card')
      if (options?.openReview) {
        await review.openReviewItemInWorkspace(revId)
      }
    },

    async openLinkedReviewFromSelected(): Promise<boolean> {
      const card = this.selectedCard
      if (!card?.reviewItemId) {
        return false
      }
      const review = useReviewStore()
      await review.initialize()
      await review.openReviewItemInWorkspace(card.reviewItemId)
      return true
    },
  },
})
