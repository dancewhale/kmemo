import { onUnmounted, ref } from 'vue'

export interface UsePaneResizeOptions {
  onMove: (delta: number) => void
  onEnd?: () => void
}

export function usePaneResize(options: UsePaneResizeOptions) {
  const isDragging = ref(false)
  let startCoord = 0

  function onPointerMoveVertical(e: PointerEvent) {
    if (!isDragging.value) {
      return
    }
    const delta = e.clientX - startCoord
    startCoord = e.clientX
    options.onMove(delta)
  }

  function onPointerMoveHorizontal(e: PointerEvent) {
    if (!isDragging.value) {
      return
    }
    const delta = e.clientY - startCoord
    startCoord = e.clientY
    options.onMove(delta)
  }

  function endVertical() {
    window.removeEventListener('pointermove', onPointerMoveVertical)
    window.removeEventListener('pointerup', endVertical)
    window.removeEventListener('pointercancel', endVertical)
    isDragging.value = false
    options.onEnd?.()
  }

  function endHorizontal() {
    window.removeEventListener('pointermove', onPointerMoveHorizontal)
    window.removeEventListener('pointerup', endHorizontal)
    window.removeEventListener('pointercancel', endHorizontal)
    isDragging.value = false
    options.onEnd?.()
  }

  function startVertical(e: PointerEvent) {
    if (isDragging.value) {
      return
    }
    e.preventDefault()
    isDragging.value = true
    startCoord = e.clientX
    window.addEventListener('pointermove', onPointerMoveVertical)
    window.addEventListener('pointerup', endVertical)
    window.addEventListener('pointercancel', endVertical)
  }

  function startHorizontal(e: PointerEvent) {
    if (isDragging.value) {
      return
    }
    e.preventDefault()
    isDragging.value = true
    startCoord = e.clientY
    window.addEventListener('pointermove', onPointerMoveHorizontal)
    window.addEventListener('pointerup', endHorizontal)
    window.addEventListener('pointercancel', endHorizontal)
  }

  onUnmounted(() => {
    endVertical()
    endHorizontal()
  })

  return {
    isDragging,
    startVertical,
    startHorizontal,
  }
}
