export interface CreateKnowledgeRequest {
  name: string
  description: string
  parentId: string | null
}

export interface UpdateKnowledgeRequest {
  name: string
  description: string
}

export interface CreateCardRequest {
  knowledgeId: string
  sourceDocumentId?: string | null
  parentId?: string | null
  title: string
  cardType: string
  htmlContent: string
  tagIds?: string[]
}

export interface UpdateCardRequest {
  title: string
  htmlContent: string
  status: string
}

export interface MoveCardRequest {
  cardId: string
  targetParentId: string | null
  targetIndex: number
}

export interface ReorderCardChildrenRequest {
  knowledgeId: string
  parentId: string | null
  orderedChildIds: string[]
}

export interface ReviewRequest {
  cardId: string
  rating: number
}

export interface CreateTagRequest {
  name: string
  slug: string
  color: string
  icon: string
  description: string
}

export interface UpdateTagRequest {
  name: string
  color: string
  icon: string
  description: string
}

export interface UpdateDefaultFSRSParameterRequest {
  parametersJson: string
  desiredRetention: number | null
  maximumInterval: number | null
}
