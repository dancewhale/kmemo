<template>
  <li>
    <div class="tree-row">
      <button v-if="canToggle" class="toggle-button" type="button" @click="$emit('toggle', node.id)">
        <span v-if="loading">…</span>
        <span v-else>{{ expanded ? '−' : '+' }}</span>
      </button>
      <span v-else class="toggle-placeholder"></span>
      <button class="tree-node" :class="{ active: node.id === selectedId }" type="button" @click="$emit('select', node.id)">
        <span class="tree-node-main">{{ node.title }}</span>
        <span class="tree-node-meta">{{ node.cardType }}</span>
      </button>
    </div>
    <div v-if="error" class="tree-error-block">
      <p class="tree-error">{{ error }}</p>
      <button class="tree-retry-button" type="button" @click="$emit('toggle', node.id)">Retry</button>
    </div>
    <ul v-if="expanded && children.length > 0" class="tree-children">
      <CardTreeNodeItem
        v-for="child in children"
        :key="child.id"
        :node="child"
        :children="treeChildren[child.id] ?? []"
        :error="treeErrors[child.id] ?? null"
        :expanded="Boolean(treeExpanded[child.id])"
        :loading="Boolean(treeLoading[child.id])"
        :selected-id="selectedId"
        :tree-children="treeChildren"
        :tree-errors="treeErrors"
        :tree-expanded="treeExpanded"
        :tree-loading="treeLoading"
        @select="$emit('select', $event)"
        @toggle="$emit('toggle', $event)"
      />
    </ul>
  </li>
</template>

<script setup lang="ts">
import { computed } from "vue";

import type { CardDTO } from "../../types/dto";

const props = defineProps<{
  node: CardDTO;
  children: CardDTO[];
  error: string | null;
  expanded: boolean;
  loading: boolean;
  selectedId: string | null;
  treeChildren: Record<string, CardDTO[]>;
  treeErrors: Record<string, string | null>;
  treeExpanded: Record<string, boolean>;
  treeLoading: Record<string, boolean>;
}>();

defineEmits<{
  (event: "select", id: string): void;
  (event: "toggle", id: string): void;
}>();

const canToggle = computed(() => {
  if (props.loading || props.error) return true;
  const cachedChildren = props.treeChildren[props.node.id];
  if (!cachedChildren) return true;
  return cachedChildren.length > 0;
});
</script>

<style scoped>
.tree-row {
  display: grid;
  grid-template-columns: 32px minmax(0, 1fr);
  gap: 8px;
}

.toggle-button,
.toggle-placeholder {
  width: 32px;
  min-height: 32px;
}

.toggle-button {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.45);
  color: var(--text-secondary);
  cursor: pointer;
}

.toggle-placeholder {
  display: inline-block;
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

.tree-error-block {
  margin: 6px 0 0 40px;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.tree-error {
  margin: 0;
  color: #fca5a5;
  font-size: 12px;
}

.tree-retry-button {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-primary);
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
}

.tree-children {
  list-style: none;
  margin: 8px 0 0 18px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
