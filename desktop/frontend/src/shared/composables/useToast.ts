import { ref } from 'vue'
import { UI_TOAST_DURATION_MS } from '@/shared/constants/ui'

type ToastKind = 'success' | 'info' | 'warning' | 'error'

export interface ToastItem {
  id: string
  kind: ToastKind
  message: string
}

const toasts = ref<ToastItem[]>([])

function removeToast(id: string) {
  toasts.value = toasts.value.filter((item) => item.id !== id)
}

function pushToast(kind: ToastKind, message: string, duration = UI_TOAST_DURATION_MS) {
  const id = `toast-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 7)}`
  toasts.value.push({ id, kind, message })
  window.setTimeout(() => removeToast(id), duration)
}

export function useToast() {
  return {
    toasts,
    success: (message: string, duration?: number) => pushToast('success', message, duration),
    info: (message: string, duration?: number) => pushToast('info', message, duration),
    warning: (message: string, duration?: number) => pushToast('warning', message, duration),
    error: (message: string, duration?: number) => pushToast('error', message, duration),
    removeToast,
  }
}
