/**
 * Wails host bridge — all Go `Desktop` methods are invoked via `window.go.main.App`.
 * UI layers must not import Wails bindings directly; route calls through this module.
 */

import { isWailsAvailable } from '@/api/core/api-guard'

export { isWailsAvailable, invokeApp, safeWailsCall } from '@/api/core/api-guard'
export * from '@/api/wails/knowledge'
export * from '@/api/wails/card'
export * from '@/api/wails/review'
export * from '@/api/wails/tag'
export * from '@/api/wails/system'

export type WailsReady = boolean

/** True when running inside Wails and `window.go.main.App` is bound. */
export function getWailsHostReady(): WailsReady {
  return isWailsAvailable()
}

export async function whenWailsReady(): Promise<void> {
  if (typeof window === 'undefined') {
    return
  }
  await new Promise<void>((resolve) => {
    if (isWailsAvailable()) {
      resolve()
      return
    }
    const max = 40
    let i = 0
    const t = window.setInterval(() => {
      i += 1
      if (isWailsAvailable() || i >= max) {
        window.clearInterval(t)
        resolve()
      }
    }, 50)
  })
}
