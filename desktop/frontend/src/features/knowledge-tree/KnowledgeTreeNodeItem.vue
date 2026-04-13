<template>
  <li>
    <button
      class="tree-node"
      :class="{ active: node.id === selectedId }"
      type="button"
      @click="$emit('select', node.id)"
    >
      <span class="tree-node-main">{{ node.name }}</span>
      <span class="tree-node-meta">{{ node.cardCount }}</span>
    </button>
    <ul v-if="(node.children ?? []).length > 0" class="tree-children">
      <KnowledgeTreeNodeItem
        v-for="child in node.children ?? []"
        :key="child.id"
        :node="child"
        :selected-id="selectedId"
        @select="$emit('select', $event)"
      />
    </ul>
  </li>
</template>

<script setup lang="ts">
import type { KnowledgeTreeNode } from "../../types/dto";

defineProps<{
  node: KnowledgeTreeNode;
  selectedId: string | null;
}>();

defineEmits<{
  (event: "select", id: string): void;
}>();
</script>

<style scoped>
.tree-node {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 10px;
  background: transparent;
  color: var(--text-primary);
  cursor: pointer;
}

.tree-node.active {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.tree-node-meta {
  color: var(--text-secondary);
  font-size: 12px;
}

.tree-children {
  list-style: none;
  margin: 6px 0 0 16px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
</style>
