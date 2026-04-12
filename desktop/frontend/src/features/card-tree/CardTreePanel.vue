<template>
  <section class="panel-block">
    <div class="card-tree-header">
      <h2 class="panel-title">Card Tree</h2>
      <p class="panel-description">Expandable lazy-loaded tree in the current knowledge scope.</p>
    </div>

    <p v-if="loading" class="panel-description">Loading card tree...</p>
    <div v-else-if="error" class="tree-feedback-block">
      <p class="panel-description">{{ error }}</p>
      <button class="tree-retry-button" type="button" @click="$emit('retry')">Retry</button>
    </div>
    <p v-else-if="items.length === 0" class="panel-description">No root cards available.</p>

    <ul v-else class="tree-list">
      <CardTreeNodeItem
        v-for="item in items"
        :key="item.id"
        :node="item"
        :children="treeChildren[item.id] ?? []"
        :error="treeErrors[item.id] ?? null"
        :expanded="Boolean(treeExpanded[item.id])"
        :loading="Boolean(treeLoading[item.id])"
        :selected-id="selectedId"
        :tree-children="treeChildren"
        :tree-errors="treeErrors"
        :tree-expanded="treeExpanded"
        :tree-loading="treeLoading"
        @select="$emit('select', $event)"
        @toggle="$emit('toggle', $event)"
      />
    </ul>
  </section>
</template>

<script setup lang="ts">
import CardTreeNodeItem from "./CardTreeNodeItem.vue";
import type { CardDTO } from "../../types/dto";

defineProps<{
  items: CardDTO[];
  error: string | null;
  loading: boolean;
  selectedId: string | null;
  treeChildren: Record<string, CardDTO[]>;
  treeErrors: Record<string, string | null>;
  treeExpanded: Record<string, boolean>;
  treeLoading: Record<string, boolean>;
}>();

defineEmits<{
  (event: "retry"): void;
  (event: "select", id: string): void;
  (event: "toggle", id: string): void;
}>();
</script>

<style scoped>
.card-tree-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tree-feedback-block {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.tree-retry-button {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-primary);
  padding: 8px 12px;
  cursor: pointer;
}

.tree-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
