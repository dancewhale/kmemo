<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '../stores/settings.store'
import PreferenceField from './PreferenceField.vue'
import SettingsSection from './SettingsSection.vue'
import ThemeModeSelect from './ThemeModeSelect.vue'
import type { ThemeMode } from '../types'

const settings = useSettingsStore()
const { appearancePreferences } = storeToRefs(settings)

function onThemeChange(mode: ThemeMode) {
  settings.updateAppearancePreference('themeMode', mode)
}
</script>

<template>
  <SettingsSection title="Appearance" description="Theme mode for the entire workbench.">
    <PreferenceField label="Theme mode">
      <ThemeModeSelect :model-value="appearancePreferences.themeMode" @update:model-value="onThemeChange" />
    </PreferenceField>
  </SettingsSection>
</template>
