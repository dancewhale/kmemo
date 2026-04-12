<template>
  <form class="editor-form" @submit.prevent="handleSubmit">
    <div>
      <h2 class="panel-title">Create Card</h2>
      <p class="panel-description">Create a new card inside the selected knowledge node.</p>
    </div>

    <p v-if="error" class="editor-error">{{ error }}</p>

    <label class="field">
      <span>Knowledge ID</span>
      <input :value="knowledgeId ?? ''" disabled type="text" />
    </label>

    <label class="field">
      <span>Title</span>
      <input v-model="title" required type="text" />
    </label>

    <label class="field">
      <span>Card Type</span>
      <select v-model="cardType">
        <option value="article">article</option>
        <option value="excerpt">excerpt</option>
        <option value="qa">qa</option>
        <option value="cloze">cloze</option>
        <option value="note">note</option>
      </select>
    </label>

    <label class="field">
      <span>Tags</span>
      <select v-model="selectedTagIds" multiple>
        <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
      </select>
    </label>

    <label class="field field-grow">
      <span>HTML Content</span>
      <textarea v-model="htmlContent" rows="12"></textarea>
    </label>

    <div class="actions">
      <button class="ghost-button" type="button" @click="handleCancel">Cancel</button>
      <button class="primary-button" :disabled="creating || !knowledgeId" type="submit">
        {{ creating ? 'Creating...' : 'Create Card' }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

import type { CreateCardRequest, TagDTO } from "../../types/dto";

const props = defineProps<{
  knowledgeId: string | null;
  tags: TagDTO[];
  creating: boolean;
  error: string | null;
}>();

const emit = defineEmits<{
  (event: "cancel"): void;
  (event: "submit", payload: CreateCardRequest): void;
}>();

const title = ref("");
const cardType = ref("article");
const htmlContent = ref("");
const selectedTagIds = ref<string[]>([]);

function resetForm() {
  title.value = "";
  cardType.value = "article";
  htmlContent.value = "";
  selectedTagIds.value = [];
}

function handleCancel() {
  resetForm();
  emit("cancel");
}

function handleSubmit() {
  if (!props.knowledgeId) return;
  emit("submit", {
    knowledgeId: props.knowledgeId,
    title: title.value.trim(),
    cardType: cardType.value,
    htmlContent: htmlContent.value,
    tagIds: selectedTagIds.value,
  });
}

watch(
  () => props.knowledgeId,
  () => {
    resetForm();
  },
);
</script>

<style scoped>
.editor-form {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.editor-error {
  margin: 0;
  color: #fca5a5;
  font-size: 13px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field input,
.field textarea,
.field select {
  width: 100%;
  border: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.7);
  color: var(--text-primary);
  border-radius: 10px;
  padding: 10px 12px;
}

.field-grow {
  flex: 1;
}

.field-grow textarea {
  min-height: 220px;
  resize: vertical;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.primary-button,
.ghost-button {
  border-radius: 10px;
  padding: 10px 14px;
  cursor: pointer;
}

.primary-button {
  border: 1px solid var(--accent);
  background: var(--accent-soft);
  color: var(--text-primary);
}

.ghost-button {
  border: 1px solid var(--border-color);
  background: transparent;
  color: var(--text-secondary);
}
</style>
