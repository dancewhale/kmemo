<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '../stores/settings.store'
import PreferenceField from './PreferenceField.vue'
import SettingsSection from './SettingsSection.vue'
import type { EditorFontSizeMode, ReadingWidthMode } from '../types'

const settings = useSettingsStore()
const { readingPreferences } = storeToRefs(settings)

const widthModes: ReadingWidthMode[] = ['narrow', 'medium', 'wide']
const fontModes: EditorFontSizeMode[] = ['small', 'medium', 'large']
</script>

<template>
  <SettingsSection title="Reading & Editor" description="Article editing density and supporting areas.">
    <PreferenceField label="Content width">
      <select
        class="pref-select"
        :value="readingPreferences.contentWidth"
        @change="
          settings.updateReadingPreference(
            'contentWidth',
            ($event.target as HTMLSelectElement).value as ReadingWidthMode,
          )
        "
      >
        <option v-for="mode in widthModes" :key="mode" :value="mode">{{ mode }}</option>
      </select>
    </PreferenceField>

    <PreferenceField label="Editor font size">
      <select
        class="pref-select"
        :value="readingPreferences.editorFontSize"
        @change="
          settings.updateReadingPreference(
            'editorFontSize',
            ($event.target as HTMLSelectElement).value as EditorFontSizeMode,
          )
        "
      >
        <option v-for="mode in fontModes" :key="mode" :value="mode">{{ mode }}</option>
      </select>
    </PreferenceField>

    <PreferenceField label="Show related extracts">
      <input
        type="checkbox"
        :checked="readingPreferences.showRelatedExtracts"
        @change="
          settings.updateReadingPreference(
            'showRelatedExtracts',
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
</style>
