import { COMMAND_IDS, type CommandId } from '@/shared/constants/commands'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import { router } from '@/app/router'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useEditorStore } from '@/modules/editor/stores/editor.store'
import { useSearchStore } from '@/modules/search/stores/search.store'
import { useCaptureStore } from '@/modules/inbox-capture/stores/capture.store'
import { useObjectCreationStore } from '@/modules/object-creation/stores/object-creation.store'
import { useCardStore } from '@/modules/card/stores/card.store'
import { buildPendingExtractFromSelection } from '@/modules/extract/services/extract.mapper'

async function goContext(context: 'inbox' | 'reading' | 'knowledge' | 'review' | 'search') {
  const workspace = useWorkspaceStore()
  workspace.setContext(context)
  await router.push({ name: ROUTE_NAMES[context] })
}

export function isCommandEnabled(commandId: CommandId): boolean {
  const workspace = useWorkspaceStore()
  const reader = useReaderStore()
  const tree = useTreeStore()
  const extract = useExtractStore()
  const review = useReviewStore()
  const editor = useEditorStore()

  switch (commandId) {
    case COMMAND_IDS.createExtract:
      return Boolean(editor.currentDocument && editor.selection?.text?.trim())
    case COMMAND_IDS.openSelectedArticle:
      return Boolean(reader.selectedArticle ?? workspace.selectedArticleId)
    case COMMAND_IDS.openSelectedExtract:
      return Boolean(extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId)
    case COMMAND_IDS.backToSourceArticle:
      return Boolean(extract.selectedExtract?.sourceArticleId)
    case COMMAND_IDS.revealCurrentExtractInTree:
      return Boolean(extract.selectedExtract?.treeNodeId)
    case COMMAND_IDS.addCurrentCardToReview:
      return workspace.currentContext === 'knowledge' && tree.selectedNode?.type === 'card'
    case COMMAND_IDS.openCurrentCardReview:
      return workspace.currentContext === 'knowledge' && tree.selectedNode?.type === 'card'
    case COMMAND_IDS.addCurrentExtractToReview:
      return (
        workspace.currentContext === 'knowledge' &&
        Boolean(extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId)
      )
    case COMMAND_IDS.openCurrentExtractReview:
      return (
        workspace.currentContext === 'knowledge' &&
        Boolean(extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId)
      )
    case COMMAND_IDS.openFirstReviewItem:
      return review.initialized ? review.queueItems.length > 0 : true
    default:
      return true
  }
}

export async function executeCommand(commandId: CommandId): Promise<boolean> {
  const workspace = useWorkspaceStore()
  const reader = useReaderStore()
  const tree = useTreeStore()
  const extract = useExtractStore()
  const review = useReviewStore()
  const editor = useEditorStore()
  const search = useSearchStore()
  const capture = useCaptureStore()
  const creation = useObjectCreationStore()
  const card = useCardStore()

  if (!isCommandEnabled(commandId)) {
    return false
  }

  switch (commandId) {
    case COMMAND_IDS.goInbox:
      await goContext('inbox')
      return true
    case COMMAND_IDS.captureBlankArticle:
      capture.openDialog('blank')
      return true
    case COMMAND_IDS.capturePasteText:
      capture.openDialog('paste-text')
      return true
    case COMMAND_IDS.captureUrl:
      capture.openDialog('url')
      return true
    case COMMAND_IDS.createNode:
      await tree.initialize()
      await goContext('knowledge')
      creation.openDialog('node')
      return true
    case COMMAND_IDS.createCard:
      await tree.initialize()
      await goContext('knowledge')
      creation.openDialog('card')
      return true
    case COMMAND_IDS.createManualExtract:
      await reader.initialize()
      creation.openDialog('manual-extract')
      return true
    case COMMAND_IDS.goReading:
      await goContext('reading')
      return true
    case COMMAND_IDS.goKnowledge:
      await goContext('knowledge')
      return true
    case COMMAND_IDS.goReview:
      await goContext('review')
      return true
    case COMMAND_IDS.goSearch:
      await goContext('search')
      return true
    case COMMAND_IDS.focusSearch:
      await goContext('search')
      await search.initialize()
      search.requestFocus()
      return true
    case COMMAND_IDS.openFirstReviewItem:
      await review.initialize()
      review.openFirstAvailable()
      await goContext('review')
      return true
    case COMMAND_IDS.openSelectedArticle: {
      await reader.initialize()
      const id = reader.selectedArticle?.id ?? workspace.selectedArticleId
      if (!id) {
        return false
      }
      reader.openArticleById(id)
      const article = reader.getArticleById(id)
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
      return true
    }
    case COMMAND_IDS.createExtract: {
      if (!editor.currentDocument || !editor.selection) {
        return false
      }
      const payload = buildPendingExtractFromSelection({
        sourceArticleId: editor.currentDocument.id,
        sourceArticleTitle: editor.currentDocument.title,
        selection: editor.selection,
        defaultParentNodeId: tree.selectedNodeId ?? 'kn-reading',
      })
      extract.openCreateDialog(payload)
      return true
    }
    case COMMAND_IDS.openSelectedExtract: {
      await extract.initialize()
      const id = extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId ?? null
      if (!id) {
        return false
      }
      extract.openExtract(id)
      const selected = extract.getExtractById(id)
      if (selected?.treeNodeId) {
        tree.setSelectedNode(selected.treeNodeId)
      }
      await goContext('knowledge')
      return true
    }
    case COMMAND_IDS.backToSourceArticle: {
      const sourceId = extract.selectedExtract?.sourceArticleId
      if (!sourceId) {
        return false
      }
      await reader.initialize()
      reader.openArticleById(sourceId)
      const article = reader.getArticleById(sourceId)
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
      return true
    }
    case COMMAND_IDS.revealCurrentExtractInTree: {
      const nodeId = extract.selectedExtract?.treeNodeId
      if (!nodeId) {
        return false
      }
      tree.setSelectedNode(nodeId)
      await goContext('knowledge')
      return true
    }
    case COMMAND_IDS.addCurrentCardToReview: {
      const node = tree.selectedNode
      if (node?.type !== 'card') {
        return false
      }
      await tree.initialize()
      await card.initialize()
      card.openCardByNodeId(node.id)
      await card.addOrNavigateReviewFromSelected()
      return true
    }
    case COMMAND_IDS.openCurrentCardReview: {
      const node = tree.selectedNode
      if (node?.type !== 'card') {
        return false
      }
      await tree.initialize()
      await card.initialize()
      card.openCardByNodeId(node.id)
      const ok = await card.openLinkedReviewFromSelected()
      if (!ok) {
        const c = card.getCardByNodeId(node.id)
        if (c) {
          await review.initialize()
          const item = review.getReviewItemByCardId(c.id, c.nodeId)
          if (item) {
            if (!c.reviewItemId) {
              card.attachReviewItem(c.id, item.id)
            }
            await review.openReviewItemInWorkspace(item.id)
          }
        }
      }
      return true
    }
    case COMMAND_IDS.addCurrentExtractToReview: {
      await extract.initialize()
      const extractId = extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId ?? null
      if (!extractId) {
        return false
      }
      extract.openExtract(extractId)
      await extract.addOrOpenReviewFromSelectedExtract()
      return true
    }
    case COMMAND_IDS.openCurrentExtractReview: {
      await extract.initialize()
      const extractId = extract.selectedExtractId ?? tree.selectedNode?.sourceExtractId ?? null
      if (!extractId) {
        return false
      }
      extract.openExtract(extractId)
      await review.initialize()
      const ex = extract.getExtractById(extractId)
      if (!ex) {
        return false
      }
      if (ex.reviewItemId) {
        await review.openReviewItemInWorkspace(ex.reviewItemId)
        return true
      }
      const item = review.getReviewItemByExtractId(ex.id)
      if (item) {
        extract.attachReviewItem(ex.id, item.id)
        await review.openReviewItemInWorkspace(item.id)
        return true
      }
      return false
    }
    case COMMAND_IDS.toggleCommandPalette:
    case COMMAND_IDS.closeCommandPalette:
      // handled by command store / shortcuts
      return true
    default:
      return false
  }
}
