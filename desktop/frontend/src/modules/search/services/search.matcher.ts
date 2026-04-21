import type { SearchDocument, SearchFilterType, SearchResultItem } from '../types'

function include(haystack: string | undefined | null, needle: string): boolean {
  if (!haystack) {
    return false
  }
  return haystack.toLowerCase().includes(needle)
}

function pickSnippet(content: string, query: string): string {
  const text = content.replace(/\s+/g, ' ').trim()
  if (!text) {
    return ''
  }
  if (!query) {
    return text.slice(0, 128)
  }
  const lower = text.toLowerCase()
  const idx = lower.indexOf(query)
  if (idx < 0) {
    return text.slice(0, 128)
  }
  const start = Math.max(0, idx - 36)
  const end = Math.min(text.length, idx + query.length + 68)
  return `${start > 0 ? '…' : ''}${text.slice(start, end)}${end < text.length ? '…' : ''}`
}

function sortByUpdatedDesc(items: SearchResultItem[]): SearchResultItem[] {
  return [...items].sort((a, b) => {
    const ta = a.updatedAt ? new Date(a.updatedAt).getTime() : 0
    const tb = b.updatedAt ? new Date(b.updatedAt).getTime() : 0
    return tb - ta || a.title.localeCompare(b.title)
  })
}

export function matchSearchDocuments(
  docs: SearchDocument[],
  query: string,
  filter: SearchFilterType,
): SearchResultItem[] {
  const q = query.trim().toLowerCase()
  const pool = docs.filter((doc) => filter === 'all' || doc.entityType === filter)

  if (!q) {
    return sortByUpdatedDesc(
      pool.map((doc) => ({
        id: doc.id,
        entityType: doc.entityType,
        title: doc.title,
        snippet: pickSnippet(doc.content, ''),
        subtitle: doc.subtitle,
        score: 0,
        sourceId: doc.sourceId,
        sourceTitle: doc.sourceTitle,
        nodeId: doc.nodeId,
        articleId: doc.articleId,
        extractId: doc.extractId,
        reviewId: doc.reviewId,
        cardId: doc.cardId,
        updatedAt: doc.updatedAt,
      })),
    ).slice(0, 28)
  }

  const out: SearchResultItem[] = []
  for (const doc of pool) {
    let score = 0
    if (include(doc.title, q)) {
      score += 6
    }
    if (include(doc.subtitle, q)) {
      score += 3
    }
    if (include(doc.content, q)) {
      score += 1
    }
    if ((doc.keywords ?? []).some((k) => include(k, q))) {
      score += 2
    }
    if (score > 0) {
      out.push({
        id: doc.id,
        entityType: doc.entityType,
        title: doc.title,
        snippet: pickSnippet(doc.content, q),
        subtitle: doc.subtitle,
        score,
        sourceId: doc.sourceId,
        sourceTitle: doc.sourceTitle,
        nodeId: doc.nodeId,
        articleId: doc.articleId,
        extractId: doc.extractId,
        reviewId: doc.reviewId,
        cardId: doc.cardId,
        updatedAt: doc.updatedAt,
      })
    }
  }

  return out.sort((a, b) => b.score - a.score || a.title.localeCompare(b.title))
}
