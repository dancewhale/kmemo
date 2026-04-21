import { COMMAND_IDS } from '@/shared/constants/commands'
import { useCommandStore } from '@/modules/command-center/stores/command.store'

let installed = false
let pendingGoPrefix = false
let pendingGoTimer: number | null = null

function isEditableTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false
  }
  const tag = target.tagName.toLowerCase()
  if (tag === 'input' || tag === 'textarea' || tag === 'select') {
    return true
  }
  return target.isContentEditable
}

function clearGoPrefixTimer() {
  if (pendingGoTimer != null) {
    window.clearTimeout(pendingGoTimer)
    pendingGoTimer = null
  }
}

function armGoPrefix() {
  pendingGoPrefix = true
  clearGoPrefixTimer()
  pendingGoTimer = window.setTimeout(() => {
    pendingGoPrefix = false
    pendingGoTimer = null
  }, 900)
}

function clearGoPrefix() {
  pendingGoPrefix = false
  clearGoPrefixTimer()
}

export function setupShortcuts(): void {
  if (installed) {
    return
  }
  installed = true
  window.addEventListener('keydown', onKeyDown)
}

async function onKeyDown(e: KeyboardEvent): Promise<void> {
  const commandStore = useCommandStore()
  const key = e.key.toLowerCase()
  const isMeta = e.metaKey || e.ctrlKey
  const inInput = isEditableTarget(e.target)

  if (isMeta && key === 'k') {
    e.preventDefault()
    await commandStore.executeById(COMMAND_IDS.toggleCommandPalette)
    return
  }
  if (isMeta && key === 'f') {
    e.preventDefault()
    await commandStore.executeById(COMMAND_IDS.focusSearch)
    return
  }

  if (commandStore.isOpen) {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      commandStore.moveDown()
      return
    }
    if (e.key === 'ArrowUp') {
      e.preventDefault()
      commandStore.moveUp()
      return
    }
    if (e.key === 'Enter') {
      e.preventDefault()
      await commandStore.executeActiveCommand()
      return
    }
    if (e.key === 'Escape') {
      e.preventDefault()
      commandStore.close()
      return
    }
  }

  if (inInput) {
    return
  }

  if (e.altKey && !isMeta && !e.shiftKey) {
    if (key === '1') {
      e.preventDefault()
      await commandStore.executeById(COMMAND_IDS.goInbox)
      return
    }
    if (key === '2') {
      e.preventDefault()
      await commandStore.executeById(COMMAND_IDS.goReading)
      return
    }
    if (key === '3') {
      e.preventDefault()
      await commandStore.executeById(COMMAND_IDS.goKnowledge)
      return
    }
    if (key === '4') {
      e.preventDefault()
      await commandStore.executeById(COMMAND_IDS.goReview)
      return
    }
    if (key === '5') {
      e.preventDefault()
      await commandStore.executeById(COMMAND_IDS.goSearch)
      return
    }
  }

  if (isMeta && e.shiftKey && key === 'r') {
    e.preventDefault()
    await commandStore.executeById(COMMAND_IDS.openFirstReviewItem)
    return
  }

  if (!e.altKey && !e.shiftKey && !isMeta && key === 'g') {
    armGoPrefix()
    return
  }

  if (pendingGoPrefix && !e.altKey && !e.shiftKey && !isMeta) {
    if (key === 'r') {
      e.preventDefault()
      clearGoPrefix()
      await commandStore.executeById(COMMAND_IDS.goReading)
      return
    }
    if (key === 'k') {
      e.preventDefault()
      clearGoPrefix()
      await commandStore.executeById(COMMAND_IDS.goKnowledge)
      return
    }
    if (key === 'v') {
      e.preventDefault()
      clearGoPrefix()
      await commandStore.executeById(COMMAND_IDS.goReview)
      return
    }
    clearGoPrefix()
  }
}
