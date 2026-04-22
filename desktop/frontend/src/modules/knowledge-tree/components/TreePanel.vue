<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import AppActionBar from '@/shared/components/AppActionBar.vue'
import AppSectionTitle from '@/shared/components/AppSectionTitle.vue'
import { useTreeStore } from '../stores/tree.store'
import TreeNodeItem from './TreeNodeItem.vue'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import AppLoading from '@/shared/components/AppLoading.vue'
import CreateQuickEntry from '@/modules/object-creation/components/CreateQuickEntry.vue'

const tree = useTreeStore()
const { visibleNodes, selectedNodeId, loading, searchKeyword } = storeToRefs(tree)

onMounted(async () => {
  await tree.initialize()
})

function isExpanded(id: string) {
  return tree.effectiveExpandedNodeIds.includes(id)
}
</script>

<template>
  <div class="tree-panel">
    <AppActionBar class="tree-panel__head">
      <div class="tree-panel__title-wrap">
        <AppSectionTitle title="Knowledge Tree" />
      </div>
      <div class="tree-panel__toolbar">
        <CreateQuickEntry :compact="true" class="tree-panel__new" />
        <el-input
          :model-value="searchKeyword"
          size="small"
          clearable
          placeholder="Filter titles…"
          class="tree-panel__search"
          @update:model-value="tree.setSearchKeyword(String($event ?? ''))"
        />
      </div>
    </AppActionBar>
    <div class="tree-panel__body">
      <div v-if="loading" class="tree-panel__loading">
        <AppLoading label="Loading tree…" />
      </div>
      <AppEmpty v-else-if="!visibleNodes.length" :message="tree.treeListEmptyMessage" />
      <div v-else class="tree-panel__list" role="tree">
        <TreeNodeItem
          v-for="n in visibleNodes"
          :key="n.id"
          :node="n"
          :selected="selectedNodeId === n.id"
          :expanded="isExpanded(n.id)"
          @select="tree.setSelectedNode($event)"
          @toggle-expand="tree.toggleNode($event)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.tree-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: $color-bg-pane;
}

.tree-panel__head {
  justify-content: space-between;
}

.tree-panel__title-wrap {
  min-width: 0;
}

.tree-panel__toolbar {
  display: flex;
  align-items: center;
  gap: $space-xs;
}

.tree-panel__new {
  flex-shrink: 0;
}

.tree-panel__search {
  flex: 1 1 auto;
  min-width: 80px;
}

.tree-panel__body {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: $space-xs 0 $space-sm;
}

.tree-panel__loading {
  display: flex;
  justify-content: center;
  padding: $space-xl 0;
}

.tree-panel__list {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.tree-panel__body :deep(.app-empty) {
  padding: $space-xl $space-md;
}
</style>
