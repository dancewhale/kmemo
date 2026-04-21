import { invokeApp, safeWailsCall } from '@/api/core/api-guard'
import type { KnowledgeDTO, KnowledgeTreeNode } from '@/api/types/dto'
import type { CreateKnowledgeRequest, UpdateKnowledgeRequest } from '@/api/types/requests'
import type { ApiResult } from '@/api/core/api-result'

export function getKnowledgeTree(rootID: string | null): Promise<ApiResult<KnowledgeTreeNode[]>> {
  return safeWailsCall(() => invokeApp<KnowledgeTreeNode[]>('GetKnowledgeTree', rootID))
}

export function listKnowledge(parentID: string | null): Promise<ApiResult<KnowledgeDTO[]>> {
  return safeWailsCall(() => invokeApp<KnowledgeDTO[]>('ListKnowledge', parentID))
}

export function createKnowledge(req: CreateKnowledgeRequest): Promise<ApiResult<string>> {
  return safeWailsCall(() => invokeApp<string>('CreateKnowledge', req))
}

export function updateKnowledge(id: string, req: UpdateKnowledgeRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('UpdateKnowledge', id, req)
  })
}

export function deleteKnowledge(id: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('DeleteKnowledge', id)
  })
}

export function moveKnowledge(id: string, newParentID: string | null): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('MoveKnowledge', id, newParentID)
  })
}

export function archiveKnowledge(id: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('ArchiveKnowledge', id)
  })
}

export function unarchiveKnowledge(id: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('UnarchiveKnowledge', id)
  })
}

export function getKnowledge(id: string): Promise<ApiResult<KnowledgeDTO>> {
  return safeWailsCall(() => invokeApp<KnowledgeDTO>('GetKnowledge', id))
}
