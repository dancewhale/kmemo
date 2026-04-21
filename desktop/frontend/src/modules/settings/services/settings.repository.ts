import { getDefaultFSRSParameter, updateDefaultFSRSParameter } from '@/api/wails/system'
import type { ApiResult } from '@/api/core/api-result'
import type { UpdateDefaultFSRSParameterRequest } from '@/api/types/requests'
import { mapFsrsParameterDtoToView, type HostFsrsParameterView } from './settings.mapper'

export async function fetchDefaultFsrsParameter(): Promise<ApiResult<HostFsrsParameterView>> {
  const res = await getDefaultFSRSParameter()
  if (!res.ok) {
    return res
  }
  return { ok: true, data: mapFsrsParameterDtoToView(res.data) }
}

export async function persistDefaultFsrsParameter(
  req: UpdateDefaultFSRSParameterRequest,
): Promise<ApiResult<HostFsrsParameterView>> {
  const res = await updateDefaultFSRSParameter(req)
  if (!res.ok) {
    return res
  }
  return { ok: true, data: mapFsrsParameterDtoToView(res.data) }
}
