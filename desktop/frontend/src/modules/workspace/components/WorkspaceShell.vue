<script setup lang="ts">
import { watch } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import type { WorkspaceContext } from '../types'
import LeftSidebar from './LeftSidebar.vue'
import CenterPane from './CenterPane.vue'
import RightPane from './RightPane.vue'
import BottomStatusBar from './BottomStatusBar.vue'
import AppSplitter from '@/shared/components/AppSplitter.vue'
import { mockArticles } from '@/mock/articles'
import { mockKnowledgeNodes } from '@/mock/tree'
import { mockReviewItems } from '@/mock/review'

const props = defineProps<{
  context: WorkspaceContext
}>()

const store = useWorkspaceStore()

const COLLAPSED_RAIL = 48

function applyContext(ctx: WorkspaceContext) {
  store.setContext(ctx)
  if (ctx === 'reading') {
    const cur = mockArticles.find((a) => a.id === store.selectedArticleId)
    if (!cur) {
      store.selectArticle(mockArticles[0]?.id ?? null)
    }
  } else if (ctx === 'inbox') {
    const pending = mockArticles.filter((a) => a.status === 'pending')
    const cur = pending.find((a) => a.id === store.selectedArticleId)
    if (!cur) {
      store.selectArticle(pending[0]?.id ?? null)
    }
  } else if (ctx === 'knowledge') {
    const nodes = mockKnowledgeNodes.filter((n) => n.type !== 'folder')
    const cur = nodes.find((n) => n.id === store.selectedNodeId)
    if (!cur) {
      store.selectNode(nodes[0]?.id ?? null)
    }
  } else if (ctx === 'review') {
    const cur = mockReviewItems.find((r) => r.id === store.selectedReviewId)
    if (!cur) {
      store.selectReview(mockReviewItems[0]?.id ?? null)
    }
  } else if (ctx === 'search') {
    store.selectArticle(null)
    store.selectNode(null)
    store.selectReview(null)
  }
}

watch(
  () => props.context,
  (c) => applyContext(c),
  { immediate: true },
)
</script>

<template>
  <div class="workspace-shell">
    <div class="workspace-shell__row">
      <LeftSidebar
        class="workspace-shell__left"
        :style="{
          width: (store.isLeftCollapsed ? COLLAPSED_RAIL : store.leftPaneWidth) + 'px',
        }"
      />
      <AppSplitter orientation="vertical" @drag="store.bumpLeftWidth($event)" />
      <CenterPane class="workspace-shell__center" />
      <template v-if="!store.isRightCollapsed">
        <AppSplitter orientation="vertical" @drag="store.bumpRightWidth($event)" />
        <RightPane class="workspace-shell__right" :style="{ width: store.rightPaneWidth + 'px' }" />
      </template>
    </div>
    <AppSplitter orientation="horizontal" @drag="store.bumpBottomHeight($event)" />
    <BottomStatusBar
      class="workspace-shell__status"
      :style="{
        height: store.bottomPaneHeight + 'px',
        minHeight: store.bottomPaneHeight + 'px',
      }"
    />
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

.workspace-shell__left {
  flex: 0 0 auto;
  min-height: 0;
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
