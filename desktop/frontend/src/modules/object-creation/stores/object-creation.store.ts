import { defineStore } from 'pinia'
import { router } from '@/app/router'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { createEntityId } from '@/shared/utils/id'
import { useToast } from '@/shared/composables/useToast'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { isWailsAvailable } from '@/api/wails'
import * as knowledgeRepository from '@/modules/knowledge-tree/services/knowledge.repository'
import * as cardRepository from '@/modules/card/services/card.repository'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { buildReviewItemFromCard } from '@/modules/review/services/review.mapper'
import { useCardStore } from '@/modules/card/stores/card.store'
import type { CardItem } from '@/modules/card/types'
import {
  buildExtractFromManualExtractInput,
  buildTreeNodeFromCreateCardInput,
  buildTreeNodeFromCreateNodeInput,
  buildTreeNodeFromManualExtract,
} from '../services/creation.mapper'
import { validateCardInput, validateManualExtractInput, validateNodeInput } from '../services/creation.validator'
import type {
  CreateCardInput,
  CreateManualExtractInput,
  CreateNodeInput,
  CreateObjectDraftState,
  CreateObjectType,
} from '../types'

interface ObjectCreationState {
  isDialogOpen: boolean
  activeType: CreateObjectType
  submitting: boolean
  draft: CreateObjectDraftState
}

function emptyDraft(): CreateObjectDraftState {
  return {
    type: 'node',
    node: { title: '', content: '', parentNodeId: null },
    card: { title: '', prompt: '', answer: '', parentNodeId: null, addToReview: true },
    manualExtract: {
      title: '',
      quote: '',
      note: '',
      parentNodeId: null,
      sourceArticleId: null,
      sourceArticleTitle: null,
    },
  }
}

export const useObjectCreationStore = defineStore('object-creation', {
  state: (): ObjectCreationState => ({
    isDialogOpen: false,
    activeType: 'node',
    submitting: false,
    draft: emptyDraft(),
  }),

  getters: {
    activeDraft(): CreateNodeInput | CreateCardInput | CreateManualExtractInput {
      switch (this.activeType) {
        case 'node':
          return this.draft.node
        case 'card':
          return this.draft.card
        case 'manual-extract':
          return this.draft.manualExtract
        default:
          return this.draft.node
      }
    },

    dialogTitle(): string {
      return 'Create'
    },

    canSubmit(): boolean {
      switch (this.activeType) {
        case 'node':
          return validateNodeInput(this.draft.node).valid
        case 'card':
          return validateCardInput(this.draft.card).valid
        case 'manual-extract':
          return validateManualExtractInput(this.draft.manualExtract).valid
        default:
          return false
      }
    },
  },

  actions: {
    openDialog(
      type: CreateObjectType = 'node',
      opts?: {
        parentNodeId?: string | null
        sourceArticleId?: string | null
        sourceArticleTitle?: string | null
      },
    ) {
      const parent = opts?.parentNodeId ?? null
      this.activeType = type
      this.draft.type = type
      this.updateNodeDraft({ parentNodeId: parent })
      this.updateCardDraft({ parentNodeId: parent })
      this.updateManualExtractDraft({
        parentNodeId: parent,
        sourceArticleId: opts?.sourceArticleId ?? null,
        sourceArticleTitle: opts?.sourceArticleTitle ?? null,
      })

      const workspace = useWorkspaceStore()
      const reader = useReaderStore()
      if (
        type === 'manual-extract' &&
        !opts?.sourceArticleId &&
        workspace.currentContext === 'reading' &&
        reader.selectedArticle
      ) {
        const a = reader.selectedArticle
        this.updateManualExtractDraft({
          sourceArticleId: a.id,
          sourceArticleTitle: a.title,
        })
      }

      this.isDialogOpen = true
    },

    closeDialog() {
      this.isDialogOpen = false
    },

    setActiveType(type: CreateObjectType) {
      this.activeType = type
      this.draft.type = type
    },

    updateNodeDraft(payload: Partial<CreateNodeInput>) {
      this.draft.node = { ...this.draft.node, ...payload }
    },

    updateCardDraft(payload: Partial<CreateCardInput>) {
      this.draft.card = { ...this.draft.card, ...payload }
    },

    updateManualExtractDraft(payload: Partial<CreateManualExtractInput>) {
      this.draft.manualExtract = { ...this.draft.manualExtract, ...payload }
    },

    resetDraft(type?: CreateObjectType) {
      if (type === 'node') {
        this.draft.node = { title: '', content: '', parentNodeId: null }
      } else if (type === 'card') {
        this.draft.card = { title: '', prompt: '', answer: '', parentNodeId: null, addToReview: true }
      } else if (type === 'manual-extract') {
        this.draft.manualExtract = {
          title: '',
          quote: '',
          note: '',
          parentNodeId: null,
          sourceArticleId: null,
          sourceArticleTitle: null,
        }
      } else {
        this.draft = emptyDraft()
        this.draft.type = this.activeType
      }
    },

    async submit() {
      if (this.activeType === 'node') {
        await this.submitNode()
      } else if (this.activeType === 'card') {
        await this.submitCard()
      } else {
        await this.submitManualExtract()
      }
    },

    async goKnowledge() {
      const workspace = useWorkspaceStore()
      workspace.setContext('knowledge')
      await router.push({ name: ROUTE_NAMES.knowledge })
    },

    async submitNode() {
      const toast = useToast()
      const v = validateNodeInput(this.draft.node)
      if (!v.valid) {
        toast.warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const tree = useTreeStore()
        await tree.initialize()
        if (isWailsAvailable()) {
          const res = await knowledgeRepository.createKnowledgeSpace(
            this.draft.node.title.trim(),
            this.draft.node.content.trim(),
            this.draft.node.parentNodeId,
          )
          if (!res.ok) {
            toast.warning(res.error.message ?? '创建知识库失败')
            return
          }
          await tree.initialize({ force: true })
          tree.setSelectedNode(res.data)
        } else {
          const id = createEntityId('kn-topic')
          const node = buildTreeNodeFromCreateNodeInput(this.draft.node, id)
          tree.addNode(node)
        }
        await this.goKnowledge()
        toast.success('Node created')
        this.closeDialog()
        this.resetDraft()
      } finally {
        this.submitting = false
      }
    },

    async submitCard() {
      const toast = useToast()
      const v = validateCardInput(this.draft.card)
      if (!v.valid) {
        toast.warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const tree = useTreeStore()
        const review = useReviewStore()
        const cardStore = useCardStore()
        await Promise.all([tree.initialize(), cardStore.initialize()])

        if (isWailsAvailable()) {
          const knowledgeId =
            this.draft.card.parentNodeId ??
            tree.rawNodes.find((n) => n.parentId == null)?.id ??
            null
          if (!knowledgeId) {
            toast.warning('请选择父级知识库节点')
            return
          }
          const res = await cardRepository.createCardFromInput(this.draft.card, knowledgeId)
          if (!res.ok) {
            toast.warning(res.error.message ?? '创建卡片失败')
            return
          }
          const newId = res.data
          await tree.initialize({ force: true })
          const displayTitle =
            this.draft.card.title.trim() ||
            (this.draft.card.prompt.trim() ? this.draft.card.prompt.trim().slice(0, 48) : 'Card')
          tree.appendHostCardNode({
            id: newId,
            parentId: this.draft.card.parentNodeId ?? knowledgeId,
            title: displayTitle,
          })
          const cardItem: CardItem = {
            id: newId,
            nodeId: newId,
            title: displayTitle.startsWith('Card:') ? displayTitle : `Card: ${displayTitle}`,
            prompt: this.draft.card.prompt.trim(),
            answer: this.draft.card.answer.trim(),
            parentNodeId: this.draft.card.parentNodeId ?? knowledgeId,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            reviewItemId: null,
          }
          cardStore.addCard(cardItem)
          if (this.draft.card.addToReview) {
            toast.success('Card created — FSRS will schedule review')
          } else {
            toast.success('Card created')
          }
        } else {
          const nodeId = createEntityId('kn-card')
          const kn = buildTreeNodeFromCreateCardInput(this.draft.card, nodeId)
          tree.addNode(kn)
          const cardId = createEntityId('card')
          const cardItem: CardItem = {
            id: cardId,
            nodeId: kn.id,
            title: kn.title,
            prompt: this.draft.card.prompt.trim(),
            answer: this.draft.card.answer.trim(),
            parentNodeId: this.draft.card.parentNodeId,
            createdAt: kn.createdAt ?? new Date().toISOString(),
            updatedAt: kn.updatedAt ?? new Date().toISOString(),
            reviewItemId: null,
          }
          cardStore.addCard(cardItem)
          if (this.draft.card.addToReview) {
            await review.initialize()
            const revId = createEntityId('rev')
            const item = buildReviewItemFromCard(cardItem, revId)
            review.addReviewItem(item, { select: false })
            cardStore.attachReviewItem(cardItem.id, revId)
            toast.success('Card added to review')
          } else {
            toast.success('Card created')
          }
        }
        await this.goKnowledge()
        this.closeDialog()
        this.resetDraft()
      } finally {
        this.submitting = false
      }
    },

    async submitManualExtract() {
      const toast = useToast()
      const v = validateManualExtractInput(this.draft.manualExtract)
      if (!v.valid) {
        toast.warning(v.message ?? 'Invalid input')
        return
      }
      if (this.submitting) {
        return
      }
      this.submitting = true
      try {
        const tree = useTreeStore()
        const extract = useExtractStore()
        await Promise.all([tree.initialize(), extract.initialize()])
        const extId = createEntityId('ext')
        const treeNodeId = `kn-ext-${extId}`
        const item = buildExtractFromManualExtractInput(this.draft.manualExtract, extId, treeNodeId)
        const kn = buildTreeNodeFromManualExtract(item)
        extract.addExtract(item)
        tree.addNode(kn)
        await this.goKnowledge()
        toast.success('Manual extract created')
        this.closeDialog()
        this.resetDraft()
      } finally {
        this.submitting = false
      }
    },
  },
})
