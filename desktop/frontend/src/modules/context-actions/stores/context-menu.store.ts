import { defineStore } from 'pinia'
import { buildContextMenuItems } from '../services/context-menu.builder'
import { executeContextAction } from '../services/context-menu.executor'
import type { ContextActionItem, ContextMenuContext } from '../types'

interface ContextMenuStoreState {
  isOpen: boolean
  x: number
  y: number
  context: ContextMenuContext | null
}

export const useContextMenuStore = defineStore('context-menu', {
  state: (): ContextMenuStoreState => ({
    isOpen: false,
    x: 0,
    y: 0,
    context: null,
  }),

  getters: {
    hasContext(): boolean {
      return this.context != null
    },
    menuItems(): ContextActionItem[] {
      if (!this.context) {
        return []
      }
      return buildContextMenuItems(this.context)
    },
  },

  actions: {
    open(context: ContextMenuContext) {
      this.context = context
      this.x = context.x
      this.y = context.y
      this.isOpen = true
    },

    close() {
      this.isOpen = false
      this.context = null
    },

    setPosition(x: number, y: number) {
      this.x = x
      this.y = y
    },

    async execute(actionId: string) {
      const ctx = this.context
      if (!ctx) {
        return
      }
      await executeContextAction(actionId, ctx)
      this.close()
    },
  },
})
