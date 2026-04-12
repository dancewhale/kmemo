import { defineStore } from "pinia";
import { computed, ref } from "vue";

import type { KnowledgeTreeNode } from "../types/dto";
import { getKnowledgeTree } from "../bridge/knowledgeService";

export const useWorkspaceStore = defineStore("workspace", () => {
  const loading = ref(false);
  const knowledgeTree = ref<KnowledgeTreeNode[]>([]);
  const selectedKnowledgeId = ref<string | null>(null);
  const selectedCardId = ref<string | null>(null);
  const mode = ref<"browse" | "create">("browse");
  const error = ref<string | null>(null);

  const selectedKnowledge = computed(() => {
    const stack = [...knowledgeTree.value];
    while (stack.length > 0) {
      const node = stack.shift();
      if (!node) continue;
      if (node.id === selectedKnowledgeId.value) return node;
      stack.push(...node.children);
    }
    return null;
  });

  async function loadKnowledgeTree() {
    loading.value = true;
    error.value = null;
    try {
      const tree = await getKnowledgeTree(null);
      knowledgeTree.value = tree;
      if (!selectedKnowledgeId.value && tree.length > 0) {
        selectedKnowledgeId.value = tree[0].id;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load knowledge tree";
    } finally {
      loading.value = false;
    }
  }

  function selectKnowledge(id: string) {
    selectedKnowledgeId.value = id;
    selectedCardId.value = null;
    mode.value = "browse";
  }

  function selectCard(id: string) {
    selectedCardId.value = id;
    mode.value = "browse";
  }

  function openCreateCard() {
    mode.value = "create";
    selectedCardId.value = null;
  }

  function closeCreateCard() {
    mode.value = "browse";
  }

  return {
    loading,
    error,
    knowledgeTree,
    selectedKnowledgeId,
    selectedCardId,
    selectedKnowledge,
    mode,
    loadKnowledgeTree,
    selectKnowledge,
    selectCard,
    openCreateCard,
    closeCreateCard,
  };
});
