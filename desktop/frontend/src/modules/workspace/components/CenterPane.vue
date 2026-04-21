<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import TreePanel from '@/modules/knowledge-tree/components/TreePanel.vue'
import ReaderList from '@/modules/reader/components/ReaderList.vue'
import ReviewQueue from '@/modules/review/components/ReviewQueue.vue'
import GlobalSearch from '@/modules/search/components/GlobalSearch.vue'
import CaptureQuickEntry from '@/modules/inbox-capture/components/CaptureQuickEntry.vue'
import { useReaderStore } from '@/modules/reader/stores/reader.store'

const store = useWorkspaceStore()
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
    return new Date(iso).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
  } catch {
    return iso
  }
}
</script>

<template>
  <AppPane :title="title" class="center-pane" :padded="'none'">
    <div v-if="store.currentContext === 'reading'" class="center-pane__reader">
      <ReaderList />
    </div>

    <div v-else-if="store.currentContext === 'inbox'" class="center-pane__inbox">
      <div class="center-pane__inbox-bar">
        <CaptureQuickEntry />
      </div>
      <div class="center-pane__list">
        <AppEmpty v-if="!inboxArticles.length" message="Inbox is clear" />
        <button
          v-for="a in inboxArticles"
          :key="a.id"
          type="button"
          class="center-pane__row"
          :class="{ 'center-pane__row--active': store.selectedArticleId === a.id }"
          @click="reader.setSelectedArticle(a.id)"
        >
        <span class="center-pane__row-title">{{ a.title }}</span>
        <span class="center-pane__meta">{{ a.status }} · {{ formatShort(a.updatedAt) }}</span>
        </button>
      </div>
    </div>

    <div v-else-if="store.currentContext === 'knowledge'" class="center-pane__tree">
      <TreePanel />
    </div>

    <div v-else-if="store.currentContext === 'review'" class="center-pane__review">
      <ReviewQueue />
    </div>

    <div v-else-if="store.currentContext === 'search'" class="center-pane__search">
      <GlobalSearch />
    </div>
  </AppPane>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.center-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.center-pane__inbox {
  display: flex;
  flex-direction: column;
  min-height: 0;
  flex: 1 1 auto;
}

.center-pane__inbox-bar {
  flex: 0 0 auto;
  display: flex;
  justify-content: flex-end;
  padding: $space-xs $space-sm;
  border-bottom: 1px solid $color-border-subtle;
}

.center-pane__list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: $space-xs;
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
}

.center-pane__tree,
.center-pane__reader,
.center-pane__review,
.center-pane__search {
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.center-pane__hint {
  padding: $space-sm $space-md;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.center-pane__row {
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

.center-pane__row:hover {
  background: $color-hover;
}

.center-pane__row--active {
  background: $color-active;
  border-color: $color-active-border;
}

.center-pane__row-title {
  font-size: $font-size-sm;
  font-weight: 600;
  line-height: $line-tight;
}

.center-pane__meta {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}
</style>
