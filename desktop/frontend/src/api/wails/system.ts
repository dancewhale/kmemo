import { invokeApp, safeWailsCall } from '@/api/core/api-guard'
import type { ApiResult } from '@/api/core/api-result'
import type { FSRSParameterDTO } from '@/api/types/dto'
import type { UpdateDefaultFSRSParameterRequest } from '@/api/types/requests'

export function listFSRSParameters(): Promise<ApiResult<FSRSParameterDTO[]>> {
  return safeWailsCall(() => invokeApp<FSRSParameterDTO[]>('ListFSRSParameters'))
}

export function getDefaultFSRSParameter(): Promise<ApiResult<FSRSParameterDTO>> {
  return safeWailsCall(() => invokeApp<FSRSParameterDTO>('GetDefaultFSRSParameter'))
}

export function updateDefaultFSRSParameter(
  req: UpdateDefaultFSRSParameterRequest,
): Promise<ApiResult<FSRSParameterDTO>> {
  return safeWailsCall(() => invokeApp<FSRSParameterDTO>('UpdateDefaultFSRSParameter', req))
}
