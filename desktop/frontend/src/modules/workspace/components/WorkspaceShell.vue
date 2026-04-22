<script setup lang="ts">
import { watch } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import type { WorkspaceContext } from '../types'
import TreePane from './TreePane.vue'
import ReadPane from './ReadPane.vue'
import BottomStatusBar from './BottomStatusBar.vue'
import AppSplitter from '@/shared/components/AppSplitter.vue'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useReaderStore } from '@/modules/reader/stores/reader.store'
import { useReviewStore } from '@/modules/review/stores/review.store'

const props = defineProps<{
  context: WorkspaceContext
}>()

const store = useWorkspaceStore()

async function applyContext(ctx: WorkspaceContext) {
  store.setContext(ctx)
  if (ctx === 'reading') {
    const reader = useReaderStore()
    await reader.ensureReadyWithSelection(store.selectedArticleId)
  } else if (ctx === 'inbox') {
    const reader = useReaderStore()
    await reader.initialize()
    const inbox = reader.inboxArticles
    const preferred = store.selectedArticleId
    const still = preferred && inbox.some((a) => a.id === preferred)
    reader.setSelectedArticle(still ? preferred : inbox[0]?.id ?? null)
  } else if (ctx === 'knowledge') {
    const tree = useTreeStore()
    await tree.ensureReadyWithSelection(store.selectedNodeId)
  } else if (ctx === 'review') {
    const review = useReviewStore()
    await review.ensureReadyWithSelection(store.selectedReviewId)
  } else if (ctx === 'search') {
    store.selectArticle(null)
    store.selectNode(null)
    store.selectReview(null)
  }
}

watch(
  () => props.context,
  (c) => {
    void applyContext(c)
  },
  { immediate: true },
)
</script>

<template>
  <div class="workspace-shell">
    <div class="workspace-shell__row">
      <TreePane class="workspace-shell__center" />
      <template v-if="!store.isRightCollapsed">
        <AppSplitter orientation="vertical" @drag="store.bumpRightWidth($event)" />
        <ReadPane class="workspace-shell__right" :style="{ width: store.rightPaneWidth + 'px' }" />
      </template>
    </div>
    <template v-if="store.showBottomStatusBar">
      <AppSplitter orientation="horizontal" @drag="store.bumpBottomHeight($event)" />
      <BottomStatusBar
        class="workspace-shell__status"
        :style="{
          height: store.bottomPaneHeight + 'px',
          minHeight: store.bottomPaneHeight + 'px',
        }"
      />
    </template>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.workspace-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: $color-bg-app;
}

.workspace-shell__row {
  display: flex;
  flex: 1 1 auto;
  min-height: 0;
  min-width: 0;
}

.workspace-shell__center {
  flex: 1 1 auto;
  min-width: 0;
  min-height: 0;
}

.workspace-shell__right {
  flex: 0 0 auto;
  min-height: 0;
}

.workspace-shell__status {
  flex: 0 0 auto;
}
</style>
