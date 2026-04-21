import type { KnowledgeNodeType } from '@/shared/types/common'

export interface MockKnowledgeNode {
  id: string
  parentId: string | null
  title: string
  type: KnowledgeNodeType
}

export const mockKnowledgeNodes: MockKnowledgeNode[] = [
  { id: 'kn-root', parentId: null, title: 'Library', type: 'folder' },
  { id: 'kn-01', parentId: 'kn-root', title: 'Incremental reading', type: 'folder' },
  { id: 'kn-02', parentId: 'kn-01', title: 'Queue discipline', type: 'note' },
  { id: 'kn-03', parentId: 'kn-01', title: 'Highlight → item', type: 'card' },
  { id: 'kn-04', parentId: 'kn-root', title: 'Engineering', type: 'folder' },
  { id: 'kn-05', parentId: 'kn-04', title: 'Wails shell', type: 'note' },
  { id: 'kn-06', parentId: 'kn-04', title: 'Storage model', type: 'card' },
  { id: 'kn-07', parentId: 'kn-04', title: 'Python worker', type: 'note' },
  { id: 'kn-08', parentId: 'kn-root', title: 'Review deck', type: 'folder' },
]
