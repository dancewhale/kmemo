import { CONTEXT_ACTION_IDS } from '@/shared/constants/context-actions'
import type { ContextEntityType, ContextMenuRegistryRow } from '../types'

const T = CONTEXT_ACTION_IDS

const registry: Record<ContextEntityType, ContextMenuRegistryRow[]> = {
  'tree-node': [
    { id: T.tree.openNode, label: 'Open', icon: 'FolderOpened' },
    { id: T.tree.renameNode, label: 'Rename…', icon: 'EditPen', dividerBefore: true },
    { id: T.tree.newChildNode, label: 'New Child Node…', icon: 'Plus' },
    { id: T.tree.newCardUnderNode, label: 'New Card Under Node', icon: 'Postcard' },
    { id: T.tree.newManualExtractUnderNode, label: 'New Manual Extract Under Node', icon: 'EditPen' },
    { id: T.tree.createReview, label: 'Create Review Item', icon: 'Notebook' },
    { id: T.tree.revealNode, label: 'Reveal in Knowledge', icon: 'View' },
    { id: T.tree.deleteNode, label: 'Delete Node', icon: 'Delete', danger: true, dividerBefore: true },
  ],
  'tree-card': [
    { id: T.tree.openNode, label: 'Open', icon: 'FolderOpened' },
    { id: T.tree.cardAddToReview, label: 'Add to Review', icon: 'Notebook', dividerBefore: true },
    { id: T.tree.cardOpenReview, label: 'Open Review', icon: 'Notebook' },
    { id: T.tree.revealNode, label: 'Reveal in Knowledge', icon: 'View' },
  ],
  'tree-extract': [
    { id: T.extract.openExtract, label: 'Open', icon: 'EditPen' },
    { id: T.extract.addToReview, label: 'Add to Review', icon: 'Notebook', dividerBefore: true },
    { id: T.extract.openReviewItem, label: 'Open Review Item', icon: 'Notebook' },
    { id: T.tree.revealNode, label: 'Reveal in Knowledge', icon: 'View' },
  ],
  article: [
    { id: T.article.openArticle, label: 'Open Article', icon: 'Document' },
    { id: T.article.openInReading, label: 'Open in Reading', icon: 'Reading' },
    {
      id: T.article.createExtractFromSelection,
      label: 'Create Extract From Selection',
      icon: 'EditPen',
      dividerBefore: true,
    },
    { id: T.article.markReading, label: 'Mark as Reading', icon: 'Clock' },
    { id: T.article.markProcessed, label: 'Mark as Processed', icon: 'CircleCheck' },
    { id: T.article.revealRelatedExtracts, label: 'Reveal Related Extracts', icon: 'Collection' },
    { id: T.article.createManualFromArticle, label: 'New Manual Extract', icon: 'EditPen', dividerBefore: true },
    { id: T.article.createCardFromArticle, label: 'New Card From Article', icon: 'Postcard' },
  ],
  extract: [
    { id: T.extract.openExtract, label: 'Open Extract', icon: 'EditPen' },
    { id: T.extract.backToSource, label: 'Back to Source Article', icon: 'ArrowLeft' },
    { id: T.extract.revealInTree, label: 'Reveal in Tree', icon: 'FolderOpened' },
    { id: T.extract.addToReview, label: 'Add to Review', icon: 'Notebook', dividerBefore: true },
    { id: T.extract.openReviewItem, label: 'Open Review Item', icon: 'Notebook' },
    { id: T.extract.copyQuote, label: 'Copy Quote', icon: 'DocumentCopy' },
    { id: T.extract.deleteExtract, label: 'Delete Extract', icon: 'Delete', danger: true, dividerBefore: true },
  ],
}

export function getRegistryRows(entityType: ContextEntityType): ContextMenuRegistryRow[] {
  return registry[entityType].map((row) => ({ ...row }))
}
