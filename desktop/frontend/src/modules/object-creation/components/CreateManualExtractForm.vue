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
  return creation.draft.manualExtract.parentNodeId ?? ''
}

function setParent(v: string) {
  creation.updateManualExtractDraft({ parentNodeId: v === '' ? null : v })
}

const sourceLine = computed(() => {
  const id = draft.value.manualExtract.sourceArticleId
  const t = draft.value.manualExtract.sourceArticleTitle
  if (!id) {
    return ''
  }
  return t ? `${t} (${id})` : id
})
</script>

<template>
  <div class="create-manual-extract-form">
    <p v-if="sourceLine" class="create-manual-extract-form__source">
      Source article: <span class="kmono">{{ sourceLine }}</span>
    </p>
    <label class="create-manual-extract-form__label">Quote</label>
    <el-input
      :model-value="draft.manualExtract.quote"
      type="textarea"
      :autosize="{ minRows: 5, maxRows: 14 }"
      size="small"
      placeholder="Primary excerpt text…"
      @update:model-value="creation.updateManualExtractDraft({ quote: String($event ?? '') })"
    />
    <label class="create-manual-extract-form__label">Title <span class="create-manual-extract-form__hint">optional</span></label>
    <el-input
      :model-value="draft.manualExtract.title"
      size="small"
      placeholder="From quote if empty"
      @update:model-value="creation.updateManualExtractDraft({ title: String($event ?? '') })"
    />
    <label class="create-manual-extract-form__label">Note <span class="create-manual-extract-form__hint">optional</span></label>
    <el-input
      :model-value="draft.manualExtract.note"
      type="textarea"
      :autosize="{ minRows: 2, maxRows: 6 }"
      size="small"
      placeholder="Interpretation or context…"
      @update:model-value="creation.updateManualExtractDraft({ note: String($event ?? '') })"
    />
    <label class="create-manual-extract-form__label">Parent</label>
    <el-select
      :model-value="parentModel()"
      size="small"
      class="create-manual-extract-form__select"
      @update:model-value="setParent(String($event ?? ''))"
    >
      <el-option v-for="o in parentOptions" :key="o.value || 'root'" :label="o.label" :value="o.value" />
    </el-select>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-manual-extract-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.create-manual-extract-form__source {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.create-manual-extract-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.create-manual-extract-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}

.create-manual-extract-form__select {
  width: 100%;
}
</style>
