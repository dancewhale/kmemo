import { CONTEXT_ACTION_IDS } from '@/shared/constants/context-actions'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useCardStore } from '@/modules/card/stores/card.store'
import { buildCardItemFromTreeNode } from '@/modules/card/services/card.mapper'
import type { ContextActionItem, ContextMenuContext } from '../types'
import { getRegistryRows } from './context-menu.registry'

const A = CONTEXT_ACTION_IDS

function cloneRows(context: ContextMenuContext): ContextActionItem[] {
  return getRegistryRows(context.entityType).map((r) => ({
    id: r.id,
    label: r.label,
    icon: r.icon,
    danger: r.danger,
    dividerBefore: r.dividerBefore,
    disabled: false,
  }))
}

export function buildContextMenuItems(context: ContextMenuContext): ContextActionItem[] {
  const items = cloneRows(context)
  const tree = useTreeStore()
  const reader = useReaderStore()
  const extract = useExtractStore()
  const review = useReviewStore()
  const editor = useEditorStore()

  if (context.entityType === 'tree-node') {
    const nodeId = context.nodeId ?? context.entityId
    const node = tree.findNodeById(nodeId)
    const isRoot = node?.parentId == null
    for (const row of items) {
      if (row.id === A.tree.deleteNode) {
        row.disabled = isRoot || !node
      }
      if (row.id === A.tree.renameNode) {
        row.disabled = !node
      }
    }
  }

  if (context.entityType === 'tree-card') {
    const nodeId = context.nodeId ?? context.entityId
    const node = tree.findNodeById(nodeId)
    const cardStore = useCardStore()
    const stored = node?.type === 'card' ? cardStore.getCardByNodeId(nodeId) : null
    const card = stored ?? (node?.type === 'card' ? buildCardItemFromTreeNode(node) : null)
    const linked = card ? review.getReviewItemByCardId(card.id, card.nodeId) : null
    const canOpenReview = Boolean(card && ((stored?.reviewItemId ?? null) != null || linked))
    for (const row of items) {
      if (row.id === A.tree.cardOpenReview) {
        row.disabled = !canOpenReview
      }
    }
  }

  if (context.entityType === 'tree-extract') {
    const extractId = context.extractId ?? context.entityId
    const ex = extract.getExtractById(extractId)
    const linked = ex ? review.getReviewItemByExtractId(ex.id) : null
    for (const row of items) {
      if (!ex) {
        row.disabled = true
        continue
      }
      if (row.id === A.extract.openReviewItem) {
        row.disabled = !(ex.reviewItemId || linked)
      }
    }
  }

  if (context.entityType === 'article') {
    const articleId = context.articleId ?? context.entityId
    const article = reader.getArticleById(articleId)
    const sel = editor.selection
    const doc = editor.currentDocument
    const canExtract =
      !!article &&
      doc?.id === articleId &&
      !!sel?.text?.trim() &&
      doc.contentType === 'article'
    for (const row of items) {
      if (!article) {
        row.disabled = true
        continue
      }
      if (row.id === A.article.createExtractFromSelection) {
        row.disabled = !canExtract
      }
    }
  }

  if (context.entityType === 'extract') {
    const extractId = context.extractId ?? context.entityId
    const ex = extract.getExtractById(extractId)
    for (const row of items) {
      if (!ex) {
        row.disabled = true
        continue
      }
      if (row.id === A.extract.backToSource) {
        row.disabled = !ex.sourceArticleId
      }
      if (row.id === A.extract.revealInTree) {
        row.disabled = !ex.treeNodeId
      }
      if (row.id === A.extract.openReviewItem) {
        const linked = review.getReviewItemByExtractId(ex.id)
        row.disabled = !(ex.reviewItemId || linked)
      }
    }
  }

  return items
}
