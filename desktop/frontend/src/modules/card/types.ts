export interface CardItem {
  id: string
  nodeId: string
  title: string
  prompt: string
  answer: string
  parentNodeId: string | null
  createdAt: string
  updatedAt: string
  reviewItemId?: string | null
}

export interface CardUpdatePayload {
  title?: string
  prompt?: string
  answer?: string
}
