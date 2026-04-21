<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useExtractStore } from '@/modules/extract/stores/extract.store'
import { useTreeStore } from '@/modules/knowledge-tree/stores/tree.store'
import { useToast } from '@/shared/composables/useToast'
import type { UITreeNode } from '@/modules/knowledge-tree/types'

interface ParentOption {
  value: string | null
  label: string
}

const extractStore = useExtractStore()
const treeStore = useTreeStore()
const toast = useToast()
const { pending, isCreateDialogOpen, creating } = storeToRefs(extractStore)

onMounted(async () => {
  await treeStore.initialize()
})

function flattenNodes(nodes: UITreeNode[], out: ParentOption[] = []): ParentOption[] {
  for (const n of nodes) {
    if (n.type !== 'card') {
      out.push({
        value: n.id,
        label: `${'  '.repeat(n.level)}${n.title}`,
      })
    }
    if (n.children.length > 0) {
      flattenNodes(n.children, out)
    }
  }
  return out
}

const parentOptions = computed<ParentOption[]>(() => {
  return [{ value: null, label: 'Top level' }, ...flattenNodes(treeStore.rootNodes)]
})

const canCreate = computed(() => {
  if (!pending.value) {
    return false
  }
  return pending.value.title.trim().length > 0 && pending.value.quote.trim().length > 0
})

async function onCreate() {
  const created = await extractStore.createExtract()
  if (created) {
    toast.success('Extract created')
  } else {
    toast.error('Unable to create extract')
  }
}
</script>

<template>
  <el-dialog
    :model-value="isCreateDialogOpen"
    width="620px"
    top="10vh"
    title="Create Extract"
    @close="extractStore.closeCreateDialog()"
  >
    <div v-if="pending" class="extract-dialog">
      <label class="extract-dialog__label">Title</label>
      <el-input
        :model-value="pending.title"
        size="small"
        maxlength="80"
        @update:model-value="extractStore.updatePendingField('title', String($event ?? ''))"
      />

      <label class="extract-dialog__label">Quote</label>
      <el-input
        type="textarea"
        :model-value="pending.quote"
        :autosize="{ minRows: 4, maxRows: 8 }"
        @update:model-value="extractStore.updatePendingField('quote', String($event ?? ''))"
      />

      <label class="extract-dialog__label">Parent node</label>
      <el-select
        :model-value="pending.parentNodeId"
        size="small"
        style="width: 100%"
        @update:model-value="extractStore.updatePendingField('parentNodeId', $event as string | null)"
      >
        <el-option v-for="opt in parentOptions" :key="String(opt.value)" :label="opt.label" :value="opt.value" />
      </el-select>

      <label class="extract-dialog__label">Note (optional)</label>
      <el-input
        type="textarea"
        :model-value="pending.note"
        :autosize="{ minRows: 2, maxRows: 4 }"
        @update:model-value="extractStore.updatePendingField('note', String($event ?? ''))"
      />
    </div>
    <template #footer>
      <div class="extract-dialog__footer">
        <el-button size="small" @click="extractStore.closeCreateDialog()">Cancel</el-button>
        <el-button size="small" type="primary" :loading="creating" :disabled="!canCreate" @click="onCreate">
          Create Extract
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-dialog {
  display: flex;
  flex-direction: column;
  gap: $space-sm;
}

.extract-dialog__label {
  margin-top: $space-xs;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.extract-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: $space-sm;
}
</style>
