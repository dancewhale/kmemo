export type WorkspaceContext = 'inbox' | 'reading' | 'knowledge' | 'review' | 'search'

export interface WorkspaceState {
  currentContext: WorkspaceContext
  leftPaneWidth: number
  rightPaneWidth: number
  bottomPaneHeight: number
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
  selectedArticleId: string | null
  selectedNodeId: string | null
  selectedReviewId: string | null
  syncStatus: 'idle' | 'syncing' | 'saved'
}

export interface PersistedWorkspaceLayout {
  leftPaneWidth: number
  rightPaneWidth: number
  bottomPaneHeight: number
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
}
