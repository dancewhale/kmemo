<script setup lang="ts">
import { Collection, Document, EditPen, Postcard } from '@element-plus/icons-vue'
import AppStatusTag from '@/shared/components/AppStatusTag.vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import type { KnowledgeNodeType, UITreeNode } from '../types'

const contextMenu = useContextMenuStore()

const props = defineProps<{
  node: UITreeNode
  selected: boolean
  expanded: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
  toggleExpand: [id: string]
}>()

function typeIcon(t: KnowledgeNodeType) {
  switch (t) {
    case 'topic':
      return Collection
    case 'article':
      return Document
    case 'extract':
      return EditPen
    case 'card':
      return Postcard
    default:
      return Document
  }
}

function typeShort(t: KnowledgeNodeType): string {
  switch (t) {
    case 'topic':
      return 'topic'
    case 'article':
      return 'article'
    case 'extract':
      return 'extract'
    case 'card':
      return 'card'
    default:
      return t
  }
}

function onRowClick(e: MouseEvent) {
  const t = e.target as HTMLElement
  if (t.closest('.tree-node-item__exp')) {
    return
  }
  emit('select', props.node.id)
}

function onExpandClick(e: MouseEvent) {
  e.stopPropagation()
  emit('toggleExpand', props.node.id)
}

function onContextMenu(e: MouseEvent) {
  e.preventDefault()
  const base = {
    entityId: props.node.id,
    nodeId: props.node.id,
    x: e.clientX,
    y: e.clientY,
  }
  if (props.node.type === 'card') {
    contextMenu.open({ ...base, entityType: 'tree-card' })
    return
  }
  if (props.node.type === 'extract' && props.node.sourceExtractId) {
    contextMenu.open({
      ...base,
      entityType: 'tree-extract',
      extractId: props.node.sourceExtractId,
    })
    return
  }
  contextMenu.open({ ...base, entityType: 'tree-node' })
}
</script>

<template>
  <div
    class="tree-node-item"
    :class="{ 'tree-node-item--selected': selected }"
    :style="{ paddingLeft: 8 + node.level * 14 + 'px' }"
    role="treeitem"
    :aria-selected="selected"
    :aria-expanded="node.hasChildren ? expanded : undefined"
    @click="onRowClick"
    @contextmenu="onContextMenu"
  >
    <span class="tree-node-item__exp">
      <button
        v-if="node.hasChildren"
        type="button"
        class="tree-node-item__exp-btn"
        :title="expanded ? 'Collapse' : 'Expand'"
        @click="onExpandClick"
      >
        <span class="tree-node-item__chev" :class="{ 'tree-node-item__chev--open': expanded }">›</span>
      </button>
      <span v-else class="tree-node-item__exp-spacer" />
    </span>
    <span class="tree-node-item__ico" :data-type="node.type">
      <el-icon><component :is="typeIcon(node.type)" /></el-icon>
    </span>
    <span class="tree-node-item__title">{{ node.title }}</span>
    <span v-if="node.childCount" class="tree-node-item__count">{{ node.childCount }}</span>
    <AppStatusTag class="tree-node-item__tag" kind="neutral">{{ typeShort(node.type) }}</AppStatusTag>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.tree-node-item {
  display: flex;
  align-items: center;
  gap: 4px;
  min-height: 24px;
  padding: 2px 6px 2px 0;
  margin: 0 2px;
  border-radius: $radius-sm;
  border: 1px solid transparent;
  font-size: $font-size-xs;
  color: $color-text;
  cursor: default;
  user-select: none;
}

.tree-node-item:hover {
  background: $color-hover;
}

.tree-node-item--selected {
  background: $color-active;
  border-color: $color-active-border;
}

.tree-node-item__exp {
  flex: 0 0 auto;
  width: 18px;
  display: flex;
  justify-content: center;
}

.tree-node-item__exp-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: $radius-sm;
  color: $color-text-secondary;
}

.tree-node-item__exp-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: $color-text;
}

.tree-node-item__exp-spacer {
  display: inline-block;
  width: 18px;
  height: 18px;
}

.tree-node-item__chev {
  display: block;
  font-size: 14px;
  line-height: 1;
  transform: rotate(0deg);
  transition: transform 0.12s ease;
}

.tree-node-item__chev--open {
  transform: rotate(90deg);
}

.tree-node-item__ico {
  flex: 0 0 auto;
  font-size: 14px;
  color: $color-text-secondary;
}

.tree-node-item__ico[data-type='topic'] {
  color: #8eb4d9;
}

.tree-node-item__ico[data-type='article'] {
  color: #b8c4d4;
}

.tree-node-item__ico[data-type='extract'] {
  color: #c9b87a;
}

.tree-node-item__ico[data-type='card'] {
  color: #9bc99b;
}

.tree-node-item__title {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: $line-tight;
}

.tree-node-item__count {
  flex: 0 0 auto;
  font-size: 10px;
  padding: 0 4px;
  border-radius: 2px;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
}

.tree-node-item__tag {
  flex: 0 0 auto;
  max-width: 52px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
