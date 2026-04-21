<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCaptureStore } from '../stores/capture.store'

const capture = useCaptureStore()
const { draft } = storeToRefs(capture)

const tagLine = computed({
  get: () => (draft.value.pasteText.tags ?? []).join(', '),
  set: (v: string) => capture.applyTagsLine('paste-text', v),
})
</script>

<template>
  <div class="paste-text-form">
    <label class="paste-text-form__label">Content</label>
    <el-input
      :model-value="draft.pasteText.content"
      type="textarea"
      :autosize="{ minRows: 10, maxRows: 22 }"
      size="small"
      placeholder="Paste or type the text you want to capture…"
      @update:model-value="capture.updatePasteTextDraft({ content: String($event ?? '') })"
    />
    <label class="paste-text-form__label">Title <span class="paste-text-form__hint">optional; generated if empty</span></label>
    <el-input
      :model-value="draft.pasteText.title"
      size="small"
      placeholder="Override auto title…"
      @update:model-value="capture.updatePasteTextDraft({ title: String($event ?? '') })"
    />
    <label class="paste-text-form__label">Tags <span class="paste-text-form__hint">comma-separated</span></label>
    <el-input v-model="tagLine" size="small" placeholder="optional" />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.paste-text-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.paste-text-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.paste-text-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}
</style>
