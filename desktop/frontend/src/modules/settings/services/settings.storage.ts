import { PREFERENCES_STORAGE_KEY } from '@/shared/constants/preferences'
import { readLocalStorageJSON, writeLocalStorageJSON } from '@/shared/utils/storage'
import { DEFAULT_PREFERENCES } from './settings.defaults'
import { mergeWithDefaults } from './settings.mapper'
import type { AppPreferences } from '../types'

export function loadPreferences(): AppPreferences {
  const saved = readLocalStorageJSON<Partial<AppPreferences>>(PREFERENCES_STORAGE_KEY)
  return mergeWithDefaults(saved)
}

export function savePreferences(preferences: AppPreferences): void {
  writeLocalStorageJSON(PREFERENCES_STORAGE_KEY, preferences)
}

export function resetPreferences(): AppPreferences {
  const payload = mergeWithDefaults(DEFAULT_PREFERENCES)
  savePreferences(payload)
  return payload
}

export function exportPreferences(preferences: AppPreferences): string {
  return JSON.stringify(preferences, null, 2)
}
