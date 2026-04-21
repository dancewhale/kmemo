<script setup lang="ts">
import { computed } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { mockArticles } from '@/mock/articles'
import { mockKnowledgeNodes } from '@/mock/tree'
import { mockReviewItems } from '@/mock/review'

const store = useWorkspaceStore()

const title = computed(() => {
  switch (store.currentContext) {
    case 'inbox':
      return 'Inbox'
    case 'reading':
      return 'Reading queue'
    case 'knowledge':
      return 'Knowledge outline'
    case 'review':
      return 'Review queue'
    case 'search':
      return 'Search'
    default:
      return ''
  }
})

const inboxArticles = computed(() => mockArticles.filter((a) => a.status === 'pending'))

const searchHits = computed(() => mockArticles.slice(0, 4))

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
    <div v-if="store.currentContext === 'reading'" class="center-pane__list">
      <button
        v-for="a in mockArticles"
        :key="a.id"
        type="button"
        class="center-pane__row"
        :class="{ 'center-pane__row--active': store.selectedArticleId === a.id }"
        @click="store.selectArticle(a.id)"
      >
        <span class="center-pane__row-title">{{ a.title }}</span>
        <span class="center-pane__meta">{{ a.sourceType }} · {{ formatShort(a.updatedAt) }}</span>
      </button>
    </div>

    <div v-else-if="store.currentContext === 'inbox'" class="center-pane__list">
      <AppEmpty v-if="!inboxArticles.length" message="Inbox is clear" />
      <button
        v-for="a in inboxArticles"
        :key="a.id"
        type="button"
        class="center-pane__row"
        :class="{ 'center-pane__row--active': store.selectedArticleId === a.id }"
        @click="store.selectArticle(a.id)"
      >
        <span class="center-pane__row-title">{{ a.title }}</span>
        <span class="center-pane__meta">{{ a.status }} · {{ formatShort(a.updatedAt) }}</span>
      </button>
    </div>

    <div v-else-if="store.currentContext === 'knowledge'" class="center-pane__list">
      <button
        v-for="n in mockKnowledgeNodes"
        :key="n.id"
        type="button"
        class="center-pane__row"
        :class="{ 'center-pane__row--active': store.selectedNodeId === n.id }"
        @click="store.selectNode(n.id)"
      >
        <span class="center-pane__row-title">{{ n.title }}</span>
        <span class="center-pane__meta">{{ n.type }}{{ n.parentId ? ` · parent ${n.parentId}` : '' }}</span>
      </button>
    </div>

    <div v-else-if="store.currentContext === 'review'" class="center-pane__list">
      <button
        v-for="r in mockReviewItems"
        :key="r.id"
        type="button"
        class="center-pane__row"
        :class="{ 'center-pane__row--active': store.selectedReviewId === r.id }"
        @click="store.selectReview(r.id)"
      >
        <span class="center-pane__row-title">{{ r.title }}</span>
        <span class="center-pane__meta">{{ r.status }} · due {{ formatShort(r.dueAt) }}</span>
      </button>
    </div>

    <div v-else-if="store.currentContext === 'search'" class="center-pane__list">
      <div class="center-pane__hint">Mock results for “incremental”</div>
      <button
        v-for="a in searchHits"
        :key="a.id"
        type="button"
        class="center-pane__row"
        :class="{ 'center-pane__row--active': store.selectedArticleId === a.id }"
        @click="store.selectArticle(a.id)"
      >
        <span class="center-pane__row-title">{{ a.title }}</span>
        <span class="center-pane__meta">{{ a.summary.slice(0, 72) }}…</span>
      </button>
    </div>
  </AppPane>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.center-pane {
  height: 100%;
}

.center-pane__list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: $space-xs;
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
