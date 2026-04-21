/**
 * Global shortcut registration placeholder.
 * Keep handlers thin; route to stores/services when wiring the host.
 */

export function setupShortcuts(): void {
  window.addEventListener('keydown', onKeyDown)
}

function onKeyDown(e: KeyboardEvent): void {
  if ((e.metaKey || e.ctrlKey) && e.key === '1') {
    e.preventDefault()
  }
}
