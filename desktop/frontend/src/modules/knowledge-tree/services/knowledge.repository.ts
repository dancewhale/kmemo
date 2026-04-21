import {
  getKnowledgeTree,
  createKnowledge as createKnowledgeWails,
  updateKnowledge as updateKnowledgeWails,
  deleteKnowledge as deleteKnowledgeWails,
} from '@/api/wails/knowledge'
import type { ApiResult } from '@/api/core/api-result'
import type { CreateKnowledgeRequest, UpdateKnowledgeRequest } from '@/api/types/requests'
import type { KnowledgeNode } from '../types'
import { knowledgeForestToFlatNodes } from './knowledge.mapper'

export async function fetchKnowledgeTree(rootId: string | null): Promise<ApiResult<KnowledgeNode[]>> {
  const res = await getKnowledgeTree(rootId)
  if (!res.ok) {
    return res
  }
  return { ok: true, data: knowledgeForestToFlatNodes(res.data) }
}

export async function createKnowledgeSpace(
  name: string,
  description: string,
  parentId: string | null,
): Promise<ApiResult<string>> {
  const req: CreateKnowledgeRequest = {
    name: name.trim(),
    description: description.trim(),
    parentId,
  }
  return createKnowledgeWails(req)
}

export async function updateKnowledgeSpace(
  id: string,
  name: string,
  description: string,
): Promise<ApiResult<void>> {
  const req: UpdateKnowledgeRequest = { name: name.trim(), description: description.trim() }
  return updateKnowledgeWails(id, req)
}

export async function removeKnowledgeSpace(id: string): Promise<ApiResult<void>> {
  return deleteKnowledgeWails(id)
}
