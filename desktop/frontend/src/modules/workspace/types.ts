export type WorkspaceContext = 'inbox' | 'reading' | 'knowledge' | 'review' | 'search'

export interface WorkspaceState {
  currentContext: WorkspaceContext
  showLeftPane: boolean
  leftPaneWidth: number
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
  showLeftPane?: boolean
  leftPaneWidth: number
  rightPaneWidth: number
  bottomPaneHeight: number
  showBottomStatusBar?: boolean
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
}
