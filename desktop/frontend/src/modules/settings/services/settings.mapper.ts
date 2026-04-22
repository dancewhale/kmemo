import type { FSRSParameterDTO } from '@/api/types/dto'
import { DEFAULT_PREFERENCES } from './settings.defaults'
import type { AppPreferences, WorkspacePreferences } from '../types'

function normalizeWorkspacePreferences(ws: Partial<WorkspacePreferences> | undefined): WorkspacePreferences {
  const merged: WorkspacePreferences = { ...DEFAULT_PREFERENCES.workspace, ...ws }
  merged.defaultContext = merged.defaultContext === 'knowledge' ? 'knowledge' : 'reading'
  return merged
}

export function mergeWithDefaults(partial: Partial<AppPreferences> | null | undefined): AppPreferences {
  const p = partial ?? {}
  return {
    workspace: normalizeWorkspacePreferences(p.workspace),
    reading: { ...DEFAULT_PREFERENCES.reading, ...p.reading },
    appearance: { ...DEFAULT_PREFERENCES.appearance, ...p.appearance },
    advanced: { ...DEFAULT_PREFERENCES.advanced, ...p.advanced },
  }
}

/** Minimal host FSRS snapshot for settings / debug UI. */
export interface HostFsrsParameterView {
  id: string
  name: string
  parametersJson: string
  desiredRetention: number | null
  maximumInterval: number | null
}

export function mapFsrsParameterDtoToView(dto: FSRSParameterDTO): HostFsrsParameterView {
  return {
    id: dto.id,
    name: dto.name,
    parametersJson: dto.parametersJson,
    desiredRetention: dto.desiredRetention,
    maximumInterval: dto.maximumInterval,
  }
}
