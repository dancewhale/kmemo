import { invokeApp, safeWailsCall } from '@/api/core/api-guard'
import type { ApiResult } from '@/api/core/api-result'
import type { CardDTO, CardDetailDTO, CardFilters, ListCardsResult } from '@/api/types/dto'
import type {
  CreateCardRequest,
  MoveCardRequest,
  ReorderCardChildrenRequest,
  UpdateCardRequest,
} from '@/api/types/requests'

export function listCards(filters: CardFilters): Promise<ApiResult<ListCardsResult>> {
  return safeWailsCall(() => invokeApp<ListCardsResult>('ListCards', filters))
}

export function getCard(id: string): Promise<ApiResult<CardDTO>> {
  return safeWailsCall(() => invokeApp<CardDTO>('GetCard', id))
}

export function getCardDetail(id: string): Promise<ApiResult<CardDetailDTO>> {
  return safeWailsCall(() => invokeApp<CardDetailDTO>('GetCardDetail', id))
}

export function createCard(req: CreateCardRequest): Promise<ApiResult<string>> {
  return safeWailsCall(() => invokeApp<string>('CreateCard', req))
}

export function updateCard(id: string, req: UpdateCardRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('UpdateCard', id, req)
  })
}

export function deleteCard(id: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('DeleteCard', id)
  })
}

export function getCardChildren(parentID: string): Promise<ApiResult<CardDTO[]>> {
  return safeWailsCall(() => invokeApp<CardDTO[]>('GetCardChildren', parentID))
}

export function moveCard(req: MoveCardRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('MoveCard', req)
  })
}

export function reorderCardChildren(req: ReorderCardChildrenRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('ReorderCardChildren', req)
  })
}

export function addCardTags(cardID: string, tagIDs: string[]): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('AddCardTags', cardID, tagIDs)
  })
}

export function removeCardTags(cardID: string, tagIDs: string[]): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('RemoveCardTags', cardID, tagIDs)
  })
}

export function suspendCard(cardID: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('SuspendCard', cardID)
  })
}

export function resumeCard(cardID: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('ResumeCard', cardID)
  })
}
