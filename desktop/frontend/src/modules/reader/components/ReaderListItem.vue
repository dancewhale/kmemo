<script setup lang="ts">
import AppStatusTag from '@/shared/components/AppStatusTag.vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import type { ReaderArticle } from '../types'

const contextMenu = useContextMenuStore()

const props = defineProps<{
  article: ReaderArticle
  selected: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

function formatUpdated(iso: string): string {
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

function summaryPreview(text: string, max = 120): string {
  const t = text.trim()
  if (t.length <= max) {
    return t
  }
  return `${t.slice(0, max).trimEnd()}…`
}

function statusKind(s: ReaderArticle['status']): 'warning' | 'info' | 'success' {
  if (s === 'inbox') {
    return 'warning'
  }
  if (s === 'reading') {
    return 'info'
  }
  return 'success'
}

function onContextMenu(e: MouseEvent) {
  e.preventDefault()
  contextMenu.open({
    entityType: 'article',
    entityId: props.article.id,
    articleId: props.article.id,
    x: e.clientX,
    y: e.clientY,
  })
}
</script>

<template>
  <button
    type="button"
    class="reader-list-item"
    :class="{ 'reader-list-item--selected': selected }"
    @click="emit('select', props.article.id)"
    @contextmenu="onContextMenu"
  >
    <div class="reader-list-item__row1">
      <span class="reader-list-item__title">{{ article.title }}</span>
      <AppStatusTag :kind="statusKind(article.status)">
        {{ article.status }}
      </AppStatusTag>
    </div>
    <p class="reader-list-item__summary">{{ summaryPreview(article.summary) }}</p>
    <div class="reader-list-item__row3 kmono">
      <span>{{ article.sourceType }}</span>
      <span class="reader-list-item__dot">·</span>
      <span>{{ formatUpdated(article.updatedAt) }}</span>
      <template v-if="article.tags?.length">
        <span class="reader-list-item__dot">·</span>
        <span class="reader-list-item__tags">{{ article.tags.join(', ') }}</span>
      </template>
    </div>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.reader-list-item {
  display: block;
  width: calc(100% - 4px);
  margin: 0 2px 1px;
  padding: $space-sm $space-md;
  text-align: left;
  border-radius: $radius-sm;
  border: 1px solid transparent;
  background: transparent;
  color: $color-text;
  cursor: default;
}

.reader-list-item:hover {
  background: $color-hover;
}

.reader-list-item--selected {
  background: $color-active;
  border-color: $color-active-border;
}

.reader-list-item__row1 {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $space-sm;
  margin-bottom: 3px;
}

.reader-list-item__title {
  font-size: $font-size-sm;
  font-weight: 600;
  line-height: $line-tight;
  flex: 1 1 auto;
  min-width: 0;
}

.reader-list-item__summary {
  margin: 0 0 4px;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.reader-list-item__row3 {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 2px;
  font-size: 10px;
  color: $color-text-secondary;
}

.reader-list-item__dot {
  opacity: 0.55;
}

.reader-list-item__tags {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}
</style>
