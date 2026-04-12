<template>
  <section class="panel-block">
    <template v-if="mode === 'create'">
      <CardEditorPanel
        :creating="creating"
        :error="createError"
        :knowledge-id="knowledgeId"
        :tags="tags"
        @cancel="$emit('cancel-create')"
        @submit="$emit('submit-create', $event)"
      />
    </template>

    <template v-else-if="loading">
      <h2 class="panel-title">Card Detail</h2>
      <p class="panel-description">Loading card detail...</p>
    </template>

    <template v-else-if="error">
      <h2 class="panel-title">Card Detail</h2>
      <div class="detail-feedback-block">
        <p class="panel-description">{{ error }}</p>
        <button class="detail-retry-button" type="button" @click="$emit('retry-detail')">Retry</button>
      </div>
    </template>

    <template v-else-if="card">
      <h2 class="panel-title">Card Detail</h2>
      <p class="panel-description">{{ card.cardType }} · {{ card.status }}</p>
      <div class="detail-section">
        <strong>Title</strong>
        <span>{{ card.title }}</span>
      </div>
      <div class="detail-section">
        <strong>Knowledge</strong>
        <span>{{ card.knowledgeName || card.knowledgeId }}</span>
      </div>
      <div class="detail-section" v-if="card.parent">
        <strong>Parent</strong>
        <button class="related-card-button" type="button" @click="$emit('select-parent', card.parent.id)">
          <span>{{ card.parent.title }}</span>
          <span>{{ card.parent.cardType }}</span>
        </button>
      </div>
      <div class="detail-section" v-if="card.tags?.length">
        <strong>Tags</strong>
        <div class="tag-row">
          <span v-for="tag in card.tags" :key="tag.id" class="tag-chip">{{ tag.name }}</span>
        </div>
      </div>
      <div class="detail-section" v-if="card.srs">
        <strong>SRS</strong>
        <span>{{ card.srs.fsrsState }} · reps {{ card.srs.reps }}</span>
      </div>
      <div class="detail-section" v-if="children.length">
        <strong>Children</strong>
        <div class="child-list">
          <button v-for="child in children" :key="child.id" class="related-card-button" type="button" @click="$emit('select-child', child.id)">
            <span>{{ child.title }}</span>
            <span>{{ child.cardType }}</span>
          </button>
        </div>
      </div>
      <div class="detail-section">
        <strong>HTML</strong>
        <pre class="html-preview">{{ card.htmlContent || card.htmlPath || 'HTML content is not available yet.' }}</pre>
      </div>
    </template>

    <template v-else>
      <h2 class="panel-title">Card Detail</h2>
      <p class="panel-description">Select a card from the middle pane, or create a new one from the sidebar.</p>
    </template>
  </section>
</template>

<script setup lang="ts">
import CardEditorPanel from "../card-editor/CardEditorPanel.vue";
import type { CardDTO, CardDetailDTO, CreateCardRequest, TagDTO } from "../../types/dto";

defineProps<{
  card: CardDetailDTO | null;
  children: CardDTO[];
  tags: TagDTO[];
  knowledgeId: string | null;
  mode: "browse" | "create";
  loading: boolean;
  error: string | null;
  createError: string | null;
  creating: boolean;
}>();

defineEmits<{
  (event: "cancel-create"): void;
  (event: "retry-detail"): void;
  (event: "submit-create", payload: CreateCardRequest): void;
  (event: "select-parent", id: string): void;
  (event: "select-child", id: string): void;
}>();
</script>

<style scoped>
.detail-feedback-block {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.detail-retry-button {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-primary);
  padding: 8px 12px;
  cursor: pointer;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tag-row,
.child-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-chip {
  border-radius: 999px;
  padding: 4px 10px;
  background: rgba(15, 23, 42, 0.65);
  color: var(--text-secondary);
  font-size: 12px;
}

.related-card-button {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 8px 10px;
  background: rgba(15, 23, 42, 0.55);
  color: var(--text-primary);
  cursor: pointer;
  display: inline-flex;
  gap: 8px;
}

.html-preview {
  margin: 0;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.7);
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
