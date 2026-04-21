import type { KnowledgeNode } from '@/modules/knowledge-tree/types'
import type { EditorSelection } from '@/modules/editor/types'
import type { ExtractItem, PendingExtractPayload } from '../types'

function compactText(text: string): string {
  return text.replace(/\s+/g, ' ').trim()
}

function defaultExtractTitle(sourceTitle: string, quote: string): string {
  const normalized = compactText(quote)
  if (normalized.length > 0) {
    return normalized.slice(0, 36)
  }
  return `Extract from ${sourceTitle}`
}

export function buildPendingExtractFromSelection(input: {
  sourceArticleId: string
  sourceArticleTitle: string
  selection: EditorSelection
  defaultParentNodeId?: string | null
}): PendingExtractPayload {
  const quote = compactText(input.selection.text)
  return {
    sourceArticleId: input.sourceArticleId,
    sourceArticleTitle: input.sourceArticleTitle,
    quote,
    note: '',
    title: defaultExtractTitle(input.sourceArticleTitle, quote),
    parentNodeId: input.defaultParentNodeId ?? null,
    selection: {
      from: input.selection.from,
      to: input.selection.to,
      text: quote,
    },
  }
}

export function buildTreeNodeFromExtract(extract: ExtractItem): KnowledgeNode {
  return {
    id: extract.treeNodeId ?? `kn-ext-${extract.id}`,
    parentId: extract.parentNodeId,
    title: extract.title,
    type: 'extract',
    sourceArticleId: extract.sourceArticleId,
    sourceExtractId: extract.id,
    createdAt: extract.createdAt,
    updatedAt: extract.updatedAt,
  }
}
