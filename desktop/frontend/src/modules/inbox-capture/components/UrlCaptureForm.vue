<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCaptureStore } from '../stores/capture.store'

const capture = useCaptureStore()
const { draft } = storeToRefs(capture)

const tagLine = computed({
  get: () => (draft.value.url.tags ?? []).join(', '),
  set: (v: string) => capture.applyTagsLine('url', v),
})
</script>

<template>
  <div class="url-capture-form">
    <p class="url-capture-form__notice">
      Phase 1 only creates a placeholder article linked to this URL. No fetching or parsing is performed.
    </p>
    <label class="url-capture-form__label">URL</label>
    <el-input
      :model-value="draft.url.url"
      size="small"
      placeholder="https://…"
      @update:model-value="capture.updateUrlDraft({ url: String($event ?? '') })"
    />
    <label class="url-capture-form__label">Title <span class="url-capture-form__hint">optional</span></label>
    <el-input
      :model-value="draft.url.title ?? ''"
      size="small"
      placeholder="Defaults from host / path"
      @update:model-value="capture.updateUrlDraft({ title: String($event ?? '') })"
    />
    <label class="url-capture-form__label">Note <span class="url-capture-form__hint">optional</span></label>
    <el-input
      :model-value="draft.url.note ?? ''"
      type="textarea"
      :autosize="{ minRows: 2, maxRows: 6 }"
      size="small"
      placeholder="Why you are saving this link…"
      @update:model-value="capture.updateUrlDraft({ note: String($event ?? '') })"
    />
    <label class="url-capture-form__label">Tags <span class="url-capture-form__hint">comma-separated</span></label>
    <el-input v-model="tagLine" size="small" placeholder="optional" />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.url-capture-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.url-capture-form__notice {
  margin: 0 0 $space-sm;
  padding: $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane;
  font-size: $font-size-xs;
  line-height: $line-normal;
  color: $color-text-secondary;
}

.url-capture-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.url-capture-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}
</style>
