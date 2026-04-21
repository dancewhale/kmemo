import { listCards } from '@/api/wails/card'
import { searchCardsByTags } from '@/api/wails/tag'
import type { ApiResult } from '@/api/core/api-result'
import type { CardDTO } from '@/api/types/dto'

export async function searchCardsByKeyword(keyword: string, limit = 50): Promise<ApiResult<CardDTO[]>> {
  const kw = keyword.trim()
  if (!kw) {
    return { ok: true, data: [] }
  }
  const res = await listCards({ keyword: kw, limit, offset: 0 })
  if (!res.ok) {
    return res
  }
  return { ok: true, data: res.data.items }
}

export async function searchCardsByTagIds(tagIDs: string[]): Promise<ApiResult<CardDTO[]>> {
  if (!tagIDs.length) {
    return { ok: true, data: [] }
  }
  return searchCardsByTags(tagIDs)
}
