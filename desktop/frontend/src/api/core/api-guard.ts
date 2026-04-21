import type { ApiErrorCode } from './api-error'
import { classifyWailsError } from './api-error'
import type { ApiFailure, ApiResult, ApiSuccess } from './api-result'

export function ok<T>(data: T): ApiSuccess<T> {
  return { ok: true, data }
}

export function fail(code: ApiErrorCode, message: string, details?: unknown): ApiFailure {
  return { ok: false, error: { code, message, details } }
}

export type WailsWindow = {
  go?: {
    main?: {
      App?: Record<string, (...args: unknown[]) => Promise<unknown>>
    }
  }
}

export function isWailsAvailable(): boolean {
  if (typeof window === 'undefined') {
    return false
  }
  const App = (window as unknown as WailsWindow).go?.main?.App
  return Boolean(App && typeof App === 'object')
}

export async function invokeApp<T>(method: string, ...args: unknown[]): Promise<T> {
  const App = (window as unknown as WailsWindow).go?.main?.App
  if (!App || typeof App[method] !== 'function') {
    throw new Error(`Wails App.${method} is not available`)
  }
  return (await App[method](...args)) as T
}

export async function safeWailsCall<T>(fn: () => Promise<T>): Promise<ApiResult<T>> {
  try {
    const data = await fn()
    return ok(data)
  } catch (err) {
    if (err instanceof Error && err.message.includes('not available')) {
      return fail('WAILS_NOT_READY', err.message, err)
    }
    return { ok: false, error: classifyWailsError(err) }
  }
}
