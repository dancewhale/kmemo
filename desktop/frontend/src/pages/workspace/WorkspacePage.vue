<template>
  <WorkspaceLayout>
    <template #sidebar>
      <KnowledgeTreePanel
        :error="workspace.error"
        :loading="workspace.loading"
        :nodes="workspace.knowledgeTree"
        :selected-id="workspace.selectedKnowledgeId"
        @create-card="handleOpenCreateCard"
        @select="handleSelectKnowledge"
      />
    </template>

    <template #main>
      <div class="workspace-main-stack">
        <CardListPanel
          :error="card.error"
          :items="card.cards"
          :keyword="keyword"
          :loading="card.loadingList"
          :selected-id="workspace.selectedCardId"
          :subtitle="workspace.selectedKnowledge?.name ?? 'All cards'"
          :total="card.total"
          @retry="handleRetryCardList"
          @select="handleSelectCard"
          @update:keyword="handleKeywordChange"
        />
        <CardTreePanel
          :knowledge-id="workspace.selectedKnowledgeId"
          :on-reload="refreshCards"
          :on-select-card="handleSelectCard"
          :items="card.rootCards"
          :error="card.rootTreeError"
          :loading="card.loadingRootTree"
          :selected-id="workspace.selectedCardId"
          :tree-children="card.treeChildren"
          :tree-errors="card.treeErrors"
          :tree-expanded="card.treeExpanded"
          :tree-loading="card.treeLoading"
          @retry="handleRetryCardTree"
          @select="handleSelectCard"
          @toggle="handleToggleTreeNode"
        />
      </div>
    </template>

    <template #detail>
      <CardDetailPanel
        :card="card.selectedCard"
        :children="card.selectedCardChildren"
        :creating="card.creating"
        :create-error="card.createError"
        :error="card.detailError"
        :knowledge-id="workspace.selectedKnowledgeId"
        :loading="card.loadingDetail"
        :mode="workspace.mode"
        :tags="card.tags"
        @cancel-create="handleCancelCreateCard"
        @retry-detail="handleRetryCardDetail"
        @select-child="handleSelectCard"
        @select-parent="handleSelectCard"
        @submit-create="handleCreateCard"
      />
    </template>
  </WorkspaceLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";

import CardDetailPanel from "../../features/card-detail/CardDetailPanel.vue";
import CardListPanel from "../../features/card-list/CardListPanel.vue";
import CardTreePanel from "../../features/card-tree/CardTreePanel.vue";
import KnowledgeTreePanel from "../../features/knowledge-tree/KnowledgeTreePanel.vue";
import WorkspaceLayout from "../../layouts/WorkspaceLayout.vue";
import { useCardStore } from "../../stores/card";
import { useWorkspaceStore } from "../../stores/workspace";
import type { CreateCardRequest } from "../../types/dto";

const workspace = useWorkspaceStore();
const card = useCardStore();
const keyword = ref("");
let keywordTimer: ReturnType<typeof setTimeout> | null = null;
let createSessionId = 0;

const cardFilters = computed(() => ({
  knowledgeId: workspace.selectedKnowledgeId,
  keyword: keyword.value.trim(),
  limit: 50,
  offset: 0,
}));

async function refreshCards() {
  if (!workspace.selectedKnowledgeId) {
    card.resetWorkspaceState();
    return;
  }
  await Promise.all([
    card.loadCards(cardFilters.value),
    card.loadRootCards(workspace.selectedKnowledgeId),
  ]);
}

async function handleSelectKnowledge(id: string) {
  if (keywordTimer) {
    clearTimeout(keywordTimer);
    keywordTimer = null;
  }
  workspace.selectKnowledge(id);
  keyword.value = "";
  card.resetSelection();
  await refreshCards();
}

async function handleSelectCard(id: string) {
  workspace.selectCard(id);
  await card.loadCardDetail(id);
}

async function handleToggleTreeNode(id: string) {
  await card.expandTreeNode(id);
}

async function handleRetryCardTree() {
  await card.loadRootCards(workspace.selectedKnowledgeId);
}

async function handleRetryCardList() {
  if (keywordTimer) {
    clearTimeout(keywordTimer);
    keywordTimer = null;
  }
  if (!workspace.selectedKnowledgeId) return;
  await card.loadCards(cardFilters.value);
}

async function handleRetryCardDetail() {
  if (!workspace.selectedCardId) return;
  await card.loadCardDetail(workspace.selectedCardId);
}

function handleOpenCreateCard() {
  createSessionId += 1;
  card.resetSelection();
  workspace.openCreateCard();
}

function handleCancelCreateCard() {
  createSessionId += 1;
  card.resetSelection();
  workspace.closeCreateCard();
}

async function handleCreateCard(payload: CreateCardRequest) {
  const submitKnowledgeId = workspace.selectedKnowledgeId;
  const submitFilters = { ...cardFilters.value };
  const submitSessionId = createSessionId;
  const id = await card.submitCreateCard(payload);
  if (!id) {
    return;
  }
  if (createSessionId !== submitSessionId) {
    return;
  }
  if (workspace.selectedKnowledgeId !== submitKnowledgeId) {
    return;
  }
  await Promise.all([
    card.loadCards(submitFilters),
    card.loadRootCards(submitKnowledgeId),
  ]);
  workspace.selectCard(id);
  await card.loadCardDetail(id);
  workspace.closeCreateCard();
}

function handleKeywordChange(value: string) {
  keyword.value = value;
}

watch(keyword, () => {
  if (keywordTimer) {
    clearTimeout(keywordTimer);
  }
  if (!workspace.selectedKnowledgeId) return;
  keywordTimer = setTimeout(async () => {
    keywordTimer = null;
    await card.loadCards(cardFilters.value);
  }, 250);
});

onBeforeUnmount(() => {
  if (keywordTimer) {
    clearTimeout(keywordTimer);
  }
});

onMounted(async () => {
  await Promise.all([workspace.loadKnowledgeTree(), card.loadTags()]);
  if (workspace.selectedKnowledgeId) {
    await refreshCards();
  }
});
</script>

<style scoped>
.workspace-main-stack {
  height: 100%;
  display: grid;
  grid-template-rows: minmax(0, 1fr) minmax(220px, 280px);
}
</style>
