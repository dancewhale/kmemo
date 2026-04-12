import { defineStore } from "pinia";
import { ref } from "vue";

import { createCard, getCardChildren, getCardDetail, listCards } from "../bridge/cardService";
import { listTags } from "../bridge/tagService";
import type { CardDTO, CardDetailDTO, CardFilters, CreateCardRequest, TagDTO } from "../types/dto";

function normalizeListCardsResult(result: [CardDTO[], number] | { 0: CardDTO[]; 1: number }) {
  if (Array.isArray(result)) {
    return { items: result[0] ?? [], total: result[1] ?? 0 };
  }
  return { items: result[0] ?? [], total: result[1] ?? 0 };
}

export const useCardStore = defineStore("card", () => {
  const cards = ref<CardDTO[]>([]);
  const rootCards = ref<CardDTO[]>([]);
  const total = ref(0);
  const selectedCard = ref<CardDetailDTO | null>(null);
  const selectedCardChildren = ref<CardDTO[]>([]);
  const treeChildren = ref<Record<string, CardDTO[]>>({});
  const treeExpanded = ref<Record<string, boolean>>({});
  const treeLoading = ref<Record<string, boolean>>({});
  const treeErrors = ref<Record<string, string | null>>({});
  const tags = ref<TagDTO[]>([]);
  const loadingList = ref(false);
  const loadingRootTree = ref(false);
  const loadingDetail = ref(false);
  const creating = ref(false);
  const error = ref<string | null>(null);
  const rootTreeError = ref<string | null>(null);
  const detailError = ref<string | null>(null);
  const createError = ref<string | null>(null);
  let listRequestId = 0;
  let rootTreeRequestId = 0;
  let treeScopeId = 0;
  let detailRequestId = 0;

  function resetSelection() {
    detailRequestId += 1;
    selectedCard.value = null;
    selectedCardChildren.value = [];
    loadingDetail.value = false;
    detailError.value = null;
    createError.value = null;
  }

  function resetWorkspaceState() {
    listRequestId += 1;
    rootTreeRequestId += 1;
    cards.value = [];
    rootCards.value = [];
    total.value = 0;
    error.value = null;
    rootTreeError.value = null;
    resetSelection();
    resetTreeState();
  }

  function resetTreeState() {
    treeScopeId += 1;
    treeChildren.value = {};
    treeExpanded.value = {};
    treeLoading.value = {};
    treeErrors.value = {};
  }

  async function loadTags() {
    tags.value = await listTags();
  }

  async function loadCards(filters: CardFilters) {
    const requestId = ++listRequestId;
    loadingList.value = true;
    error.value = null;
    try {
      const result = await listCards(filters);
      if (requestId !== listRequestId) return;
      const normalized = normalizeListCardsResult(result);
      cards.value = normalized.items;
      total.value = normalized.total;
    } catch (err) {
      if (requestId !== listRequestId) return;
      error.value = err instanceof Error ? err.message : "Failed to load cards";
    } finally {
      if (requestId === listRequestId) {
        loadingList.value = false;
      }
    }
  }

  async function loadRootCards(knowledgeId: string | null) {
    const requestId = ++rootTreeRequestId;
    loadingRootTree.value = true;
    rootTreeError.value = null;
    if (!knowledgeId) {
      if (requestId !== rootTreeRequestId) return;
      rootCards.value = [];
      resetTreeState();
      loadingRootTree.value = false;
      return;
    }
    try {
      const result = await listCards({
        knowledgeId,
        isRoot: true,
        limit: 200,
        offset: 0,
        orderBy: "sort_order",
        orderDesc: false,
      });
      if (requestId !== rootTreeRequestId) return;
      rootCards.value = normalizeListCardsResult(result).items;
      resetTreeState();
    } catch (err) {
      if (requestId !== rootTreeRequestId) return;
      rootCards.value = [];
      resetTreeState();
      rootTreeError.value = err instanceof Error ? err.message : "Failed to load card tree";
    } finally {
      if (requestId === rootTreeRequestId) {
        loadingRootTree.value = false;
      }
    }
  }

  async function expandTreeNode(cardId: string) {
    if (treeLoading.value[cardId]) return;
    const cachedChildren = treeChildren.value[cardId];
    if (cachedChildren) {
      if (cachedChildren.length === 0) return;
      if (treeExpanded.value[cardId]) {
        treeExpanded.value = { ...treeExpanded.value, [cardId]: false };
        return;
      }
      treeExpanded.value = { ...treeExpanded.value, [cardId]: true };
      return;
    }

    const scopeId = treeScopeId;
    treeLoading.value = { ...treeLoading.value, [cardId]: true };
    treeErrors.value = { ...treeErrors.value, [cardId]: null };
    try {
      const children = await getCardChildren(cardId);
      if (scopeId !== treeScopeId) return;
      treeChildren.value = { ...treeChildren.value, [cardId]: children };
      treeExpanded.value = { ...treeExpanded.value, [cardId]: children.length > 0 };
    } catch (err) {
      if (scopeId !== treeScopeId) return;
      treeErrors.value = {
        ...treeErrors.value,
        [cardId]: err instanceof Error ? err.message : "Failed to load child cards",
      };
    } finally {
      if (scopeId === treeScopeId) {
        treeLoading.value = { ...treeLoading.value, [cardId]: false };
      }
    }
  }

  async function loadCardDetail(cardId: string) {
    const requestId = ++detailRequestId;
    loadingDetail.value = true;
    detailError.value = null;
    selectedCard.value = null;
    selectedCardChildren.value = [];
    try {
      const [detail, children] = await Promise.all([
        getCardDetail(cardId),
        getCardChildren(cardId),
      ]);
      if (requestId !== detailRequestId) return;
      selectedCard.value = detail;
      selectedCardChildren.value = children;
    } catch (err) {
      if (requestId !== detailRequestId) return;
      selectedCard.value = null;
      selectedCardChildren.value = [];
      detailError.value = err instanceof Error ? err.message : "Failed to load card detail";
    } finally {
      if (requestId === detailRequestId) {
        loadingDetail.value = false;
      }
    }
  }

  async function submitCreateCard(payload: CreateCardRequest) {
    creating.value = true;
    createError.value = null;
    try {
      return await createCard(payload);
    } catch (err) {
      createError.value = err instanceof Error ? err.message : "Failed to create card";
      throw err;
    } finally {
      creating.value = false;
    }
  }

  return {
    cards,
    rootCards,
    total,
    selectedCard,
    selectedCardChildren,
    treeChildren,
    treeExpanded,
    treeLoading,
    treeErrors,
    tags,
    loadingList,
    loadingRootTree,
    loadingDetail,
    creating,
    error,
    rootTreeError,
    detailError,
    createError,
    resetSelection,
    resetWorkspaceState,
    loadTags,
    loadCards,
    loadRootCards,
    expandTreeNode,
    loadCardDetail,
    submitCreateCard,
  };
});
