<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useObjectCreationStore } from '../stores/object-creation.store'

const creation = useObjectCreationStore()
const tree = useTreeStore()
const { draft } = storeToRefs(creation)

onMounted(() => {
  void tree.initialize()
})

const parentOptions = computed(() => {
  const roots = [{ value: '' as string, label: '(root)' }]
  const rows = tree.rawNodes.map((n) => ({
    value: n.id,
    label: n.title,
  }))
  return [...roots, ...rows]
})

function parentModel(): string {
  return creation.draft.card.parentNodeId ?? ''
}

function setParent(v: string) {
  creation.updateCardDraft({ parentNodeId: v === '' ? null : v })
}
</script>

<template>
  <div class="create-card-form">
    <label class="create-card-form__label">Title <span class="create-card-form__hint">optional</span></label>
    <el-input
      :model-value="draft.card.title"
      size="small"
      placeholder="From prompt if empty"
      @update:model-value="creation.updateCardDraft({ title: String($event ?? '') })"
    />
    <label class="create-card-form__label">Prompt</label>
    <el-input
      :model-value="draft.card.prompt"
      type="textarea"
      :autosize="{ minRows: 3, maxRows: 8 }"
      size="small"
      placeholder="Front / question…"
      @update:model-value="creation.updateCardDraft({ prompt: String($event ?? '') })"
    />
    <label class="create-card-form__label">Answer</label>
    <el-input
      :model-value="draft.card.answer"
      type="textarea"
      :autosize="{ minRows: 3, maxRows: 8 }"
      size="small"
      placeholder="Back / answer…"
      @update:model-value="creation.updateCardDraft({ answer: String($event ?? '') })"
    />
    <label class="create-card-form__label">Parent</label>
    <el-select :model-value="parentModel()" size="small" class="create-card-form__select" @update:model-value="setParent(String($event ?? ''))">
      <el-option v-for="o in parentOptions" :key="o.value || 'root'" :label="o.label" :value="o.value" />
    </el-select>
    <div class="create-card-form__row">
      <span class="create-card-form__label create-card-form__label--inline">Add to review queue</span>
      <el-switch :model-value="draft.card.addToReview" @update:model-value="creation.updateCardDraft({ addToReview: Boolean($event) })" />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-card-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.create-card-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.create-card-form__label--inline {
  margin-top: 0;
  text-transform: none;
  letter-spacing: 0;
}

.create-card-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}

.create-card-form__select {
  width: 100%;
}

.create-card-form__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-md;
  margin-top: $space-sm;
  padding-top: $space-xs;
  border-top: 1px solid $color-border-subtle;
}
</style>
