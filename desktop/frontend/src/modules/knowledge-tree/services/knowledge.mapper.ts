import type { KnowledgeTreeNode } from '@/api/types/dto'
import type { KnowledgeNode } from '../types'

function toIso(value: string | Date): string {
  if (typeof value === 'string') {
    return value
  }
  return value.toISOString()
}

/**
 * Map one backend knowledge tree node to a flat `KnowledgeNode`.
 * Uses `type: 'topic'` so existing tree UI keeps working.
 */
export function knowledgeTreeNodeToKnowledgeNode(node: KnowledgeTreeNode): KnowledgeNode {
  return {
    id: node.id,
    parentId: node.parentId,
    title: node.name,
    type: 'topic',
    description: node.description?.trim() ? node.description : undefined,
    createdAt: toIso(node.createdAt),
    updatedAt: toIso(node.updatedAt),
  }
}

function walk(nodes: KnowledgeTreeNode[], out: KnowledgeNode[]): void {
  for (const n of nodes) {
    out.push(knowledgeTreeNodeToKnowledgeNode(n))
    if (n.children?.length) {
      walk(n.children, out)
    }
  }
}

/** Depth-first flatten of backend forest → flat list for `buildTree`. */
export function knowledgeForestToFlatNodes(roots: KnowledgeTreeNode[]): KnowledgeNode[] {
  const out: KnowledgeNode[] = []
  walk(roots, out)
  return out
}
