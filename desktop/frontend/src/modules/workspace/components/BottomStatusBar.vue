<script setup lang="ts">
import { computed } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'

const store = useWorkspaceStore()
const tree = useTreeStore()
const reader = useReaderStore()
const settings = useSettingsStore()

const selectionLabel = computed(() => {
  if (store.currentContext === 'knowledge') {
    const n = tree.selectedNode
    return n ? `${n.id} · ${n.title}` : '—'
  }
  if (store.currentContext === 'reading') {
    const a = reader.selectedArticle
    return a ? `${a.id} · ${a.title}` : '—'
  }
  return '—'
})

const syncLabel = computed(() => {
  if (store.syncStatus === 'syncing') {
    return 'syncing…'
  }
  if (store.syncStatus === 'saved') {
    return 'saved (mock)'
  }
  return 'idle'
})

function cycleSync() {
  if (store.syncStatus === 'idle') {
    store.setSyncStatus('syncing')
    window.setTimeout(() => store.setSyncStatus('saved'), 600)
    window.setTimeout(() => store.setSyncStatus('idle'), 1600)
  }
}

function toggleRight() {
  store.setRightCollapsed(!store.isRightCollapsed)
  store.persistLayout()
}
</script>

<template>
  <footer class="bottom-status">
    <span class="bottom-status__seg">Mode: <strong>{{ store.currentContext }}</strong></span>
    <span class="bottom-status__sep" />
    <span class="bottom-status__seg" :title="selectionLabel">Sel: {{ selectionLabel }}</span>
    <span class="bottom-status__sep" />
    <span class="bottom-status__seg">Sync: {{ syncLabel }}</span>
    <button type="button" class="bottom-status__btn" @click="cycleSync">Mock sync</button>
    <span class="bottom-status__grow" />
    <span v-if="settings.workspacePreferences.showShortcutHints" class="bottom-status__hints kmono">
      ⌘/Ctrl+K commands · Alt+1/2 navigate · g r/k jump
    </span>
    <button type="button" class="bottom-status__btn" @click="toggleRight">
      {{ store.isRightCollapsed ? 'Show detail' : 'Hide detail' }}
    </button>
  </footer>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.bottom-status {
  display: flex;
  align-items: center;
  gap: $space-sm;
  padding: 0 $space-md;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  background: $color-bg-subtle;
  border-top: 1px solid $color-border-subtle;
  min-width: 0;
}

.bottom-status__seg {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 42vw;
}

.bottom-status__seg strong {
  color: $color-text;
  font-weight: 600;
}

.bottom-status__sep {
  width: 1px;
  height: 12px;
  background: $color-border;
  flex-shrink: 0;
}

.bottom-status__grow {
  flex: 1 1 auto;
}

.bottom-status__hints {
  flex-shrink: 0;
}

.bottom-status__btn {
  flex-shrink: 0;
  font-size: $font-size-xs;
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.bottom-status__btn:hover {
  border-color: $color-active-border;
  color: $color-text;
}
</style>
