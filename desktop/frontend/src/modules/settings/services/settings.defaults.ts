import type { AppPreferences } from '../types'

export const DEFAULT_PREFERENCES: AppPreferences = {
  workspace: {
    defaultContext: 'reading',
    rightPaneWidth: 300,
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
