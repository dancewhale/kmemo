<script setup lang="ts">
import { computed } from 'vue'
import { useWorkspaceStore } from '../stores/workspace.store'
import AppPane from '@/shared/components/AppPane.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import { mockArticles } from '@/mock/articles'
import { mockKnowledgeNodes } from '@/mock/tree'
import { mockReviewItems } from '@/mock/review'

const store = useWorkspaceStore()

const article = computed(() => mockArticles.find((a) => a.id === store.selectedArticleId) ?? null)
const node = computed(() => mockKnowledgeNodes.find((n) => n.id === store.selectedNodeId) ?? null)
const review = computed(() => mockReviewItems.find((r) => r.id === store.selectedReviewId) ?? null)

const title = computed(() => {
  switch (store.currentContext) {
    case 'inbox':
      return 'Inbox item'
    case 'reading':
      return 'Article'
    case 'knowledge':
      return 'Node'
    case 'review':
      return 'Review item'
    case 'search':
      return 'Result detail'
    default:
      return 'Detail'
  }
})
</script>

<template>
  <AppPane :title="title" class="right-pane" :scrollable="true">
    <template v-if="store.currentContext === 'reading' || store.currentContext === 'inbox'">
      <AppEmpty v-if="!article" message="Select an article" />
      <div v-else class="right-pane__block">
        <h2 class="right-pane__h">{{ article.title }}</h2>
        <p class="right-pane__p">{{ article.summary }}</p>
        <dl class="right-pane__dl">
          <div><dt>Id</dt><dd>{{ article.id }}</dd></div>
          <div><dt>Source</dt><dd>{{ article.sourceType }}</dd></div>
          <div><dt>Status</dt><dd>{{ article.status }}</dd></div>
          <div><dt>Updated</dt><dd>{{ article.updatedAt }}</dd></div>
        </dl>
      </div>
    </template>

    <template v-else-if="store.currentContext === 'knowledge'">
      <AppEmpty v-if="!node" message="Select a node" />
      <div v-else class="right-pane__block">
        <h2 class="right-pane__h">{{ node.title }}</h2>
        <dl class="right-pane__dl">
          <div><dt>Id</dt><dd>{{ node.id }}</dd></div>
          <div><dt>Type</dt><dd>{{ node.type }}</dd></div>
          <div><dt>Parent</dt><dd>{{ node.parentId ?? '—' }}</dd></div>
        </dl>
      </div>
    </template>

    <template v-else-if="store.currentContext === 'review'">
      <AppEmpty v-if="!review" message="Select a review item" />
      <div v-else class="right-pane__block">
        <h2 class="right-pane__h">{{ review.title }}</h2>
        <p class="right-pane__label">Prompt</p>
        <p class="right-pane__p">{{ review.prompt }}</p>
        <dl class="right-pane__dl">
          <div><dt>Id</dt><dd>{{ review.id }}</dd></div>
          <div><dt>Due</dt><dd>{{ review.dueAt }}</dd></div>
          <div><dt>Status</dt><dd>{{ review.status }}</dd></div>
        </dl>
      </div>
    </template>

    <template v-else-if="store.currentContext === 'search'">
      <div v-if="article" class="right-pane__block">
        <h2 class="right-pane__h">{{ article.title }}</h2>
        <p class="right-pane__p">{{ article.summary }}</p>
      </div>
      <div v-else class="right-pane__block">
        <p class="right-pane__p">
          Search is mock-only. Pick a result on the left to preview metadata here; host-backed index
          comes later.
        </p>
      </div>
    </template>
  </AppPane>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.right-pane {
  height: 100%;
}

.right-pane__block {
  padding: $space-md;
}

.right-pane__h {
  margin: 0 0 $space-md;
  font-size: $font-size-lg;
  font-weight: 600;
  line-height: $line-tight;
}

.right-pane__label {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.right-pane__p {
  margin: 0 0 $space-md;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.right-pane__dl {
  margin: 0;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.right-pane__dl > div {
  display: grid;
  grid-template-columns: 72px 1fr;
  gap: $space-sm;
  padding: $space-xs 0;
  border-top: 1px solid $color-border-subtle;
}

.right-pane__dl dt {
  margin: 0;
  color: $color-text-secondary;
}

.right-pane__dl dd {
  margin: 0;
  color: $color-text;
  font-family: $font-mono;
  word-break: break-all;
}
</style>
