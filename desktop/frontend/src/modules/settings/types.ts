export type WorkspaceDefaultContext = 'inbox' | 'reading' | 'knowledge' | 'review' | 'search'

export type ReadingWidthMode = 'narrow' | 'medium' | 'wide'
export type EditorFontSizeMode = 'small' | 'medium' | 'large'
export type ThemeMode = 'light' | 'dark' | 'system'

export interface WorkspacePreferences {
  defaultContext: WorkspaceDefaultContext
  rightPaneWidth: number
  showRightPane: boolean
  showBottomStatusBar: boolean
  restoreLastContextOnLaunch: boolean
  showShortcutHints: boolean
  commandPaletteShowGroup: boolean
}

export interface ReadingPreferences {
  contentWidth: ReadingWidthMode
  editorFontSize: EditorFontSizeMode
  showRelatedExtracts: boolean
}

export interface AppearancePreferences {
  themeMode: ThemeMode
}

export interface AdvancedPreferences {
  exportDebugInfo: boolean
}

export interface AppPreferences {
  workspace: WorkspacePreferences
  reading: ReadingPreferences
  appearance: AppearancePreferences
  advanced: AdvancedPreferences
}
