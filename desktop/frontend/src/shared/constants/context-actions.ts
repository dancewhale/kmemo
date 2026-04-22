/** Action ids for context menu executor / registry (single source of truth). */
export const CONTEXT_ACTION_IDS = {
  tree: {
    openNode: 'open-node',
    renameNode: 'rename-node',
    newChildNode: 'new-child-node',
    newCardUnderNode: 'new-card-under-node',
    newManualExtractUnderNode: 'new-manual-extract-under-node',
    createReview: 'create-review-from-node',
    revealNode: 'reveal-node',
    deleteNode: 'delete-node',
    cardAddToReview: 'tree-card-add-to-review',
    cardOpenReview: 'tree-card-open-review',
  },
  article: {
    openArticle: 'open-article',
    openInReading: 'open-article-reading',
    createExtractFromSelection: 'create-extract-from-selection',
    markReading: 'mark-article-reading',
    markProcessed: 'mark-article-processed',
    revealRelatedExtracts: 'reveal-related-extracts',
    createManualFromArticle: 'article-create-manual-extract',
    createCardFromArticle: 'article-create-card',
  },
  extract: {
    openExtract: 'open-extract',
    backToSource: 'back-to-source-article',
    revealInTree: 'reveal-extract-in-tree',
    addToReview: 'add-extract-to-review',
    openReviewItem: 'open-extract-review-item',
    copyQuote: 'copy-extract-quote',
    deleteExtract: 'delete-extract',
  },
} as const
