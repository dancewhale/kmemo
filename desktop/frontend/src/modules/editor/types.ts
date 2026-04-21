export type EditorDocumentType = 'article' | 'node' | 'extract'

export interface EditorSelection {
  from: number
  to: number
  text: string
}

export interface EditorDocument {
  id: string
  title: string
  content: string
  contentType: EditorDocumentType
  updatedAt: string
  sourceUrl?: string
  tags?: string[]
}

export interface EditorStateSnapshot {
  originalContent: string
  currentContent: string
  dirty: boolean
  saving: boolean
}
