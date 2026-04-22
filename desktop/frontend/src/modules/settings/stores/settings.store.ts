import { defineStore } from 'pinia'
import { WORKSPACE_LAST_CONTEXT_STORAGE_KEY } from '@/shared/constants/preferences'
import { readLocalStorageJSON, writeLocalStorageJSON } from '@/shared/utils/storage'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { DEFAULT_PREFERENCES } from '../services/settings.defaults'
import { isWailsAvailable } from '@/api/wails'
import {
  exportPreferences as exportPreferencesText,
  loadPreferences,
  resetPreferences as resetPreferenceStorage,
  savePreferences,
} from '../services/settings.storage'
import * as hostSettingsRepository from '../services/settings.repository'
import type { HostFsrsParameterView } from '../services/settings.mapper'
import type {
  AdvancedPreferences,
  AppearancePreferences,
  AppPreferences,
  ReadingPreferences,
  ThemeMode,
  WorkspacePreferences,
} from '../types'

interface SettingsState {
  preferences: AppPreferences
  initialized: boolean
  /** Default FSRS row from Go host when running under Wails. */
  hostFsrsParameter: HostFsrsParameterView | null
}

function resolveThemeMode(themeMode: ThemeMode): 'light' | 'dark' {
  if (themeMode === 'system') {
    if (typeof window === 'undefined' || !window.matchMedia) {
      return 'dark'
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
  return themeMode
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    preferences: DEFAULT_PREFERENCES,
    initialized: false,
    hostFsrsParameter: null,
  }),
  getters: {
    workspacePreferences(): WorkspacePreferences {
      return this.preferences.workspace
    },
    readingPreferences(): ReadingPreferences {
      return this.preferences.reading
    },
    appearancePreferences(): AppearancePreferences {
      return this.preferences.appearance
    },
    advancedPreferences(): AdvancedPreferences {
      return this.preferences.advanced
    },
    themeMode(): ThemeMode {
      return this.preferences.appearance.themeMode
    },
    editorFontSizeClass(): string {
      return `pref-editor-font-${this.preferences.reading.editorFontSize}`
    },
    contentWidthClass(): string {
      return `pref-content-width-${this.preferences.reading.contentWidth}`
    },
  },
  actions: {
    initialize() {
      if (this.initialized) {
        this.applyPreferences()
        void this.loadHostFsrsIfAvailable()
        return
      }
      this.preferences = loadPreferences()
      this.initialized = true
      this.applyPreferences()
      void this.loadHostFsrsIfAvailable()
    },

    async loadHostFsrsIfAvailable() {
      if (!isWailsAvailable()) {
        return
      }
      const res = await hostSettingsRepository.fetchDefaultFsrsParameter()
      if (res.ok) {
        this.hostFsrsParameter = res.data
      }
    },
    persist() {
      savePreferences(this.preferences)
    },
    updateWorkspacePreference<K extends keyof WorkspacePreferences>(
      key: K,
      value: WorkspacePreferences[K],
    ) {
      this.preferences.workspace = {
        ...this.preferences.workspace,
        [key]: value,
      }
      this.persist()
      this.applyPreferences()
    },
    updateReadingPreference<K extends keyof ReadingPreferences>(
      key: K,
      value: ReadingPreferences[K],
    ) {
      this.preferences.reading = {
        ...this.preferences.reading,
        [key]: value,
      }
      this.persist()
      this.applyPreferences()
    },
    updateAppearancePreference<K extends keyof AppearancePreferences>(
      key: K,
      value: AppearancePreferences[K],
    ) {
      this.preferences.appearance = {
        ...this.preferences.appearance,
        [key]: value,
      }
      this.persist()
      this.applyPreferences()
    },
    updateAdvancedPreference<K extends keyof AdvancedPreferences>(
      key: K,
      value: AdvancedPreferences[K],
    ) {
      this.preferences.advanced = {
        ...this.preferences.advanced,
        [key]: value,
      }
      this.persist()
      this.applyPreferences()
    },
    resetToDefaults() {
      this.preferences = resetPreferenceStorage()
      this.applyPreferences()
    },
    exportPreferences() {
      return exportPreferencesText(this.preferences)
    },
    saveLastContext(context: WorkspacePreferences['defaultContext']) {
      writeLocalStorageJSON(WORKSPACE_LAST_CONTEXT_STORAGE_KEY, context)
    },
    getStartupContext(): WorkspacePreferences['defaultContext'] {
      const fallback = this.preferences.workspace.defaultContext
      if (!this.preferences.workspace.restoreLastContextOnLaunch) {
        return fallback
      }
      const saved = readLocalStorageJSON<WorkspacePreferences['defaultContext']>(
        WORKSPACE_LAST_CONTEXT_STORAGE_KEY,
      )
      return saved ?? fallback
    },
    applyPreferences() {
      const workspace = useWorkspaceStore()
      const pref = this.preferences.workspace
      workspace.setRightPaneWidth(pref.rightPaneWidth)
      workspace.setRightCollapsed(!pref.showRightPane)
      workspace.setShowBottomStatusBar(pref.showBottomStatusBar)
      workspace.persistLayout()

      const resolved = resolveThemeMode(this.preferences.appearance.themeMode)
      document.documentElement.dataset.theme = resolved
    },
  },
})
