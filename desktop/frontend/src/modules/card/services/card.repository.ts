import { createCard, getCard, listCards, updateCard } from '@/api/wails/card'
import type { ApiResult } from '@/api/core/api-result'
import type { CardFilters } from '@/api/types/dto'
import type { CreateCardInput } from '@/modules/object-creation/types'
import type { CardItem } from '../types'
import type { CardUpdatePayload } from '../types'
import { buildCardHtmlContent, mapCardDtoToCardItem, mapCreateCardInputToRequest } from './card.api-mapper'

export async function fetchCards(filters: CardFilters): Promise<ApiResult<{ items: CardItem[]; total: number }>> {
  const res = await listCards(filters)
  if (!res.ok) {
    return res
  }
  return {
    ok: true,
    data: {
      items: res.data.items.map((x) => mapCardDtoToCardItem(x)),
      total: res.data.total,
    },
  }
}

export async function fetchCard(id: string): Promise<ApiResult<CardItem>> {
  const res = await getCard(id)
  if (!res.ok) {
    return res
  }
  return { ok: true, data: mapCardDtoToCardItem(res.data) }
}

export async function createCardFromInput(
  input: CreateCardInput,
  knowledgeId: string,
): Promise<ApiResult<string>> {
  const req = mapCreateCardInputToRequest(input, knowledgeId)
  return createCard(req)
}

export async function persistCardUpdate(
  id: string,
  payload: CardUpdatePayload,
  current: CardItem,
): Promise<ApiResult<void>> {
  const title = payload.title ?? current.title
  const prompt = payload.prompt ?? current.prompt
  const answer = payload.answer ?? current.answer
  return updateCard(id, {
    title,
    htmlContent: buildCardHtmlContent(prompt, answer),
    status: 'active',
  })
}
