<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '../stores/settings.store'
import PreferenceField from './PreferenceField.vue'
import SettingsSection from './SettingsSection.vue'

const settings = useSettingsStore()
const { advancedPreferences } = storeToRefs(settings)
</script>

<template>
  <SettingsSection title="Advanced" description="Small debugging and export-related flags.">
    <PreferenceField label="Export debug info">
      <input
        type="checkbox"
        :checked="advancedPreferences.exportDebugInfo"
        @change="
          settings.updateAdvancedPreference(
            'exportDebugInfo',
            ($event.target as HTMLInputElement).checked,
          )
        "
      />
    </PreferenceField>
  </SettingsSection>
</template>
