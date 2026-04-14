<template>
  <Draggable
    v-model="treeData"
    class="he-tree"
    tree-line
    @before-drag-start="handleBeforeDragStart"
    @after-drop="handleAfterDrop"
  >
    <template #default="{ node, stat }">
      <div class="tree-row">
        <button
          v-if="node.hasChildrenHint || stat.children.length > 0"
          class="toggle-button"
          type="button"
          :disabled="isNodeLoading(node.id)"
          @click="handleToggle(node, stat)"
        >
          <span v-if="isNodeLoading(node.id)">…</span>
          <span v-else>{{ stat.open ? "−" : "+" }}</span>
        </button>
        <span v-else class="toggle-placeholder"></span>

        <CardTreeNodeContent
          :node="node"
          :selected="node.id === selectedId"
          @select="$emit('node-select', node.id)"
          @create-child="props.onNodeCreateChild(node.source)"
          @rename="props.onNodeRename(node.source)"
          @delete="props.onNodeDelete(node.source)"
        />
      </div>
      <p v-if="treeErrors[node.id]" class="tree-error">{{ treeErrors[node.id] }}</p>
    </template>
  </Draggable>
</template>

<script setup lang="ts">
import { Draggable } from "@he-tree/vue";
import "@he-tree/vue/style/default.css";
import { nextTick, ref, watch } from "vue";

import type { CardDTO } from "../../types/dto";
import CardTreeNodeContent from "./CardTreeNodeContent.vue";
import type { CardTreeDropPayload, CardTreeNodeVM } from "./types/cardTree";

const props = defineProps<{
  rootItems: CardDTO[];
  selectedId: string | null;
  treeChildren: Record<string, CardDTO[]>;
  treeExpanded: Record<string, boolean>;
  treeLoading: Record<string, boolean>;
  treeErrors: Record<string, string | null>;
  onNodeCreateChild: (node: CardDTO) => Promise<void>;
  onNodeRename: (node: CardDTO) => Promise<void>;
  onNodeDelete: (node: CardDTO) => Promise<void>;
  onNodeDrop: (payload: CardTreeDropPayload) => Promise<void>;
}>();

const emit = defineEmits<{
  (event: "node-select", id: string): void;
  (event: "node-toggle", id: string): void;
}>();

const treeData = ref<CardTreeNodeVM[]>([]);
const lastStableTree = ref<CardTreeNodeVM[]>([]);
const syncingFromProps = ref(false);
const applyingDrop = ref(false);
const dragging = ref(false);

function cloneTree(input: CardTreeNodeVM[]): CardTreeNodeVM[] {
  return input.map((node) => ({
    id: node.id,
    parentId: node.parentId,
    open: node.open,
    title: node.title,
    cardType: node.cardType,
    status: node.status,
    hasChildrenHint: node.hasChildrenHint,
    depth: node.depth,
    draggable: node.draggable,
    droppable: node.droppable,
    // 保持 source 引用，避免克隆 Vue proxy 触发 DataCloneError。
    source: node.source,
    children: cloneTree(node.children),
  }));
}

function buildNodes(items: CardDTO[], depth: number, parentId: string | null): CardTreeNodeVM[] {
  return items.map((item) => {
    const children = props.treeChildren[item.id] ?? [];
    const hasChildrenHint = children.length > 0 || !Object.prototype.hasOwnProperty.call(props.treeChildren, item.id);
    return {
      id: item.id,
      parentId,
      open: Boolean(props.treeExpanded[item.id]),
      title: item.title,
      cardType: item.cardType,
      status: item.status,
      hasChildrenHint,
      depth,
      draggable: true,
      droppable: true,
      source: item,
      children: props.treeExpanded[item.id] ? buildNodes(children, depth + 1, item.id) : [],
    };
  });
}

function flattenParentIndex(nodes: CardTreeNodeVM[], parentId: string | null, map: Map<string, { parentId: string | null; index: number }>) {
  nodes.forEach((node, index) => {
    map.set(node.id, { parentId, index });
    flattenParentIndex(node.children, node.id, map);
  });
}

function buildSiblingMap(nodes: CardTreeNodeVM[], map: Map<string | null, string[]>, parentId: string | null = null) {
  map.set(parentId, nodes.map((node) => node.id));
  nodes.forEach((node) => buildSiblingMap(node.children, map, node.id));
}

function syncFromProps() {
  syncingFromProps.value = true;
  const nextTree = buildNodes(props.rootItems, 0, null);
  treeData.value = nextTree;
  lastStableTree.value = cloneTree(nextTree);
  nextTick(() => {
    syncingFromProps.value = false;
  });
}

function isNodeLoading(id: string): boolean {
  return Boolean(props.treeLoading[id]);
}

function handleToggle(node: CardTreeNodeVM, stat: { open: boolean; children: CardTreeNodeVM[] }) {
  if (isNodeLoading(node.id)) return;
  stat.open = !stat.open;
  emit("node-toggle", node.id);
}

function findMovedNode(nextTree: CardTreeNodeVM[]) {
  const before = new Map<string, { parentId: string | null; index: number }>();
  const after = new Map<string, { parentId: string | null; index: number }>();
  flattenParentIndex(lastStableTree.value, null, before);
  flattenParentIndex(nextTree, null, after);
  const moved = [...after.entries()].find(([id, pos]) => {
    const previous = before.get(id);
    return previous && (previous.parentId !== pos.parentId || previous.index !== pos.index);
  });
  return { before, moved };
}

async function persistDroppedTree(nextTree: CardTreeNodeVM[]) {
  const { before, moved } = findMovedNode(nextTree);
  if (!moved) {
    lastStableTree.value = cloneTree(nextTree);
    return;
  }
  const [cardId, targetPos] = moved;
  const sourcePos = before.get(cardId);
  const siblingMap = new Map<string | null, string[]>();
  buildSiblingMap(nextTree, siblingMap);
  const payload: CardTreeDropPayload = {
    cardId,
    sourceParentId: sourcePos?.parentId ?? null,
    targetParentId: targetPos.parentId,
    targetIndex: targetPos.index,
    siblingIds: siblingMap.get(targetPos.parentId) ?? [],
  };

  applyingDrop.value = true;
  try {
    await props.onNodeDrop(payload);
    lastStableTree.value = cloneTree(nextTree);
  } catch {
    treeData.value = cloneTree(lastStableTree.value);
  } finally {
    nextTick(() => {
      applyingDrop.value = false;
    });
  }
}

function handleBeforeDragStart() {
  dragging.value = true;
}

async function handleAfterDrop() {
  dragging.value = false;
  await persistDroppedTree(treeData.value);
}

watch(
  () => [props.rootItems, props.treeChildren, props.treeExpanded] as const,
  () => {
    if (applyingDrop.value || dragging.value) return;
    syncFromProps();
  },
  { deep: true, immediate: true },
);

watch(
  treeData,
  async (nextTree) => {
    if (syncingFromProps.value || applyingDrop.value || dragging.value) return;
    // 非拖拽场景（如库内部格式化）只同步本地快照，不触发后端写操作。
    lastStableTree.value = cloneTree(nextTree);
  },
  { deep: true },
);
</script>

<style scoped>
.he-tree {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

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

.toggle-button:disabled {
  opacity: 0.65;
  cursor: wait;
}

.toggle-placeholder {
  display: inline-block;
}

.tree-error {
  margin: 4px 0 0 40px;
  color: #fca5a5;
  font-size: 12px;
}
</style>
