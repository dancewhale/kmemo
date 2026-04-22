<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '../stores/settings.store'
import PreferenceField from './PreferenceField.vue'
import SettingsSection from './SettingsSection.vue'
import type { WorkspaceDefaultContext } from '../types'

const settings = useSettingsStore()
const { workspacePreferences } = storeToRefs(settings)

const contexts: WorkspaceDefaultContext[] = [
  'inbox',
  'reading',
  'knowledge',
  'review',
  'search',
]
</script>

<template>
  <SettingsSection title="Workspace" description="Layout defaults and command center behavior.">
    <PreferenceField label="Default workspace" description="Used when no last context is restored.">
      <select
        class="pref-select"
        :value="workspacePreferences.defaultContext"
        @change="
          settings.updateWorkspacePreference(
            'defaultContext',
            ($event.target as HTMLSelectElement).value as WorkspaceDefaultContext,
          )
        "
      >
        <option v-for="context in contexts" :key="context" :value="context">
          {{ context }}
        </option>
      </select>
    </PreferenceField>

    <PreferenceField label="Right pane width" description="Default detail pane width.">
      <input
        class="pref-range"
        type="range"
        min="200"
        max="520"
        :value="workspacePreferences.rightPaneWidth"
        @input="
          settings.updateWorkspacePreference(
            'rightPaneWidth',
            Number(($event.target as HTMLInputElement).value),
          )
        "
      />
      <span class="pref-value">{{ workspacePreferences.rightPaneWidth }}px</span>
    </PreferenceField>

    <PreferenceField label="Show right pane">
      <input
        type="checkbox"
        :checked="workspacePreferences.showRightPane"
        @change="
          settings.updateWorkspacePreference(
            'showRightPane',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>

    <PreferenceField label="Show bottom status bar">
      <input
        type="checkbox"
        :checked="workspacePreferences.showBottomStatusBar"
        @change="
          settings.updateWorkspacePreference(
            'showBottomStatusBar',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>

    <PreferenceField label="Restore last context on launch">
      <input
        type="checkbox"
        :checked="workspacePreferences.restoreLastContextOnLaunch"
        @change="
          settings.updateWorkspacePreference(
            'restoreLastContextOnLaunch',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>

    <PreferenceField label="Show shortcut hints">
      <input
        type="checkbox"
        :checked="workspacePreferences.showShortcutHints"
        @change="
          settings.updateWorkspacePreference(
            'showShortcutHints',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>

    <PreferenceField label="Show command groups">
      <input
        type="checkbox"
        :checked="workspacePreferences.commandPaletteShowGroup"
        @change="
          settings.updateWorkspacePreference(
            'commandPaletteShowGroup',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>
  </SettingsSection>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.pref-select {
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  background: $color-bg-pane-elevated;
  color: $color-text;
  padding: 3px 8px;
  font-size: $font-size-xs;
}

.pref-range {
  width: 180px;
}

.pref-value {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  min-width: 48px;
}
</style>
