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

    <CardTreeHeTree
      v-else
      :root-items="items"
      :selected-id="selectedId"
      :tree-children="treeChildren"
      :tree-errors="treeErrors"
      :tree-expanded="treeExpanded"
      :tree-loading="treeLoading"
      :on-node-create-child="commands.onNodeCreateChild"
      :on-node-rename="commands.onNodeRename"
      :on-node-delete="commands.onNodeDelete"
      :on-node-drop="commands.onNodeDrop"
      @node-select="$emit('select', $event)"
      @node-toggle="$emit('toggle', $event)"
    />
  </section>
</template>

<script setup lang="ts">
import CardTreeHeTree from "./CardTreeHeTree.vue";
import { useCardTreeCommands } from "./composables/useCardTreeCommands";
import type { CardDTO } from "../../types/dto";

const props = defineProps<{
  knowledgeId: string | null;
  onReload: () => Promise<void>;
  onSelectCard: (id: string) => Promise<void>;
  items: CardDTO[];
  error: string | null;
  loading: boolean;
  selectedId: string | null;
  treeChildren: Record<string, CardDTO[]>;
  treeErrors: Record<string, string | null>;
  treeExpanded: Record<string, boolean>;
  treeLoading: Record<string, boolean>;
}>();

const emit = defineEmits<{
  (event: "retry"): void;
  (event: "select", id: string): void;
  (event: "toggle", id: string): void;
  (event: "reload"): Promise<void>;
}>();

const commands = useCardTreeCommands({
  getKnowledgeId: () => props.knowledgeId,
  onReload: props.onReload,
  onSelect: props.onSelectCard,
});
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

</style>
