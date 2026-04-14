<template>
  <div class="tree-node-content">
    <button class="tree-node" :class="{ active: selected }" type="button" @click="$emit('select')">
      <span class="tree-node-main">{{ node.title }}</span>
      <span class="tree-node-meta">{{ node.cardType }}</span>
    </button>
    <div class="tree-actions">
      <button type="button" class="action-btn" @click.stop="$emit('create-child')">+</button>
      <button type="button" class="action-btn" @click.stop="$emit('rename')">R</button>
      <button type="button" class="action-btn danger" @click.stop="$emit('delete')">D</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CardTreeNodeVM } from "./types/cardTree";

defineProps<{
  node: CardTreeNodeVM;
  selected: boolean;
}>();

defineEmits<{
  (event: "select"): void;
  (event: "create-child"): void;
  (event: "rename"): void;
  (event: "delete"): void;
}>();
</script>

<style scoped>
.tree-node-content {
  display: flex;
  gap: 8px;
  align-items: center;
}

.tree-node {
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.45);
  color: var(--text-primary);
  padding: 10px 12px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.tree-node.active {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.tree-node-meta {
  color: var(--text-secondary);
  font-size: 12px;
}

.tree-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-secondary);
  font-size: 11px;
  min-width: 24px;
  min-height: 24px;
  cursor: pointer;
}

.action-btn.danger {
  color: #fca5a5;
}
</style>
