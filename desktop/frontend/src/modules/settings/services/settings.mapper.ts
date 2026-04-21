import { DEFAULT_PREFERENCES } from './settings.defaults'
import type {
  AdvancedPreferences,
  AppearancePreferences,
  AppPreferences,
  ReadingPreferences,
  ThemeMode,
  WorkspacePreferences,
} from '../types'

const VALID_THEME_MODE = new Set<ThemeMode>(['light', 'dark', 'system'])

function toNumber(value: unknown, fallback: number): number {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return fallback
  }
  return value
}

function toBoolean(value: unknown, fallback: boolean): boolean {
  if (typeof value !== 'boolean') {
    return fallback
  }
  return value
}

export function mergeWithDefaults(input: unknown): AppPreferences {
  const raw = (input ?? {}) as Partial<AppPreferences>
  const workspace = (raw.workspace ?? {}) as Partial<WorkspacePreferences>
  const reading = (raw.reading ?? {}) as Partial<ReadingPreferences>
  const appearance = (raw.appearance ?? {}) as Partial<AppearancePreferences>
  const advanced = (raw.advanced ?? {}) as Partial<AdvancedPreferences>

  return {
    workspace: {
      defaultContext:
        workspace.defaultContext ?? DEFAULT_PREFERENCES.workspace.defaultContext,
      leftPaneWidth: toNumber(
        workspace.leftPaneWidth,
        DEFAULT_PREFERENCES.workspace.leftPaneWidth,
      ),
      rightPaneWidth: toNumber(
        workspace.rightPaneWidth,
        DEFAULT_PREFERENCES.workspace.rightPaneWidth,
      ),
      showLeftPane: toBoolean(
        workspace.showLeftPane,
        DEFAULT_PREFERENCES.workspace.showLeftPane,
      ),
      showRightPane: toBoolean(
        workspace.showRightPane,
        DEFAULT_PREFERENCES.workspace.showRightPane,
      ),
      showBottomStatusBar: toBoolean(
        workspace.showBottomStatusBar,
        DEFAULT_PREFERENCES.workspace.showBottomStatusBar,
      ),
      restoreLastContextOnLaunch: toBoolean(
        workspace.restoreLastContextOnLaunch,
        DEFAULT_PREFERENCES.workspace.restoreLastContextOnLaunch,
      ),
      showShortcutHints: toBoolean(
        workspace.showShortcutHints,
        DEFAULT_PREFERENCES.workspace.showShortcutHints,
      ),
      commandPaletteShowGroup: toBoolean(
        workspace.commandPaletteShowGroup,
        DEFAULT_PREFERENCES.workspace.commandPaletteShowGroup,
      ),
    },
    reading: {
      contentWidth: reading.contentWidth ?? DEFAULT_PREFERENCES.reading.contentWidth,
      editorFontSize: reading.editorFontSize ?? DEFAULT_PREFERENCES.reading.editorFontSize,
      showRelatedExtracts: toBoolean(
        reading.showRelatedExtracts,
        DEFAULT_PREFERENCES.reading.showRelatedExtracts,
      ),
    },
    appearance: {
      themeMode: VALID_THEME_MODE.has(appearance.themeMode as ThemeMode)
        ? (appearance.themeMode as ThemeMode)
        : DEFAULT_PREFERENCES.appearance.themeMode,
    },
    advanced: {
      exportDebugInfo: toBoolean(
        advanced.exportDebugInfo,
        DEFAULT_PREFERENCES.advanced.exportDebugInfo,
      ),
    },
  }
}
