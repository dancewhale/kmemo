export type CreateObjectType = 'node' | 'card' | 'manual-extract'

export interface CreateNodeInput {
  title: string
  content: string
  parentNodeId: string | null
}

export interface CreateCardInput {
  title: string
  prompt: string
  answer: string
  parentNodeId: string | null
  addToReview: boolean
}

export interface CreateManualExtractInput {
  title: string
  quote: string
  note: string
  parentNodeId: string | null
  sourceArticleId?: string | null
  sourceArticleTitle?: string | null
}

export interface CreateObjectDraftState {
  type: CreateObjectType
  node: CreateNodeInput
  card: CreateCardInput
  manualExtract: CreateManualExtractInput
}

export interface CreateObjectResult {
  objectType: CreateObjectType
  id: string
  nodeId?: string | null
  extractId?: string | null
  reviewItemId?: string | null
}

export interface CreateObjectValidationResult {
  valid: boolean
  message?: string
}
