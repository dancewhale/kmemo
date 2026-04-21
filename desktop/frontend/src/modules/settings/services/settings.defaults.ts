import type { AppPreferences } from '../types'

export const DEFAULT_PREFERENCES: AppPreferences = {
  workspace: {
    defaultContext: 'reading',
    leftPaneWidth: 196,
    rightPaneWidth: 300,
    showLeftPane: true,
    showRightPane: true,
    showBottomStatusBar: true,
    restoreLastContextOnLaunch: true,
    showShortcutHints: true,
    commandPaletteShowGroup: true,
  },
  reading: {
    contentWidth: 'medium',
    editorFontSize: 'medium',
    showRelatedExtracts: true,
  },
  appearance: {
    themeMode: 'dark',
  },
  advanced: {
    exportDebugInfo: false,
  },
}
