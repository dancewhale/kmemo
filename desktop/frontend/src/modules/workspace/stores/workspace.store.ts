import { defineStore } from 'pinia'
import {
  DEFAULT_BOTTOM_PANE_HEIGHT,
  DEFAULT_LEFT_PANE_WIDTH,
  DEFAULT_RIGHT_PANE_WIDTH,
  LAYOUT_STORAGE_KEY,
  MAX_BOTTOM_PANE_HEIGHT,
  MAX_LEFT_PANE_WIDTH,
  MAX_RIGHT_PANE_WIDTH,
  MIN_BOTTOM_PANE_HEIGHT,
  MIN_LEFT_PANE_WIDTH,
  MIN_RIGHT_PANE_WIDTH,
} from '@/shared/constants/layout'
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
    leftPaneWidth: DEFAULT_LEFT_PANE_WIDTH,
    rightPaneWidth: DEFAULT_RIGHT_PANE_WIDTH,
    bottomPaneHeight: DEFAULT_BOTTOM_PANE_HEIGHT,
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
    },

    setLeftPaneWidth(w: number) {
      this.leftPaneWidth = clamp(w, MIN_LEFT_PANE_WIDTH, MAX_LEFT_PANE_WIDTH)
    },

    setRightPaneWidth(w: number) {
      this.rightPaneWidth = clamp(w, MIN_RIGHT_PANE_WIDTH, MAX_RIGHT_PANE_WIDTH)
    },

    setBottomPaneHeight(h: number) {
      this.bottomPaneHeight = clamp(h, MIN_BOTTOM_PANE_HEIGHT, MAX_BOTTOM_PANE_HEIGHT)
    },

    setLeftCollapsed(v: boolean) {
      this.isLeftCollapsed = v
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
      if (typeof p.leftPaneWidth === 'number') {
        this.setLeftPaneWidth(p.leftPaneWidth)
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
    },

    persistLayout() {
      const payload: PersistedWorkspaceLayout = {
        leftPaneWidth: this.leftPaneWidth,
        rightPaneWidth: this.rightPaneWidth,
        bottomPaneHeight: this.bottomPaneHeight,
        isLeftCollapsed: this.isLeftCollapsed,
        isRightCollapsed: this.isRightCollapsed,
      }
      writeLocalStorageJSON(LAYOUT_STORAGE_KEY, payload)
    },

    bumpLeftWidth(delta: number) {
      this.setLeftPaneWidth(this.leftPaneWidth + delta)
      this.persistLayout()
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
