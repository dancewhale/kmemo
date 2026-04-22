<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { Setting } from '@element-plus/icons-vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import AppIconButton from '@/shared/components/AppIconButton.vue'
import TreePanel from '@/modules/knowledge-tree/components/TreePanel.vue'
import ReaderList from '@/modules/reader/components/ReaderList.vue'
import ReviewQueue from '@/modules/review/components/ReviewQueue.vue'
import GlobalSearch from '@/modules/search/components/GlobalSearch.vue'
import CaptureQuickEntry from '@/modules/inbox-capture/components/CaptureQuickEntry.vue'
import CreateQuickEntry from '@/modules/object-creation/components/CreateQuickEntry.vue'
import { ROUTE_PATHS } from '@/shared/constants/routes'
import { useReaderStore } from '@/modules/reader/stores/reader.store'

const store = useWorkspaceStore()
const router = useRouter()
const reader = useReaderStore()
const { inboxArticles } = storeToRefs(reader)

onMounted(() => {
  void reader.initialize()
})

const title = computed(() => {
  switch (store.currentContext) {
    case 'inbox':
      return 'Inbox'
    case 'reading':
      return ''
    case 'knowledge':
      return ''
    case 'review':
      return ''
    case 'search':
      return 'Search'
    default:
      return ''
  }
})

function formatShort(iso: string) {
  try {
    return new Date(iso).toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return iso
  }
}

function goSettings() {
  void router.push(ROUTE_PATHS.settings)
}
</script>

<template>
  <div class="tree-pane-root">
    <AppPane class="tree-pane-root__pane" :padded="'none'" :scrollable="true">
    <template #header>
      <div class="tree-pane__header">
        <span v-if="title" class="tree-pane__header-title">{{ title }}</span>
        <div class="tree-pane__toolbar">
          <div class="tree-pane__toolbar-main">
            <CaptureQuickEntry :compact="true" />
            <CreateQuickEntry :compact="true" />
            <span v-if="!store.isLeftCollapsed" class="tree-pane__toolbar-hint">Inbox / New</span>
          </div>
          <div class="tree-pane__toolbar-actions">
            <AppIconButton
              :label="store.isLeftCollapsed ? 'Expand toolbar' : 'Collapse toolbar'"
              @click="store.setLeftCollapsed(!store.isLeftCollapsed)"
            >
              <span class="tree-pane__chev">{{ store.isLeftCollapsed ? '»' : '«' }}</span>
            </AppIconButton>
            <AppIconButton v-if="!store.isLeftCollapsed" label="Settings" @click="goSettings">
              <el-icon><Setting /></el-icon>
            </AppIconButton>
            <button
              v-else
              type="button"
              class="tree-pane__icon-only"
              title="Settings"
              @click="goSettings"
            >
              <el-icon><Setting /></el-icon>
            </button>
          </div>
        </div>
      </div>
    </template>

    <div v-if="store.currentContext === 'reading'" class="tree-pane__reader">
      <ReaderList />
    </div>

    <div v-else-if="store.currentContext === 'inbox'" class="tree-pane__inbox">
      <div class="tree-pane__list">
        <AppEmpty v-if="!inboxArticles.length" message="Inbox is clear" />
        <button
          v-for="a in inboxArticles"
          :key="a.id"
          type="button"
          class="tree-pane__row"
          :class="{ 'tree-pane__row--active': store.selectedArticleId === a.id }"
          @click="reader.setSelectedArticle(a.id)"
        >
          <span class="tree-pane__row-title">{{ a.title }}</span>
          <span class="tree-pane__meta">{{ a.status }} · {{ formatShort(a.updatedAt) }}</span>
        </button>
      </div>
    </div>

    <div v-else-if="store.currentContext === 'knowledge'" class="tree-pane__tree">
      <TreePanel />
    </div>

    <div v-else-if="store.currentContext === 'review'" class="tree-pane__review">
      <ReviewQueue />
    </div>

    <div v-else-if="store.currentContext === 'search'" class="tree-pane__search">
      <GlobalSearch />
    </div>
    </AppPane>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.tree-pane-root {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.tree-pane-root :deep(.app-pane) {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.tree-pane-root :deep(.app-pane__body) {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.tree-pane__header {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: $space-xs;
  width: 100%;
}

.tree-pane__header-title {
  font-size: $font-size-sm;
  font-weight: 600;
  color: $color-text-secondary;
}

.tree-pane__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
  flex-wrap: wrap;
  width: 100%;
}

.tree-pane__toolbar-main {
  display: flex;
  align-items: center;
  gap: $space-xs;
  flex: 1 1 auto;
  min-width: 0;
}

.tree-pane__toolbar-hint {
  font-size: $font-size-xs;
  font-weight: 600;
  color: $color-text-secondary;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tree-pane__toolbar-actions {
  display: flex;
  align-items: center;
  gap: $space-xs;
  flex-shrink: 0;
}

.tree-pane__chev {
  font-size: $font-size-md;
  line-height: 1;
}

.tree-pane__icon-only {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: $radius-sm;
  color: $color-text-secondary;
  border: none;
  background: transparent;
  cursor: pointer;
}

.tree-pane__icon-only:hover {
  background: $color-hover;
  color: $color-text;
}

.tree-pane__inbox {
  display: flex;
  flex-direction: column;
  min-height: 0;
  flex: 1 1 auto;
}

.tree-pane__list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: $space-xs;
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
}

.tree-pane__tree,
.tree-pane__reader,
.tree-pane__review,
.tree-pane__search {
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.tree-pane__row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  text-align: left;
  width: 100%;
  padding: $space-sm $space-md;
  border-radius: $radius-sm;
  border: 1px solid transparent;
  background: transparent;
  color: $color-text;
}

.tree-pane__row:hover {
  background: $color-hover;
}

.tree-pane__row--active {
  background: $color-active;
  border-color: $color-active-border;
}

.tree-pane__row-title {
  font-size: $font-size-sm;
  font-weight: 600;
  line-height: $line-tight;
}

.tree-pane__meta {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}
</style>
