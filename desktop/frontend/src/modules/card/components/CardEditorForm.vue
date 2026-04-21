<script setup lang="ts">
import type { CardItem } from '../types'

const props = defineProps<{
  card: CardItem
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:title': [value: string]
  'update:prompt': [value: string]
  'update:answer': [value: string]
}>()
</script>

<template>
  <div class="card-editor-form">
    <label class="card-editor-form__label">Title</label>
    <el-input
      :model-value="props.card.title"
      size="small"
      maxlength="160"
      :disabled="props.disabled"
      @update:model-value="emit('update:title', String($event ?? ''))"
    />
    <label class="card-editor-form__label">Prompt</label>
    <el-input
      :model-value="props.card.prompt"
      type="textarea"
      :autosize="{ minRows: 3, maxRows: 10 }"
      size="small"
      :disabled="props.disabled"
      @update:model-value="emit('update:prompt', String($event ?? ''))"
    />
    <label class="card-editor-form__label">Answer</label>
    <el-input
      :model-value="props.card.answer"
      type="textarea"
      :autosize="{ minRows: 3, maxRows: 12 }"
      size="small"
      :disabled="props.disabled"
      @update:model-value="emit('update:answer', String($event ?? ''))"
    />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.card-editor-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.card-editor-form__label {
  margin-top: $space-sm;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
}

.card-editor-form__label:first-of-type {
  margin-top: 0;
}
</style>
