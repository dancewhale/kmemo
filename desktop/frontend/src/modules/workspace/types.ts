export type WorkspaceContext = 'reading' | 'knowledge'

export interface WorkspaceState {
  currentContext: WorkspaceContext
  rightPaneWidth: number
  bottomPaneHeight: number
  showBottomStatusBar: boolean
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
  selectedArticleId: string | null
  selectedNodeId: string | null
  selectedReviewId: string | null
  syncStatus: 'idle' | 'syncing' | 'saved'
}

export interface PersistedWorkspaceLayout {
  rightPaneWidth: number
  bottomPaneHeight: number
  showBottomStatusBar?: boolean
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
}
