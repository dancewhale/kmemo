import { invokeApp, safeWailsCall } from '@/api/core/api-guard'
import type { ApiResult } from '@/api/core/api-result'
import type { CardDTO, TagDTO } from '@/api/types/dto'
import type { CreateTagRequest, UpdateTagRequest } from '@/api/types/requests'

export function listTags(): Promise<ApiResult<TagDTO[]>> {
  return safeWailsCall(() => invokeApp<TagDTO[]>('ListTags'))
}

export function createTag(req: CreateTagRequest): Promise<ApiResult<string>> {
  return safeWailsCall(() => invokeApp<string>('CreateTag', req))
}

export function updateTag(id: string, req: UpdateTagRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('UpdateTag', id, req)
  })
}

export function deleteTag(id: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('DeleteTag', id)
  })
}

export function getTag(id: string): Promise<ApiResult<TagDTO>> {
  return safeWailsCall(() => invokeApp<TagDTO>('GetTag', id))
}

export function searchCardsByTags(tagIDs: string[]): Promise<ApiResult<CardDTO[]>> {
  return safeWailsCall(() => invokeApp<CardDTO[]>('SearchCardsByTags', tagIDs))
}
