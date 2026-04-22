import { defineStore } from 'pinia'
import {
  DEFAULT_BOTTOM_PANE_HEIGHT,
  DEFAULT_RIGHT_PANE_WIDTH,
  LAYOUT_STORAGE_KEY,
  MAX_BOTTOM_PANE_HEIGHT,
  MIN_BOTTOM_PANE_HEIGHT,
  MIN_RIGHT_PANE_WIDTH,
} from '@/shared/constants/layout'
import { WORKSPACE_LAST_CONTEXT_STORAGE_KEY } from '@/shared/constants/preferences'
import { readLocalStorageJSON, writeLocalStorageJSON } from '@/shared/utils/storage'
import type { PersistedWorkspaceLayout, WorkspaceContext } from '../types'

function clamp(n: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, n))
}

function loadPersistedLayout(): PersistedWorkspaceLayout | null {
  return readLocalStorageJSON<PersistedWorkspaceLayout>(LAYOUT_STORAGE_KEY)
}

export const useWorkspaceStore = defineStore('workspace', {
  state: () => ({
    currentContext: 'reading' as WorkspaceContext,
    rightPaneWidth: DEFAULT_RIGHT_PANE_WIDTH,
    bottomPaneHeight: DEFAULT_BOTTOM_PANE_HEIGHT,
    showBottomStatusBar: true,
    isLeftCollapsed: false,
    isRightCollapsed: false,
    selectedArticleId: null as string | null,
    selectedNodeId: null as string | null,
    selectedReviewId: null as string | null,
    syncStatus: 'idle' as 'idle' | 'syncing' | 'saved',
  }),
  actions: {
    setContext(ctx: WorkspaceContext) {
      this.currentContext = ctx
      writeLocalStorageJSON(WORKSPACE_LAST_CONTEXT_STORAGE_KEY, ctx)
    },

    setRightPaneWidth(w: number) {
      const v = Number.isFinite(w) ? w : MIN_RIGHT_PANE_WIDTH
      this.rightPaneWidth = Math.max(MIN_RIGHT_PANE_WIDTH, v)
    },

    setBottomPaneHeight(h: number) {
      this.bottomPaneHeight = clamp(h, MIN_BOTTOM_PANE_HEIGHT, MAX_BOTTOM_PANE_HEIGHT)
    },

    setShowBottomStatusBar(v: boolean) {
      this.showBottomStatusBar = v
    },

    setLeftCollapsed(v: boolean) {
      this.isLeftCollapsed = v
      this.persistLayout()
    },

    setRightCollapsed(v: boolean) {
      this.isRightCollapsed = v
    },

    selectArticle(id: string | null) {
      this.selectedArticleId = id
    },

    selectNode(id: string | null) {
      this.selectedNodeId = id
    },

    selectReview(id: string | null) {
      this.selectedReviewId = id
    },

    setSyncStatus(status: 'idle' | 'syncing' | 'saved') {
      this.syncStatus = status
    },

    restoreLayoutFromStorage() {
      const p = loadPersistedLayout()
      if (!p) {
        return
      }
      if (typeof p.rightPaneWidth === 'number') {
        this.setRightPaneWidth(p.rightPaneWidth)
      }
      if (typeof p.bottomPaneHeight === 'number') {
        this.setBottomPaneHeight(p.bottomPaneHeight)
      }
      if (typeof p.isLeftCollapsed === 'boolean') {
        this.isLeftCollapsed = p.isLeftCollapsed
      }
      if (typeof p.isRightCollapsed === 'boolean') {
        this.isRightCollapsed = p.isRightCollapsed
      }
      if (typeof p.showBottomStatusBar === 'boolean') {
        this.showBottomStatusBar = p.showBottomStatusBar
      }
    },

    persistLayout() {
      const payload: PersistedWorkspaceLayout = {
        rightPaneWidth: this.rightPaneWidth,
        bottomPaneHeight: this.bottomPaneHeight,
        showBottomStatusBar: this.showBottomStatusBar,
        isLeftCollapsed: this.isLeftCollapsed,
        isRightCollapsed: this.isRightCollapsed,
      }
      writeLocalStorageJSON(LAYOUT_STORAGE_KEY, payload)
    },

    bumpRightWidth(delta: number) {
      this.setRightPaneWidth(this.rightPaneWidth - delta)
      this.persistLayout()
    },

    bumpBottomHeight(delta: number) {
      this.setBottomPaneHeight(this.bottomPaneHeight + delta)
      this.persistLayout()
    },
  },
})
