import type { CardDTO } from '@/api/types/dto'
import type { CreateCardRequest } from '@/api/types/requests'
import type { CreateCardInput } from '@/modules/object-creation/types'
import type { CardItem } from '../types'

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

/** Build minimal HTML body for backend storage from prompt + answer. */
export function buildCardHtmlContent(prompt: string, answer: string): string {
  const p = escapeHtml(prompt.trim())
  const a = escapeHtml(answer.trim())
  if (a) {
    return `<section class="kmemo-card"><p class="prompt">${p}</p><p class="answer">${a}</p></section>`
  }
  return `<section class="kmemo-card"><p class="prompt">${p}</p></section>`
}

function htmlToPlainText(html: string): string {
  if (!html.trim()) {
    return ''
  }
  if (typeof document === 'undefined') {
    return html.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim()
  }
  const div = document.createElement('div')
  div.innerHTML = html
  return div.textContent?.replace(/\s+/g, ' ').trim() ?? ''
}

export function mapCardDtoToCardItem(dto: CardDTO): CardItem {
  const plain = htmlToPlainText(dto.htmlContent)
  let prompt = plain
  let answer = ''
  const mid = Math.floor(plain.length / 2)
  if (plain.length > 120) {
    prompt = plain.slice(0, mid).trim()
    answer = plain.slice(mid).trim()
  }
  return {
    id: dto.id,
    nodeId: dto.id,
    title: dto.title,
    prompt,
    answer,
    parentNodeId: dto.parentId ?? dto.knowledgeId,
    createdAt: typeof dto.createdAt === 'string' ? dto.createdAt : String(dto.createdAt),
    updatedAt: typeof dto.updatedAt === 'string' ? dto.updatedAt : String(dto.updatedAt),
    reviewItemId: null,
  }
}

export function mapCreateCardInputToRequest(
  input: CreateCardInput,
  knowledgeId: string,
  options?: { cardType?: string },
): CreateCardRequest {
  return {
    knowledgeId,
    parentId: null,
    title: input.title.trim() || 'Untitled card',
    cardType: options?.cardType ?? 'qa',
    htmlContent: buildCardHtmlContent(input.prompt, input.answer),
    tagIds: [],
  }
}
