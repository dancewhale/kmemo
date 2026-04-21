import { INBOX_CAPTURE_DEFAULTS } from '@/shared/constants/inbox'
import { normalizeCaptureUrl } from '@/shared/utils/url'
import type { ReaderArticle } from '@/modules/reader/types'
import type { BlankArticleInput, PasteTextCaptureInput, UrlCaptureInput } from '../types'

function excerptSummary(text: string, max = 160): string {
  const t = text.replace(/\s+/g, ' ').trim()
  if (t.length <= max) {
    return t
  }
  return `${t.slice(0, max).trimEnd()}…`
}

function titleFromPasteBody(content: string): string {
  const line = content.replace(/\s+/g, ' ').trim()
  if (!line) {
    return INBOX_CAPTURE_DEFAULTS.untitledBlankTitle
  }
  const slice = line.slice(0, 36).trimEnd()
  return slice.length < line.length ? `${slice}…` : slice
}

function titleFromUrl(url: string, userTitle?: string): string {
  const trimmed = userTitle?.trim()
  if (trimmed) {
    return trimmed
  }
  try {
    const u = new URL(normalizeCaptureUrl(url))
    const path = u.pathname === '/' ? '' : u.pathname
    const tail = `${u.hostname}${path}`.slice(0, 72)
    return tail || INBOX_CAPTURE_DEFAULTS.capturedFromUrlTitle
  } catch {
    return INBOX_CAPTURE_DEFAULTS.capturedFromUrlTitle
  }
}

export function buildArticleFromBlankInput(input: BlankArticleInput, id: string): ReaderArticle {
  const now = new Date().toISOString()
  const title = input.title.trim() || INBOX_CAPTURE_DEFAULTS.untitledBlankTitle
  const content = input.content
  return {
    id,
    title,
    summary: excerptSummary(content || title),
    content,
    sourceType: 'text',
    status: 'inbox',
    updatedAt: now,
    tags: input.tags?.length ? [...input.tags] : undefined,
  }
}

export function buildArticleFromPasteTextInput(input: PasteTextCaptureInput, id: string): ReaderArticle {
  const now = new Date().toISOString()
  const content = input.content.trim()
  const title = input.title.trim() || titleFromPasteBody(content)
  return {
    id,
    title,
    summary: excerptSummary(content),
    content,
    sourceType: 'text',
    status: 'inbox',
    updatedAt: now,
    tags: input.tags?.length ? [...input.tags] : undefined,
  }
}

export function buildArticleFromUrlInput(input: UrlCaptureInput, id: string): ReaderArticle {
  const now = new Date().toISOString()
  const sourceUrl = normalizeCaptureUrl(input.url)
  const title = titleFromUrl(input.url, input.title)
  const note = input.note?.trim() ?? ''
  const content = [
    `Source URL: ${sourceUrl}`,
    note ? `Capture note: ${note}` : 'Capture note: —',
    '',
    INBOX_CAPTURE_DEFAULTS.urlPlaceholderLead,
  ].join('\n')
  const summary = note ? excerptSummary(note, 120) : excerptSummary(`URL capture: ${sourceUrl}`, 120)
  return {
    id,
    title,
    summary,
    content,
    sourceType: 'web',
    sourceUrl,
    status: 'inbox',
    updatedAt: now,
    tags: input.tags?.length ? [...input.tags] : undefined,
  }
}
