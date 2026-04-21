import { OBJECT_CREATION_DEFAULTS } from '@/shared/constants/object-creation'
import type { KnowledgeNode } from '@/modules/knowledge-tree/types'
import type { ExtractItem } from '@/modules/extract/types'
import type { CreateCardInput, CreateManualExtractInput, CreateNodeInput } from '../types'

function compactLine(text: string, max: number): string {
  const t = text.replace(/\s+/g, ' ').trim()
  if (t.length <= max) {
    return t
  }
  return `${t.slice(0, max).trimEnd()}…`
}

function titleFromPrompt(prompt: string): string {
  const t = prompt.replace(/\s+/g, ' ').trim()
  if (!t) {
    return OBJECT_CREATION_DEFAULTS.untitledCard
  }
  return compactLine(t, 40)
}

function titleFromQuote(quote: string): string {
  const t = quote.replace(/\s+/g, ' ').trim()
  if (!t) {
    return 'Untitled extract'
  }
  return compactLine(t, 40)
}

export function buildTreeNodeFromCreateNodeInput(input: CreateNodeInput, nodeId: string): KnowledgeNode {
  const now = new Date().toISOString()
  const title = input.title.trim() || OBJECT_CREATION_DEFAULTS.untitledNode
  const desc = input.content.trim()
  return {
    id: nodeId,
    parentId: input.parentNodeId,
    title,
    type: 'topic',
    description: desc || undefined,
    createdAt: now,
    updatedAt: now,
  }
}

export function buildTreeNodeFromCreateCardInput(input: CreateCardInput, nodeId: string): KnowledgeNode {
  const now = new Date().toISOString()
  const title = input.title.trim() || titleFromPrompt(input.prompt)
  return {
    id: nodeId,
    parentId: input.parentNodeId,
    title,
    type: 'card',
    createdAt: now,
    updatedAt: now,
  }
}

export function buildExtractFromManualExtractInput(
  input: CreateManualExtractInput,
  extractId: string,
  treeNodeId: string,
): ExtractItem {
  const now = new Date().toISOString()
  const quote = input.quote.trim()
  const title = input.title.trim() || titleFromQuote(quote)
  const sid = input.sourceArticleId?.trim() || ''
  const stitle = input.sourceArticleTitle?.trim() || (sid ? 'Linked article' : 'Manual extract')
  return {
    id: extractId,
    sourceArticleId: sid || 'manual',
    sourceArticleTitle: stitle,
    title,
    quote,
    note: input.note.trim(),
    parentNodeId: input.parentNodeId,
    treeNodeId,
    createdAt: now,
    updatedAt: now,
  }
}

export function buildTreeNodeFromManualExtract(extract: ExtractItem): KnowledgeNode {
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

