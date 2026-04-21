<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCaptureStore } from '../stores/capture.store'

const capture = useCaptureStore()
const { draft } = storeToRefs(capture)

const tagLine = computed({
  get: () => (draft.value.blank.tags ?? []).join(', '),
  set: (v: string) => capture.applyTagsLine('blank', v),
})
</script>

<template>
  <div class="new-article-form">
    <label class="new-article-form__label">Title <span class="new-article-form__hint">optional</span></label>
    <el-input
      :model-value="draft.blank.title"
      size="small"
      placeholder="Untitled is fine"
      @update:model-value="capture.updateBlankDraft({ title: String($event ?? '') })"
    />
    <label class="new-article-form__label">Content</label>
    <el-input
      :model-value="draft.blank.content"
      type="textarea"
      :autosize="{ minRows: 6, maxRows: 16 }"
      size="small"
      placeholder="Start typing or leave empty…"
      @update:model-value="capture.updateBlankDraft({ content: String($event ?? '') })"
    />
    <label class="new-article-form__label">Tags <span class="new-article-form__hint">comma-separated</span></label>
    <el-input v-model="tagLine" size="small" placeholder="e.g. work, reference" />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.new-article-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.new-article-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.new-article-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}
</style>
