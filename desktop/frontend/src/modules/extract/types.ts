export interface TextSelectionRange {
  from: number
  to: number
  text: string
}

export interface PendingExtractPayload {
  sourceArticleId: string
  sourceArticleTitle: string
  quote: string
  note: string
  title: string
  parentNodeId: string | null
  selection: TextSelectionRange
}

export interface ExtractItem {
  id: string
  sourceArticleId: string
  sourceArticleTitle: string
  title: string
  quote: string
  note: string
  parentNodeId: string | null
  treeNodeId?: string | null
  reviewItemId?: string | null
  createdAt: string
  updatedAt: string
}

export interface ExtractUpdatePayload {
  title?: string
  note?: string
  quote?: string
}
