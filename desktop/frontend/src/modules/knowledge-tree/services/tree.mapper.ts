import type { KnowledgeNode, UITreeNode } from '../types'

function sortByTitle(a: UITreeNode, b: UITreeNode): number {
  return a.title.localeCompare(b.title, undefined, { sensitivity: 'base' })
}

function assignLevels(nodes: UITreeNode[], level: number): void {
  for (const n of nodes) {
    n.level = level
    if (n.children.length) {
      assignLevels(n.children, level + 1)
    }
  }
}

/**
 * Pure: flat list → hierarchical UITreeNode roots (children populated, levels set).
 */
export function buildTree(nodes: KnowledgeNode[]): UITreeNode[] {
  const byId = new Map<string, UITreeNode>()
  for (const n of nodes) {
    byId.set(n.id, {
      ...n,
      level: 0,
      children: [],
      hasChildren: false,
      childCount: 0,
    })
  }

  const roots: UITreeNode[] = []
  for (const n of nodes) {
    const ui = byId.get(n.id)
    if (!ui) {
      continue
    }
    if (n.parentId == null || !byId.has(n.parentId)) {
      roots.push(ui)
    } else {
      byId.get(n.parentId)!.children.push(ui)
    }
  }

  function finalize(list: UITreeNode[]): void {
    list.sort(sortByTitle)
    for (const item of list) {
      item.childCount = item.children.length
      item.hasChildren = item.children.length > 0
      if (item.children.length) {
        finalize(item.children)
      }
    }
  }

  finalize(roots)
  assignLevels(roots, 0)
  return roots
}

/**
 * Pure: depth-first visible rows given expansion set (roots always shown).
 */
export function flattenVisibleTree(roots: UITreeNode[], expandedNodeIds: string[]): UITreeNode[] {
  const out: UITreeNode[] = []

  function walk(list: UITreeNode[]): void {
    for (const n of list) {
      out.push(n)
      if (expandedNodeIds.includes(n.id) && n.children.length) {
        walk(n.children)
      }
    }
  }

  walk(roots)
  return out
}
