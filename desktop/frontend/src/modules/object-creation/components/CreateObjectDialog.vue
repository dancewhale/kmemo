<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useObjectCreationStore } from '../stores/object-creation.store'
import CreateObjectTypeTabs from './CreateObjectTypeTabs.vue'
import CreateNodeForm from './CreateNodeForm.vue'
import CreateCardForm from './CreateCardForm.vue'
import CreateManualExtractForm from './CreateManualExtractForm.vue'

const creation = useObjectCreationStore()
const { isDialogOpen, activeType, submitting, dialogTitle } = storeToRefs(creation)

const visible = computed({
  get: () => isDialogOpen.value,
  set: (v: boolean) => {
    if (!v) {
      creation.closeDialog()
    }
  },
})

function onCancel() {
  creation.closeDialog()
}

function onSubmit() {
  void creation.submit()
}
</script>

<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="520px"
    class="create-object-dialog"
    append-to-body
    :close-on-click-modal="true"
  >
    <div class="create-object-dialog__body">
      <CreateObjectTypeTabs />
      <div class="create-object-dialog__form">
        <CreateNodeForm v-if="activeType === 'node'" />
        <CreateCardForm v-else-if="activeType === 'card'" />
        <CreateManualExtractForm v-else />
      </div>
    </div>
    <template #footer>
      <div class="create-object-dialog__footer">
        <button type="button" class="create-object-dialog__btn create-object-dialog__btn--ghost" @click="onCancel">Cancel</button>
        <el-button type="primary" size="small" :loading="submitting" :disabled="!creation.canSubmit" @click="onSubmit">
          Create
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-object-dialog__body {
  display: flex;
  flex-direction: column;
  gap: $space-md;
}

.create-object-dialog__form {
  min-height: 200px;
}

.create-object-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: $space-sm;
}

.create-object-dialog__btn {
  font-size: $font-size-xs;
  padding: 6px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: transparent;
}

.create-object-dialog__btn--ghost:hover {
  border-color: $color-active-border;
  color: $color-text;
}
</style>
