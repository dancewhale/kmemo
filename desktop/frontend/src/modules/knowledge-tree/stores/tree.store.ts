import { defineStore } from 'pinia'
import { mockKnowledgeNodes } from '@/mock/tree'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { buildTree, flattenVisibleTree } from '../services/tree.mapper'
import type { KnowledgeNode, UITreeNode } from '../types'

interface TreeState {
  rawNodes: KnowledgeNode[]
  selectedNodeId: string | null
  expandedNodeIds: string[]
  searchKeyword: string
  loading: boolean
  initialized: boolean
}

function keywordMatch(node: KnowledgeNode, keyword: string): boolean {
  const k = keyword.trim().toLowerCase()
  if (!k) {
    return true
  }
  return node.title.toLowerCase().includes(k)
}

function filterRawByKeyword(nodes: KnowledgeNode[], keyword: string): KnowledgeNode[] {
  const k = keyword.trim()
  if (!k) {
    return nodes
  }
  const matched = new Set<string>()
  for (const n of nodes) {
    if (keywordMatch(n, k)) {
      matched.add(n.id)
      let p = n.parentId
      while (p) {
        matched.add(p)
        const parent = nodes.find((x) => x.id === p)
        p = parent?.parentId ?? null
      }
    }
  }
  return nodes.filter((n) => matched.has(n.id))
}

export const useTreeStore = defineStore('knowledge-tree', {
  state: (): TreeState => ({
    rawNodes: [],
    selectedNodeId: null,
    expandedNodeIds: [],
    searchKeyword: '',
    loading: false,
    initialized: false,
  }),

  getters: {
    treeNodes(): UITreeNode[] {
      const filtered = filterRawByKeyword(this.rawNodes, this.searchKeyword)
      return buildTree(filtered)
    },

    rootNodes(): UITreeNode[] {
      return this.treeNodes
    },

    /** When filtering, force-expand ancestor chains so matches stay visible. */
    effectiveExpandedNodeIds(): string[] {
      const k = this.searchKeyword.trim()
      if (!k) {
        return [...this.expandedNodeIds]
      }
      const filtered = filterRawByKeyword(this.rawNodes, this.searchKeyword)
      const s = new Set<string>()
      for (const n of filtered) {
        let p = n.parentId
        while (p) {
          s.add(p)
          const parent = this.rawNodes.find((r) => r.id === p)
          p = parent?.parentId ?? null
        }
      }
      return [...s]
    },

    visibleNodes(): UITreeNode[] {
      return flattenVisibleTree(this.treeNodes, this.effectiveExpandedNodeIds)
    },

    selectedNode(): KnowledgeNode | null {
      if (!this.selectedNodeId) {
        return null
      }
      return this.rawNodes.find((n) => n.id === this.selectedNodeId) ?? null
    },

    selectedNodeChildCount(): number {
      if (!this.selectedNodeId) {
        return 0
      }
      return this.rawNodes.filter((n) => n.parentId === this.selectedNodeId).length
    },
    findNodeById(): (id: string | null) => KnowledgeNode | null {
      return (id: string | null) => {
        if (!id) {
          return null
        }
        return this.rawNodes.find((n) => n.id === id) ?? null
      }
    },
    isNodeRoot(): (nodeId: string) => boolean {
      return (nodeId: string) => {
        const n = this.rawNodes.find((x) => x.id === nodeId)
        return n?.parentId == null
      }
    },
    isSelectedExtractNode(): boolean {
      return this.selectedNode?.type === 'extract'
    },
  },

  actions: {
    async initialize() {
      if (this.initialized) {
        return
      }
      this.loading = true
      await Promise.resolve()
      this.rawNodes = mockKnowledgeNodes.map((n) => ({ ...n }))
      this.expandedNodeIds = [
        'kn-reading',
        'kn-knowledge',
        'kn-mem-topic',
        'kn-rw-topic',
        'kn-sr-topic',
      ]
      this.initialized = true
      this.loading = false
    },

    setSelectedNode(id: string | null) {
      this.selectedNodeId = id
      useWorkspaceStore().selectNode(id)
    },

    toggleNode(id: string) {
      const i = this.expandedNodeIds.indexOf(id)
      if (i >= 0) {
        this.expandedNodeIds.splice(i, 1)
      } else {
        this.expandedNodeIds.push(id)
      }
    },

    expandNode(id: string) {
      if (!this.expandedNodeIds.includes(id)) {
        this.expandedNodeIds.push(id)
      }
    },

    collapseNode(id: string) {
      const i = this.expandedNodeIds.indexOf(id)
      if (i >= 0) {
        this.expandedNodeIds.splice(i, 1)
      }
    },

    setSearchKeyword(keyword: string) {
      this.searchKeyword = keyword
    },

    updateNodeTitle(nodeId: string, title: string) {
      const node = this.rawNodes.find((item) => item.id === nodeId)
      if (!node) {
        return
      }
      node.title = title
      node.updatedAt = new Date().toISOString()
    },

    addExtractNode(payload: {
      parentNodeId: string | null
      title: string
      sourceArticleId: string
      sourceExtractId?: string
      nodeId?: string
    }) {
      const now = new Date().toISOString()
      const id = payload.nodeId ?? `kn-ext-${Date.now().toString(36)}`
      const parentExists =
        payload.parentNodeId == null ||
        this.rawNodes.some((node) => node.id === payload.parentNodeId)
      const parentId = parentExists ? payload.parentNodeId : null

      this.rawNodes.push({
        id,
        parentId,
        title: payload.title,
        type: 'extract',
        sourceArticleId: payload.sourceArticleId,
        sourceExtractId: payload.sourceExtractId,
        createdAt: now,
        updatedAt: now,
      })

      if (parentId) {
        this.expandNode(parentId)
      }
      this.setSelectedNode(id)
    },

    async ensureReadyWithSelection(preferredId: string | null) {
      await this.initialize()
      const exists = preferredId && this.rawNodes.some((n) => n.id === preferredId)
      if (exists) {
        this.setSelectedNode(preferredId)
        return
      }
      const first = this.visibleNodes[0] ?? this.rawNodes[0]
      this.setSelectedNode(first?.id ?? null)
    },

    addChildTopicNode(parentId: string): string | null {
      const parent = this.rawNodes.find((n) => n.id === parentId)
      if (!parent) {
        return null
      }
      const id = `kn-new-${Date.now().toString(36)}`
      const now = new Date().toISOString()
      this.rawNodes.push({
        id,
        parentId: parentId,
        title: 'New Node',
        type: 'topic',
        createdAt: now,
        updatedAt: now,
      })
      this.expandNode(parentId)
      this.setSelectedNode(id)
      return id
    },

    addNode(node: KnowledgeNode) {
      let next = { ...node }
      if (next.parentId != null && !this.rawNodes.some((n) => n.id === next.parentId)) {
        next = { ...next, parentId: null }
      }
      this.rawNodes.push(next)
      if (next.parentId) {
        this.expandNode(next.parentId)
      }
      this.setSelectedNode(next.id)
    },

    removeNodeSubtree(nodeId: string): boolean {
      const node = this.rawNodes.find((n) => n.id === nodeId)
      if (!node || node.parentId == null) {
        return false
      }
      const toRemove = new Set<string>()
      const stack = [nodeId]
      while (stack.length > 0) {
        const id = stack.pop()!
        toRemove.add(id)
        for (const n of this.rawNodes) {
          if (n.parentId === id) {
            stack.push(n.id)
          }
        }
      }
      this.rawNodes = this.rawNodes.filter((n) => !toRemove.has(n.id))
      if (this.selectedNodeId && toRemove.has(this.selectedNodeId)) {
        const next = this.rawNodes[0]?.id ?? null
        this.setSelectedNode(next)
      }
      return true
    },
  },
})
