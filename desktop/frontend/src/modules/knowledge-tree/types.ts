export type KnowledgeNodeType = 'topic' | 'article' | 'extract' | 'card'

export interface KnowledgeNode {
  id: string
  parentId: string | null
  title: string
  type: KnowledgeNodeType
  /** Optional body text for topic/card nodes (phase-1 creation). */
  description?: string
  sourceArticleId?: string
  sourceExtractId?: string
  expanded?: boolean
  selected?: boolean
  hasChildren?: boolean
  childCount?: number
  createdAt?: string
  updatedAt?: string
}

export interface UITreeNode extends KnowledgeNode {
  level: number
  children: UITreeNode[]
}
