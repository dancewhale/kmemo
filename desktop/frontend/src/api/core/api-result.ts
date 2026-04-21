import type { ApiError } from './api-error'

export type ApiSuccess<T> = { ok: true; data: T }
export type ApiFailure = { ok: false; error: ApiError }
export type ApiResult<T> = ApiSuccess<T> | ApiFailure
