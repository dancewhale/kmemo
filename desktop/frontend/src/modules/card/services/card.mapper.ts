import type { KnowledgeNode } from '@/modules/knowledge-tree/types'
import type { CardItem } from '../types'

/**
 * Build a CardItem from a knowledge tree node of type `card`.
 * Title-only mock nodes get prompt seeded from "Card: …" title when present.
 */
export function buildCardItemFromTreeNode(node: KnowledgeNode): CardItem {
  if (node.type !== 'card') {
    throw new Error('buildCardItemFromTreeNode: node type must be card')
  }
  const now = node.updatedAt ?? node.createdAt ?? new Date().toISOString()
  let prompt = ''
  let answer = ''
  const title = node.title
  const m = title.match(/^Card:\s*(.+)$/i)
  if (m) {
    prompt = m[1].trim()
  }
  return {
    id: `card-${node.id}`,
    nodeId: node.id,
    title,
    prompt,
    answer,
    parentNodeId: node.parentId ?? null,
    createdAt: node.createdAt ?? now,
    updatedAt: now,
    reviewItemId: null,
  }
}
