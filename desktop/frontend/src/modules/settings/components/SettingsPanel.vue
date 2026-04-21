<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ROUTE_PATHS } from '@/shared/constants/routes'
import { useToast } from '@/shared/composables/useToast'
import { useSettingsStore } from '../stores/settings.store'
import WorkspacePreferenceSection from './WorkspacePreferenceSection.vue'
import ReadingPreferenceSection from './ReadingPreferenceSection.vue'
import AppearancePreferenceSection from './AppearancePreferenceSection.vue'
import AdvancedPreferenceSection from './AdvancedPreferenceSection.vue'

const settings = useSettingsStore()
const toast = useToast()
const exportDialogOpen = ref(false)
const exportText = computed(() => settings.exportPreferences())

function onResetDefaults() {
  const confirmed = window.confirm('Reset all preferences to defaults?')
  if (!confirmed) {
    return
  }
  settings.resetToDefaults()
  toast.success('Preferences reset to defaults')
}

function onOpenExport() {
  exportDialogOpen.value = true
}

async function onCopyExport() {
  try {
    await navigator.clipboard.writeText(exportText.value)
    toast.success('Preferences JSON copied')
  } catch {
    toast.error('Unable to copy preferences')
  }
}
</script>

<template>
  <div class="settings-panel">
    <header class="settings-panel__head">
      <div>
        <h1 class="settings-panel__title">Preferences</h1>
        <p class="settings-panel__desc">
          Local workbench settings are saved on this device and applied immediately.
        </p>
      </div>
      <RouterLink :to="ROUTE_PATHS.reading" class="settings-panel__back">
        Back to workspace
      </RouterLink>
    </header>

    <div class="settings-panel__sections">
      <WorkspacePreferenceSection />
      <ReadingPreferenceSection />
      <AppearancePreferenceSection />
      <AdvancedPreferenceSection />
    </div>

    <footer class="settings-panel__actions">
      <button type="button" class="settings-panel__btn" @click="onOpenExport">Export JSON</button>
      <button type="button" class="settings-panel__btn settings-panel__btn--danger" @click="onResetDefaults">
        Reset to defaults
      </button>
    </footer>

    <el-dialog v-model="exportDialogOpen" title="Export Preferences" width="720px">
      <pre class="settings-panel__export">{{ exportText }}</pre>
      <template #footer>
        <button type="button" class="settings-panel__btn" @click="onCopyExport">Copy</button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.settings-panel {
  max-width: 920px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: $space-lg;
}

.settings-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $space-lg;
}

.settings-panel__title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.settings-panel__desc {
  margin: $space-xs 0 0;
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.settings-panel__back {
  padding: 4px 8px;
  border-radius: $radius-sm;
  border: 1px solid $color-border;
  color: $color-text-secondary;
  font-size: $font-size-xs;
  text-decoration: none;
}

.settings-panel__back:hover {
  color: $color-text;
  border-color: $color-active-border;
}

.settings-panel__sections {
  display: grid;
  gap: $space-md;
}

.settings-panel__actions {
  display: flex;
  gap: $space-sm;
}

.settings-panel__btn {
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  background: $color-bg-pane-elevated;
  color: $color-text;
  font-size: $font-size-xs;
  padding: 5px 10px;
}

.settings-panel__btn:hover {
  border-color: $color-active-border;
}

.settings-panel__btn--danger {
  color: #f3b9b9;
  border-color: rgba(185, 89, 89, 0.55);
}

.settings-panel__export {
  margin: 0;
  max-height: 360px;
  overflow: auto;
  padding: $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane;
  color: $color-text;
  font-size: $font-size-xs;
  line-height: 1.45;
}
</style>
