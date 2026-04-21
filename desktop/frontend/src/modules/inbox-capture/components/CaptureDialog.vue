<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCaptureStore } from '../stores/capture.store'
import CaptureTypeTabs from './CaptureTypeTabs.vue'
import NewArticleForm from './NewArticleForm.vue'
import PasteTextCaptureForm from './PasteTextCaptureForm.vue'
import UrlCaptureForm from './UrlCaptureForm.vue'

const capture = useCaptureStore()
const { isDialogOpen, activeType, submitting, dialogTitle } = storeToRefs(capture)

const visible = computed({
  get: () => isDialogOpen.value,
  set: (v: boolean) => {
    if (!v) {
      capture.closeDialog()
    }
  },
})

function onCancel() {
  capture.closeDialog()
}

function onSubmit() {
  void capture.submit()
}
</script>

<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="520px"
    class="capture-dialog"
    append-to-body
    :close-on-click-modal="true"
  >
    <div class="capture-dialog__body">
      <CaptureTypeTabs />
      <div class="capture-dialog__form">
        <NewArticleForm v-if="activeType === 'blank'" />
        <PasteTextCaptureForm v-else-if="activeType === 'paste-text'" />
        <UrlCaptureForm v-else />
      </div>
    </div>
    <template #footer>
      <div class="capture-dialog__footer">
        <button type="button" class="capture-dialog__btn capture-dialog__btn--ghost" @click="onCancel">Cancel</button>
        <el-button type="primary" size="small" :loading="submitting" :disabled="!capture.canSubmit" @click="onSubmit">
          Add to Inbox
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.capture-dialog__body {
  display: flex;
  flex-direction: column;
  gap: $space-md;
}

.capture-dialog__form {
  min-height: 200px;
}

.capture-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: $space-sm;
}

.capture-dialog__btn {
  font-size: $font-size-xs;
  padding: 6px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: transparent;
}

.capture-dialog__btn--ghost:hover {
  border-color: $color-active-border;
  color: $color-text;
}

:deep(.capture-dialog .el-dialog__header) {
  padding-bottom: $space-sm;
  border-bottom: 1px solid $color-border-subtle;
}

:deep(.capture-dialog .el-dialog__body) {
  padding-top: $space-md;
}
</style>
