export type ApiErrorCode =
  | 'UNKNOWN'
  | 'NOT_FOUND'
  | 'VALIDATION_ERROR'
  | 'CONFLICT'
  | 'EMPTY'
  | 'BACKEND_ERROR'
  | 'WAILS_NOT_READY'

export interface ApiError {
  code: ApiErrorCode
  message: string
  details?: unknown
}

export function classifyWailsError(err: unknown): ApiError {
  const msg = err instanceof Error ? err.message : String(err)
  const lower = msg.toLowerCase()
  if (lower.includes('not found') || lower.includes('record not found')) {
    return { code: 'NOT_FOUND', message: msg, details: err }
  }
  if (lower.includes('invalid input') || lower.includes('validation') || lower.includes('invalid')) {
    return { code: 'VALIDATION_ERROR', message: msg, details: err }
  }
  if (lower.includes('conflict')) {
    return { code: 'CONFLICT', message: msg, details: err }
  }
  return { code: 'BACKEND_ERROR', message: msg, details: err }
}
