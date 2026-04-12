<template>
  <section class="panel-block">
    <div class="tree-header">
      <h2 class="panel-title">Knowledge Tree</h2>
      <button class="primary-button" type="button" @click="$emit('create-card')">New Card</button>
    </div>

    <p v-if="loading" class="panel-description">Loading knowledge tree...</p>
    <p v-else-if="error" class="panel-description">{{ error }}</p>

    <ul v-else class="tree-list">
      <KnowledgeTreeNodeItem
        v-for="node in nodes"
        :key="node.id"
        :node="node"
        :selected-id="selectedId"
        @select="$emit('select', $event)"
      />
    </ul>
  </section>
</template>

<script setup lang="ts">
import KnowledgeTreeNodeItem from "./KnowledgeTreeNodeItem.vue";
import type { KnowledgeTreeNode } from "../../types/dto";

defineProps<{
  nodes: KnowledgeTreeNode[];
  selectedId: string | null;
  loading: boolean;
  error: string | null;
}>();

defineEmits<{
  (event: "select", id: string): void;
  (event: "create-card"): void;
}>();
</script>

<style scoped>
.tree-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.tree-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.primary-button {
  border: 1px solid var(--accent);
  background: var(--accent-soft);
  color: var(--text-primary);
  border-radius: 10px;
  padding: 8px 12px;
  cursor: pointer;
}
</style>
