/**
 * Wails bridge — reserved for later host calls.
 * UI layers must not import generated Wails bindings directly; route calls through this module.
 */

export type WailsReady = false

export const wails: WailsReady = false

export async function whenWailsReady(): Promise<void> {
  /* no-op until host integration */
}
