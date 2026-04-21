<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useExtractStore } from '../stores/extract.store'

const contextMenu = useContextMenuStore()

const props = defineProps<{
  articleId: string
}>()

const workspace = useWorkspaceStore()
const tree = useTreeStore()
const extract = useExtractStore()

onMounted(async () => {
  await extract.initialize()
})

const items = computed(() => extract.getExtractsForArticle(props.articleId))

function openExtract(extractId: string) {
  const item = extract.getExtractById(extractId)
  if (!item) {
    return
  }
  extract.setSelectedExtract(item.id)
  if (item.treeNodeId) {
    tree.setSelectedNode(item.treeNodeId)
  }
  workspace.setContext('knowledge')
}

function onExtractContextMenu(e: MouseEvent, extractId: string) {
  e.preventDefault()
  e.stopPropagation()
  contextMenu.open({
    entityType: 'extract',
    entityId: extractId,
    extractId,
    articleId: props.articleId,
    x: e.clientX,
    y: e.clientY,
  })
}
</script>

<template>
  <section class="article-extracts">
    <header class="article-extracts__head">
      <h3 class="article-extracts__title">Related Extracts</h3>
      <span class="article-extracts__count kmono">{{ items.length }}</span>
    </header>
    <p v-if="items.length === 0" class="article-extracts__empty">No extracts from this article yet.</p>
    <div v-else>
      <button
        v-for="item in items"
        :key="item.id"
        type="button"
        class="article-extracts__item"
        @click="openExtract(item.id)"
        @contextmenu="onExtractContextMenu($event, item.id)"
      >
        <span class="article-extracts__item-title">{{ item.title }}</span>
        <span class="article-extracts__item-quote"
          >{{ item.quote.slice(0, 88) }}{{ item.quote.length > 88 ? '…' : '' }}</span
        >
        <span class="article-extracts__item-time kmono">{{ item.updatedAt }}</span>
      </button>
    </div>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.article-extracts {
  border-top: 1px solid $color-border-subtle;
  padding: $space-sm;
  background: $color-bg-pane;
}

.article-extracts__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $space-xs;
}

.article-extracts__title {
  margin: 0;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.article-extracts__count {
  font-size: 10px;
  color: $color-text-secondary;
}

.article-extracts__empty {
  margin: 0;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.article-extracts__item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  width: 100%;
  margin-top: $space-xs;
  padding: $space-xs $space-sm;
  border: 1px solid transparent;
  border-radius: $radius-sm;
  text-align: left;
  color: $color-text;
}

.article-extracts__item:hover {
  background: $color-hover;
  border-color: $color-border-subtle;
}

.article-extracts__item-title {
  font-size: $font-size-sm;
  font-weight: 600;
}

.article-extracts__item-quote {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.article-extracts__item-time {
  font-size: 10px;
  color: $color-text-secondary;
}
</style>
