import { router } from '@/app/router'
import { CONTEXT_ACTION_IDS } from '@/shared/constants/context-actions'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { useToast } from '@/shared/composables/useToast'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useSearchStore } from '@/modules/search/stores/search.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useObjectCreationStore } from '@/modules/object-creation/stores/object-creation.store'
import { useCardStore } from '@/modules/card/stores/card.store'
import type { ContextMenuContext } from '../types'

const A = CONTEXT_ACTION_IDS

async function goContext(ctx: 'reading' | 'knowledge' | 'review' | 'search') {
  useWorkspaceStore().setContext(ctx)
  await router.push({ name: ROUTE_NAMES[ctx] })
}

async function copyToClipboard(text: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch {
    return false
  }
}

export async function executeContextAction(
  actionId: string,
  context: ContextMenuContext,
): Promise<void> {
  const toast = useToast()
  const tree = useTreeStore()
  const reader = useReaderStore()
  const extract = useExtractStore()
  const review = useReviewStore()
  const search = useSearchStore()
  const editor = useEditorStore()
  const workspace = useWorkspaceStore()

  switch (actionId) {
    case A.tree.openNode: {
      const id = context.nodeId ?? context.entityId
      await tree.initialize()
      tree.setSelectedNode(id)
      await goContext('knowledge')
      return
    }
    case A.tree.newChildNode: {
      await tree.initialize()
      const parentId = context.nodeId ?? context.entityId
      useObjectCreationStore().openDialog('node', { parentNodeId: parentId })
      return
    }
    case A.tree.newCardUnderNode: {
      await tree.initialize()
      const parentId = context.nodeId ?? context.entityId
      useObjectCreationStore().openDialog('card', { parentNodeId: parentId })
      return
    }
    case A.tree.newManualExtractUnderNode: {
      await tree.initialize()
      const parentId = context.nodeId ?? context.entityId
      useObjectCreationStore().openDialog('manual-extract', { parentNodeId: parentId })
      return
    }
    case A.tree.revealNode: {
      const id = context.nodeId ?? context.entityId
      await tree.initialize()
      tree.setSelectedNode(id)
      await goContext('knowledge')
      return
    }
    case A.tree.renameNode: {
      await tree.initialize()
      const id = context.nodeId ?? context.entityId
      const node = tree.findNodeById(id)
      if (!node) {
        return
      }
      const next = window.prompt('节点标题', node.title)
      if (next == null || !next.trim()) {
        return
      }
      tree.updateNodeTitle(id, next.trim())
      toast.success('标题已更新')
      return
    }
    case A.tree.deleteNode: {
      await tree.initialize()
      const id = context.nodeId ?? context.entityId
      if (tree.isNodeRoot(id)) {
        toast.warning('根节点不能删除')
        return
      }
      if (tree.removeNodeSubtree(id)) {
        toast.success('节点已删除')
      }
      return
    }
    case A.tree.cardAddToReview: {
      const nodeId = context.nodeId ?? context.entityId
      await tree.initialize()
      const cardStore = useCardStore()
      await cardStore.initialize()
      cardStore.openCardByNodeId(nodeId)
      await cardStore.addOrNavigateReviewFromSelected()
      return
    }
    case A.tree.cardOpenReview: {
      const nodeId = context.nodeId ?? context.entityId
      const cardStore = useCardStore()
      await cardStore.initialize()
      cardStore.openCardByNodeId(nodeId)
      const opened = await cardStore.openLinkedReviewFromSelected()
      if (!opened) {
        const c = cardStore.getCardByNodeId(nodeId)
        if (c) {
          await review.initialize()
          const item = review.getReviewItemByCardId(c.id, c.nodeId)
          if (item) {
            if (!c.reviewItemId) {
              cardStore.attachReviewItem(c.id, item.id)
            }
            await review.openReviewItemInWorkspace(item.id)
          } else {
            toast.warning('No review item linked for this card')
          }
        }
      }
      return
    }
    case A.tree.createReview: {
      await tree.initialize()
      const id = context.nodeId ?? context.entityId
      const node = tree.findNodeById(id)
      if (!node) {
        return
      }
      await review.initialize()
      review.addMockReviewFromTreeNode(id, node.title)
      toast.success('已加入复习队列（演示）')
      return
    }

    case A.article.openArticle:
    case A.article.openInReading: {
      await reader.initialize()
      const articleId = context.articleId ?? context.entityId
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
      workspace.setContext('reading')
      await router.push({ name: ROUTE_NAMES.reading })
      return
    }
    case A.article.createExtractFromSelection: {
      const articleId = context.articleId ?? context.entityId
      const article = reader.getArticleById(articleId)
      const sel = editor.selection
      if (!article || !sel?.text?.trim()) {
        toast.warning('请先选中一段文字')
        return
      }
      await extract.initialize()
      extract.openCreateDialog({
        sourceArticleId: article.id,
        sourceArticleTitle: article.title,
        quote: sel.text,
        note: '',
        title: '',
        parentNodeId: tree.selectedNodeId,
        selection: { from: sel.from, to: sel.to, text: sel.text },
      })
      toast.info('已打开摘录对话框')
      return
    }
    case A.article.markReading: {
      await reader.initialize()
      const articleId = context.articleId ?? context.entityId
      reader.setArticleStatus(articleId, 'reading')
      toast.success('已标为阅读中')
      return
    }
    case A.article.markProcessed: {
      await reader.initialize()
      const articleId = context.articleId ?? context.entityId
      reader.setArticleStatus(articleId, 'processed')
      toast.success('已标为已处理')
      return
    }
    case A.article.revealRelatedExtracts: {
      const articleId = context.articleId ?? context.entityId
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
      workspace.setContext('reading')
      await router.push({ name: ROUTE_NAMES.reading })
      toast.info('请在下方「Related Extracts」查看摘录')
      return
    }
    case A.article.createManualFromArticle: {
      const articleId = context.articleId ?? context.entityId
      await reader.initialize()
      const article = reader.getArticleById(articleId)
      if (!article) {
        return
      }
      useObjectCreationStore().openDialog('manual-extract', {
        sourceArticleId: article.id,
        sourceArticleTitle: article.title,
      })
      return
    }
    case A.article.createCardFromArticle: {
      const articleId = context.articleId ?? context.entityId
      await reader.initialize()
      const article = reader.getArticleById(articleId)
      const creation = useObjectCreationStore()
      creation.openDialog('card')
      if (article) {
        creation.updateCardDraft({
          prompt: `Article: ${article.title}\n\n`,
        })
      }
      return
    }

    case A.extract.openExtract: {
      const extractId = context.extractId ?? context.entityId
      await Promise.all([extract.initialize(), tree.initialize()])
      extract.openExtract(extractId)
      const ex = extract.getExtractById(extractId)
      if (ex?.treeNodeId) {
        tree.setSelectedNode(ex.treeNodeId)
      }
      await goContext('knowledge')
      return
    }
    case A.extract.backToSource: {
      const extractId = context.extractId ?? context.entityId
      await extract.initialize()
      const ex = extract.getExtractById(extractId)
      if (!ex?.sourceArticleId) {
        return
      }
      await reader.initialize()
      reader.openArticleById(ex.sourceArticleId)
      const article = reader.getArticleById(ex.sourceArticleId)
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
      workspace.setContext('reading')
      await router.push({ name: ROUTE_NAMES.reading })
      return
    }
    case A.extract.revealInTree: {
      const extractId = context.extractId ?? context.entityId
      await Promise.all([extract.initialize(), tree.initialize()])
      const ex = extract.getExtractById(extractId)
      if (!ex?.treeNodeId) {
        return
      }
      tree.setSelectedNode(ex.treeNodeId)
      extract.setSelectedExtract(extractId)
      await goContext('knowledge')
      return
    }
    case A.extract.addToReview: {
      const extractId = context.extractId ?? context.entityId
      await extract.initialize()
      extract.openExtract(extractId)
      await extract.addOrOpenReviewFromSelectedExtract()
      return
    }
    case A.extract.openReviewItem: {
      const extractId = context.extractId ?? context.entityId
      await extract.initialize()
      extract.openExtract(extractId)
      await review.initialize()
      const ex = extract.getExtractById(extractId)
      if (!ex) {
        return
      }
      if (ex.reviewItemId) {
        await review.openReviewItemInWorkspace(ex.reviewItemId)
        return
      }
      const item = review.getReviewItemByExtractId(ex.id)
      if (item) {
        extract.attachReviewItem(ex.id, item.id)
        await review.openReviewItemInWorkspace(item.id)
      } else {
        toast.warning('No review item linked for this extract')
      }
      return
    }
    case A.extract.copyQuote: {
      const extractId = context.extractId ?? context.entityId
      await extract.initialize()
      const ex = extract.getExtractById(extractId)
      if (!ex) {
        return
      }
      const ok = await copyToClipboard(ex.quote)
      if (ok) {
        toast.success('摘录已复制到剪贴板')
      } else {
        toast.error('复制失败')
      }
      return
    }
    case A.extract.deleteExtract: {
      const extractId = context.extractId ?? context.entityId
      await Promise.all([extract.initialize(), tree.initialize()])
      extract.removeExtractById(extractId)
      toast.success('摘录已删除')
      return
    }

    case A.search.openResult:
    case A.search.revealResult: {
      const rid = context.searchResultId ?? context.entityId
      const hit = search.results.find((r) => r.id === rid)
      if (!hit) {
        return
      }
      await search.openResult(hit)
      return
    }
    case A.search.copyTitle: {
      const rid = context.searchResultId ?? context.entityId
      const hit = search.results.find((r) => r.id === rid)
      if (!hit?.title) {
        return
      }
      const ok = await copyToClipboard(hit.title)
      if (ok) {
        toast.success('标题已复制')
      } else {
        toast.error('复制失败')
      }
      return
    }

    case A.review.openItem: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      review.setSelectedItem(rid)
      await goContext('review')
      return
    }
    case A.review.openSourceArticle: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      const item = review.getItemById(rid)
      if (!item?.sourceArticleId) {
        return
      }
      await reader.initialize()
      reader.openArticleById(item.sourceArticleId)
      const article = reader.getArticleById(item.sourceArticleId)
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
      workspace.setContext('reading')
      await router.push({ name: ROUTE_NAMES.reading })
      return
    }
    case A.review.openExtract: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      const item = review.getItemById(rid)
      if (!item?.extractId) {
        return
      }
      await Promise.all([extract.initialize(), tree.initialize()])
      extract.openExtract(item.extractId)
      const ex = extract.getExtractById(item.extractId)
      if (ex?.treeNodeId) {
        tree.setSelectedNode(ex.treeNodeId)
      }
      await goContext('knowledge')
      return
    }
    case A.review.openCard: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      const item = review.getItemById(rid)
      if (!item || item.type !== 'card') {
        return
      }
      const cardStore = useCardStore()
      await Promise.all([cardStore.initialize(), tree.initialize()])
      if (item.nodeId) {
        cardStore.openCardByNodeId(item.nodeId)
      } else if (item.cardId) {
        cardStore.setSelectedCard(item.cardId)
        const c = cardStore.getCardById(item.cardId)
        if (c?.nodeId) {
          tree.setSelectedNode(c.nodeId)
        }
      }
      await goContext('knowledge')
      return
    }
    case A.review.markReviewed: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      review.removeItemFromSessionQueue(rid)
      toast.success('已标记为已复习（演示）')
      return
    }
    case A.review.removeFromQueue: {
      const rid = context.reviewId ?? context.entityId
      await review.initialize()
      review.removeItemFromSessionQueue(rid)
      toast.success('已从本轮队列移除')
      return
    }

    default:
      return
  }
}
