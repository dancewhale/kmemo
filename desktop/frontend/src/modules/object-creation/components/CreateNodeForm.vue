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
  return creation.draft.node.parentNodeId ?? ''
}

function setParent(v: string) {
  creation.updateNodeDraft({ parentNodeId: v === '' ? null : v })
}
</script>

<template>
  <div class="create-node-form">
    <label class="create-node-form__label">Title <span class="create-node-form__hint">optional</span></label>
    <el-input
      :model-value="draft.node.title"
      size="small"
      placeholder="Untitled Node"
      @update:model-value="creation.updateNodeDraft({ title: String($event ?? '') })"
    />
    <label class="create-node-form__label">Parent</label>
    <el-select :model-value="parentModel()" size="small" class="create-node-form__select" @update:model-value="setParent(String($event ?? ''))">
      <el-option v-for="o in parentOptions" :key="o.value || 'root'" :label="o.label" :value="o.value" />
    </el-select>
    <label class="create-node-form__label">Content <span class="create-node-form__hint">optional note</span></label>
    <el-input
      :model-value="draft.node.content"
      type="textarea"
      :autosize="{ minRows: 4, maxRows: 12 }"
      size="small"
      placeholder="Optional body stored on the node…"
      @update:model-value="creation.updateNodeDraft({ content: String($event ?? '') })"
    />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-node-form {
  display: flex;
  flex-direction: column;
  gap: $space-xs;
}

.create-node-form__label {
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: $color-text-secondary;
  margin-top: $space-xs;
}

.create-node-form__hint {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 400;
  opacity: 0.75;
}

.create-node-form__select {
  width: 100%;
}
</style>
